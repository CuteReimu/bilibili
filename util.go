package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type paramHandler func(map[string]string) error

func fillCsrf(c *Client) paramHandler {
	return func(m map[string]string) error {
		csrf := c.getCookie("bili_jct")
		if len(csrf) == 0 {
			return errors.New("B站登录过期")
		}
		m["csrf"] = csrf
		return nil
	}
}

func fillParam(key, value string) paramHandler {
	return func(params map[string]string) error {
		params[key] = value
		return nil
	}
}

// execute 发起请求
func execute[Out any](c *Client, method, url string, in any, handlers ...paramHandler) (out Out, err error) {
	r := c.resty.R()
	if in != nil {
		var params map[string]string
		params, err = structToMap(in, handlers...)
		if err != nil {
			return
		}
		r = r.SetQueryParams(params)
	}
	fmt.Println(r.Header)
	resp, err := r.Execute(method, url)
	if err != nil {
		return out, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return out, errors.Errorf("status code: %d", resp.StatusCode())
	}
	c.resty.SetCookies(resp.Cookies())
	var cr commonResp[Out]
	if err = json.Unmarshal(resp.Body(), &cr); err != nil {
		return out, errors.WithStack(err)
	}
	if cr.Code != 0 {
		return out, errors.Errorf("错误码: %d, 错误信息: %s", cr.Code, cr.Message)
	}
	return cr.Data, errors.WithStack(err)
}

type commonResp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func structToMap(s any, handlers ...paramHandler) (map[string]string, error) {
	// TODO 后续优化一下这个算法，改成使用反射的方式实现
	buf, err := json.Marshal(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var m map[string]any
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	m2 := make(map[string]string, len(m))
	for k, v := range m {
		m2[k] = cast.ToString(v)
	}
	for _, handler := range handlers {
		if err = handler(m2); err != nil {
			return nil, err
		}
	}
	return m2, nil
}
