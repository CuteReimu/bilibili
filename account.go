package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// AccountInformation  自己账号相关的简单信息
type AccountInformation struct {
	Mid      int64  `json:"mid"`       //我的mid
	Uname    string `json:"uname"`     //我的昵称
	Userid   string `json:"userid"`    //我的用户名
	Sign     string `json:"sign"`      //我的签名
	Birthday string `json:"birthday"`  //我的生日
	Sex      string `json:"sex"`       //我的性别
	NickFree bool   `json:"nick_free"` //false：设置过昵称 true：未设置昵称
	Rank     string `json:"rank"`      //我的会员等级
}

func (c *Client) GetAccountInformation() (*AccountInformation, error) {
	request := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParam("version", "1")

	var accountInfo AccountInformation

	resp, err := request.Get("https://api.bilibili.com/x/member/web/account")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取我的信息")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &accountInfo)

	return &accountInfo, err
}

// GetAccountInformation 获取我的信息 无参数
func GetAccountInformation() (*AccountInformation, error) {
	return std.GetAccountInformation()
}
