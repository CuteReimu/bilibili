package bilibili

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-resty/resty/v2"
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

func withParams(r *resty.Request, in any) error {
	if in == nil {
		return nil
	}

	inType := reflect.TypeOf(in)
	inValue := reflect.ValueOf(in)

	switch inType.Kind() {
	case reflect.Ptr:
		// 如果是空指针，直接返回
		if inValue.IsNil() {
			return nil
		}
		inType = inType.Elem()
		inValue = inValue.Elem()
	case reflect.Struct:
	default:
		return errors.New("参数类型错误")
	}

	bodyMap := make(map[string]interface{}, 4)
	contentType := ""
	for i := 0; i < inType.NumField(); i++ {
		fieldType := inType.Field(i)
		fieldValue := inValue.Field(i)
		tValue := fieldType.Tag.Get("request")

		if tValue == "" || tValue == "-" {
			continue
		}

		// 获取字段名
		var fieldName string

		tagMap := parseTag(tValue)
		if name, ok := tagMap["field"]; ok {
			fieldName = name
		} else {
			fieldName = toSnakeCase(fieldType.Name)
		}

		var realVal any
		if !fieldValue.IsZero() {
			realVal = fieldValue.Interface()
		} else {
			// 设置了 omitempty 代表不传
			if _, ok := tagMap["omitempty"]; ok {
				continue
			}
			// 设置了 default 代表使用默认值
			if v, ok := tagMap["default"]; ok {
				realVal = v
			} else {
				// 否则使用零值
				realVal = fieldValue.Interface()
			}
		}

		for name := range tagMap {
			switch name {
			case "query":
				contentType = "application/x-www-form-urlencoded"
				r.SetQueryParam(fieldName, cast.ToString(realVal))
			case "json":
				contentType = "application/json"
				bodyMap[fieldName] = realVal
			case "form-data":
				contentType = "multipart/form-data"
				bodyMap[fieldName] = realVal
			}
		}
	}

	r.SetHeader("Content-Type", contentType)
	if len(bodyMap) > 0 {
		r.SetBody(bodyMap)
	}

	return nil
}

func parseTag(tag string) map[string]string {
	parts := strings.Split(tag, ",")

	pMap := make(map[string]string, 10)
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) == 1 {
			pMap[kv[0]] = ""
		} else {
			pMap[kv[0]] = kv[1]
		}
	}

	return pMap
}

func toSnakeCase(s string) string {
	var result strings.Builder
	result.Grow(len(s) * 2)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
