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
	cookies       []*http.Cookie
	cookiesString string
	timeout       time.Duration
	logger        resty.Logger
}

// New 返回一个 bilibili.Client
func New() *Client {
	return &Client{timeout: 20 * time.Second}
}

var std = New()

// GetTimeout 获取http请求超时时间，默认20秒
func GetTimeout() time.Duration {
	return std.GetTimeout()
}
func (c *Client) GetTimeout() time.Duration {
	if c.timeout == 0 {
		return time.Second * 20
	}
	return c.timeout
}

// SetLogger 设置logger
func SetLogger(logger resty.Logger) {
	std.SetLogger(logger)
}
func (c *Client) SetLogger(logger resty.Logger) {
	c.logger = logger
}

// GetLogger 获取logger，默认使用resty默认的logger
func GetLogger() resty.Logger {
	return std.GetLogger()
}
func (c *Client) GetLogger() resty.Logger {
	return c.logger
}

// GetCookiesString 获取字符串格式的cookies，方便自行存储后下次使用。配合下面的 SetCookiesString 使用。
func GetCookiesString() string {
	return std.cookiesString
}
func (c *Client) GetCookiesString() string {
	return c.cookiesString
}

// SetCookiesString 设置Cookies，但是是字符串格式，配合 GetCookiesString 使用。有些功能必须登录或设置Cookies后才能使用。
func SetCookiesString(cookiesString string) {
	std.SetCookiesString(cookiesString)
}
func (c *Client) SetCookiesString(cookiesString string) {
	c.cookiesString = cookiesString
	c.cookies = (&resty.Response{RawResponse: &http.Response{Header: http.Header{
		"Set-Cookie": strings.Split(cookiesString, "\n"),
	}}}).Cookies()
}

// GetCookies 获取Cookies。配合下面的SetCookies使用。
func GetCookies() []*http.Cookie {
	return std.GetCookies()
}
func (c *Client) GetCookies() []*http.Cookie {
	return c.cookies
}

// SetCookies 设置Cookies。有些功能必须登录之后才能使用，设置Cookies可以代替登录。
func SetCookies(cookies []*http.Cookie) {
	std.SetCookies(cookies)
}
func (c *Client) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
	var cookieStrings []string
	for _, cookie := range c.cookies {
		cookieStrings = append(cookieStrings, cookie.String())
	}
	c.cookiesString = strings.Join(cookieStrings, "\n")
}

// 获取resty的一个request
func (c *Client) resty() *resty.Client {
	client := resty.New().SetTimeout(c.GetTimeout()).SetHeader("user-agent", "go")
	if c.logger != nil {
		client.SetLogger(c.logger)
	}
	if c.cookies != nil {
		client.SetCookies(c.cookies)
	}
	return client
}

// 根据key获取指定的cookie值
func (c *Client) getCookie(name string) string {
	now := time.Now()
	for _, cookie := range c.cookies {
		if cookie.Name == name && now.Before(cookie.Expires) {
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
