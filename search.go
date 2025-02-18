package bilibili

import "github.com/go-resty/resty/v2"

// https://socialsisteryi.github.io/bilibili-API-collect/docs/search/*

type SearchParam struct {
	Keyword string `json:"keyword" request:"query"` //	需要搜索的关键词
	// Page     int    `json:"page" request:"query"`      // 从1开始
	// PageSize int    `json:"page_size" request:"query"` // 默认42
}

type SearchRespData struct {
	Seid           string             `json:"seid"`             // 搜索id
	Page           int                `json:"page"`             // 页数，固定为1
	PageSize       int                `json:"page_size"`        // 每页条数，固定为20
	NumResults     int                `json:"numResults"`       // 总条数，最大值为1000
	NumPages       int                `json:"numPages"`         // 分页数，最大值为50
	SuggestKeyword string             `json:"suggest_keyword"`  // 作用尚不明确
	RqtType        string             `json:"rqt_type"`         // 作用尚不明确
	CostTime       SearchRespCostTime `json:"cost_time"`        // 详细搜索用时
	ExpList        map[string]any     `json:"exp_list"`         // 作用尚不明确
	EggHit         int                `json:"egg_hit"`          // 作用尚不明确
	PageInfo       SearchRespPageInfo `json:"pageinfo"`         // 分类页数信息
	TopTList       SearchRespTopTList `json:"top_tlist"`        // 分类结果数目信息
	ShowColumn     int                `json:"show_column"`      // 作用尚不明确
	ShowModuleList []string           `json:"show_module_list"` // 返回结果类型列表
	Result         []SeachRespResult  `json:"result"`           // 结果列表
}

type SearchRespCostTime struct {
	ParamsCheck      string `json:"params_check"`         // 参数检查时间
	IllegalHandler   string `json:"illegal_handler"`      // 违规处理时间
	AsResponseFormat string `json:"as_response_format"`   // 响应格式化时间
	AsRequest        string `json:"as_request"`           // 请求时间
	SaveCache        string `json:"save_cache"`           // 缓存保存时间
	DeserializeResp  string `json:"deserialize_response"` // 反序列化时间
	AsRequestFormat  string `json:"as_request_format"`    // 请求格式化时间
	Total            string `json:"total"`                // 总耗时
	MainHandler      string `json:"main_handler"`         // 主处理时间
}

type SearchRespPageInfo struct {
	Pgc           SearchRespItemCount `json:"pgc"`            // PGC相关内容数
	LiveRoom      SearchRespItemCount `json:"live_room"`      // 直播间数
	Photo         SearchRespItemCount `json:"photo"`          // 相簿数
	Topic         SearchRespItemCount `json:"topic"`          // 话题数
	Video         SearchRespItemCount `json:"video"`          // 视频数
	User          SearchRespItemCount `json:"user"`           // 用户数
	BiliUser      SearchRespItemCount `json:"bili_user"`      // B站用户数
	MediaFt       SearchRespItemCount `json:"media_ft"`       // 电影数
	Article       SearchRespItemCount `json:"article"`        // 专栏数
	MediaBangumi  SearchRespItemCount `json:"media_bangumi"`  // 番剧数
	Special       SearchRespItemCount `json:"special"`        // 特别项目数
	OperationCard SearchRespItemCount `json:"operation_card"` // 运营卡片数
	UpUser        SearchRespItemCount `json:"upuser"`         // UP主数
	Movie         SearchRespItemCount `json:"movie"`          // 电影数
	LiveAll       SearchRespItemCount `json:"live_all"`       // 直播相关数
	Tv            SearchRespItemCount `json:"tv"`             // 电视数
	Live          SearchRespItemCount `json:"live"`           // 直播间数
	Bangumi       SearchRespItemCount `json:"bangumi"`        // 番剧数
	Activity      SearchRespItemCount `json:"activity"`       // 活动数
	LiveMaster    SearchRespItemCount `json:"live_master"`    // 主播数
	LiveUser      SearchRespItemCount `json:"live_user"`      // 直播用户数
}

type SearchRespItemCount struct {
	NumResults int `json:"numResults"` // 结果总数
	Total      int `json:"total"`      // 总计数量
	Pages      int `json:"pages"`      // 分页数量
}

type SearchRespTopTList struct {
	Pgc           int `json:"pgc"`            // PGC内容数
	LiveRoom      int `json:"live_room"`      // 直播数
	Photo         int `json:"photo"`          // 相簿数
	Topic         int `json:"topic"`          // 话题数
	Video         int `json:"video"`          // 视频数
	User          int `json:"user"`           // 用户数
	BiliUser      int `json:"bili_user"`      // B站用户数
	MediaFt       int `json:"media_ft"`       // 电影数
	Article       int `json:"article"`        // 专栏数
	MediaBangumi  int `json:"media_bangumi"`  // 番剧数
	Card          int `json:"card"`           // 作用不明确
	OperationCard int `json:"operation_card"` // 运营卡片数
	UpUser        int `json:"upuser"`         // UP主数
	Movie         int `json:"movie"`          // 电影数
	LiveAll       int `json:"live_all"`       // 直播相关数
	Tv            int `json:"tv"`             // 电视数
	Live          int `json:"live"`           // 直播间数
	Special       int `json:"special"`        // 特别项目数
	Bangumi       int `json:"bangumi"`        // 番剧数
	Activity      int `json:"activity"`       // 活动数
	LiveMaster    int `json:"live_master"`    // 主播数
	LiveUser      int `json:"live_user"`      // 直播用户数
}

type SeachRespResult struct {
	ResultType string           `json:"result_type"` // 结果类型 与result数组对应的项相同
	Data       []map[string]any `json:"data"`        // 具体结果数据 结果为该项所对应的对象条目格式
}

type ResultData struct {
}

// IntergratedSearch 综合搜索（web端）
// https://socialsisteryi.github.io/bilibili-API-collect/docs/search/search_request.html#%E7%BB%BC%E5%90%88%E6%90%9C%E7%B4%A2-web%E7%AB%AF
func (c *Client) IntergratedSearch(param SearchParam) (*SearchRespData, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/wbi/search/all/v2"
	)
	return execute[*SearchRespData](c, method, url, param, fillWbiHandler(c.wbi, c.GetCookies()))
}
