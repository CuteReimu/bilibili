package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"math/rand"
	"strconv"
	"time"
)

type UnreadMessage struct {
	At     int `json:"at"`      // 未读at数
	Chat   int `json:"chat"`    // 固定值0，作用尚不明确
	Like   int `json:"like"`    // 未读点赞数
	Reply  int `json:"reply"`   // 未读回复数
	SysMsg int `json:"sys_msg"` // 未读系统通知数
	Up     int `json:"up"`      // UP主助手信息数
}

// GetUnreadMessage 获取未读消息数
func GetUnreadMessage() (*UnreadMessage, error) {
	return std.GetUnreadMessage()
}
func (c *Client) GetUnreadMessage() (*UnreadMessage, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get("https://api.bilibili.com/x/msgfeed/unread")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取未读消息数")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *UnreadMessage
	err = json.Unmarshal(data, &ret)
	return ret, err
}

type UnreadPrivateMessage struct {
	UnfollowUnread int `json:"unfollow_unread"` // 未关注用户未读私信数
	FollowUnread   int `json:"follow_unread"`   // 已关注用户未读私信数
	Gt             int `json:"_gt_"`            // 固定值0，作用尚不明确
}

// GetUnreadPrivateMessage 获取未读私信数
func GetUnreadPrivateMessage() (*UnreadPrivateMessage, error) {
	return std.GetUnreadPrivateMessage()
}
func (c *Client) GetUnreadPrivateMessage() (*UnreadPrivateMessage, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get("https://api.vc.bilibili.com/session_svr/v1/session_svr/single_unread")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取未读私信数")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *UnreadPrivateMessage
	err = json.Unmarshal(data, &ret)
	return ret, err
}

var deviceId string

func init() {
	b := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}
	s := []byte("xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx")
	randBytes := make([]byte, len(s))
	_, _ = rand.Read(randBytes)
	for i := range s {
		if '-' == s[i] || '4' == s[i] {
			continue
		}
		j := randBytes[i] % 16
		if 'x' == s[i] {
			s[i] = b[j]
		} else {
			s[i] = b[3&j|8]
		}
	}
	deviceId = string(s)
}

// SendPrivateMessageText 发送私信（文字消息）
func SendPrivateMessageText(senderUid, receiverId int, content string) (int, string, error) {
	return std.SendPrivateMessageText(senderUid, receiverId, content)
}
func (c *Client) SendPrivateMessageText(senderUid, receiverId int, content string) (int, string, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return 0, "", errors.New("B站登录过期")
	}
	contentBytes, err := json.Marshal(map[string]interface{}{"content": content})
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"msg[sender_uid]":    strconv.Itoa(senderUid),
		"msg[receiver_id]":   strconv.Itoa(receiverId),
		"msg[receiver_type]": "1",
		"msg[msg_type]":      "1",
		"msg[dev_id]":        deviceId,
		"msg[timestamp]":     strconv.FormatInt(time.Now().Unix(), 10),
		"msg[content]":       string(contentBytes),
		"csrf":               biliJct,
	}).Post("https://api.vc.bilibili.com/web_im/v1/web_im/send_msg")
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return 0, "", errors.Errorf("发送私信失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return 0, "", errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return 0, "", errors.Errorf("发送私信失败，返回值：%d，返回信息：%s", code, res.Get("message").String())
	}
	return int(res.Get("data.msg_key").Int()), res.Get("data.msg_content").String(), err
}

// SendPrivateMessageImage 发送私信（图片消息）
func SendPrivateMessageImage(senderUid, receiverId int, url string) (int, string, error) {
	return std.SendPrivateMessageImage(senderUid, receiverId, url)
}
func (c *Client) SendPrivateMessageImage(senderUid, receiverId int, url string) (int, string, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return 0, "", errors.New("B站登录过期")
	}
	contentBytes, err := json.Marshal(map[string]interface{}{"url": url})
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"msg[sender_uid]":    strconv.Itoa(senderUid),
		"msg[receiver_id]":   strconv.Itoa(receiverId),
		"msg[receiver_type]": "1",
		"msg[msg_type]":      "2",
		"msg[dev_id]":        deviceId,
		"msg[timestamp]":     strconv.FormatInt(time.Now().Unix(), 10),
		"msg[content]":       string(contentBytes),
		"csrf":               biliJct,
	}).Post("https://api.vc.bilibili.com/web_im/v1/web_im/send_msg")
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return 0, "", errors.Errorf("发送私信失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return 0, "", errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return 0, "", errors.Errorf("发送私信失败，返回值：%d，返回信息：%s", code, res.Get("message").String())
	}
	return int(res.Get("data.msg_key").Int()), res.Get("data.msg_content").String(), err
}

