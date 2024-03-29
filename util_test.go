package bilibili

import (
	"maps"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestStructToMap(t *testing.T) {
	type GetVideoCommentParams struct {
		AccessKey string `json:"access_key,omitempty"` // APP 登录 Token，不是APP方式可以填空
		Type      int    `json:"type"`                 // 评论区类型代码
		Oid       int    `json:"oid"`                  // 目标评论区 id
		Sort      int    `json:"sort,omitempty"`       // 排序方式
		Nohot     int    `json:"nohot,omitempty"`      // 是否不显示热评
		Ps        int    `json:"ps,omitempty"`         // 每页项数
		Pn        int    `json:"pn,omitempty"`         // 页码
	}

	params := GetVideoCommentParams{
		AccessKey: "abc",
		Type:      1,
		Oid:       2,
		Sort:      3,
	}

	m, err := structToMap(params)
	if err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(m, map[string]string{
		"access_key": "abc",
		"type":       "1",
		"oid":        "2",
		"sort":       "3",
	}) {
		t.Fatal("structToMap result not correct ", m)
	}
}

func TestWithParams(t *testing.T) {
	type GetVideoCommentParams struct {
		AccessKey string `json:"access_key,omitempty" request:"query"`       // APP 登录 Token，不是APP方式可以填空
		Type      int    `json:"type" request:"query,default=1"`             // 评论区类型代码
		Oid       int    `json:"oid" request:"query,default=2"`              // 目标评论区 id
		Sort      int    `json:"sort,omitempty" request:"query,default=3"`   // 排序方式
		Nohot     int    `json:"nohot,omitempty" request:"query,default=-1"` // 是否不显示热评
		Ps        int    `json:"ps,omitempty" request:"query,default=20"`    // 每页项数
		Pn        int    `json:"pn,omitempty" request:"query,default=1"`     // 页码

		TestA string `json:"test_a,omitempty" request:"json,default=1"`
		TestB string `json:"test_b,omitempty" request:"json,omitempty"`
		TestC string `json:"test_c,omitempty" request:"form-data,field=TC"`
		TestD string `json:"test_d,omitempty" request:"-"`
		TestE string `json:"test_e,omitempty" request:"json"`
	}

	params := GetVideoCommentParams{
		AccessKey: "abc",
		Type:      1,
		Oid:       2,
		Sort:      3,

		TestC: "test_c",
		TestD: "test_d",
	}

	r := resty.New().R()
	err := withParams(r, params)

	if err != nil {
		t.Fatal(err)
		return
	}

	query := make(map[string]string)
	for k := range r.QueryParam {
		query[k] = r.QueryParam.Get(k)
	}
	if !maps.Equal(query, map[string]string{
		"access_key": "abc",
		"type":       "1",
		"oid":        "2",
		"sort":       "3",
		"nohot":      "-1",
		"ps":         "20",
		"pn":         "1",
	}) {
		t.Fatal("withParams query result not correct ", r.QueryParam)
	}

	if !maps.Equal(r.Body.(map[string]interface{}), map[string]interface{}{
		"test_a": "1",
		"TC":     "test_c",
		"test_e": "",
	}) {
		t.Fatal("withParams body result not correct ", r.Body)
	}
}

func TestWithParams2(t *testing.T) {
	type GetVideoCommentParams struct {
		AccessKey string `json:"access_key,omitempty" request:"query"`       // APP 登录 Token，不是APP方式可以填空
		Type      int    `json:"type" request:"query,default=1"`             // 评论区类型代码
		Oid       int    `json:"oid" request:"query,default=2"`              // 目标评论区 id
		Sort      int    `json:"sort,omitempty" request:"query,default=3"`   // 排序方式
		Nohot     int    `json:"nohot,omitempty" request:"query,default=-1"` // 是否不显示热评
		Ps        int    `json:"ps,omitempty" request:"query,default=20"`    // 每页项数
		Pn        int    `json:"pn,omitempty" request:"query,default=1"`     // 页码

		TestA string `json:"test_a,omitempty" request:"json,default=1"`
		TestB string `json:"test_b,omitempty" request:"json,omitempty"`
		TestC string `json:"test_c,omitempty" request:"form-data,field=TC"`
		TestD string `json:"test_d,omitempty" request:"-"`
		TestE string `json:"test_e,omitempty" request:"json"`
	}

	params := &GetVideoCommentParams{
		AccessKey: "abc",
		Type:      1,
		Oid:       2,
		Sort:      3,

		TestC: "test_c",
		TestD: "test_d",
	}

	r := resty.New().R()
	err := withParams(r, params)

	if err != nil {
		t.Fatal(err)
		return
	}

	query := make(map[string]string)
	for k := range r.QueryParam {
		query[k] = r.QueryParam.Get(k)
	}
	if !maps.Equal(query, map[string]string{
		"access_key": "abc",
		"type":       "1",
		"oid":        "2",
		"sort":       "3",
		"nohot":      "-1",
		"ps":         "20",
		"pn":         "1",
	}) {
		t.Fatal("withParams query result not correct ", r.QueryParam)
	}

	if !maps.Equal(r.Body.(map[string]interface{}), map[string]interface{}{
		"test_a": "1",
		"TC":     "test_c",
		"test_e": "",
	}) {
		t.Fatal("withParams body result not correct ", r.Body)
	}
}

func TestWithParams3(t *testing.T) {
	r := resty.New().R()
	err := withParams(r, []int{1, 2, 3})

	if err == nil || err.Error() != "参数类型错误" {
		t.Fatal(err)
	}

	err = withParams(r, nil)
	if err != nil {
		t.Fatal("nil params should not return error")
	}

	if len(r.QueryParam) != 0 {
		t.Fatal("withParams query result not correct ", r.QueryParam)
	}

	if r.Body != nil && len(r.Body.(map[string]interface{})) != 0 {
		t.Fatal("withParams body result not correct ", r.Body)
	}
}
