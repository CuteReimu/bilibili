package bilibili

import (
	"crypto/rand"
	"github.com/go-resty/resty/v2"
)

type UnreadMessage struct {
	At     int `json:"at"`      // 未读at数
	Chat   int `json:"chat"`    // 0。作用尚不明确
	Like   int `json:"like"`    // 未读点赞数
	Reply  int `json:"reply"`   // 未读回复数
	SysMsg int `json:"sys_msg"` // 未读系统通知数
	Up     int `json:"up"`      // UP主助手信息数
}

// GetUnreadMessage 获取未读消息数
func (c *Client) GetUnreadMessage() (*UnreadMessage, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/msgfeed/unread"
	)
	return execute[*UnreadMessage](c, method, url, nil)
}

type UnreadPrivateMessage struct {
	UnfollowUnread int `json:"unfollow_unread"` // 未关注用户未读私信数
	FollowUnread   int `json:"follow_unread"`   // 已关注用户未读私信数
	Gt             int `json:"_gt_"`            // 0
}

// GetUnreadPrivateMessage 获取未读私信数
func (c *Client) GetUnreadPrivateMessage() (*UnreadPrivateMessage, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.vc.bilibili.com/session_svr/v1/session_svr/single_unread"
	)
	return execute[*UnreadPrivateMessage](c, method, url, nil)
}

var deviceId string

func init() {
	b := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	s := []byte("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	randBytes := make([]byte, len(s))
	_, _ = rand.Read(randBytes)
	for i := range s {
		if s[i] == '-' || s[i] == '4' {
			continue
		}
		j := randBytes[i] % 16
		if s[i] == 'x' {
			s[i] = b[j]
		} else {
			s[i] = b[3&j|8]
		}
	}
	deviceId = string(s)
}

type SendPrivateMessageParam struct {
	SenderUid      int `json:"msg[sender_uid]"`                                           // 发送者mid
	ReceiverId     int `json:"msg[receiver_id]"`                                          // 接收者mid
	ReceiverType   int `json:"msg[receiver_type]"`                                        // 1。固定为1
	MsgType        int `json:"msg[msg_type]"`                                             // 消息类型。1:发送文字。2:发送图片。5:撤回消息
	MsgStatus      int `json:"msg[msg_status],omitempty" request:"query,omitempty"`       // 0
	Timestamp      int `json:"msg[timestamp]"`                                            // 时间戳（秒）
	NewFaceVersion int `json:"msg[new_face_version],omitempty" request:"query,omitempty"` // 表情包版本
	Content        any `json:"msg[content]"`                                              // 消息内容。发送文字时：str<br />撤回消息时：num
}

type SendPrivateMessageResult struct {
	MsgKey      int    `json:"msg_key"`       // 消息唯一id
	MsgContent  string `json:"msg_content"`   // 发送的消息
	KeyHitInfos any    `json:"key_hit_infos"` // 作用尚不明确
}

// SendPrivateMessage 发送私信（文字消息）
func (c *Client) SendPrivateMessage(param SendPrivateMessageParam) (*SendPrivateMessageResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.vc.bilibili.com/web_im/v1/web_im/send_msg"
	)
	return execute[*SendPrivateMessageResult](c, method, url, param, fillCsrf(c), func(request *resty.Request) error {
		request.SetQueryParam("msg[dev_id]", deviceId)
		return nil
	})
}

type GetPrivateMessageRecordsParam struct {
	TalkerId       int    `json:"talker_id"`                                            // 聊天对象的uid
	SenderDeviceId int    `json:"sender_device_id,omitempty" request:"query,omitempty"` // 发送者设备。1
	SessionType    int    `json:"session_type"`                                         // 聊天对象的类型。1为用户，2为粉丝团
	Size           int    `json:"size,omitempty" request:"query,omitempty"`             // 列出消息条数。默认是20，最大为200
	Build          int    `json:"build,omitempty" request:"query,omitempty"`            // 未知。默认是0
	MobiApp        string `json:"mobi_app,omitempty" request:"query,omitempty"`         // 设备。web
	BeginSeqno     int    `json:"begin_seqno,omitempty" request:"query,omitempty"`      // 开始的序列号。默认0为全部
	EndSeqno       int    `json:"end_seqno,omitempty" request:"query,omitempty"`        // 结束的序列号。默认0为全部
}

