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
