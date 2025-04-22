package bilibili

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"time"
)

type Geetest struct {
	Gt        string `json:"gt"`        // 极验id。一般为固定值
	Challenge string `json:"challenge"` // 极验KEY。由B站后端产生用于人机验证
}

type CaptchaResult struct {
	Geetest Geetest `json:"geetest"` // 极验captcha数据
	Tencent any     `json:"tencent"` // (?)。**作用尚不明确**
	Token   string  `json:"token"`   // 登录 API token。与 captcha 无关，与登录接口有关
	Type    string  `json:"type"`    // 验证方式。用于判断使用哪一种验证方式，目前所见只有极验。geetest：极验
}

// Captcha 申请验证码参数
func (c *Client) Captcha() (*CaptchaResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/x/passport-login/captcha"
	)
	return execute[*CaptchaResult](c, method, url, nil, fillParam("source", "main_web"))
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
	pk, _ := publicKeyInterface.(*rsa.PublicKey)
	// 加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte(data))
	if err != nil {
		return "", errors.WithStack(err)
	}
	// base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

type LoginWithPasswordParam struct {
	Username  string `json:"username"`                                   // 用户登录账号。手机号或邮箱地址
	Password  string `json:"password"`                                   // 参数传入原密码，下文会自动转为加密后的带盐密码
	Keep      int    `json:"keep"`                                       // 0
	Token     string `json:"token"`                                      // 登录 API token。使用 Captcha() 方法获取
	Challenge string `json:"challenge"`                                  // 极验 challenge。使用 Captcha() 方法获取
	Validate  string `json:"validate"`                                   // 极验 result。极验验证后得到
	Seccode   string `json:"seccode"`                                    // 极验 result +jordan。极验验证后得到
	GoUrl     string `json:"go_url,omitempty" request:"query,omitempty"` // 跳转 url。默认为 https://www.bilibili.com
	Source    string `json:"source,omitempty" request:"query,omitempty"` // 登录来源。main_web：独立登录页。main_mini：小窗登录
}

type LoginWithPasswordResult struct {
	Message      string `json:"message"`       // 扫码状态信息
	RefreshToken string `json:"refresh_token"` // 刷新refresh_token
	Status       int    `json:"status"`        // 成功为0
	Timestamp    int    `json:"timestamp"`     // 登录时间。未登录为0。时间戳 单位为毫秒
	Url          string `json:"url"`           // 游戏分站跨域登录 url
}

// LoginWithPassword 账号密码登录，其中validate, seccode字段需要在极验人机验证后获取
func (c *Client) LoginWithPassword(param LoginWithPasswordParam) (*LoginWithPasswordResult, error) {
	type getKeyResult struct {
		Hash string `json:"hash"` // 密码盐值。有效时间为 20s。恒为 16 字符。需要拼接在明文密码之前
		Key  string `json:"key"`  // rsa 公钥。PEM 格式编码。加密密码时需要使用
	}

	// 获取密码盐值
	const (
		method1 = resty.MethodGet
		url1    = "https://passport.bilibili.com/x/passport-login/web/key"
	)
	getKey, err := execute[*getKeyResult](c, method1, url1, nil)
	if err != nil {
		return nil, err
	}

	// 密码加盐
	encryptPwd, err := encrypt(getKey.Key, getKey.Hash+param.Password)
	if err != nil {
		return nil, err
	}
	param.Password = encryptPwd

	// 登录操作(web端)
	const (
		method2 = resty.MethodPost
		url2    = "https://passport.bilibili.com/x/passport-login/web/login"
	)
	return execute[*LoginWithPasswordResult](c, method2, url2, param)
}

type CountryCrown struct {
	Id        int    `json:"id"`         // 国际代码值
	Cname     string `json:"cname"`      // 国家或地区名
	CountryId string `json:"country_id"` // 国家或地区区号
}

type GetCountryCrownResult struct {
	Common []CountryCrown `json:"common"` // 常用国家&地区
	Others []CountryCrown `json:"others"` // 其他国家&地区
}

// GetCountryCrown 获取国际冠字码
func (c *Client) GetCountryCrown() (*GetCountryCrownResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/web/generic/country/list"
	)
	return execute[*GetCountryCrownResult](c, method, url, nil)
}

type SendSMSParam struct {
	Cid       int    `json:"cid"`       // 国际冠字码。使用 GetCountryCrown() 方法获取
	Tel       int    `json:"tel"`       // 手机号码
	Source    string `json:"source"`    // 登录来源。main_web：独立登录页。main_mini：小窗登录
	Token     string `json:"token"`     // 登录 API token。使用 Captcha() 方法获取
	Challenge string `json:"challenge"` // 极验 challenge。使用 Captcha() 方法获取
	Validate  string `json:"validate"`  // 极验 result。极验验证后得到
	Seccode   string `json:"seccode"`   // 极验 result +jordan。极验验证后得到
}

