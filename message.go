package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
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
