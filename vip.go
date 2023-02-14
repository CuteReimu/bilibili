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

type VipCenterInfo struct {
	User struct { // 用户信息
		Account struct { // 账号基本信息
			Mid            int    `json:"mid"`      // 用户mid
			Name           string `json:"name"`     // 昵称
			Sex            string `json:"sex"`      // 性别，男/女/保密
			Face           string `json:"face"`     // 头像url
			Sign           string `json:"sign"`     // 签名
			Rank           int    `json:"rank"`     // 等级
			Birthday       int    `json:"birthday"` // 生日时间戳，单位：秒
			IsFakeAccount  int    `json:"is_fake_account"`
			IsDeleted      int    `json:"is_deleted"`       // 是否注销，0：正常，1：注销
			InRegAudit     int    `json:"in_reg_audit"`     // 是否注册审核，0：正常，1：审核
			IsSeniorMember int    `json:"is_senior_member"` // 是否转正，0：未转正，1：正式会员
		} `json:"account"`
		Vip struct { // 账号会员信息
			Mid        int      `json:"mid"`          // 用户mid
			VipType    int      `json:"vip_type"`     // 会员类型，0：无，1：月大会员，2：年度及以上大会员
			VipStatus  int      `json:"vip_status"`   // 会员状态，0：无，1：有
			VipDueDate int64    `json:"vip_due_date"` // 会员过期时间戳，单位：毫秒
			VipPayType int      `json:"vip_pay_type"` // 支付类型，0：未支付（常见于官方账号），1：已支付（以正常渠道获取的大会员均为此值）
			ThemeType  int      `json:"theme_type"`
			Label      struct { // 会员标签
				Text                  string `json:"text"`          // 会员类型文案，大会员，年度大会员，十年大会员，百年大会员，最强绿鲤鱼
				LabelTheme            string `json:"label_theme"`   // 会员标签，vip，annual_vip，ten_annual_vip，hundred_annual_vip，fools_day_hundred_annual_vip
				TextColor             string `json:"text_color"`    // 会员标签文本颜色
				BgStyle               int    `json:"bg_style"`      // 固定值1
				BgColor               string `json:"bg_color"`      // 会员标签背景颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
				BorderColor           string `json:"border_color"`  // 会员标签边框颜色
				UseImgLabel           bool   `json:"use_img_label"` // 固定值true
				ImgLabelUriHans       string `json:"img_label_uri_hans"`
				ImgLabelUriHant       string `json:"img_label_uri_hant"`
				ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"` // 大会员牌子图片简体版
				ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"` // 大会员牌子图片繁体版
			} `json:"label"`
			AvatarSubscript    int         `json:"avatar_subscript"` // 是否显示会员图标，0：不显示，1：显示
			AvatarSubscriptUrl string      `json:"avatar_subscript_url"`
			NicknameColor      string      `json:"nickname_color"` // 会员昵称颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
			IsNewUser          bool        `json:"is_new_user"`
			TipMaterial        interface{} `json:"tip_material"`
		} `json:"vip"`
		TV struct { // 电视会员信息
			Type       int   `json:"type"`         // 电视大会员类型，0：无，1：月大会员，2：年度及以上大会员
			VipPayType int   `json:"vip_pay_type"` // 电视大支付类型，0：未支付（常见于官方账号），1：已支付（以正常渠道获取的大会员均为此值）
			Status     int   `json:"status"`       // 电视大会员状态，0：无，1：有
			DueDate    int64 `json:"due_date"`     // 电视大会员过期时间戳，单位：毫秒
		} `json:"tv"`
		BackgroundImageSmall string   `json:"background_image_small"`
		BackgroundImageBig   string   `json:"background_image_big"`
		PanelTitle           string   `json:"panel_title"` // 用户昵称
		PanelSubTitle        string   `json:"panel_sub_title"`
		AvatarPendant        struct { // 用户头像框信息
			Image             string `json:"image"`               // 头像框url
			ImageEnhance      string `json:"image_enhance"`       // 头像框url，动态图
			ImageEnhanceFrame string `json:"image_enhance_frame"` // 动态头像框帧波普版url
		} `json:"avatar_pendant"`
		VipOverdueExplain    string `json:"vip_overdue_explain"` // 大会员提示文案
		TvOverdueExplain     string `json:"tv_overdue_explain"`  // 电视大会员提示文案
		AccountExceptionText string `json:"account_exception_text"`
		IsAutoRenew          bool   `json:"is_auto_renew"`    // 是否自动续费
		IsTvAutoRenew        bool   `json:"is_tv_auto_renew"` // 是否自动续费电视大会员
		SurplusSeconds       int    `json:"surplus_seconds"`  // 大会员到期剩余时间，单位：秒
		VipKeepTime          int    `json:"vip_keep_time"`    // 持续开通大会员时间，单位：秒
		Renew                struct {
			Text string `json:"text"`
			Link string `json:"link"`
		} `json:"renew"`
		Notice struct {
			Text             string `json:"text"`
			TvText           string `json:"tv_text"`
			Type             int    `json:"type"`
			CanClose         bool   `json:"can_close"`
			SurplusSeconds   int    `json:"surplus_seconds"`
			TvSurplusSeconds int    `json:"tv_surplus_seconds"`
		} `json:"notice"`
	} `json:"user"`
	Wallet struct { // 钱包信息
		Coupon            int  `json:"coupon"` // 当前B币券
		Point             int  `json:"point"`
		PrivilegeReceived bool `json:"privilege_received"`
	} `json:"wallet"`
	Privileges []struct { // 会员特权信息列表
		Id              int        `json:"id"`   // 特权父类id
		Name            string     `json:"name"` // 类型名称
		ChildPrivileges []struct { // 特权子类列表
			FirstId            int    `json:"first_id"`             // 特权父类id
			ReportId           string `json:"report_id"`            // 上报id
			Name               string `json:"name"`                 // 特权名称
			Desc               string `json:"desc"`                 // 特权简介文案
			Explain            string `json:"explain"`              // 特权介绍正文
			IconUrl            string `json:"icon_url"`             // 特权图标url
			IconGrayUrl        string `json:"icon_gray_url"`        // 特权图标灰色主题url
			BackgroundImageUrl string `json:"background_image_url"` // 背景图片url
			Link               string `json:"link"`                 // 特权介绍页url
			ImageUrl           string `json:"image_url"`            // 特权示例图url
			Type               int    `json:"type"`                 // 类型，目前为0
			HotType            int    `json:"hot_type"`             // 是否热门特权，0：普通特权，1：热门特权
			NewType            int    `json:"new_type"`             // 是否新特权，0：普通特权，1：新特权
			Id                 int    `json:"id"`                   // 特权子类id
		} `json:"child_privileges"`
	} `json:"privileges"`
	Welfare struct { // 福利信息
		Count int        `json:"count"` // 福利数
		List  []struct { // 福利项目列表
			Id          int    `json:"id"`           // 福利id
			Name        string `json:"name"`         // 福利名称
			HomepageUri string `json:"homepage_uri"` // 福利图片url
			BackdropUri string `json:"backdrop_uri"` // 福利图片banner url
			Tid         int    `json:"tid"`          // 目前为0
			Rank        int    `json:"rank"`         // 排列顺序
			ReceiveUri  string `json:"receive_uri"`  // 福利跳转页url
		} `json:"list"`
	} `json:"welfare"`
	RecommendPendants struct { // 推荐头像框信息
		List []struct { // 推荐头像框列表
			Id      int    `json:"id"`       // 头像框id
			Name    string `json:"name"`     // 头像框名称
			Image   string `json:"image"`    // 头像框图片url
			JumpUrl string `json:"jump_url"` // 头像框页面url
		} `json:"list"`
		JumpUrl string `json:"jump_url"` // 头像框商城页面跳转url
	} `json:"recommend_pendants"`
	RecommendCards struct { // 推荐装扮信息
		List []struct { // 推荐个性装扮列表
			Id      int    `json:"id"`       // 个性装扮id
			Name    string `json:"name"`     // 个性装扮名称
			Image   string `json:"image"`    // 个性装扮图标url
			JumpUrl string `json:"jump_url"` // 个性装扮页面url
		} `json:"list"`
		JumpUrl string `json:"jump_url"` // 推荐个性装扮商城页面跳转url
	} `json:"recommend_cards"`
	Sort []struct {
		Key  string `json:"key"`  // 扩展row字段名
		Sort int    `json:"sort"` // 排列顺序
	} `json:"sort"`
	InReview bool     `json:"in_review"`
	BigPoint struct { // 大积分信息
		PointInfo struct { // 点数信息
			Point       int `json:"point"`        // 当前拥有大积分数量
			ExpirePoint int `json:"expire_point"` // 失效积分，目前为0
			ExpireTime  int `json:"expire_time"`  // 失效时间，目前为0
			ExpireDays  int `json:"expire_days"`  // 失效天数，目前为0
		} `json:"point_info"`
		SignInfo struct { // 签到信息
			SignRemind   bool `json:"sign_remind"`
			Benefit      int  `json:"benefit"` // 签到收益，单位为积分
			BonusBenefit int  `json:"bonus_benefit"`
			NormalRemind bool `json:"normal_remind"`
			MuggleTask   bool `json:"muggle_task"`
		} `json:"sign_info"`
		SkuInfo struct { // 大积分商品预览
			Skus []GoodsSku `json:"skus"`
		} `json:"sku_info"`
		PointSwitchOff bool `json:"point_switch_off"`
		Tips           []struct {
			Content string `json:"content"`
		} `json:"tips"`
	} `json:"big_point"`
	HotList struct { // 热门榜单类型信息
		Taps []struct { // 热门榜单类型信息
			Oid       string `json:"oid"`        // 分类数据，类似json
			RankId    int    `json:"rank_id"`    // 分类id
			RankTitle string `json:"rank_title"` // 分类名称
		} `json:"taps"`
	} `json:"hot_list"`
}

