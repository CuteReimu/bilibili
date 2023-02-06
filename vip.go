package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

// ReceiveVipPrivilege 兑换大会员卡券，1：B币券，2：会员购优惠券，3：漫画福利券，4：会员购包邮券，5：漫画商城优惠券
func ReceiveVipPrivilege(privilegeType int) error {
	return std.ReceiveVipPrivilege(privilegeType)
}
func (c *Client) ReceiveVipPrivilege(privilegeType int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"type": strconv.Itoa(privilegeType),
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/vip/privilege/receive")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "兑换大会员卡券")
	return err
}

// SignVipScore 大积分签到
func SignVipScore() error {
	return std.SignVipScore()
}
func (c *Client) SignVipScore() error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      "https://www.bilibili.com",
	}).SetQueryParam("csrf", biliJct).Post("https://api.bilibili.com/pgc/activity/score/task/sign")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "大积分签到")
	return err
}

type VipPrivilege struct {
	List []struct { // 卡券信息列表
		Type            int `json:"type"`              // 卡券类型，1：B币券，2：会员购优惠券，3：漫画福利券，4：会员购包邮券，5：漫画商城优惠券
		State           int `json:"state"`             // 兑换状态，0：当月未兑换，1：已兑换
		ExpireTime      int `json:"expire_time"`       // 本轮卡券过期时间戳（秒）
		VipType         int `json:"vip_type"`          // 可兑换的会员类型，2：年度大会员
		NextReceiveDays int `json:"next_receive_days"` // 距下一轮兑换剩余天数
		PeriodEndUnix   int `json:"period_end_unix"`   // 下一轮兑换开始时间戳（秒）
	} `json:"list"`
	IsShortVip    bool `json:"is_short_vip"`
	IsFreightOpen bool `json:"is_freight_open"`
}

// GetVipPrivilege 卡券状态查询
func GetVipPrivilege() (*VipPrivilege, error) {
	return std.GetVipPrivilege()
}
func (c *Client) GetVipPrivilege() (*VipPrivilege, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get("https://api.bilibili.com/x/vip/privilege/my")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "卡券状态查询")
	if err != nil {
		return nil, err
	}
	var ret *VipPrivilege
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
