package bilibili

import (
	"github.com/go-resty/resty/v2"
)

// ClearHistory 清空历史记录
func (c *Client) ClearHistory() error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v2/history/clear"
	)
	_, err := execute[any](c, method, url, nil, fillCsrf(c))
	return err
}
