package bilibili

import "github.com/go-resty/resty/v2"

type ZoneVideoRankListParam struct {
	Tid  int    `json:"tid,omitempty" request:"query,omitempty"`  // 目标分区tid，可不填。可调用 GetAllZoneInfos 获取，或者直接自行查阅 video_zone.csv
	Type string `json:"type,omitempty" request:"query,omitempty"` // 未知。默认为：all，且为目前唯一已知值。怀疑为稿件类型，但没有找到其他值佐证。
}

type ZoneVideoRankList struct {
	Note string      `json:"note"` // “根据稿件内容质量、近期的数据综合展示，动态更新”
	List []VideoInfo `json:"list"` // 视频列表
}

// GetZoneVideoRankList 获取分区视频排行榜列表
func (c *Client) GetZoneVideoRankList(param ZoneVideoRankListParam) (*ZoneVideoRankList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/ranking/v2"
	)
	return execute[*ZoneVideoRankList](c, method, url, param)
}

type GetZoneVideoListNewParam struct {
	Pn  int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为1
	Ps  int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为14, 留空为5
	Rid int `json:"rid"`                                    // 目标分区tid
}
type ZoneVideoListInfo struct {
	Archives []VideoInfo   `json:"archives"` // 视频列表
	Page     ZoneVideoPage `json:"page"`     // 页面信息
}
type ZoneVideoPage struct {
	Count int `json:"count"` // 总计视频数
	Num   int `json:"num"`   // 当前页码
	Size  int `json:"size"`  // 每页项数
}

// GetZoneVideoListNew 获取分区最新视频列表
func (c *Client) GetZoneVideoListNew(param GetZoneVideoListNewParam) (*ZoneVideoListInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/dynamic/region"
	)
	return execute[*ZoneVideoListInfo](c, method, url, param)
}

type GetZoneVideoListWithTagParam struct {
	Ps    int `json:"ps,omitempty" request:"query,omitempty"` // 视频数。默认为14, 留空为5
	Pn    int `json:"pn,omitempty" request:"query,omitempty"` // 列数。留空为1
	Rid   int `json:"rid"`                                    // 目标分区id。参见[视频分区一览](../video/video_zone.md)
	TagId int `json:"tag_id"`                                 // 目标标签id
}

// GetZoneVideoListWithTag 获取分区标签近期互动列表
func (c *Client) GetZoneVideoListWithTag(param GetZoneVideoListWithTagParam) (*ZoneVideoListInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/dynamic/tag"
	)
	return execute[*ZoneVideoListInfo](c, method, url, param)
}

type GetZoneVideoListRecentParam struct {
	Ps   int `json:"ps,omitempty" request:"query,omitempty"`   // 视频数。默认为14, 留空为5
	Pn   int `json:"pn,omitempty" request:"query,omitempty"`   // 页码。默认为1
	Rid  int `json:"rid,omitempty" request:"query,omitempty"`  // 目标分区id。参见[视频分区一览](../video/video_zone.md)
	Type int `json:"type,omitempty" request:"query,omitempty"` // 类型?。默认为0
}

// GetZoneVideoListRecent 获取分区近期投稿列表
func (c *Client) GetZoneVideoListRecent(param GetZoneVideoListRecentParam) (*ZoneVideoListInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/newlist"
	)
	return execute[*ZoneVideoListInfo](c, method, url, param)
}

type GetZoneVideoListByOrderParam struct {
	MainVer    string `json:"main_ver,omitempty" request:"query,omitempty"`    // 主页版本。默认为 v3
	SearchType string `json:"search_type"`                                     // 搜索类型。默认为 video
	ViewType   string `json:"view_type"`                                       // 查看类型?。默认为 hot_rank
	CopyRight  int    `json:"copy_right,omitempty" request:"query,omitempty"`  // 版权?。默认为 -1
	NewWebTag  int    `json:"new_web_tag,omitempty" request:"query,omitempty"` // 标签?。默认为 1
	Order      string `json:"order,omitempty" request:"query,omitempty"`       // 排序方式。click: 按播放排序(默认)。scores: 按评论数排序。stow: 按收藏排序。coin: 按硬币数排序。dm: 按弹幕数排序
	CateId     int    `json:"cate_id"`                                         // 分区id。留空会导致响应中data中result为null, 参见[视频分区一览](../video/video_zone.md)
	Page       int    `json:"page,omitempty" request:"query,omitempty"`        // 页码。默认以 1 开始
	Pagesize   int    `json:"pagesize"`                                        // 视频数。默认为 30, 留空会导致 -500
	TimeFrom   int    `json:"time_from"`                                       // 起始时间。yyyyMMdd, 默认为 time_to - 7
	TimeTo     int    `json:"time_to"`                                         // 结束时间。yyyyMMdd, 默认为当前时间(大于起始时间)
}

