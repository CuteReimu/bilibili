package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type GetWebCookieRefreshInfoResult struct {
	Refresh   bool  `json:"refresh"`   // 是否应该刷新 Cookie。true-需要刷新，false-不需要刷新
	Timestamp int64 `json:"timestamp"` // 用于获取 refresh_csrf 的毫秒时间戳
}

// GetWebCookieRefreshInfo 获取web端cookie刷新信息
func (c *Client) GetWebCookieRefreshInfo() (*GetWebCookieRefreshInfoResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/x/passport-login/web/cookie/info"
	)

	return execute[*GetWebCookieRefreshInfoResult](c, method, url, nil)
}

type (
	GetWebCookieRefreshCsrfParam struct {
		Timestamp int64 `json:"timestamp"` // 毫秒时间戳
	}
	GetWebCookieRefreshCsrfResult struct {
		RefreshCsrf string `json:"refresh_csrf"` // 实时刷新口令
	}
)

// 正则匹配 <div id="1-name">RefreshCsrf</div> 中的刷新口令
var refreshCsrfRegex = regexp.MustCompile(`<div\s+id="1-name"\s*>(.*?)</div>`)

// GetWebCookieRefreshCsrf 获取web端cookie刷新口令
func (c *Client) GetWebCookieRefreshCsrf(param GetWebCookieRefreshCsrfParam) (*GetWebCookieRefreshCsrfResult, error) {
	correspondPath, err := getCorrespondPath(param.Timestamp)
	if err != nil {
		return nil, errors.Errorf("getCorrespondPath failed: %v", err)
	}

	url := "https://www.bilibili.com/correspond/1/" + correspondPath
	response, err := resty.New().R().SetCookies(c.resty.Cookies).Get(url)
	if err != nil || response == nil || !response.IsSuccess() {
		return nil, errors.Errorf("Request RefreshCsrf failed: %v", err)
	}

	matches := refreshCsrfRegex.FindStringSubmatch(response.String())
	if len(matches) < 2 {
		return nil, errors.Errorf("Failed to match RefreshCsrf")
	}

	return &GetWebCookieRefreshCsrfResult{RefreshCsrf: matches[1]}, nil
}

func init() {
	const publicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----
`
	pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	pubInterface, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	correspondPathPublicKey = pubInterface.(*rsa.PublicKey) //nolint:forcetypeassert // 在init中，直接panic，无需处理
}

var correspondPathPublicKey *rsa.PublicKey

// 生成CorrespondPath 算法，参数：GetWebCookieRefreshInfoResult.Timestamp
func getCorrespondPath(timestamp int64) (string, error) {
	var (
		hash   = sha256.New()
		random = rand.Reader
		msg    = []byte(fmt.Sprintf("refresh_%d", timestamp))
	)
	encryptedData, err := rsa.EncryptOAEP(hash, random, correspondPathPublicKey, msg, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(encryptedData), nil
}