type Message struct {
	SenderUid      int    `json:"sender_uid"`       // 发送者uid。注意名称是sender_uid
	ReceiverType   int    `json:"receiver_type"`    // 与session_type对应。1为用户，2为粉丝团
	ReceiverId     int    `json:"receiver_id"`      // 接收者uid。注意名称是receiver_id
	MsgType        int    `json:"msg_type"`         // 消息类型。1:文字消息。2:图片消息。5:撤回的消息。12、13:通知
	Content        string `json:"content"`          // 消息内容。此处存在设计缺陷
	MsgSeqno       int    `json:"msg_seqno"`        // 消息序列号，保证按照时间顺序从小到大
	Timestamp      int    `json:"timestamp"`        // 消息发送时间戳
	AtUids         []int  `json:"at_uids"`          // 未知
	MsgKey         int    `json:"msg_key"`          // 未知
	MsgStatus      int    `json:"msg_status"`       // 消息状态。0
	NotifyCode     string `json:"notify_code"`      // 未知
	NewFaceVersion int    `json:"new_face_version"` // 表情包版本。0或者没有是旧版，此时b站会自动转换成新版表情包，例如[doge] -> [tv_doge]；1是新版
}

type EInfo struct {
	Text string `json:"text"` // 表情名称
	Uri  string `json:"uri"`  // 表情链接
	Size int    `json:"size"` // 表情尺寸。1
}

type PrivateMessageRecords struct {
	Messages []Message `json:"messages"`  // 聊天记录列表
	HasMore  int       `json:"has_more"`  // 0
	MinSeqno uint64    `json:"min_seqno"` // 所有消息最小的序列号（最早）
	MaxSeqno uint64    `json:"max_seqno"` // 所有消息最大的序列号（最晚）
	EInfos   []EInfo   `json:"e_infos"`   // 聊天表情列表
}

// GetPrivateMessageRecords 获取与聊天对象的私信消息记录
func (c *Client) GetPrivateMessageRecords(param GetPrivateMessageRecordsParam) (*PrivateMessageRecords, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.vc.bilibili.com/svr_sync/v1/svr_sync/fetch_session_msgs"
	)
	return execute[*PrivateMessageRecords](c, method, url, param)
}

type GetPrivateMessageListParam struct {
	SessionType int    `json:"session_type"`                                 // 1：系统，2：用户，3：应援团
	MobiApp     string `json:"mobi_app,omitempty" request:"query,omitempty"` // 设备
}

type PrivateMessageList struct {
	SessionList []struct {
		TalkerId    int64  `json:"talker_id"`
		SessionType int    `json:"session_type"`
		AtSeqno     int    `json:"at_seqno"`
		TopTs       int    `json:"top_ts"`
		GroupName   string `json:"group_name"`
		GroupCover  string `json:"group_cover"`
		IsFollow    int    `json:"is_follow"`
		IsDnd       int    `json:"is_dnd"`
		AckSeqno    int64  `json:"ack_seqno"`
		AckTs       int64  `json:"ack_ts"`
		SessionTs   int64  `json:"session_ts"`
		UnreadCount int    `json:"unread_count"`
		LastMsg     struct {
			SenderUid      int64  `json:"sender_uid"`
			ReceiverType   int    `json:"receiver_type"`
			ReceiverId     int    `json:"receiver_id"`
			MsgType        int    `json:"msg_type"`
			Content        string `json:"content"`
			MsgSeqno       int64  `json:"msg_seqno"`
			Timestamp      int    `json:"timestamp"`
			MsgKey         int64  `json:"msg_key"`
			MsgStatus      int    `json:"msg_status"`
			NotifyCode     string `json:"notify_code"`
			NewFaceVersion int    `json:"new_face_version,omitempty"`
		} `json:"last_msg"`
		GroupType         int   `json:"group_type"`
		CanFold           int   `json:"can_fold"`
		Status            int   `json:"status"`
		MaxSeqno          int64 `json:"max_seqno"`
		NewPushMsg        int   `json:"new_push_msg"`
		Setting           int   `json:"setting"`
		IsGuardian        int   `json:"is_guardian"`
		IsIntercept       int   `json:"is_intercept"`
		IsTrust           int   `json:"is_trust"`
		SystemMsgType     int   `json:"system_msg_type"`
		LiveStatus        int   `json:"live_status"`
		BizMsgUnreadCount int   `json:"biz_msg_unread_count"`
		AccountInfo       struct {
			Name   string `json:"name"`
			PicUrl string `json:"pic_url"`
		} `json:"account_info,omitempty"`
	} `json:"session_list"`
	HasMore             int              `json:"has_more"`
	AntiDisturbCleaning bool             `json:"anti_disturb_cleaning"`
	IsAddressListEmpty  int              `json:"is_address_list_empty"`
	SystemMsg           map[string]int64 `json:"system_msg"`
	ShowLevel           bool             `json:"show_level"`
}

// GetPrivateMessageList 获取消息列表 session_type，1：系统，2：用户，3：应援团
//
// 参照 https://github.com/CuteReimu/bilibili/issues/8
func (c *Client) GetPrivateMessageList(param GetPrivateMessageListParam) (*PrivateMessageList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.vc.bilibili.com/session_svr/v1/session_svr/get_sessions"
	)
	return execute[*PrivateMessageList](c, method, url, param)
}
