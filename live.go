package bilibili

import (
	"github.com/go-resty/resty/v2"
)

type GetLiveRoomInfoParam struct {
	RoomId int `json:"room_id"` // 直播间号。可以为短号
}

type Frame struct {
	Name       string `json:"name"`         // 名称
	Value      string `json:"value"`        // 值
	Position   int    `json:"position"`     // 位置
	Desc       string `json:"desc"`         // 描述
	Area       int    `json:"area"`         // 分区
	AreaOld    int    `json:"area_old"`     // 旧分区
	BgColor    string `json:"bg_color"`     // 背景色
	BgPic      string `json:"bg_pic"`       // 背景图
	UseOldArea bool   `json:"use_old_area"` // 是否旧分区号
}

type Badge struct {
	Name     string `json:"name"`     // 类型。v_person: 个人认证(黄) 。 v_company: 企业认证(蓝)
	Position int    `json:"position"` // 位置
	Value    string `json:"value"`    // 值
	Desc     string `json:"desc"`     // 描述
}

type NewPendants struct {
	Frame       *Frame `json:"frame"`        // 头像框
	MobileFrame *Frame `json:"mobile_frame"` // 同上。手机版, 结构一致, 可能null
	Badge       *Badge `json:"badge"`        // 大v
	MobileBadge *Badge `json:"mobile_badge"` // 同上。手机版, 结构一致, 可能null
}

type StudioInfo struct {
	Status     int   `json:"status"`
	MasterList []any `json:"master_list"`
}

type LiveRoomInfo struct {
	Uid                  int         `json:"uid"`                // 主播mid
	RoomId               int         `json:"room_id"`            // 直播间长号
	ShortId              int         `json:"short_id"`           // 直播间短号。为0是无短号
	Attention            int         `json:"attention"`          // 关注数量
	Online               int         `json:"online"`             // 观看人数
	IsPortrait           bool        `json:"is_portrait"`        // 是否竖屏
	Description          string      `json:"description"`        // 描述
	LiveStatus           int         `json:"live_status"`        // 直播状态。0：未开播。1：直播中。2：轮播中
	AreaId               int         `json:"area_id"`            // 分区id
	ParentAreaId         int         `json:"parent_area_id"`     // 父分区id
	ParentAreaName       string      `json:"parent_area_name"`   // 父分区名称
	OldAreaId            int         `json:"old_area_id"`        // 旧版分区id
	Background           string      `json:"background"`         // 背景图片链接
	Title                string      `json:"title"`              // 标题
	UserCover            string      `json:"user_cover"`         // 封面
	Keyframe             string      `json:"keyframe"`           // 关键帧。用于网页端悬浮展示
	IsStrictRoom         bool        `json:"is_strict_room"`     // 未知。未知
	LiveTime             string      `json:"live_time"`          // 直播开始时间。YYYY-MM-DD HH:mm:ss
	Tags                 string      `json:"tags"`               // 标签。','分隔
	IsAnchor             int         `json:"is_anchor"`          // 未知。未知
	RoomSilentType       string      `json:"room_silent_type"`   // 禁言状态
	RoomSilentLevel      int         `json:"room_silent_level"`  // 禁言等级
	RoomSilentSecond     int         `json:"room_silent_second"` // 禁言时间。单位是秒
	AreaName             string      `json:"area_name"`          // 分区名称
	Pardants             string      `json:"pardants"`           // 未知。未知
	AreaPardants         string      `json:"area_pardants"`      // 未知。未知
	HotWords             []string    `json:"hot_words"`          // 热词
	HotWordsStatus       int         `json:"hot_words_status"`   // 热词状态
	Verify               string      `json:"verify"`             // 未知。未知
	NewPendants          NewPendants `json:"new_pendants"`       // 头像框\大v
	UpSession            string      `json:"up_session"`         // 未知
	PkStatus             int         `json:"pk_status"`          // pk状态
	PkId                 int         `json:"pk_id"`              // pk id
	BattleId             int         `json:"battle_id"`          // 未知
	AllowChangeAreaTime  int         `json:"allow_change_area_time"`
	AllowUploadCoverTime int         `json:"allow_upload_cover_time"`
	StudioInfo           StudioInfo  `json:"studio_info"`
}

// GetLiveRoomInfo 获取直播间信息
func (c *Client) GetLiveRoomInfo(param GetLiveRoomInfoParam) (*LiveRoomInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.live.bilibili.com/room/v1/Room/get_info"
	)
	return execute[*LiveRoomInfo](c, method, url, param)
}

