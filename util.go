package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type paramHandler func(*Client, map[string]string)

func fillCsrf(c *Client, params map[string]string) {
	params["csrf"] = c.getCookie("bili_jct")
}

// execute 发起请求
func execute[In, Out any](c *Client, method, url string, in In, handlers ...paramHandler) (Out, error) {
	var out Out
	params, err := structToMap(in)
	if err != nil {
		return out, err
	}
	for _, handler := range handlers {
		handler(c, params)
	}
	resp, err := c.resty().R().SetQueryParams(params).Execute(method, url)
	if err != nil {
		return out, errors.WithStack(err)
	}
	var cr commonResp[Out]
	if err = json.Unmarshal(resp.Body(), &cr); err != nil {
		return out, errors.WithStack(err)
	}
	return cr.Data, errors.WithStack(err)
}

type commonResp[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func structToMap(s any) (map[string]string, error) {
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
	return m2, nil
}
