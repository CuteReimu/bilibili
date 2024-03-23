package bilibili

import (
	"maps"
	"testing"
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
