package bilibili

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	resty *resty.Client
}

// New 返回一个默认的 bilibili.Client
func New() *Client {
	restyClient := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetTimeout(20*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Origin", "https://www.bilibili.com").
		SetHeader("Referer", "https://www.bilibili.com/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	return NewWithClient(restyClient)
}

// NewWithClient 接收一个自定义的*resty.Client为参数
func NewWithClient(restyClient *resty.Client) *Client {
	return &Client{resty: restyClient}
}

func (c *Client) Resty() *resty.Client {
	return c.resty
}

// GetCookiesString 获取字符串格式的cookies，方便自行存储后下次使用。配合下面的 SetCookiesString 使用。
func (c *Client) GetCookiesString() string {
	cookies := c.resty.Cookies
	cookieStrings := make([]string, 0, len(cookies))
	for _, cookie := range c.resty.Cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}
	return strings.Join(cookieStrings, "\n")
}

// SetCookiesString 设置Cookies，但是是字符串格式，配合 GetCookiesString 使用。有些功能必须登录或设置Cookies后才能使用。
func (c *Client) SetCookiesString(cookiesString string) {
	c.resty.SetCookies((&resty.Response{RawResponse: &http.Response{Header: http.Header{
		"Set-Cookie": strings.Split(cookiesString, "\n"),
	}}}).Cookies())
}

// 根据key获取指定的cookie值
func (c *Client) getCookie(name string) string {
	now := time.Now()
	// 查找指定name的cookie
	for _, cookie := range c.resty.Cookies {
		if cookie.Name == name && (cookie.Expires.IsZero() || cookie.Expires.After(now)) {
			return cookie.Value
		}
	}
	return ""
}

func formatError(prefix string, code int64, message ...string) error {
	for _, m := range message {
		if len(m) > 0 {
			return errors.New(prefix + "失败，返回值：" + strconv.FormatInt(code, 10) + "，返回信息：" + m)
		}
	}
	return errors.New(prefix + "失败，返回值：" + strconv.FormatInt(code, 10))
}

func getRespData(resp *resty.Response, prefix string) ([]byte, error) {
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf(prefix+"失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return nil, errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return nil, formatError(prefix, code, res.Get("message").String(), res.Get("msg").String())
	}
	return []byte(res.Get("data").Raw), nil
}
