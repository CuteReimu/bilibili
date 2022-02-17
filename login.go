package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strings"
)

type CaptchaResult struct {
	Code    int      // 返回值，0是成功
	Message string   // 返回信息
	TTL     int      `json:"ttl"` // 固定值1，作用尚不明确
	Data    struct { // 信息本体
		Geetest struct {
			Gt        string // 极验id
			Challenge string // 极验KEY
		}
		Tencent struct { // 作用不明确
			Appid string
		}
		Token string // 极验token
		Type  string // 验证方式
	}
}

// Captcha 申请验证码参数
func Captcha() (*CaptchaResult, error) {
	return std.Captcha()
}
func (c *Client) Captcha() (*CaptchaResult, error) {
	resp, err := c.resty().R().SetQueryParam("source", "main_web").Get("https://passport.bilibili.com/web/captcha/combine")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("status code: %d", resp.StatusCode())
	}
	var result *CaptchaResult
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, err
}

func encrypt(publicKey, data string) (string, error) {
	// pem解码
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("failed to decode public key")
	}
	// x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", errors.WithStack(err)
	}
	pk := publicKeyInterface.(*rsa.PublicKey)
	// 加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte(data))
	if err != nil {
		return "", errors.WithStack(err)
	}
	// base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// LoginWithPassword 账号密码登录
func LoginWithPassword(userName, password, key, challenge, validate, seccode string) error {
	return std.LoginWithPassword(userName, password, key, challenge, validate, seccode)
}
func (c *Client) LoginWithPassword(userName, password, key, challenge, validate, seccode string) error {
	client := c.resty()
	resp, err := client.R().SetQueryParam("act", "getkey").Get("https://passport.bilibili.com/login")
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return errors.Errorf("获取密码盐值失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return errors.Errorf("json invalid: %s", resp.String())
	}
	getKeyResp := gjson.ParseBytes(resp.Body())
	encryptPwd, err := encrypt(getKeyResp.Get("key").String(), getKeyResp.Get("hash").String()+password)
	if err != nil {
		return errors.WithStack(err)
	}
	resp, err = client.R().SetQueryParams(map[string]string{
		"captchaType": "6",
		"username":    userName,
		"password":    encryptPwd,
		"keep":        "true",
		"key":         key,
		"challenge":   challenge,
		"validate":    validate,
		"seccode":     seccode,
	}).Post("https://passport.bilibili.com/web/login/v2")
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return errors.New("登录bilibili失败")
	}
	c.cookies = resp.Cookies()
	var cookieStrings []string
	for _, cookie := range c.cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}
	c.cookiesString = strings.Join(cookieStrings, "\n")
	return nil
}