// GetVipCenterInfo 获取大会员中心信息
func GetVipCenterInfo() (*VipCenterInfo, error) {
	return std.GetVipCenterInfo()
}
func (c *Client) GetVipCenterInfo() (*VipCenterInfo, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get("https://api.bilibili.com/x/vip/web/vip_center/combine")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取大会员中心信息")
	if err != nil {
		return nil, err
	}
	var ret *VipCenterInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type GoodsSku struct {
	Base struct {
		Token            string   `json:"token"`             // 商品token
		Title            string   `json:"title"`             // 商品名称
		Picture          string   `json:"picture"`           // 商品图片url
		RotationPictures []string `json:"rotation_pictures"` // 商品图片组
		Price            struct { // 价格信息
			Origin    int       `json:"origin"` // 商品原价，单位为积分
			Promotion *struct { // 折扣信息
				Price    int    `json:"price"`    // 折后价格，单位为积分
				Type     int    `json:"type"`     // 折扣类型，1：普通折扣，2：秒杀
				Discount int    `json:"discount"` // 折扣力度
				Label    string `json:"label"`    // 标签文案
			} `json:"promotion"`
		} `json:"price"`
		Inventory struct { // 库存信息
			AvailableNum int `json:"available_num"` // 库存总量
			UsedNum      int `json:"used_num"`      // 已售数量
			SurplusNum   int `json:"surplus_num"`   // 剩余数量
		} `json:"inventory"`
		UserType          int `json:"user_type"`
		ExchangeLimitType int `json:"exchange_limit_type"`
		ExchangeLimitNum  int `json:"exchange_limit_num"` // 限购数量
		StartTime         int `json:"start_time"`         // 起售时间戳，单位：秒
		EndTime           int `json:"end_time"`           // 止售时间，单位：秒
		State             int `json:"state"`              // 状态，固定值2
	} `json:"base"`
}
