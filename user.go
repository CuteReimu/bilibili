package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type OrderType string

const (
	OrderPubDate OrderType = "pubdate"
	OrderClick   OrderType = "click"
	OrderStow    OrderType = "stow"
)

type Video struct {
	Aid          int    `json:"aid"`            // 稿件avid
	Author       string `json:"author"`         // 视频UP主，不一定为目标用户（合作视频）
	Bvid         string `json:"bvid"`           // 稿件bvid
	Comment      int    `json:"comment"`        // 视频评论数
	Copyright    string `json:"copyright"`      // 空，作用尚不明确
	Created      int64  `json:"created"`        // 投稿时间戳
	Description  string `json:"description"`    // 视频简介
	HideClick    bool   `json:"hide_click"`     // 固定值false，作用尚不明确
	IsPay        int    `json:"is_pay"`         // 固定值0，作用尚不明确
	IsUnionVideo int    `json:"is_union_video"` // 是否为合作视频，0：否，1：是
	Length       string `json:"length"`         // 视频长度，MM:SS
	Mid          int    `json:"mid"`            // 视频UP主mid，不一定为目标用户（合作视频）
	Pic          string `json:"pic"`            // 视频封面
	Play         int    `json:"play"`           // 视频播放次数
	Review       int    `json:"review"`         // 固定值0，作用尚不明确
	Subtitle     string `json:"subtitle"`       // 固定值空，作用尚不明确
	Title        string `json:"title"`          // 视频标题
	Typeid       int    `json:"typeid"`         // 视频分区tid
	VideoReview  int    `json:"video_review"`   // 视频弹幕数
}

type GetUserVideosResult struct {
	List struct { // 列表信息
		Tlist map[int]struct { // 投稿视频分区索引
			Count int    `json:"count"` // 投稿至该分区的视频数
			Name  string `json:"name"`  // 该分区名称
			Tid   int    `json:"tid"`   // 该分区tid
		} `json:"tlist"`
		Vlist []Video `json:"vlist"` // 投稿视频列表
	} `json:"list"`
	Page struct { // 页面信息
		Count int `json:"count"` // 总计稿件数
		Pn    int `json:"pn"`    // 当前页码
		Ps    int `json:"ps"`    // 每页项数
	} `json:"page"`
	EpisodicButton struct { // “播放全部“按钮
		Text string `json:"text"` // 按钮文字
		Uri  string `json:"uri"`  // 全部播放页url
	} `json:"episodic_button"`
}

// GetUserVideos 获取用户投稿视频明细
func GetUserVideos(mid int, order OrderType, tid int, keyword string, pn int, ps int) (*GetUserVideosResult, error) {
	return std.GetUserVideos(mid, order, tid, keyword, pn, ps)
}
func (c *Client) GetUserVideos(mid int, order OrderType, tid int, keyword string, pn int, ps int) (*GetUserVideosResult, error) {
	postData := map[string]string{
		"mid": strconv.Itoa(mid),
		"pn":  strconv.Itoa(pn),
		"ps":  strconv.Itoa(ps),
	}
	if len(order) > 0 {
		postData["order"] = string(order)
	}
	if tid != 0 {
		postData["tid"] = strconv.Itoa(tid)
	}
	if len(keyword) > 0 {
		postData["keyword"] = keyword
	}
	resp, err := c.resty().R().SetQueryParams(postData).Get("https://api.bilibili.com/x/space/arc/search")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取用户视频")
	if err != nil {
		return nil, err
	}
	var ret *GetUserVideosResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
