package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type SearchDynamicAtResult struct {
	Groups []struct { // 内容分组
		GroupType int        `json:"group_type"` // 2：我的关注，4：其他
		GroupName string     `json:"group_name"` // 分组名字
		Items     []struct { // 用户信息
			Uid                int    `json:"uid"`   // 用户id
			Uname              string `json:"uname"` // 用户昵称
			Face               string `json:"face"`  // 用户头像url
			Fans               int    `json:"fans"`  // 用户粉丝数
			OfficialVerifyType int    `json:"official_verify_type"`
		} `json:"items"`
	} `json:"groups"`
	Gt int `json:"_gt_"` // 固定值0
}

// SearchDynamicAt 根据关键字搜索用户(at别人时的填充列表)
func SearchDynamicAt(uid int, keyword string) (*SearchDynamicAtResult, error) {
	return std.SearchDynamicAt(uid, keyword)
}
func (c *Client) SearchDynamicAt(uid int, keyword string) (*SearchDynamicAtResult, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"uid":     strconv.Itoa(uid),
		"keyword": keyword,
	}).Get("https://api.vc.bilibili.com/dynamic_mix/v1/dynamic_mix/at_search")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "根据关键字搜索用户")
	if err != nil {
		return nil, err
	}
	var ret *SearchDynamicAtResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
