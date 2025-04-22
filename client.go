package bilibili

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	wbi   *WBI
	resty *resty.Client
}

// New 返回一个默认的 bilibili.Client
func New() *Client {
	restyClient := resty.New().
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		SetTimeout(20*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Origin", "https://www.bilibili.com").
		SetHeader("Referer", "https://www.bilibili.com/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	return NewWithClient(restyClient)
}

// NewAnonymousClient 返回一个带有游客cookie的 bilibili.Client
func NewAnonymousClient() *Client {
	url := "https://www.bilibili.com/"
	method := resty.MethodGet

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("Sec-Ch-Ua", "Not")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", "Windows")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0")
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer func() { _ = res.Body.Close() }()

	bili_client := New()
	bili_client.SetCookies(res.Cookies())
	return bili_client
}

// NewWithClient 接收一个自定义的*resty.Client为参数
func NewWithClient(restyClient *resty.Client) *Client {
	return &Client{
		wbi:   NewDefaultWbi(),
		resty: restyClient,
	}
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
	c.SetCookies((&resty.Response{RawResponse: &http.Response{Header: http.Header{
		"Set-Cookie": strings.Split(cookiesString, "\n"),
	}}}).Cookies())
}

// SetRawCookies 如果你是从浏览器request的header中直接复制出来的cookies，调用这个函数。
func (c *Client) SetRawCookies(rawCookies string) {
	header := http.Header{}
	header.Add("Cookie", rawCookies)
	req := http.Request{Header: header}

	c.SetCookies(req.Cookies())
}

// SetCookie 设置单个cookie
func (c *Client) SetCookie(cookie *http.Cookie) {
	for i, c0 := range c.resty.Cookies {
		if c0.Name == cookie.Name {
			c.resty.Cookies[i] = cookie
			return
		}
	}
	c.resty.Cookies = append(c.resty.Cookies, cookie)
}

// SetCookies 设置cookies
func (c *Client) SetCookies(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		c.SetCookie(cookie)
	}
}

// GetCookies 获取当前的cookies
func (c *Client) GetCookies() []*http.Cookie {
	return c.resty.Cookies
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
