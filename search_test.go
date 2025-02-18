package bilibili_test

import (
	"testing"

	"github.com/CuteReimu/bilibili/v2"
)

func TestSearch(t *testing.T) {
	c := bilibili.NewAnonymousClient()
	// c_str := c.GetCookiesString()
	// fmt.Print(c_str)
	res, err := c.IntergratedSearch(bilibili.SearchParam{
		Keyword: "东方幻想乡",
		// Page:     2,
		// PageSize: 42,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Seid == "" {
		t.Fatal("search id is nil. search failed")
	}

	// 测试搜索结果
	if len(res.Result) == 0 {
		t.Fatal("search result is nil")
	}
	// type
	for _, ele := range res.Result {
		if ele.ResultType == "video" {
			return
		}
	}
	t.Fatal("no available video result")
}
