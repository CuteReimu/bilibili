package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type GetRoomInfoResult struct {
	Code    int      `json:"code"`    // 返回值，0表示成功
	Msg     string   `json:"msg"`     // 错误信息
	Message string   `json:"message"` // 错误信息
	Data    struct { // 信息本体
		Uid              int      `json:"uid"`
		RoomId           int      `json:"room_id"` // 直播间id
		ShortId          int      `json:"short_id"`
		Attention        int      `json:"attention"`
		Online           int      `json:"online"` // 直播间人气，值为上次直播时刷新
		IsPortrait       bool     `json:"is_portrait"`
		Description      string   `json:"description"`
		LiveStatus       int      `json:"live_status"` // 直播状态，0：未开播，1：直播中
		AreaId           int      `json:"area_id"`
		ParentAreaId     int      `json:"parent_area_id"`
		ParentAreaName   string   `json:"parent_area_name"`
		OldAreaId        int      `json:"old_area_id"`
		Background       string   `json:"background"`
		Title            string   `json:"title"` // 直播间标题
		UserCover        string   `json:"user_cover"`
		Keyframe         string   `json:"keyframe"`
		IsStrictRoom     bool     `json:"is_strict_room"`
		LiveTime         string   `json:"live_time"`
		Tags             string   `json:"tags"`
		IsAnchor         int      `json:"is_anchor"`
		RoomSilentType   string   `json:"room_silent_type"`
		RoomSilentLevel  int      `json:"room_silent_level"`
		RoomSilentSecond int      `json:"room_silent_second"`
		AreaName         string   `json:"area_name"`
		Pendants         string   `json:"pendants"`
		AreaPendants     string   `json:"area_pendants"`
		HotWords         []string `json:"hot_words"`
		HotWordsStatus   int      `json:"hot_words_status"`
		Verify           string   `json:"verify"`
		NewPendants      struct {
			Frame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"frame"`
			//Badge       interface{} `json:"badge"`
			MobileFrame struct {
				Name       string `json:"name"`
				Value      string `json:"value"`
				Position   int    `json:"position"`
				Desc       string `json:"desc"`
				Area       int    `json:"area"`
				AreaOld    int    `json:"area_old"`
				BgColor    string `json:"bg_color"`
				BgPic      string `json:"bg_pic"`
				UseOldArea bool   `json:"use_old_area"`
			} `json:"mobile_frame"`
			//MobileBadge interface{} `json:"mobile_badge"`
		} `json:"new_pendants"`
		UpSession            string `json:"up_session"`
		PkStatus             int    `json:"pk_status"`
		PkId                 int    `json:"pk_id"`
		BattleId             int    `json:"battle_id"`
		AllowChangeAreaTime  int    `json:"allow_change_area_time"`
		AllowUploadCoverTime int    `json:"allow_upload_cover_time"`
		StudioInfo           struct {
			Status int `json:"status"`
			//MasterList []interface{} `json:"master_list"`
		} `json:"studio_info"`
	} `json:"data"`
}

// GetRoomInfo 获取直播间状态
func GetRoomInfo(roomId int) (*GetRoomInfoResult, error) {
	return std.GetRoomInfo(roomId)
}
func (c *Client) GetRoomInfo(roomId int) (*GetRoomInfoResult, error) {
	resp, err := c.resty().R().SetQueryParam("id", strconv.Itoa(roomId)).Get("https://api.live.bilibili.com/room/v1/Room/get_info")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取直播间状态失败，status code：%d", resp.StatusCode())
	}
	var ret *GetRoomInfoResult
	err = json.Unmarshal(resp.Body(), &ret)
	return ret, errors.WithStack(err)
}

type UpdateLiveResult struct {
	Code    int    `json:"code"`    // 返回值，0表示成功
	Msg     string `json:"msg"`     // 错误信息
	Message string `json:"message"` // 错误信息
}

// UpdateLive 更新直播间标题
func UpdateLive(roomId int, title string) (*UpdateLiveResult, error) {
	return std.UpdateLive(roomId, title)
}
func (c *Client) UpdateLive(roomId int, title string) (*UpdateLiveResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"room_id": strconv.Itoa(roomId),
		"title":   title,
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/update")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("修改直播间标题失败，status code：%d", resp.StatusCode())
	}
	var ret *UpdateLiveResult
	err = json.Unmarshal(resp.Body(), &ret)
	return ret, errors.WithStack(err)
}

type StartLiveResult struct {
	Code    int      `json:"code"`    // 返回值，0表示成功
	Msg     string   `json:"msg"`     // 错误信息
	Message string   `json:"message"` // 错误信息
	Data    struct { // 信息本体
		Change   int      `json:"change"`    // 是否改变状态，0：未改变，1：改变
		Status   string   `json:"status"`    // 固定值LIVE
		RoomType int      `json:"room_type"` // 固定值0，作用尚不明确
		Rtmp     struct { // RTMP推流地址信息
			Addr     string `json:"addr"`     // RTMP推流（发送）地址，重要
			Code     string `json:"code"`     // RTMP推流参数（密钥），重要
			NewLink  string `json:"new_link"` // 获取CDN推流ip地址重定向信息的url，没啥用
			Provider string `json:"provider"` // 作用尚不明确
		} `json:"rtmp"`
		Protocols []struct { // 作用尚不明确
			Protocol string `json:"protocol"` // 固定值RTMP，作用尚不明确
			Addr     string `json:"addr"`     // RTMP推流（发送）地址
			Code     string `json:"code"`     // RTMP推流参数（密钥）
			NewLink  string `json:"new_link"` // 获取CDN推流ip地址重定向信息的url
			Provider string `json:"provider"` // 固定值txy，作用尚不明确
		} `json:"protocols"`
		TryTime string   `json:"try_time"` // 作用尚不明确
		LiveKey string   `json:"live_key"` // 作用尚不明确
		Notice  struct { // 作用尚不明确
			Type       int    `json:"type"`        // 固定值1，作用尚不明确
			Status     int    `json:"status"`      // 固定值0，作用尚不明确
			Title      string `json:"title"`       // 固定值空，作用尚不明确
			Msg        string `json:"msg"`         // 固定值空，作用尚不明确
			ButtonText string `json:"button_text"` // 固定值空，作用尚不明确
			ButtonUrl  string `json:"button_url"`  // 固定值空，作用尚不明确
		} `json:"notice"`
	} `json:"data"`
}

// StartLive 开始直播
func StartLive(roomId, area int) (*StartLiveResult, error) {
	return std.StartLive(roomId, area)
}
func (c *Client) StartLive(roomId, area int) (*StartLiveResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"room_id":  strconv.Itoa(roomId),
		"platform": "pc",
		"area_v2":  strconv.Itoa(area),
		"csrf":     biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/startLive")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("开启直播间失败，status code：%d", resp.StatusCode())
	}
	var ret *StartLiveResult
	err = json.Unmarshal(resp.Body(), &ret)
	return ret, errors.WithStack(err)
}

type StopLiveResult struct {
	Code    int      `json:"code"`    // 返回值，0表示成功
	Msg     string   `json:"msg"`     // 错误信息
	Message string   `json:"message"` // 错误信息
	Data    struct { // 信息本体
		Change int    `json:"change"` // 是否改变状态，0：未改变，1：改变
		Status string `json:"status"` // 固定值PREPARING
	} `json:"data"`
}

// StopLive 关闭直播
func StopLive(roomId int) (*StopLiveResult, error) {
	return std.StopLive(roomId)
}
func (c *Client) StopLive(roomId int) (*StopLiveResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetQueryParams(map[string]string{
		"room_id": strconv.Itoa(roomId),
		"csrf":    biliJct,
	}).Post("https://api.live.bilibili.com/room/v1/Room/stopLive")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("关闭直播间失败，status code：%d", resp.StatusCode())
	}
	var ret *StopLiveResult
	err = json.Unmarshal(resp.Body(), &ret)
	return ret, errors.WithStack(err)
}
