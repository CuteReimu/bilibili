package bilibili

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type GetWebCookieRefreshInfoResult struct {
	Refresh   bool  `json:"refresh"`   // 是否应该刷新 Cookie。true-需要刷新，false-不需要刷新
	Timestamp int64 `json:"timestamp"` // 用于获取 refresh_csrf 的毫秒时间戳
}

// GetWebCookieRefreshInfo 获取web端cookie刷新信息
func (c *Client) GetWebCookieRefreshInfo() (*GetWebCookieRefreshInfoResult, error) {
	var url = fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/cookie/info?bili_jct=%s", c.getCookie("bili_jct"))

	return execute[*GetWebCookieRefreshInfoResult](c, resty.MethodGet, url, nil)
}
