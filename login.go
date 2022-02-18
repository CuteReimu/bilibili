package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/tidwall/gjson"
	"strconv"
)

type CaptchaResult struct {
	Code    int      `json:"code"`    // 返回值，0是成功
	Message string   `json:"message"` // 返回信息
	TTL     int      `json:"ttl"`     // 固定值1，作用尚不明确
	Data    struct { // 信息本体
		Geetest struct {
			Gt        string `json:"gt"`        // 极验id
			Challenge string `json:"challenge"` // 极验KEY
		} `json:"geetest"`
		Tencent struct { // 作用不明确
			Appid string `json:"appid"`
		} `json:"tencent"`
		Token string `json:"token"` // 极验token
		Type  string `json:"type"`  // 验证方式
	} `json:"data"`
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
		return nil, errors.Errorf("申请验证码参数失败，status code: %d", resp.StatusCode())
	}
	var result *CaptchaResult
	err = json.Unmarshal(resp.Body(), &result)
	return result, errors.WithStack(err)
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
func LoginWithPassword(userName, password string, captchaResult *CaptchaResult, validate, seccode string) error {
	return std.LoginWithPassword(userName, password, captchaResult, validate, seccode)
}
func (c *Client) LoginWithPassword(userName, password string, captchaResult *CaptchaResult, validate, seccode string) error {
	if captchaResult == nil {
		return errors.New("请先进行极验人机验证")
	}
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
		"key":         captchaResult.Data.Token,
		"challenge":   captchaResult.Data.Geetest.Challenge,
		"validate":    validate,
		"seccode":     seccode,
	}).Post("https://passport.bilibili.com/web/login/v2")
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return errors.New("登录bilibili失败")
	}
	if !gjson.ValidBytes(resp.Body()) {
		return errors.Errorf("json invalid: %s", resp.String())
	}
	code := gjson.GetBytes(resp.Body(), "code").Int()
	if code != 0 {
		return errors.Errorf("登录bilibili失败，错误码：%d，错误信息：%s", code, gjson.GetBytes(resp.Body(), "message").String())
	}
	c.SetCookies(resp.Cookies())
	return nil
}

type CountryInfo struct {
	Id        int    `json:"id"`         // 国际代码值
	Cname     string `json:"cname"`      // 国家或地区名
	CountryId string `json:"country_id"` // 国家或地区区号
}

type ListCountryResult struct {
	Code int      `json:"code"` // 返回值，0表示成功
	Data struct { // 数据本体
		Common []*CountryInfo `json:"common"` // 常用国家或地区
		Others []*CountryInfo `json:"others"` // 其他国家或地区
	} `json:"data"`
}

// ListCountry 获取国际地区代码
func ListCountry() (*ListCountryResult, error) {
	return std.ListCountry()
}
func (c *Client) ListCountry() (*ListCountryResult, error) {
	resp, err := c.resty().R().Get("https://passport.bilibili.com/web/generic/country/list")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取国际地区代码失败，status code: %d", resp.StatusCode())
	}
	var result *ListCountryResult
	err = json.Unmarshal(resp.Body(), &result)
	return result, errors.WithStack(err)
}

type SendSMSResult struct {
	Code    int      `json:"code"`    // 返回值，0表示成功
	Message string   `json:"message"` // 错误信息
	Data    struct { // 数据
		CaptchaKey string `json:"captcha_key"`
	} `json:"data"`
}

// SendSMS 发送短信验证码
func SendSMS(tel, cid int, captchaResult *CaptchaResult, validate, seccode string) (*SendSMSResult, error) {
	return std.SendSMS(tel, cid, captchaResult, validate, seccode)
}
func (c *Client) SendSMS(tel, cid int, captchaResult *CaptchaResult, validate, seccode string) (*SendSMSResult, error) {
	if captchaResult == nil {
		return nil, errors.New("请先进行极验人机验证")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"tel":       strconv.Itoa(tel),
		"cid":       strconv.Itoa(cid),
		"source":    "main_web",
		"token":     captchaResult.Data.Token,
		"challenge": captchaResult.Data.Geetest.Challenge,
		"validate":  validate,
		"seccode":   seccode,
	}).Post("https://passport.bilibili.com/x/passport-login/web/sms/send")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("发送短信验证码失败，status code: %d", resp.StatusCode())
	}
	var result *SendSMSResult
	err = json.Unmarshal(resp.Body(), &result)
	return result, errors.WithStack(err)
}

