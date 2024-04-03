package bilibili

import (
	"github.com/go-resty/resty/v2"
)

type ReceiveVipPrivilegeParam struct {
	Type int `json:"type"` // 兑换类型。1：B币券。2：会员购优惠券。3：漫画福利券。4：会员购包邮券。5：漫画商城优惠券。6：装扮体验卡。7：课堂优惠券
}

// ReceiveVipPrivilege 兑换大会员卡券，1：B币券，2：会员购优惠券，3：漫画福利券，4：会员购包邮券，5：漫画商城优惠券
func (c *Client) ReceiveVipPrivilege(param ReceiveVipPrivilegeParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/vip/privilege/receive"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

// SignVipScore 大积分签到
func (c *Client) SignVipScore() error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/pgc/activity/score/task/sign"
	)
	_, err := execute[any](c, method, url, nil, fillCsrf(c))
	return err
}

type VipPrivilegeInfo struct {
	Type            int `json:"type"`              // 卡券类型。详见 list 数组表格中的 type 项
	State           int `json:"state"`             // 兑换状态。0：未兑换。1：已兑换。2：未完成（若需要完成）
	ExpireTime      int `json:"expire_time"`       // 本轮卡券过期时间戳。当月月底/当日24点
	VipType         int `json:"vip_type"`          // 当前用户的大会员状态。2：年度大会员
	NextReceiveDays int `json:"next_receive_days"` // 距下一轮兑换剩余天数。无权限时，每月任务固定为 0，每日固定为 1
	PeriodEndUnix   int `json:"period_end_unix"`   // 下一轮兑换开始时间戳。秒级时间戳
}

type VipPrivilege struct {
	List           []VipPrivilegeInfo `json:"list"`             // 卡券信息列表
	IsShortVip     bool               `json:"is_short_vip"`     // (?)
	IsFreightOpen  bool               `json:"is_freight_open"`  // (?)
	Level          int                `json:"level"`            // 当前等级
	CurExp         int                `json:"cur_exp"`          // 当前拥有经验值
	NextExp        int                `json:"next_exp"`         // 升级所需经验值。满级时为 -1
	IsVip          bool               `json:"is_vip"`           // 是否为大会员
	IsSeniorMember int                `json:"is_senior_member"` // (?)
	Format060102   int                `json:"format060102"`     // (?)
}

// GetVipPrivilege 卡券状态查询
func (c *Client) GetVipPrivilege() (*VipPrivilege, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/vip/privilege/my"
	)
	return execute[*VipPrivilege](c, method, url, nil)
}

type GetVipCenterInfoParam struct {
	AccessKey string `json:"access_key,omitempty" request:"query,omitempty"` // APP登录Token
	Platform  string `json:"platform,omitempty" request:"query,omitempty"`   // 平台表示。web端：web。安卓APP：android
	MobiApp   string `json:"mobi_app,omitempty" request:"query,omitempty"`   // APP 名称。安卓APP：android
	Build     int    `json:"build,omitempty" request:"query,omitempty"`      // 构建 id
}

type VipUserAccount struct {
	Mid            int    `json:"mid"`              // 用户 mid
	Name           string `json:"name"`             // 昵称
	Sex            string `json:"sex"`              // 性别。男 / 女 / 保密
	Face           string `json:"face"`             // 头像 url
	Sign           string `json:"sign"`             // 签名
	Rank           int    `json:"rank"`             // 等级
	Birthday       int    `json:"birthday"`         // 生日。秒时间戳
	IsFakeAccount  int    `json:"is_fake_account"`  // (?)
	IsDeleted      int    `json:"is_deleted"`       // 是否注销。0：正常。1：注销
	InRegAudit     int    `json:"in_reg_audit"`     // 是否注册审核。0：正常。1：审核
	IsSeniorMember int    `json:"is_senior_member"` // 是否转正。0：未转正。1：正式会员
}

type VipUserVip struct {
	Mid             int    `json:"mid"`              // 用户 mid
	VipType         int    `json:"vip_type"`         // 会员类型。0：无。1：月大会员。2：年度及以上大会员
	VipStatus       int    `json:"vip_status"`       // 会员状态。0：无。1：有
	VipDueDate      int    `json:"vip_due_date"`     // 会员过期时间。毫秒时间戳
	VipPayType      int    `json:"vip_pay_type"`     // 支付类型。0：未支付（常见于官方账号）。1：已支付（以正常渠道获取的大会员均为此值）
	ThemeType       int    `json:"theme_type"`       // (?)
	Label           Label  `json:"label"`            // 会员标签
	AvatarSubscript int    `json:"avatar_subscript"` // 是否显示会员图标。0：不显示。1：显示
	NicknameColor   string `json:"nickname_color"`   // 会员昵称颜色。颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
	IsNewUser       bool   `json:"is_new_user"`      // (?)
	TipMaterial     any    `json:"tip_material"`     // (?)
}

