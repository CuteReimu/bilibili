package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type paramHandler func(*resty.Request) error

func fillCsrf(c *Client) paramHandler {
	return func(r *resty.Request) error {
		csrf := c.getCookie("bili_jct")
		if len(csrf) == 0 {
			return errors.New("B站登录过期")
		}
		r.SetQueryParam("csrf", csrf)
		r.SetQueryParam("csrf_token", csrf)
		return nil
	}
}

func fillParam(key, value string) paramHandler {
	return func(r *resty.Request) error {
		r.SetQueryParam(key, value)
		return nil
	}
}

func fillWbiHandler(wbi *WBI, cookies []*http.Cookie) func(*resty.Request) error {
	return func(r *resty.Request) error {
		newQuery, err := wbi.SignQuery(r.QueryParam, time.Now())
		if err != nil {
			return err
		}

		r.QueryParam = newQuery
		r.Cookies = cookies
		r.Header.Del("Referer")
		return nil
	}
}

// execute 发起请求
func execute[Out any](c *Client, method, url string, in any, handlers ...paramHandler) (out Out, err error) {
	r := c.resty.R()
	if err = withParams(r, in); err != nil {
		return
	}
	for _, handler := range handlers {
		if err = handler(r); err != nil {
			return
		}
	}
	resp, err := r.Execute(method, url)
	if err != nil {
		return out, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return out, errors.Errorf("status code: %d", resp.StatusCode())
	}
	c.SetCookies(resp.Cookies())
	var cr commonResp[Out]
	if err = json.Unmarshal(resp.Body(), &cr); err != nil {
		return out, errors.WithStack(err)
	}
	if cr.Code != 0 {
		return out, errors.WithStack(Error{Code: cr.Code, Message: cr.Message})
	}
	return cr.Data, errors.WithStack(err)
}

type commonResp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
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

	bodyMap := make(map[string]any, 4)
	contentType := ""
	for i := range inType.NumField() {
		fieldType := inType.Field(i)
		if !fieldType.IsExported() {
			continue
		}
		fieldValue := inValue.Field(i)
		tValue := fieldType.Tag.Get("request")

		if tValue == "-" {
			continue
		}

		// 获取字段名
		var fieldName string

		tagMap := parseTag(tValue)
		if name, ok := tagMap["field"]; ok {
			fieldName = name
		} else if jsonValue := fieldType.Tag.Get("json"); jsonValue != "" && jsonValue != "-" {
			if index := strings.Index(jsonValue, ","); index != -1 {
				jsonValue = jsonValue[:index]
			}
			fieldName = jsonValue
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

		contentType = "application/x-www-form-urlencoded"
		for name := range tagMap {
			switch name {
			case "query":
				contentType = "application/x-www-form-urlencoded"
			case "json":
				contentType = "application/json"
			case "form-data":
				contentType = "multipart/form-data"
			}
		}
		if contentType == "application/x-www-form-urlencoded" {
			_, ok1 := tagMap["json"]
			_, ok2 := tagMap["form-data"]
			if !ok1 && !ok2 {
				// 对query类型的字段进行特殊处理
				if fieldType.Type.Kind() == reflect.Slice {
					strSlice := make([]string, 0, 4)
					for i := range fieldValue.Len() {
						strSlice = append(strSlice, cast.ToString(fieldValue.Index(i).Interface()))
					}
					realVal = strings.Join(strSlice, ",")
				}
			}
			r.SetQueryParam(fieldName, cast.ToString(realVal))
		} else {
			bodyMap[fieldName] = realVal
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

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}
