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

// GetWebCookieRefreshCsrf 获取web端cookie刷新口令
func (c *Client) GetWebCookieRefreshCsrf(param GetWebCookieRefreshCsrfParam) (*GetWebCookieRefreshCsrfResult, error) {
	correspondPath, err := getCorrespondPath(param.Timestamp)
	if err != nil {
		return nil, errors.Errorf("getCorrespondPath falied: %v", err)
	}

	url := fmt.Sprintf("https://www.bilibili.com/correspond/1/%s", correspondPath)
	response, err := resty.New().R().SetCookies(c.resty.Cookies).Get(url)
	if err != nil || response == nil || !response.IsSuccess() {
		return nil, errors.Errorf("Request RefreshCsrf failed: %v", err)
	}

	// 正则匹配 <div id="1-name">RefreshCsrf</div> 中的刷新口令
	re := regexp.MustCompile(`<div\s+id="1-name"\s*>(.*?)</div>`)
	matches := re.FindStringSubmatch(response.String())
	if len(matches) < 2 {
		return nil, errors.Errorf("Failed to match RefreshCsrf")
	}

	return &GetWebCookieRefreshCsrfResult{RefreshCsrf: matches[1]}, nil
}

// 生成CorrespondPath 算法，参数：GetWebCookieRefreshInfoResult.Timestamp
func getCorrespondPath(timestamp int64) (string, error) {
	const publicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----
`
	pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	msg := []byte(fmt.Sprintf("refresh_%d", timestamp))
	var (
		hash   = sha256.New()
		random = rand.Reader
		pub    *rsa.PublicKey
	)
	pubInterface, parseErr := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if parseErr != nil {
		return "", parseErr
	}

	pub = pubInterface.(*rsa.PublicKey)
	encryptedData, encryptErr := rsa.EncryptOAEP(hash, random, pub, msg, nil)
	if encryptErr != nil {
		return "", encryptErr
	}

	return hex.EncodeToString(encryptedData), nil
}
