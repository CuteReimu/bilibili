package bilibili

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

const (
	cacheImgKey = "imgKey"
	cacheSubKey = "subKey"
)

var (
	_defaultMixinKeyEncTab = []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}

	_defaultStorage = &MemoryStorage{
		data: make(map[string]interface{}, 15),
	}
)

type Storage interface {
	Set(key string, value interface{})
	Get(key string) (v interface{}, isSet bool)
}

type MemoryStorage struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func (impl *MemoryStorage) Set(key string, value interface{}) {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	impl.data[key] = value
}

func (impl *MemoryStorage) Get(key string) (v interface{}, isSet bool) {
	impl.mu.RLock()
	defer impl.mu.RUnlock()

	if v, isSet = impl.data[key]; isSet {
		return v, true
	}
	return nil, false
}

// WBI 签名实现
// 如果希望以登录的方式获取则使用 WithCookies or WithRawCookies 设置cookie
// 如果希望以未登录的方式获取 WithCookies(nil) 设置cookie为 nil 即可, 这是 Default 行为
//
//	!!! 使用 WBI 的接口 绝对不可以 set header Referer 会导致失败 !!!
//	!!! 大部分使用 WBI 的接口都需要 set header Cookie !!!
//
// see https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/misc/sign/wbi.md
type WBI struct {
	cookies        []*http.Cookie
	mixinKeyEncTab []int

	// updateCheckerInterval is the interval to check and update wbi keys
	// default is 60 minutes. so if lastInitTime + updateCheckerInterval < now, it will update wbi keys
	updateCheckerInterval time.Duration
	lastInitTime          time.Time
	storage               Storage

	sfg singleflight.Group
}

func NewDefaultWbi() *WBI {
	return &WBI{
		cookies:        nil,
		mixinKeyEncTab: _defaultMixinKeyEncTab,

		updateCheckerInterval: 60 * time.Minute,
		storage:               _defaultStorage,
	}
}

func (wbi *WBI) WithUpdateInterval(updateInterval time.Duration) *WBI {
	wbi.updateCheckerInterval = updateInterval
	return wbi
}

func (wbi *WBI) WithCookies(cookies []*http.Cookie) *WBI {
	wbi.cookies = cookies
	return wbi
}

func (wbi *WBI) WithRawCookies(rawCookies string) *WBI {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	req := http.Request{Header: header}

	wbi.cookies = req.Cookies()
	return wbi
}

func (wbi *WBI) WithMixinKeyEncTab(mixinKeyEncTab []int) *WBI {
	wbi.mixinKeyEncTab = mixinKeyEncTab
	return wbi
}

func (wbi *WBI) WithStorage(storage Storage) *WBI {
	wbi.storage = storage
	return wbi
}

func (wbi *WBI) GetKeys() (imgKey string, subKey string, err error) {
	imgKey, subKey = wbi.getKeys()

	// 更新检查
	if imgKey == "" || subKey == "" || time.Since(wbi.lastInitTime) > wbi.updateCheckerInterval {
		if err = wbi.initWbi(); err != nil {
			return "", "", err
		}

		return wbi.GetKeys()
	}

	return imgKey, subKey, nil
}

func (wbi *WBI) getKeys() (imgKey string, subKey string) {
	if v, isSet := wbi.storage.Get(cacheImgKey); isSet {
		imgKey, _ = v.(string)
	}

	if v, isSet := wbi.storage.Get(cacheSubKey); isSet {
		subKey, _ = v.(string)
	}

	return imgKey, subKey
}

func (wbi *WBI) SetKeys(imgKey, subKey string) {
	wbi.storage.Set(cacheImgKey, imgKey)
	wbi.storage.Set(cacheSubKey, subKey)
	wbi.lastInitTime = time.Now()
}

func (wbi *WBI) GetMixinKey() (string, error) {
	imgKey, subKey, err := wbi.GetKeys()
	if err != nil {
		return "", err
	}

	return wbi.GenerateMixinKey(imgKey + subKey), nil
}

func (wbi *WBI) GenerateMixinKey(orig string) string {
	var str strings.Builder
	for _, v := range wbi.mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
		}
	}
	return str.String()[:32]
}

func (wbi *WBI) sanitizeString(s string) string {
	unwantedChars := []string{"!", "'", "(", ")", "*"}
	for _, char := range unwantedChars {
		s = strings.ReplaceAll(s, char, "")
	}
	return s
}

func (wbi *WBI) SignQuery(query url.Values, ts time.Time) (newQuery url.Values, err error) {
	payload := make(map[string]string, 10)
	for k := range query {
		payload[k] = query.Get(k)
	}

	newPayload, err := wbi.SignMap(payload, ts)
	if err != nil {
		return query, err
	}

	newQuery = url.Values{}
	for k, v := range newPayload {
		newQuery.Set(k, v)
	}

	return newQuery, nil
}

func (wbi *WBI) SignMap(payload map[string]string, ts time.Time) (newPayload map[string]string, err error) {
	newPayload = make(map[string]string, 10)
	for k, v := range payload {
		newPayload[k] = v
	}

	newPayload["wts"] = strconv.FormatInt(ts.Unix(), 10)

	// Sort keys
	keys := make([]string, 0, 10)
	for k := range newPayload {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// Remove unwanted characters
	for k, v := range newPayload {
		v = wbi.sanitizeString(v)
		newPayload[k] = v
	}

	// Build URL parameters
	signQuery := url.Values{}
	for _, k := range keys {
		signQuery.Set(k, newPayload[k])
	}
	signQueryStr := signQuery.Encode()

	// Get mixin key
	mixinKey, err := wbi.GetMixinKey()
	if err != nil {
		return payload, err
	}

	// Calculate w_rid
	hash := md5.Sum([]byte(signQueryStr + mixinKey))
	newPayload["w_rid"] = hex.EncodeToString(hash[:])

	return newPayload, nil
}

func (wbi *WBI) initWbi() error {
	_, err, _ := wbi.sfg.Do("initWbi", func() (interface{}, error) {
		return nil, wbi.doInitWbi()
	})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (wbi *WBI) doInitWbi() error {
	result := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			WbiImg struct {
				ImgUrl string `json:"img_url"`
				SubUrl string `json:"sub_url"`
			} `json:"wbi_img"`
		}
	}{}

	resp, err := resty.New().R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Origin", "https://www.bilibili.com").
		SetHeader("Referer", "https://www.bilibili.com/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0").
		SetCookies(wbi.cookies).
		SetResult(&result).
		Get("https://api.bilibili.com/x/web-interface/nav")

	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("status code: %d", resp.StatusCode())
	}

	if result.Code != 0 {
		if result.Data.WbiImg.ImgUrl == "" || result.Data.WbiImg.SubUrl == "" {
			return errors.Errorf("init wbi 失败, 错误码: %d, 错误信息: %s", result.Code, result.Message)
		}
	}

	if len(resp.Cookies()) > 0 {
		// update cookie
		wbi.cookies = resp.Cookies()
	}

	imgKey := strings.Split(strings.Split(result.Data.WbiImg.ImgUrl, "/")[len(strings.Split(result.Data.WbiImg.ImgUrl, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(result.Data.WbiImg.SubUrl, "/")[len(strings.Split(result.Data.WbiImg.SubUrl, "/"))-1], ".")[0]

	wbi.SetKeys(imgKey, subKey)
	return nil
}