type UpdateLiveRoomTitleParam struct {
	RoomId int    `json:"room_id"`                                   // 直播间id。必须为自己的直播间id
	Title  string `json:"title,omitempty" request:"query,omitempty"` // 直播间标题。最大20字符
}

// UpdateLiveRoomTitle 更新直播间标题
func (c *Client) UpdateLiveRoomTitle(param UpdateLiveRoomTitleParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.live.bilibili.com/room/v1/Room/update"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type StartLiveParam struct {
	RoomId   int    `json:"room_id"`  // 直播间id。必须为自己的直播间id
	AreaV2   int    `json:"area_v2"`  // 直播分区id（子分区id）。详见[直播分区](live_area.md)
	Platform string `json:"platform"` // 直播平台。web端：。bililink：android_link
}

type Rtmp struct {
	Addr     string `json:"addr"`     // RTMP推流（发送）地址。**重要**
	Code     string `json:"code"`     // RTMP推流参数（密钥）。**重要**
	NewLink  string `json:"new_link"` // 获取CDN推流ip地址重定向信息的url。没啥用
	Provider string `json:"provider"` // ？？？。作用尚不明确
}

type Protocol struct {
	Protocol string `json:"protocol"` // rtmp。作用尚不明确
	Addr     string `json:"addr"`     // RTMP推流（发送）地址
	Code     string `json:"code"`     // RTMP推流参数（密钥）
	NewLink  string `json:"new_link"` // 获取CDN推流ip地址重定向信息的url
	Provider string `json:"provider"` // txy。作用尚不明确
}

type Notice struct {
	Type       int    `json:"type"`        // 1。作用尚不明确
	Status     int    `json:"status"`      // 0。作用尚不明确
	Title      string `json:"title"`       // 空。作用尚不明确
	Msg        string `json:"msg"`         // 空。作用尚不明确
	ButtonText string `json:"button_text"` // 空。作用尚不明确
	ButtonUrl  string `json:"button_url"`  // 空。作用尚不明确
}

type StartLiveResult struct {
	Change    int        `json:"change"`    // 是否改变状态。0：未改变。1：改变
	Status    string     `json:"status"`    // LIVE
	RoomType  int        `json:"room_type"` // 0。作用尚不明确
	Rtmp      Rtmp       `json:"rtmp"`      // RTMP推流地址信息
	Protocols []Protocol `json:"protocols"` // ？？？。作用尚不明确
	TryTime   string     `json:"try_time"`  // ？？？。作用尚不明确
	LiveKey   string     `json:"live_key"`  // ？？？。作用尚不明确
	Notice    Notice     `json:"notice"`    // ？？？。作用尚不明确
}

// StartLive 开始直播
func (c *Client) StartLive(param StartLiveParam) (*StartLiveResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.live.bilibili.com/room/v1/Room/startLive"
	)
	return execute[*StartLiveResult](c, method, url, param, fillCsrf(c))
}

type StopLiveParam struct {
	RoomId int `json:"room_id"` // 直播间id。必须为自己的直播间id
}

type StopLiveResult struct {
	Change int    `json:"change"` // 是否改变状态。0：未改变。1：改变
	Status string `json:"status"` // PREPARING
}

// StopLive 关闭直播
func (c *Client) StopLive(param StopLiveParam) (*StopLiveResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.live.bilibili.com/room/v1/Room/stopLive"
	)
	return execute[*StopLiveResult](c, method, url, param, fillCsrf(c))
}

type SubLiveArea struct {
	Id         string `json:"id"`          // 子分区id
	ParentId   string `json:"parent_id"`   // 父分区id
	OldAreaId  string `json:"old_area_id"` // 旧分区id
	Name       string `json:"name"`        // 子分区名
	ActId      string `json:"act_id"`      // 0。**作用尚不明确**
	PkStatus   string `json:"pk_status"`   // ？？？。**作用尚不明确**
	HotStatus  int    `json:"hot_status"`  // 是否为热门分区。0：否。1：是
	LockStatus string `json:"lock_status"` // 0。**作用尚不明确**
	Pic        string `json:"pic"`         // 子分区标志图片url
	ParentName string `json:"parent_name"` // 父分区名
	AreaType   int    `json:"area_type"`
}

type LiveAreaList struct {
	Id   int           `json:"id"`   // 父分区id
	Name string        `json:"name"` // 父分区名
	List []SubLiveArea `json:"list"` // 子分区列表
}

// GetLiveAreaList 获取全部直播间分区列表
func (c *Client) GetLiveAreaList() ([]LiveAreaList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.live.bilibili.com/room/v1/Area/getList"
	)
	return execute[[]LiveAreaList](c, method, url, nil)
}
