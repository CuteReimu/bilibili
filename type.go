package bilibili

import (
	"strconv"
)

type Size struct {
	Width  int
	Height int
}

type FormatCtrl struct {
	Location int    `json:"location"` // 从全文第几个字开始
	Type     int    `json:"type"`     // 1：At
	Length   int    `json:"length"`   // 长度总共多少个字
	Data     string `json:"data"`     // 当Type为1时，这里填At的人的Uid
}

type ResourceType int

const (
	ResourceTypeVideo     ResourceType = 2  // 视频稿件
	ResourceTypeAudio     ResourceType = 12 // 音频
	ResourceTypeVideoList ResourceType = 21 // 视频合集
)

type Resource struct {
	Id   int
	Type ResourceType
}

func (r Resource) String() string {
	return strconv.Itoa(r.Id) + ":" + strconv.Itoa(int(r.Type))
}

type Pendant struct {
	Pid               int    `json:"pid"`                 // 挂件id
	Name              string `json:"name"`                // 挂件名称
	Image             string `json:"image"`               // 挂件图片url
	JumpUrl           string `json:"jump_url"`            // 挂件跳转url
	Type              string `json:"type"`                // 装扮类型。suit：一般装扮。vip_suit：vip 装扮
	Expire            int    `json:"expire"`              // 固定值0，作用尚不明确
	ImageEnhance      string `json:"image_enhance"`       // 头像框图片url
	ImageEnhanceFrame string `json:"image_enhance_frame"` // (?)
}

type Nameplate struct {
	Nid        int    `json:"nid"`         // 勋章id
	Name       string `json:"name"`        // 勋章名称
	Image      string `json:"image"`       // 勋章图标
	ImageSmall string `json:"image_small"` // 勋章图标（小）
	Level      string `json:"level"`       // 勋章等级
	Condition  string `json:"condition"`   // 获取条件
}

type OfficialVerify struct {
	Type int    `json:"type"` // 是否认证，-1：无。0：个人认证。1：机构认证
	Desc string `json:"desc"` // 认证信息，无为空
}

type VipLabel struct {
	Path        string `json:"path"`         // (?)
	Text        string `json:"text"`         // 会员类型文案
	LabelTheme  string `json:"label_theme"`  // 会员类型。vip：大会员。annual_vip：年度大会员。ten_annual_vip：十年大会员。hundred_annual_vip：百年大会员
	TextColor   string `json:"text_color"`   // 文字颜色?
	BgStyle     int    `json:"bg_style"`     // (?)
	BgColor     string `json:"bg_color"`     // 背景颜色?
	BorderColor string `json:"border_color"` // 描边颜色?
}

type Vip struct {
	Viptype            int      `json:"vipType"`              // 大会员类型。0：无。1：月会员。2：年以上会员
	Vipduedate         int      `json:"vipDueDate"`           // 大会员到期时间。毫秒 时间戳
	Dueremark          string   `json:"dueRemark"`            // (?)
	Accessstatus       int      `json:"accessStatus"`         // (?)
	Vipstatus          int      `json:"vipStatus"`            // 大会员状态。0：无。1：有
	Vipstatuswarn      string   `json:"vipStatusWarn"`        // (?)
	ThemeType          int      `json:"theme_type"`           // 会员样式 id
	Label              VipLabel `json:"label"`                // 会员铭牌样式
	AvatarSubscript    int      `json:"avatar_subscript"`     // (?)
	AvatarSubscriptUrl string   `json:"avatar_subscript_url"` // (?)
	NicknameColor      string   `json:"nickname_color"`       // 昵称颜色
}

type FansDetail struct {
	Uid          int    `json:"uid"`           // 用户 mid
	MedalId      int    `json:"medal_id"`      // 粉丝标签 id
	MedalName    string `json:"medal_name"`    // 粉丝标签名
	Score        int    `json:"score"`         // (?)
	Level        int    `json:"level"`         // 当前标签等级
	Intimacy     int    `json:"intimacy"`      // (?)
	MasterStatus int    `json:"master_status"` // (?)
	IsReceive    int    `json:"is_receive"`    // (?)
}

type Fan struct {
	IsFan   int    `json:"is_fan"`   // 是否为粉丝专属装扮。0：否。1：是
	Number  int    `json:"number"`   // 粉丝专属编号
	Color   string `json:"color"`    // 数字颜色。颜色码
	Name    string `json:"name"`     // 装扮名称
	NumDesc string `json:"num_desc"` // 粉丝专属编号。字串格式
}

type Cardbg struct {
	Id      int    `json:"id"`       // 评论条目装扮 id
	Name    string `json:"name"`     // 评论条目装扮名称
	Image   string `json:"image"`    // 评论条目装扮图片 url
	JumpUrl string `json:"jump_url"` // 评论条目装扮商城页面 url
	Fan     Fan    `json:"fan"`      // 粉丝专属信息
	Type    string `json:"type"`     // 装扮类型。suit：一般装扮。vip_suit：vip 装扮
}

type UserSailing struct {
	Pendant         *Pendant `json:"pendant"`           // 头像框信息
	Cardbg          *Cardbg  `json:"cardbg"`            // 评论卡片装扮
	CardbgWithFocus any      `json:"cardbg_with_focus"` // (?)
}

