package bilibili

import (
	"github.com/go-resty/resty/v2"
)

// ClearHistory 清空历史记录
func (c *Client) ClearHistory(param RemoveDynamicParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v2/history/delete"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}