// SendPrivateMessageRecall 发送私信（撤回消息）
func SendPrivateMessageRecall(senderUid, receiverId, msgKey int) (int, string, error) {
	return std.SendPrivateMessageRecall(senderUid, receiverId, msgKey)
}
func (c *Client) SendPrivateMessageRecall(senderUid, receiverId, msgKey int) (int, string, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return 0, "", errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"msg[sender_uid]":    strconv.Itoa(senderUid),
		"msg[receiver_id]":   strconv.Itoa(receiverId),
		"msg[receiver_type]": "1",
		"msg[msg_type]":      "5",
		"msg[dev_id]":        deviceId,
		"msg[timestamp]":     strconv.FormatInt(time.Now().Unix(), 10),
		"msg[content]":       strconv.Itoa(msgKey),
		"csrf":               biliJct,
	}).Post("https://api.vc.bilibili.com/web_im/v1/web_im/send_msg")
	if err != nil {
		return 0, "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return 0, "", errors.Errorf("发送私信（撤回消息）失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return 0, "", errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return 0, "", errors.Errorf("发送私信（撤回消息）失败，返回值：%d，返回信息：%s", code, res.Get("message").String())
	}
	return int(res.Get("data.msg_key").Int()), res.Get("data.msg_content").String(), err
}

type SessionMessages struct {
	Messages []struct { // 聊天记录列表
		SenderUid      int    `json:"sender_uid"`                 // 发送者uid
		ReceiverType   int    `json:"receiver_type"`              // 1为用户，2为粉丝团
		ReceiverId     int    `json:"receiver_id"`                // 接收者uid
		MsgType        int    `json:"msg_type"`                   // 消息类型，1:文字消息，2:图片消息，5:撤回的消息，12、13:通知
		Content        string `json:"content"`                    // 消息内容
		MsgSeqno       int64  `json:"msg_seqno"`                  // 作用尚不明确
		Timestamp      int    `json:"timestamp"`                  // 消息发送时间戳
		AtUids         []int  `json:"at_uids"`                    // 作用尚不明确
		MsgKey         int64  `json:"msg_key"`                    // 作用尚不明确
		MsgStatus      int    `json:"msg_status"`                 // 固定值0
		NotifyCode     string `json:"notify_code"`                // 作用尚不明确
		NewFaceVersion int    `json:"new_face_version,omitempty"` // 作用尚不明确
	} `json:"messages"`
	HasMore  int        `json:"has_more"`  // 固定值0
	MinSeqno int64      `json:"min_seqno"` // 作用尚不明确
	MaxSeqno int64      `json:"max_seqno"` // 作用尚不明确
	EInfos   []struct { // 聊天表情列表
		Text string `json:"text"` // 表情名称
		Url  string `json:"url"`  // 表情链接
		Size int    `json:"size"` // 表情尺寸
	} `json:"e_infos"`
}

// GetSessionMessages 获取私信消息记录
func GetSessionMessages(talkerId, sessionType, size int, mobiApp string) (*SessionMessages, error) {
	return std.GetSessionMessages(talkerId, sessionType, size, mobiApp)
}
func (c *Client) GetSessionMessages(talkerId, sessionType, size int, mobiApp string) (*SessionMessages, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"talker_id":    strconv.Itoa(talkerId),
		"session_type": strconv.Itoa(sessionType),
	})
	if size != 0 {
		r.SetQueryParam("size", strconv.Itoa(size))
	}
	if len(mobiApp) > 0 {
		r.SetQueryParam("mobi_app", mobiApp)
	}
	resp, err := r.Get("https://api.vc.bilibili.com/svr_sync/v1/svr_sync/fetch_session_msgs")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取私信消息记录")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *SessionMessages
	err = json.Unmarshal(data, &ret)
	return ret, err
}

type SessionList struct {
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

// GetSessions 获取消息列表 session_type，1：系统，2：用户，3：应援团
//
// 参照 https://github.com/CuteReimu/bilibili/issues/8
func GetSessions(sessionType int, mobiApp string) (*SessionList, error) {
	return std.GetSessions(sessionType, mobiApp)
}
func (c *Client) GetSessions(sessionType int, mobiApp string) (*SessionList, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParam("session_type", strconv.Itoa(sessionType))
	if len(mobiApp) > 0 {
		r.SetQueryParam("mobi_app", mobiApp)
	}
	resp, err := r.Get("https://api.vc.bilibili.com/session_svr/v1/session_svr/get_sessions")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取消息列表")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *SessionList
	err = json.Unmarshal(data, &ret)
	return ret, err
}