type SendSMSResult struct {
	CaptchaKey string `json:"captcha_key"` // 短信登录 token。在下方传参时需要，请备用
}

// SendSMS 发送短信验证码
func (c *Client) SendSMS(param SendSMSParam) (*SendSMSResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://passport.bilibili.com/x/passport-login/web/sms/send"
	)
	return execute[*SendSMSResult](c, method, url, param)
}

type LoginWithSMSParam struct {
	Cid        int    `json:"cid"`                                        // 国际冠字码。可以从 GetCountryCrown() 获取
	Tel        int    `json:"tel"`                                        // 手机号码
	Code       int    `json:"code"`                                       // 短信验证码。timeout 为 5min
	Source     string `json:"source"`                                     // 登录来源。main_web：独立登录页。main_mini：小窗登录
	CaptchaKey string `json:"captcha_key"`                                // 短信登录 token。从 SendSMS() 请求成功后返回
	GoUrl      string `json:"go_url,omitempty" request:"query,omitempty"` // 跳转url。默认为 https://www.bilibili.com
	Keep       bool   `json:"keep,omitempty" request:"query,omitempty"`   // 是否记住登录。true：记住登录。false：不记住登录
}

type LoginWithSMSResult struct {
	IsNew  bool   `json:"is_new"` // 是否为新注册用户。false：非新注册用户。true：新注册用户
	Status int    `json:"status"` // 0。未知，可能0就是成功吧
	Url    string `json:"url"`    // 跳转 url。默认为 https://www.bilibili.com
}

// LoginWithSMS 使用短信验证码登录
func (c *Client) LoginWithSMS(param LoginWithSMSParam) (*LoginWithSMSResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://passport.bilibili.com/x/passport-login/web/login/sms"
	)
	return execute[*LoginWithSMSResult](c, method, url, param)
}

type QRCode struct {
	Url       string `json:"url"`        // 二维码内容 (登录页面 url)
	QrcodeKey string `json:"qrcode_key"` // 扫码登录秘钥。恒为32字符
}

// Encode a QRCode and return a raw PNG image.
func (result *QRCode) Encode() ([]byte, error) {
	buf, err := qrcode.Encode(result.Url, qrcode.Medium, 256)
	return buf, errors.WithStack(err)
}

// Print the QRCode in the console
func (result *QRCode) Print() {
	front, back := qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite
	qrcodeTerminal.New2(front, back, qrcodeTerminal.QRCodeRecoveryLevels.Low).Get(result.Url).Print()
}

// GetQRCode 申请二维码
func (c *Client) GetQRCode() (*QRCode, error) {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	)
	return execute[*QRCode](c, method, url, nil)
}

type LoginWithQRCodeParam struct {
	QrcodeKey string `json:"qrcode_key"` // 扫码登录秘钥
}

type LoginWithQRCodeResult struct {
	Url          string `json:"url"`           // 游戏分站跨域登录 url。未登录为空
	RefreshToken string `json:"refresh_token"` // 刷新refresh_token。未登录为空
	Timestamp    int    `json:"timestamp"`     // 登录时间。未登录为0。时间戳 单位为毫秒
	Code         int    `json:"code"`          // 0：扫码登录成功。86038：二维码已失效。86090：二维码已扫码未确认。86101：未扫码
	Message      string `json:"message"`       // 扫码状态信息
}

// LoginWithQRCode 使用扫码登录。
//
// 该方法会阻塞直到扫码成功或者已经无法扫码。
func (c *Client) LoginWithQRCode(param LoginWithQRCodeParam) (*LoginWithQRCodeResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"
	)
	for {
		result, err := execute[*LoginWithQRCodeResult](c, method, url, param)
		if err != nil {
			return nil, err
		}
		if result.Code != 86090 && result.Code != 86101 {
			// 86090：二维码已扫码未确认
			// 86101：未扫码
			return result, nil
		}
		time.Sleep(3 * time.Second) // 主站 3s 一次请求
	}
}

type AccountInformation struct {
	Mid      int    `json:"mid"`       // 我的mid
	Uname    string `json:"uname"`     // 我的昵称
	Userid   string `json:"userid"`    // 我的用户名
	Sign     string `json:"sign"`      // 我的签名
	Birthday string `json:"birthday"`  // 我的生日。YYYY-MM-DD
	Sex      string `json:"sex"`       // 我的性别。男 女 保密
	NickFree bool   `json:"nick_free"` // 是否未设置昵称。false：设置过昵称。true：未设置昵称
	Rank     string `json:"rank"`      // 我的会员等级
}

// GetAccountInformation 获取我的信息
func (c *Client) GetAccountInformation() (*AccountInformation, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/member/web/account"
	)
	return execute[*AccountInformation](c, method, url, nil)
}