type ZoneVideoRankInfo struct {
	ExpList        *string         `json:"exp_list"`         // 作用尚不明确
	ShowModuleList []string        `json:"show_module_list"` // 显示模块列表?
	Result         []RankVideoInfo `json:"result"`           // 结果本体。失败时为null
	ShowColumn     int             `json:"show_column"`      // 0。作用尚不明确
	RqtType        string          `json:"rqt_type"`         // search。作用尚不明确
	Numpages       int             `json:"numPages"`         // 页码。失败时为0
	Numresults     int             `json:"numResults"`       // 视频数。失败时为0
	CrrQuery       *string         `json:"crr_query"`        // 空。作用尚不明确
	Pagesize       int             `json:"pagesize"`         // 视频数
	SuggestKeyword *string         `json:"suggest_keyword"`  // 空。作用尚不明确
	EggInfo        *string         `json:"egg_info"`         // 作用尚不明确
	Cache          int             `json:"cache"`            // 0。作用尚不明确
	ExpBits        int             `json:"exp_bits"`         // 1。作用尚不明确
	ExpStr         *string         `json:"exp_str"`          // 空。作用尚不明确
	Seid           string          `json:"seid"`             // 一串字符串中的数字。作用尚不明确
	Msg            string          `json:"msg"`              // 结果信息。成功时为success, 反之为as error.
	EggHit         int             `json:"egg_hit"`          // 0。作用尚不明确
	Page           int             `json:"page"`             // 页码
}
type RankVideoInfo struct {
	Pubdate      string `json:"pubdate"`        // 发布时间。格式为 yyyy-MM-dd HH:mm:ss
	Pic          string `json:"pic"`            // 封面图
	Tag          string `json:"tag"`            // 标签。用 , 分隔
	Duration     int    `json:"duration"`       // 时长。单位为秒
	Id           int    `json:"id"`             // aid
	RankScore    int    `json:"rank_score"`     // 排序分数?
	Badgepay     bool   `json:"badgepay"`       // 是否有角标?
	Senddate     int    `json:"senddate"`       // 发送时间?。UNIX 秒级时间戳
	Author       string `json:"author"`         // UP主名
	Review       int    `json:"review"`         // 评论数
	Mid          int    `json:"mid"`            // UP主mid
	IsUnionVideo int    `json:"is_union_video"` // 是否为联合投稿
	RankIndex    int    `json:"rank_index"`     // 排序索引号
	Type         string `json:"type"`           // 类型。video: 视频
	Arcrank      string `json:"arcrank"`        // 0。作用尚不明确
	Play         string `json:"play"`           // 播放数
	RankOffset   int    `json:"rank_offset"`    // 排序偏移?。与 rank_index 相同
	Description  string `json:"description"`    // 简介
	VideoReview  int    `json:"video_review"`   // 弹幕数?
	IsPay        int    `json:"is_pay"`         // 是否付费?。0: 免费。1: 付费
	Favorites    int    `json:"favorites"`      // 收藏数
	Arcurl       string `json:"arcurl"`         // 视频播放页URL
	Bvid         string `json:"bvid"`           // bvid
	Title        string `json:"title"`          // 标题
	Vt           int    `json:"vt"`             // 0。作用尚不明确
	EnableVt     int    `json:"enable_vt"`      // 0。作用尚不明确
	VtDisplay    string `json:"vt_display"`     // 空。作用尚不明确
}

// GetZoneVideoListByOrder 获取分区近期投稿列表 (带排序)
func (c *Client) GetZoneVideoListByOrder(param GetZoneVideoListByOrderParam) (*ZoneVideoRankInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/newlist_rank"
	)
	return execute[*ZoneVideoRankInfo](c, method, url, param)
}