// LoginWithSMS 使用短信验证码登录
func LoginWithSMS(tel, cid, code int, sendSMSResult *SendSMSResult) error {
	return std.LoginWithSMS(tel, cid, code, sendSMSResult)
}
func (c *Client) LoginWithSMS(tel, cid, code int, sendSMSResult *SendSMSResult) error {
	if sendSMSResult == nil {
		return errors.New("请先发送短信")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"cid":         strconv.Itoa(cid),
		"tel":         strconv.Itoa(tel),
		"code":        strconv.Itoa(code),
		"source":      "main_web",
		"captcha_key": sendSMSResult.Data.CaptchaKey,
	}).Post("https://passport.bilibili.com/x/passport-login/web/login/sms")
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return errors.New("登录bilibili失败")
	}
	if !gjson.ValidBytes(resp.Body()) {
		return errors.Errorf("json invalid: %s", resp.String())
	}
	retCode := gjson.GetBytes(resp.Body(), "code").Int()
	if retCode != 0 {
		return errors.Errorf("登录bilibili失败，错误码：%d，错误信息：%s", retCode, gjson.GetBytes(resp.Body(), "message").String())
	}
	c.SetCookies(resp.Cookies())
	return nil
}

type GetQRCodeResult struct {
	Code   int      `json:"code"`   // 返回值，0表示成功
	Status bool     `json:"status"` // 作用尚不明确
	Ts     uint32   `json:"ts"`     // 请求时间戳
	Data   struct { // 信息本体
		Url      string `json:"url"`       // 二维码内容url
		OauthKey string `json:"oauth_key"` // 扫码登录秘钥
	} `json:"data"`
}

// Encode a QRCode and return a raw PNG image.
func (result *GetQRCodeResult) Encode() ([]byte, error) {
	return qrcode.Encode(result.Data.Url, qrcode.Medium, 256)
}

// GetQRCode 申请二维码URL及扫码密钥
func GetQRCode() (*GetQRCodeResult, error) {
	return std.GetQRCode()
}
func (c *Client) GetQRCode() (*GetQRCodeResult, error) {
	resp, err := c.resty().R().Get("https://passport.bilibili.com/qrcode/getLoginUrl")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.New("申请二维码失败")
	}
	var result *GetQRCodeResult
	err = json.Unmarshal(resp.Body(), &result)
	return result, errors.WithStack(err)
}

// LoginWithQRCode 使用扫码登录
func LoginWithQRCode(getQRCodeResult *GetQRCodeResult) error {
	return std.LoginWithQRCode(getQRCodeResult)
}
func (c *Client) LoginWithQRCode(getQRCodeResult *GetQRCodeResult) error {
	if getQRCodeResult == nil {
		return errors.New("请先获取二维码")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"oauthKey": getQRCodeResult.Data.OauthKey,
		"gourl":    "https://www.bilibili.com",
	}).Post("https://passport.bilibili.com/qrcode/getLoginInfo")
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return errors.New("登录bilibili失败")
	}
	if !gjson.ValidBytes(resp.Body()) {
		return errors.Errorf("json invalid: %s", resp.String())
	}
	result := gjson.ParseBytes(resp.Body())
	retCode := result.Get("code").Int()
	if retCode != 0 {
		return errors.Errorf("登录bilibili失败，错误码：%d，错误信息：%s", retCode, gjson.GetBytes(resp.Body(), "message").String())
	} else if !result.Get("status").Bool() {
		switch result.Get("data").Int() {
		case -1:
			return errors.New("扫码登录未成功，原因：密钥错误")
		case -2:
			return errors.New("扫码登录未成功，原因：密钥超时")
		case -4:
			return errors.New("扫码登录未成功，原因：未扫描")
		case -5:
			return errors.New("扫码登录未成功，原因：未确认")
		}
		return errors.New("由于未知原因，扫码登录未成功")
	}
	c.SetCookies(resp.Cookies())
	return nil
}
