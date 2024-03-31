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

func TestQuery(t *testing.T) {
	type Test struct {
		TestA string `request:"query"`
		TestB string `request:"query,field=tb,default=1"`
		TestC string `request:"query,omitempty"`
		TestD string `request:"-"`
		TestE int    `request:"query,default=1"`
		TestF *int   `request:"query,default=1"`
		TestG *int   `request:"query,default=1"`
		TestH *int   `request:"query,omitempty"`
		TestI *int   `request:"query,omitempty"`
		TestJ *int   `request:"query"`
		TestK int    `request:"query"`
	}

	f, i := 10, 0
	params := Test{
		TestD: "test_d",
		TestE: 0,
		TestF: &f,
		TestI: &i,
	}

	r := resty.New().R()
	err := withParams(r, params)

	if err != nil {
		t.Fatal(err)
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatal("withParams content type not correct ", r.Header.Get("Content-Type"))
	}

	query := make(map[string]string)
	for k := range r.QueryParam {
		query[k] = r.QueryParam.Get(k)
	}
	if !maps.Equal(query, map[string]string{
		"test_a": "",
		"tb":     "1",
		"test_e": "1",
		"test_f": "10",
		"test_g": "1",
		"test_i": "0",
		"test_j": "",
		"test_k": "0",
	}) {
		t.Fatal("withParams query result not correct ", r.QueryParam)
	}
}

func TestJson(t *testing.T) {
	type Test struct {
		TestA string `request:"json"`
		TestB string `request:"json,field=tb,default=1"`
		TestC string `request:"json,omitempty"`
		TestD string `request:"-"`
		TestE int    `request:"json,default=1"`
		TestF *int   `request:"json,default=1"`
		TestG *int   `request:"json,default=1"`
		TestH *int   `request:"json,omitempty"`
		TestI *int   `request:"json,omitempty"`
		TestJ *int   `request:"json"`
		TestK int    `request:"json"`
	}

	f, i := 10, 0
	params := Test{
		TestD: "test_d",
		TestE: 0,
		TestF: &f,
		TestI: &i,
	}

	r := resty.New().R()
	err := withParams(r, params)

	if err != nil {
		t.Fatal(err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		t.Fatal("withParams content type not correct ", r.Header.Get("Content-Type"))
	}

	if !maps.Equal(r.Body.(map[string]interface{}), map[string]any{
		"test_a": "",
		"tb":     "1",
		"test_e": "1",
		"test_f": &f,
		"test_g": "1",
		"test_i": &i,
		"test_j": (*int)(nil),
		"test_k": 0,
	}) {
		t.Fatal("withParams body result not correct ", r.Body)
	}
}

func TestFormData(t *testing.T) {
	type Test struct {
		TestA string `request:"form-data"`
		TestB string `request:"form-data,field=tb,default=1"`
		TestC string `request:"form-data,omitempty"`
		TestD string `request:"-"`
		TestE int    `request:"form-data,default=1"`
		TestF *int   `request:"form-data,default=1"`
		TestG *int   `request:"form-data,default=1"`
		TestH *int   `request:"form-data,omitempty"`
		TestI *int   `request:"form-data,omitempty"`
		TestJ *int   `request:"form-data"`
		TestK int    `request:"form-data"`
	}

	f, i := 10, 0
	params := Test{
		TestD: "test_d",
		TestE: 0,
		TestF: &f,
		TestI: &i,
	}

	r := resty.New().R()
	err := withParams(r, params)

	if err != nil {
		t.Fatal(err)
		return
	}

	if r.Header.Get("Content-Type") != "multipart/form-data" {
		t.Fatal("withParams content type not correct ", r.Header.Get("Content-Type"))
	}

	if !maps.Equal(r.Body.(map[string]interface{}), map[string]any{
		"test_a": "",
		"tb":     "1",
		"test_e": "1",
		"test_f": &f,
		"test_g": "1",
		"test_i": &i,
		"test_j": (*int)(nil),
		"test_k": 0,
	}) {
		t.Fatal("withParams body result not correct ", r.Body)
	}
}

func TestWithParamsNil(t *testing.T) {
	r := resty.New().R()
	err := withParams(r, []int{1, 2, 3})

	if err == nil || err.Error() != "参数类型错误" {
		t.Fatal(err)
	}

	err = withParams(r, nil)
	if err != nil {
		t.Fatal("nil params should not return error")
	}

	err = withParams(r, (*int)(nil))
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
