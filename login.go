package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strconv"
	"time"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/tidwall/gjson"
)

type CaptchaResult struct {
	Geetest struct {
		Gt        string `json:"gt"`        // 极验id
		Challenge string `json:"challenge"` // 极验KEY
	} `json:"geetest"`
	Tencent struct { // 作用不明确
		Appid string `json:"appid"`
	} `json:"tencent"`
	Token string `json:"token"` // 极验token
	Type  string `json:"type"`  // 验证方式
}

// Captcha 申请验证码参数
func Captcha() (*CaptchaResult, error) {
	return std.Captcha()
}
func (c *Client) Captcha() (*CaptchaResult, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("source", "main_web").Get("https://passport.bilibili.com/x/passport-login/captcha")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "申请验证码参数")
	if err != nil {
		return nil, err
	}
	var result *CaptchaResult
	err = json.Unmarshal(data, &result)
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
	resp, err := client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("act", "getkey").Get("https://passport.bilibili.com/login")
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
	resp, err = client.R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"source":    "main_web",
		"username":  userName,
		"password":  encryptPwd,
		"keep":      "true",
		"token":     captchaResult.Token,
		"challenge": captchaResult.Geetest.Challenge,
		"validate":  validate,
		"seccode":   seccode,
	}).Post("https://passport.bilibili.com/x/passport-login/web/login")
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
		return errors.Errorf("登录bilibili失败，错误码：%d, 错误信息：%s", code, gjson.GetBytes(resp.Body(), "message").String())
	}
	status := gjson.GetBytes(resp.Body(), "data.status").Int()
	if status != 0 {
		return errors.Errorf("登录bilibili失败，错误码：%d，状态码：%d, 错误信息：%s", code, status, gjson.GetBytes(resp.Body(), "data.message").String())
	}
	c.SetCookies(resp.Cookies())
	return nil
}

type CountryInfo struct {
	Id        int    `json:"id"`         // 国际代码值
	Cname     string `json:"cname"`      // 国家或地区名
	CountryId string `json:"country_id"` // 国家或地区区号
}

// ListCountry 获取国际地区代码
func ListCountry() (common []CountryInfo, others []CountryInfo, err error) {
	return std.ListCountry()
}
func (c *Client) ListCountry() (common []CountryInfo, others []CountryInfo, err error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Get("https://passport.bilibili.com/web/generic/country/list")
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, nil, errors.Errorf("获取国际地区代码失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return nil, nil, errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return nil, nil, errors.Errorf("获取国际地区代码失败，返回值：%d", code)
	}
	err = json.Unmarshal([]byte(res.Get("data.common").Raw), &common)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	err = json.Unmarshal([]byte(res.Get("data.others").Raw), &others)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return common, others, nil
}

// SendSMS 发送短信验证码
func SendSMS(tel, cid int, captchaResult *CaptchaResult, validate, seccode string) (captchaKey string, err error) {
	return std.SendSMS(tel, cid, captchaResult, validate, seccode)
}
func (c *Client) SendSMS(tel, cid int, captchaResult *CaptchaResult, validate, seccode string) (captchaKey string, err error) {
	if captchaResult == nil {
		return "", errors.New("请先进行极验人机验证")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"tel":       strconv.Itoa(tel),
		"cid":       strconv.Itoa(cid),
		"source":    "main_web",
		"token":     captchaResult.Token,
		"challenge": captchaResult.Geetest.Challenge,
		"validate":  validate,
		"seccode":   seccode,
	}).Post("https://passport.bilibili.com/x/passport-login/web/sms/send")
	if err != nil {
		return "", errors.WithStack(err)
	}
	data, err := getRespData(resp, "发送短信验证码")
	if err != nil {
		return "", err
	}
	return gjson.GetBytes(data, "captcha_key").String(), nil
}

// LoginWithSMS 使用短信验证码登录
func LoginWithSMS(tel, cid, code int, captchaKey string) error {
	return std.LoginWithSMS(tel, cid, code, captchaKey)
}
func (c *Client) LoginWithSMS(tel, cid, code int, captchaKey string) error {
	if len(captchaKey) == 0 {
		return errors.New("请先发送短信")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"cid":         strconv.Itoa(cid),
		"tel":         strconv.Itoa(tel),
		"code":        strconv.Itoa(code),
		"source":      "main_web",
		"captcha_key": captchaKey,
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
	status := gjson.GetBytes(resp.Body(), "data.status").Int()
	if status != 0 {
		return errors.Errorf("登录bilibili失败，错误码：%d，状态码：%d, 错误信息：%s", retCode, status, gjson.GetBytes(resp.Body(), "data.message").String())
	}
	c.SetCookies(resp.Cookies())
	return nil
}

type QRCode struct {
	Url       string `json:"url"`        // 二维码内容url
	QrcodeKey string `json:"qrcode_key"` // 扫码登录秘钥
}

// Encode a QRCode and return a raw PNG image.
func (result *QRCode) Encode() ([]byte, error) {
	return qrcode.Encode(result.Url, qrcode.Medium, 256)
}

// Print the QRCode in the console
func (result *QRCode) Print() {
	front, back := qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite
	qrcodeTerminal.New2(front, back, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(result.Url).Print()
}

// GetQRCode 申请二维码URL及扫码密钥
func GetQRCode() (*QRCode, error) {
	return std.GetQRCode()
}
func (c *Client) GetQRCode() (*QRCode, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Get("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "申请二维码")
	if err != nil {
		return nil, err
	}
	var result *QRCode
	err = json.Unmarshal(data, &result)
	return result, errors.WithStack(err)
}

// LoginWithQRCode 使用扫码登录
func LoginWithQRCode(qrCode *QRCode) error {
	return std.LoginWithQRCode(qrCode)
}
func (c *Client) LoginWithQRCode(qrCode *QRCode) error {
	if qrCode == nil {
		return errors.New("请先获取二维码")
	}

	for {
		ok, err := c.qrCodeSuccess(qrCode)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		time.Sleep(3 * time.Second) // 主站 3s 一次请求
	}
}

// qrCodeSuccess 是否扫码成功
func (c *Client) qrCodeSuccess(qrCode *QRCode) (bool, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("qrcode_key", qrCode.QrcodeKey).Get("https://passport.bilibili.com/x/passport-login/web/qrcode/poll")
	if err != nil {
		return false, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return false, errors.New("登录bilibili失败")
	}
	if !gjson.ValidBytes(resp.Body()) {
		return false, errors.Errorf("json invalid: %s", resp.String())
	}
	result := gjson.ParseBytes(resp.Body())
	retCode := result.Get("code").Int()
	if retCode != 0 {
		return false, errors.Errorf("登录bilibili失败，错误码：%d，错误信息：%s", retCode, gjson.GetBytes(resp.Body(), "message").String())
	} else {
		codeValue := result.Get("data.code")
		if !codeValue.Exists() || codeValue.Type != gjson.Number {
			return false, errors.New("扫码登录未成功，返回异常")
		}
		code := codeValue.Int()
		switch code {
		case 86038: // 二维码已失效
			return false, errors.New("扫码登录未成功，原因：二维码已失效")
		case 86090: // 二维码已扫码未确认
			return false, nil
		case 86101: // 未扫码
			return false, nil
		case 0:
			c.SetCookies(resp.Cookies())
			return true, nil
		default:
			return false, errors.Errorf("由于未知原因，扫码登录未成功，错误码：%d", code)
		}
	}
}