type Member struct {
	Mid            string         `json:"mid"`             // 发送者 mid
	Uname          string         `json:"uname"`           // 发送者昵称
	Sex            string         `json:"sex"`             // 发送者性别。男 女 保密
	Sign           string         `json:"sign"`            // 发送者签名
	Avatar         string         `json:"avatar"`          // 发送者头像 url
	Rank           string         `json:"rank"`            // (?)
	Displayrank    string         `json:"DisplayRank"`     // (?)
	LevelInfo      LevelInfo      `json:"level_info"`      // 发送者等级
	Pendant        Pendant        `json:"pendant"`         // 发送者头像框信息
	Nameplate      Nameplate      `json:"nameplate"`       // 发送者勋章信息
	OfficialVerify OfficialVerify `json:"official_verify"` // 发送者认证信息
	Vip            Vip            `json:"vip"`             // 发送者会员信息
	FansDetail     *FansDetail    `json:"fans_detail"`     // 发送者粉丝标签
	Following      int            `json:"following"`       // 是否关注该用户。需要登录(Cookie或APP) 。否则恒为0。0：未关注。1：已关注
	IsFollowed     int            `json:"is_followed"`     // 是否被该用户关注。需要登录(Cookie或APP) 。否则恒为0。0：未关注。1：已关注
	UserSailing    UserSailing    `json:"user_sailing"`    // 发送者评论条目装扮信息
	IsContractor   bool           `json:"is_contractor"`   // 是否为合作用户？
	ContractDesc   string         `json:"contract_desc"`   // 合作用户说明？
}

type Upper struct {
	Mid       int    `json:"mid"`        // UP 主 mid
	Name      string `json:"name"`       // 创建者昵称
	Face      string `json:"face"`       // 创建者头像url
	Followed  bool   `json:"followed"`   // 是否已关注创建者
	VipType   int    `json:"vip_type"`   // 会员类别，0：无，1：月大会员，2：年度及以上大会员
	VipStatue int    `json:"vip_statue"` // 0：无，1：有
}

type LevelInfo struct {
	CurrentLevel int `json:"current_level"` // 用户等级
	CurrentMin   int `json:"current_min"`   // 0
	CurrentExp   int `json:"current_exp"`   // 0
	NextExp      int `json:"next_exp"`      // 0
}

type CardSpace struct {
	SImg string `json:"s_img"` // 主页头图url 小图
	LImg string `json:"l_img"` // 主页头图url 正常
}

type Label struct {
	Path                  string `json:"path"`                      // 空。作用尚不明确
	Text                  string `json:"text"`                      // 会员类型文案。大会员 年度大会员 十年大会员 百年大会员 最强绿鲤鱼
	LabelTheme            string `json:"label_theme"`               // 会员标签。vip：大会员。annual_vip：年度大会员。ten_annual_vip：十 年大会员。hundred_annual_vip：百年大会员。fools_day_hundred_annual_vip：最强绿鲤鱼
	TextColor             string `json:"text_color"`                // 会员标签
	BgStyle               int    `json:"bg_style"`                  // 1
	BgColor               string `json:"bg_color"`                  // 会员标签背景颜色。颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
	BorderColor           string `json:"border_color"`              // 会员标签边框颜色。未使用
	UseImgLabel           bool   `json:"use_img_label"`             // true
	ImgLabelUriHans       string `json:"img_label_uri_hans"`        // 空串
	ImgLabelUriHant       string `json:"img_label_uri_hant"`        // 空串
	ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"` // 大会员牌子图片。简体版
	ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"` // 大会员牌子图片。繁体版
}

type Official struct {
	Role  int    `json:"role"`  // 成员认证级别
	Title string `json:"title"` // 成员认证名。无为空
	Desc  string `json:"desc"`  // 成员认证备注。无为空
	Type  int    `json:"type"`  // 成员认证类型。-1：无。0：有
}

type ContractInfo struct {
	IsContract   bool `json:"is_contract"`   // 目标用户是否为对方的契约者
	IsContractor bool `json:"is_contractor"` // 对方是否为目标用户的契约者
	Ts           int  `json:"ts"`            // 对方成为目标用户的契约者的时间。秒级时间戳，仅当 is_contractor 项的值为 true 时才有此项
	UserAttr     int  `json:"user_attr"`     // 对方作为目标用户的契约者的属性。1：老粉。否则为原始粉丝。仅当有特殊属性时才有此项
}

// RelationUser 关系列表对象
type RelationUser struct {
	Mid            int            `json:"mid"`             // 用户 mid
	Attribute      int            `json:"attribute"`       // 对方对于自己的关系属性。0：未关注。1：悄悄关注（现已下线）。2：已关注。6：已互粉。128：已拉黑
	Mtime          int            `json:"mtime"`           // 对方关注目标用户时间。秒级时间戳。互关后刷新
	Tag            []int          `json:"tag"`             // 目标用户将对方分组到的 id
	Special        int            `json:"special"`         // 目标用户特别关注对方标识。0：否。1：是
	ContractInfo   ContractInfo   `json:"contract_info"`   // 契约计划相关信息
	Uname          string         `json:"uname"`           // 用户昵称
	Face           string         `json:"face"`            // 用户头像url
	FaceNft        int            `json:"face_nft"`        // 是否为 NFT 头像。0：非 NFT 头像。1：NFT 头像
	Sign           string         `json:"sign"`            // 用户签名
	OfficialVerify OfficialVerify `json:"official_verify"` // 认证信息
	Vip            Vip            `json:"vip"`             // 会员信息
	NftIcon        string         `json:"nft_icon"`        // （？）
	RecReason      string         `json:"rec_reason"`      // （？）
	TrackId        string         `json:"track_id"`        // （？）
}