type VipUserTv struct {
	Type       int `json:"type"`         // 电视大会员类型。0：无。1：月大会员。2：年度及以上大会员
	VipPayType int `json:"vip_pay_type"` // 电视大支付类型。0：未支付（常见于官方账号）。1：已支付（以正常渠道获取的大会员均为此值）
	Status     int `json:"status"`       // 电视大会员状态。0：无。1：有
	DueDate    int `json:"due_date"`     // 电视大会员过期时间。毫秒时间戳
}

type AvatarPendant struct {
	Image             string `json:"image"`               // 头像框 url
	ImageEnhance      string `json:"image_enhance"`       // 头像框 url。动态图
	ImageEnhanceFrame string `json:"image_enhance_frame"` // 动态头像框帧波普版 url
}

type VipUser struct {
	Account              VipUserAccount `json:"account"`                // 账号基本信息
	Vip                  VipUserVip     `json:"vip"`                    // 账号会员信息
	Tv                   VipUserTv      `json:"tv"`                     // 电视会员信息
	BackgroundImageSmall string         `json:"background_image_small"` // 空
	BackgroundImageBig   string         `json:"background_image_big"`   // 空
	PanelTitle           string         `json:"panel_title"`            // 用户昵称
	AvatarPendant        AvatarPendant  `json:"avatar_pendant"`         // 用户头像框信息
	VipOverdueExplain    string         `json:"vip_overdue_explain"`    // 大会员提示文案。有效期 / 到期
	TvOverdueExplain     string         `json:"tv_overdue_explain"`     // 电视大会员提示文案。有效期 / 到期
	AccountExceptionText string         `json:"account_exception_text"` // 空
	IsAutoRenew          bool           `json:"is_auto_renew"`          // 是否自动续费。true：是。false：否
	IsTvAutoRenew        bool           `json:"is_tv_auto_renew"`       // 是否自动续费电视大会员。true：是。false：否
	SurplusSeconds       int            `json:"surplus_seconds"`        // 大会员到期剩余时间。单位为秒
	VipKeepTime          int            `json:"vip_keep_time"`          // 持续开通大会员时间。单位为秒
	Renew                any            `json:"renew"`                  // (?)
	Notice               any            `json:"notice"`                 // (?)
}

type VipWallet struct {
	Coupon            int  `json:"coupon"`             // 当前 B 币券
	Point             int  `json:"point"`              // (?)
	PrivilegeReceived bool `json:"privilege_received"` // (?)
}

type ChildPrivilege struct {
	FirstId            int    `json:"first_id"`             // 特权父类 id
	ReportId           string `json:"report_id"`            // 上报 id。该特权的代号？
	Name               string `json:"name"`                 // 特权名称
	Desc               string `json:"desc"`                 // 特权简介文案
	Explain            string `json:"explain"`              // 特权介绍正文
	IconUrl            string `json:"icon_url"`             // 特权图标 url
	IconGrayUrl        string `json:"icon_gray_url"`        // 特权图标灰色主题 url。某些项目无此字段
	BackgroundImageUrl string `json:"background_image_url"` // 背景图片 url
	Link               string `json:"link"`                 // 特权介绍页 url
	ImageUrl           string `json:"image_url"`            // 特权示例图 url
	Type               int    `json:"type"`                 // 类型？。目前为0
	HotType            int    `json:"hot_type"`             // 是否热门特权。0：普通特权。1：热门特权
	NewType            int    `json:"new_type"`             // 是否新特权。0：普通特权。1：新特权
	Id                 int    `json:"id"`                   // 特权子类 id
}

type Privilege struct {
	Id              int              `json:"id"`               // 特权父类 id
	Name            string           `json:"name"`             // 类型名称
	ChildPrivileges []ChildPrivilege `json:"child_privileges"` // 特权子类列表
}

