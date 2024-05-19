package bilibili

import (
	"github.com/go-resty/resty/v2"
)

type HistoryParam struct {
	Max      int    `json:"max,omitempty" request:"query,omitempty"`      // 历史记录截止目标 id。默认为 0。稿件：稿件 avid。剧集（番剧 / 影视）：剧集 ssid。直播：直播间 id。文集：文集 rlid。文章：文章 cvid
	Business string `json:"business,omitempty" request:"query,omitempty"` // 历史记录截止目标业务类型。默认为空。archive：稿件。pgc：剧集（番剧 / 影视）。live：直播。article-list：文集。article：文章
	ViewAt   int    `json:"view_at,omitempty" request:"query,omitempty"`  // 历史记录截止时间。时间戳。默认为 0。0 为当前 时间
	Type     string `json:"type,omitempty" request:"query,omitempty"`     // 历史记录分类筛选。all：全部类型（默认）。archive：稿件。live：直播。article：文章
	Ps       int    `json:"ps,omitempty" request:"query,omitempty"`       // 每页项数。默认为 20，最大 30
}
type HistoryCursor struct {
	Max      int    `json:"max"`      // 最后一项目标 id。**见请求参数**
	ViewAt   int    `json:"view_at"`  // 最后一项时间节点。时间戳
	Business string `json:"business"` // 最后一项业务类型。**见请求参数**
	Ps       int    `json:"ps"`       // 每页项数
}
type HistoryTab struct {
	Type string `json:"type"` // 类型
	Name string `json:"name"` // 类型名
}
type HistoryDetail struct {
	Oid      int    `json:"oid"`      // 目标id。稿件视频&剧集（当business=archive或business=pgc时）：稿件avid。直播（当business=live时）：直播间id。文章（当business=article时）：文章cvid。文集（当business=article-list时）：文集rlid
	Epid     int    `json:"epid"`     // 剧集epid。仅用于剧集
	Bvid     string `json:"bvid"`     // 稿件bvid。仅用于稿件视频
	Page     int    `json:"page"`     // 观看到的视频分P数。仅用于稿件视频
	Cid      int    `json:"cid"`      // 观看到的对象id。稿件视频&剧集（当business=archive或business=pgc时）：视频cid。文集（当business=article-list时）：文章cvid
	Part     string `json:"part"`     // 观看到的视频分 P 标题。仅用于稿件视频
	Business string `json:"business"` // 业务类型。**见请求参数**
	Dt       int    `json:"dt"`       // 记录查看的平台代码。1 3 5 7：手机端。2：web端。4 6：pad端。33：TV端。0：其他
}
type HistoryList struct {
	Title      string        `json:"title"`       // 条目标题
	LongTitle  string        `json:"long_title"`  // 条目副标题
	Cover      string        `json:"cover"`       // 条目封面图 url。用于专栏以外的条目
	Covers     []string      `json:"covers"`      // 条目封面图组。仅用于专栏	Uri string `json:"uri"` // 重定向 url。仅用于剧集和直播
	History    HistoryDetail `json:"history"`     // 条目详细信息
	Videos     int           `json:"videos"`      // 视频分 P 数目。仅用于稿件视频
	AuthorName string        `json:"author_name"` // UP 主昵称
	AuthorFace string        `json:"author_face"` // UP 主头像 url
	AuthorMid  int           `json:"author_mid"`  // UP 主 mid
	ViewAt     int           `json:"view_at"`     // 查看时间。时间戳
	Progress   int           `json:"progress"`    // 视频观看进度。单位为秒。用于稿件视频或剧集
	Badge      string        `json:"badge"`       // 角标文案。稿件视频 / 剧集 / 笔记
	ShowTitle  string        `json:"show_title"`  // 分 P 标题。用于稿件视频或剧集
	Duration   int           `json:"duration"`    // 视频总时长。用于稿件视频或剧集
	Current    string        `json:"current"`     // (?)
	Total      int           `json:"total"`       // 总计分集数。仅用于剧集
	NewDesc    string        `json:"new_desc"`    // 最新一话 / 最新一 P 标识。用于稿件视频或剧集
	IsFinish   int           `json:"is_finish"`   // 是否已完结。仅用于剧集。0：未完结。1：已完结
	IsFav      int           `json:"is_fav"`      // 是否收藏。0：未收藏。1：已收藏
	Kid        int           `json:"kid"`         // 条目目标 id。**详细内容见参数**
	TagName    string        `json:"tag_name"`    // 子分区名。用于稿件视频和直播
	LiveStatus int           `json:"live_status"` // 直播状态。仅用于直播。0：未开播。1：已开播
}
type HistoryInfo struct {
	Cursor HistoryCursor `json:"cursor"` // 历史记录页面信息
	Tab    []HistoryTab  `json:"tab"`    // 历史记录筛选类型
	List   []HistoryList `json:"list"`   // 分段历史记录列表
}

// GetHistory 获取历史记录列表
func (c *Client) GetHistory(param HistoryParam) (*HistoryInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/history/cursor"
	)
	return execute[*HistoryInfo](c, method, url, param)
}

type DeleteHistoryParam struct {
	Kid string `json:"kid"` // 删除的目标记录，格式为{业务类型}_{目标id}详见备注。视频：archive_{稿件avid}。直播：live_{直播间id}。专栏：article_{专栏cvid}。剧集：pgc_{剧集ssid}。文集：article-list_{文集rlid}
}

// DeleteHistory 删除历史记录
func (c *Client) DeleteHistory(param DeleteHistoryParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v2/history/delete"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

// ClearHistory 清空历史记录
func (c *Client) ClearHistory() error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v2/history/clear"
	)
	_, err := execute[any](c, method, url, nil, fillCsrf(c))
	return err
}