type Banner struct {
	Id          int    `json:"id"`           // banner 卡片 id
	Index       int    `json:"index"`        // banner 卡片排序
	Image       string `json:"image"`        // banner 卡片图片 url
	Title       string `json:"title"`        // banner 卡片标题
	Uri         string `json:"uri"`          // banner 卡片跳转页 url
	TrackParams any    `json:"track_params"` // 上报参数
}

type WelfareItem struct {
	Id          int    `json:"id"`           // 福利 id
	Name        string `json:"name"`         // 福利名称
	HomepageUri string `json:"homepage_uri"` // 福利图片 url
	BackdropUri string `json:"backdrop_uri"` // 福利图片 banner url
	Tid         int    `json:"tid"`          // (?)。目前为0
	Rank        int    `json:"rank"`         // 排列顺序
	ReceiveUri  string `json:"receive_uri"`  // 福利跳转页 url
}

type Welfare struct {
	Count int           `json:"count"` // 福利数
	List  []WelfareItem `json:"list"`  // 福利项目列表
}

type RecommendPendant struct {
	Id      int    `json:"id"`       // 头像框 id
	Name    string `json:"name"`     // 头像框名称
	Image   string `json:"image"`    // 头像框图片 url
	JumpUrl string `json:"jump_url"` // 头像框页面 url
}

type RecommendPendants struct {
	JumpUrl string             `json:"jump_url"` // 头像框商城页面跳转 url
	List    []RecommendPendant `json:"list"`     // 推荐头像框列表
}

type RecommendCard struct {
	Id      int    `json:"id"`       // 个性装扮 id
	Name    string `json:"name"`     // 个性装扮名称
	Image   string `json:"image"`    // 个性装扮图标 url
	JumpUrl string `json:"jump_url"` // 个性装扮页面 url
}

type RecommendCards struct {
	JumpUrl string          `json:"jump_url"` // 推荐个性装扮商城页面跳转 url
	List    []RecommendCard `json:"list"`     // 推荐个性装扮列表
}

type Sort struct {
	Key  string `json:"key"`  // 扩展 row 字段名
	Sort int    `json:"sort"` // 排列顺序
}

type PointInfo struct {
	Point       int `json:"point"`        // 当前拥有大积分数量
	ExpirePoint int `json:"expire_point"` // 失效积分？。目前为0
	ExpireTime  int `json:"expire_time"`  // 失效时间？。目前为0
	ExpireDays  int `json:"expire_days"`  // 失效天数？。目前为0
}

type SignInfo struct {
	SignRemind   bool `json:"sign_remind"`   // (?)
	Benefit      int  `json:"benefit"`       // 签到收益。单位为积分
	BonusBenefit int  `json:"bonus_benefit"` // (?)
	NormalRemind bool `json:"normal_remind"` // (?)
	MuggleTask   bool `json:"muggle_task"`   // (?)
}

type BigPoint struct {
	PointInfo      PointInfo `json:"point_info"` // 点数信息
	SignInfo       SignInfo  `json:"sign_info"`  // 签到信息
	SkuInfo        any       `json:"sku_info"`   // 大积分商品预览
	Tips           bool      `json:"tips"`
	PointSwitchOff any       `json:"point_switch_off"`
}

type VipCenterInfo struct {
	User              VipUser           `json:"user"`               // 用户信息
	Wallet            VipWallet         `json:"wallet"`             // 钱包信息
	UnionVip          any               `json:"union_vip"`          // 联合会员信息列表，web 端：null。APP 端：array
	OtherOpenInfo     any               `json:"other_open_info"`    // 其他开通方式信息列表，web 端：null。APP 端：array
	Privileges        []Privilege       `json:"privileges"`         // 会员特权信息列表
	Banners           []Banner          `json:"banners"`            // banner 卡片列表。web 端为空
	Welfare           Welfare           `json:"welfare"`            // 福利信息
	RecommendPendants RecommendPendants `json:"recommend_pendants"` // 推荐头像框信息
	RecommendCards    RecommendCards    `json:"recommend_cards"`    // 推荐装扮信息
	Sort              []Sort            `json:"sort"`
	InReview          bool              `json:"in_review"`
	BigPoint          BigPoint          `json:"big_point"` // 大积分信息
	HotList           any               `json:"hot_list"`  // 热门榜单类型信息
}

// GetVipCenterInfo 获取大会员中心信息
func (c *Client) GetVipCenterInfo(param GetVipCenterInfoParam) (*VipCenterInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/vip/privilege/my"
	)
	return execute[*VipCenterInfo](c, method, url, param)
}
