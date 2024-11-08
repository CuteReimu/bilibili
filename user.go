package bilibili

import "github.com/go-resty/resty/v2"

type GetUserVideosParam struct {
	Mid     int    `json:"mid"`                                         // 目标用户mid
	Order   string `json:"order,omitempty" request:"query,omitempty"`   // 排序方式。默认为pubdate。最新发布：pubdate。最多播放：click。最多收藏：stow
	Tid     int    `json:"tid,omitempty" request:"query,omitempty"`     // 筛选目标分区。默认为0。0：不进行分区筛选。分区tid为所筛选的分区
	Keyword string `json:"keyword,omitempty" request:"query,omitempty"` // 关键词筛选。用于使用关键词搜索该UP主视频稿件
	Pn      int    `json:"pn,omitempty" request:"query,omitempty"`      // 页码。默认为 1
	Ps      int    `json:"ps,omitempty" request:"query,omitempty"`      // 每页项数。默认为 30
}

type VideoArea struct {
	Count int    `json:"count"` // 投稿至该分区的视频数
	Name  string `json:"name"`  // 该分区名称
	Tid   int    `json:"tid"`   // 该分区tid
}

type UserVideo struct {
	Aid          int    `json:"aid"` // 稿件avid
	Attribute    int    `json:"attribute"`
	Author       string `json:"author"`      // 视频UP主。不一定为目标用户（合作视频）
	Bvid         string `json:"bvid"`        // 稿件bvid
	Comment      int    `json:"comment"`     // 视频评论数
	Copyright    string `json:"copyright"`   // 视频版权类型
	Created      int    `json:"created"`     // 投稿时间。时间戳
	Description  string `json:"description"` // 视频简介
	EnableVt     int    `json:"enable_vt"`
	HideClick    bool   `json:"hide_click"`     // false。作用尚不明确
	IsPay        int    `json:"is_pay"`         // 0。作用尚不明确
	IsUnionVideo int    `json:"is_union_video"` // 是否为合作视频。0：否。1：是
	Length       string `json:"length"`         // 视频长度。MM:SS
	Mid          int    `json:"mid"`            // 视频UP主mid。不一定为目标用户（合作视频）
	Meta         any    `json:"meta"`           // 无数据时为 null
	Pic          string `json:"pic"`            // 视频封面
	Play         int    `json:"play"`           // 视频播放次数
	Review       int    `json:"review"`         // 0。作用尚不明确
	Subtitle     string `json:"subtitle"`       // 空。作用尚不明确
	Title        string `json:"title"`          // 视频标题
	Typeid       int    `json:"typeid"`         // 视频分区tid
	VideoReview  int    `json:"video_review"`   // 视频弹幕数
}

type UserVideosList struct {
	Tlist map[int]VideoArea `json:"tlist"` // 投稿视频分区索引
	Vlist []UserVideo       `json:"vlist"` // 投稿视频列表
}

type UserVideoPage struct {
	Count int `json:"count"` // 总计稿件数
	Pn    int `json:"pn"`    // 当前页码
	Ps    int `json:"ps"`    // 每页项数
}

type EpisodicButton struct {
	Text string `json:"text"` // 按钮文字
	Uri  string `json:"uri"`  // 全部播放页url
}

type UserVideos struct {
	List           UserVideosList `json:"list"`            // 列表信息
	Page           UserVideoPage  `json:"page"`            // 页面信息
	EpisodicButton EpisodicButton `json:"episodic_button"` // “播放全部“按钮
	IsRisk         bool           `json:"is_risk"`
	GaiaResType    int            `json:"gaia_res_type"`
	GaiaData       any            `json:"gaia_data"`
}

// GetUserVideos 查询用户投稿视频明细
func (c *Client) GetUserVideos(param GetUserVideosParam) (*UserVideos, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/space/wbi/arc/search"
	)
	return execute[*UserVideos](c, method, url, param, fillWbiHandler(c.wbi, c.GetCookies()))
}

type GetUserSpaceDetailParam struct {
	Mid int `json:"mid"` // 目标用户mid
}

type SpaceVip struct {
	Type               int    `json:"type"`                 // 会员类型。0：无。1：月大会员。2：年度及以上大会员
	Status             int    `json:"status"`               // 会员状态。0：无。1：有
	DueDate            int    `json:"due_date"`             // 会员过期时间。毫秒时间戳
	VipPayType         int    `json:"vip_pay_type"`         // 支付类型。0：未支付（常见于官方账号）。1：已支付（以正常渠道获取的大会员均为此值）
	ThemeType          int    `json:"theme_type"`           // 0。作用尚不明确
	Label              Label  `json:"label"`                // 会员标签
	AvatarSubscript    int    `json:"avatar_subscript"`     // 是否显示会员图标。0：不显示。1：显示
	NicknameColor      string `json:"nickname_color"`       // 会员昵称颜色。颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
	Role               int    `json:"role"`                 // 大角色类型。1：月度大会员。3：年度大会员。7：十年大会员。15：百年大会员
	AvatarSubscriptUrl string `json:"avatar_subscript_url"` // 大会员角标地址
	TvVipStatus        int    `json:"tv_vip_status"`        // 电视大会员状态。0：未开通
	TvVipPayType       int    `json:"tv_vip_pay_type"`      // 电视大会员支付类型
}

type Medal struct {
	Uid              int    `json:"uid"`                // 此用户mid
	TargetId         int    `json:"target_id"`          // 粉丝勋章所属UP的mid
	MedalId          int    `json:"medal_id"`           // 粉丝勋章id
	Level            int    `json:"level"`              // 粉丝勋章等级
	MedalName        string `json:"medal_name"`         // 粉丝勋章名称
	MedalColor       int    `json:"medal_color"`        // 颜色
	Intimacy         int    `json:"intimacy"`           // 当前亲密度
	NextIntimacy     int    `json:"next_intimacy"`      // 下一等级所需亲密度
	DayLimit         int    `json:"day_limit"`          // 每日亲密度获取上限
	TodayFeed        int    `json:"today_feed"`         // 今日已获得亲密度
	MedalColorStart  int    `json:"medal_color_start"`  // 粉丝勋章颜色。十进制数，可转为十六进制颜色代码
	MedalColorEnd    int    `json:"medal_color_end"`    // 粉丝勋章颜色。十进制数，可转为十六进制颜色代码
	MedalColorBorder int    `json:"medal_color_border"` // 粉丝勋章边框颜色。十进制数，可转为十六进制颜色代码
	IsLighted        int    `json:"is_lighted"`
	LightStatus      int    `json:"light_status"`
	WearingStatus    int    `json:"wearing_status"` // 当前是否佩戴。0：未佩戴。1：已佩戴
	Score            int    `json:"score"`
}

type FansMedal struct {
	Show  bool  `json:"show"`
	Wear  bool  `json:"wear"`  // 是否佩戴了粉丝勋章
	Medal Medal `json:"medal"` // 粉丝勋章信息
}

type SysNotice struct {
	Id         int    `json:"id"`          // id
	Content    string `json:"content"`     // 显示文案
	Url        string `json:"url"`         // 跳转地址
	NoticeType int    `json:"notice_type"` // 提示类型。1,2
	Icon       string `json:"icon"`        // 前缀图标
	TextColor  string `json:"text_color"`  // 文字颜色
	BgColor    string `json:"bg_color"`    // 背景颜色
}

type WatchedShow struct {
	Switch       bool   `json:"switch"` // ?
	Num          int    `json:"num"`    // total watched users
	TextSmall    string `json:"text_small"`
	TextLarge    string `json:"text_large"`
	Icon         string `json:"icon"`          // watched icon url
	IconLocation string `json:"icon_location"` // ?
	IconWeb      string `json:"icon_web"`      // watched icon url
}

type LiveRoom struct {
	Roomstatus    int         `json:"roomStatus"` // 直播间状态。0：无房间。1：有房间
	Livestatus    int         `json:"liveStatus"` // 直播状态。0：未开播。1：直播中
	Url           string      `json:"url"`        // 直播间网页 url
	Title         string      `json:"title"`      // 直播间标题
	Cover         string      `json:"cover"`      // 直播间封面 url
	WatchedShow   WatchedShow `json:"watched_show"`
	Roomid        int         `json:"roomid"`         // 直播间 id(短号)
	Roundstatus   int         `json:"roundStatus"`    // 轮播状态。0：未轮播。1：轮播
	BroadcastType int         `json:"broadcast_type"` // 0
}

type School struct {
	Name string `json:"name"` // 就读大学名称。没有则为空
}

type Profession struct {
	Name       string `json:"name"`       // 资质名称
	Department string `json:"department"` // 职位
	Title      string `json:"title"`      // 所属机构
	IsShow     int    `json:"is_show"`    // 是否显示。0：不显示。1：显示
}

type Elec struct {
	ShowInfo struct {
		Show    bool   `json:"show"`     // 是否开通了充电
		State   int    `json:"state"`    // 状态。-1：未开通。1：已开通
		Title   string `json:"title"`    // 空串
		Icon    string `json:"icon"`     // 空串
		JumpUrl string `json:"jump_url"` // 空串
	} `json:"show_info"`
}

type Contract struct {
	IsDisplay       bool `json:"is_display"`        // true/false。在页面中未使用此字段
	IsFollowDisplay bool `json:"is_follow_display"` // 是否在显示老粉计划。true：显示。false：不显示
}

type UserSpaceDetail struct {
	Mid            int       `json:"mid"`           // mid
	Name           string    `json:"name"`          // 昵称
	Sex            string    `json:"sex"`           // 性别。男/女/保密
	Face           string    `json:"face"`          // 头像链接
	FaceNft        int       `json:"face_nft"`      // 是否为 NFT 头像。0：不是 NFT 头像。1：是 NFT 头像
	FaceNftType    int       `json:"face_nft_type"` // NFT 头像类型？
	Sign           string    `json:"sign"`          // 签名
	Rank           int       `json:"rank"`          // 用户权限等级。目前应该无任何作用。5000：0级未答题。10000：普通会员。20000：字幕君。25000：VIP。30000：真·职人。32000：管理员
	Level          int       `json:"level"`         // 当前等级。0-6 级
	Jointime       int       `json:"jointime"`      // 注册时间。此接口返回恒为0
	Moral          int       `json:"moral"`         // 节操值。此接口返回恒为0
	Silence        int       `json:"silence"`       // 封禁状态。0：正常。1：被封
	Coins          int       `json:"coins"`         // 硬币数。需要登录（Cookie） 。只能查看自己的。默认为0
	FansBadge      bool      `json:"fans_badge"`    // 是否具有粉丝勋章。false：无。true：有
	FansMedal      FansMedal `json:"fans_medal"`    // 粉丝勋章信息
	Official       Official  `json:"official"`      // 认证信息
	Vip            SpaceVip  `json:"vip"`           // 会员信息
	Pendant        Pendant   `json:"pendant"`       // 头像框信息
	Nameplate      Nameplate `json:"nameplate"`     // 勋章信息
	UserHonourInfo struct {
		Mid    int    `json:"mid"`    // 0
		Colour string `json:"colour"` // null
		Tags   any    `json:"tags"`   // null
	} `json:"user_honour_info"` // （？）
	IsFollowed bool       `json:"is_followed"` // 是否关注此用户。true：已关注。false：未关注。需要登录（Cookie） 。未登录恒为false
	TopPhoto   string     `json:"top_photo"`   // 主页头图链接
	Theme      any        `json:"theme"`       // （？）
	SysNotice  SysNotice  `json:"sys_notice"`  // 系统通知。无内容则为空对象。主要用于展示如用户争议、纪念账号等等的小黄条
	LiveRoom   LiveRoom   `json:"live_room"`   // 直播间信息
	Birthday   string     `json:"birthday"`    // 生日。MM-DD。如设置隐私为空
	School     School     `json:"school"`      // 学校
	Profession Profession `json:"profession"`  // 专业资质信息
	Tags       any        `json:"tags"`        // 个人标签
	Series     struct {
		UserUpgradeStatus int  `json:"user_upgrade_status"` // (?)
		ShowUpgradeWindow bool `json:"show_upgrade_window"` // (?)
	} `json:"series"`
	IsSeniorMember int      `json:"is_senior_member"` // 是否为硬核会员。0：否。1：是
	McnInfo        any      `json:"mcn_info"`         // （？）
	GaiaResType    int      `json:"gaia_res_type"`    // （？）
	GaiaData       any      `json:"gaia_data"`        // （？）
	IsRisk         bool     `json:"is_risk"`          // （？）
	Elec           Elec     `json:"elec"`             // 充电信息
	Contract       Contract `json:"contract"`         // 是否显示老粉计划
}

// GetUserSpaceDetail 获取用户空间详细信息
func (c *Client) GetUserSpaceDetail(param GetUserSpaceDetailParam) (*UserSpaceDetail, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/space/wbi/acc/info"
	)
	return execute[*UserSpaceDetail](c, method, url, param, fillWbiHandler(c.wbi, c.GetCookies()))
}

type GetUserCardParam struct {
	Mid   int  `json:"mid"`                                       // 目标用户mid
	Photo bool `json:"photo,omitempty" request:"query,omitempty"` // 是否请求用户主页头图。true：是。false：否
}

type UserCardVip struct {
	Viptype       int    `json:"vipType"`       // 大会员类型。0：无。1：月度大会员。2：年度及以上大会员
	Dueremark     string `json:"dueRemark"`     // 空。**作用尚不明确**
	Accessstatus  int    `json:"accessStatus"`  // 0。**作用尚不明确**
	Vipstatus     int    `json:"vipStatus"`     // 大会员状态。0：无。1：有
	Vipstatuswarn string `json:"vipStatusWarn"` // 空。**作用尚不明确**
	ThemeType     int    `json:"theme_type"`    // 0。**作用尚不明确**
}

type UserCardInfo struct {
	Mid            string         `json:"mid"`             // 用户mid
	Approve        bool           `json:"approve"`         // false。**作用尚不明确**
	Name           string         `json:"name"`            // 用户昵称
	Sex            string         `json:"sex"`             // 用户性别。男 女 保密
	Face           string         `json:"face"`            // 用户头像链接
	Displayrank    string         `json:"DisplayRank"`     // 0。**作用尚不明确**
	Regtime        int            `json:"regtime"`         // 0。**作用尚不明确**
	Spacesta       int            `json:"spacesta"`        // 用户状态。0：正常。-2：被封禁
	Birthday       string         `json:"birthday"`        // 空。**作用尚不明确**
	Place          string         `json:"place"`           // 空。**作用尚不明确**
	Description    string         `json:"description"`     // 空。**作用尚不明确**
	Article        int            `json:"article"`         // 0。**作用尚不明确**
	Attentions     any            `json:"attentions"`      // 空。**作用尚不明确**
	Fans           int            `json:"fans"`            // 粉丝数
	Friend         int            `json:"friend"`          // 关注数
	Attention      int            `json:"attention"`       // 关注数
	Sign           string         `json:"sign"`            // 签名
	LevelInfo      LevelInfo      `json:"level_info"`      // 等级
	Pendant        Pendant        `json:"pendant"`         // 挂件
	Nameplate      Nameplate      `json:"nameplate"`       // 勋章
	Official       Official       `json:"Official"`        // 认证信息
	OfficialVerify OfficialVerify `json:"official_verify"` // 认证信息2
	Vip            UserCardVip    `json:"vip"`             // 大会员状态
	Space          CardSpace      `json:"space"`           // 主页头图
}

type UserCard struct {
	Card         UserCardInfo `json:"card"`          // 卡片信息
	Following    bool         `json:"following"`     // 是否关注此用户。true：已关注。false：未关注。需要登录(Cookie)。未登录为false
	ArchiveCount int          `json:"archive_count"` // 用户稿件数
	ArticleCount int          `json:"article_count"` // 0。**作用尚不明确**
	Follower     int          `json:"follower"`      // 粉丝数
	LikeNum      int          `json:"like_num"`      // 点赞数
}

// GetUserCard 获取用户用户名片 免登录
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/user/info.md#%E7%94%A8%E6%88%B7%E5%90%8D%E7%89%87%E4%BF%A1%E6%81%AF
func (c *Client) GetUserCard(param GetUserCardParam) (*UserCard, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/card"
	)
	return execute[*UserCard](c, method, url, param)
}

type MyVip struct {
	Type            int    `json:"type"`             // 会员类型。0：无。1：月大会员。2：年度及以上大会员
	Status          int    `json:"status"`           // 会员状态。0：无。1：有
	DueDate         int    `json:"due_date"`         // 会员过期时间。Unix时间戳(毫秒)
	ThemeType       int    `json:"theme_type"`       // 0。作用尚不明确
	Label           Label  `json:"label"`            // 会员标签
	AvatarSubscript int    `json:"avatar_subscript"` // 是否显示会员图标。0：不显示。1：显示
	NicknameColor   string `json:"nickname_color"`   // 会员昵称颜色。颜色码
}

type MyProfession struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ShowName string `json:"show_name"`
}

type MyUserSpaceDetail struct {
	Mid            int          `json:"mid"`             // mid
	Name           string       `json:"name"`            // 昵称
	Sex            string       `json:"sex"`             // 性别。男 女 保密
	Face           string       `json:"face"`            // 头像图片url
	Sign           string       `json:"sign"`            // 签名
	Rank           int          `json:"rank"`            // 10000。**作用尚不明确**
	Level          int          `json:"level"`           // 当前等级。0-6级
	Jointime       int          `json:"jointime"`        // 0。**作用尚不明确**
	Moral          int          `json:"moral"`           // 节操。默认70
	Silence        int          `json:"silence"`         // 封禁状态。0：正常。1：被封
	EmailStatus    int          `json:"email_status"`    // 已验证邮箱。0：未验证。1：已验证
	TelStatus      int          `json:"tel_status"`      // 已验证手机号。0：未验证。1：已验证
	Identification int          `json:"identification"`  // 1。**作用尚不明确**
	Vip            MyVip        `json:"vip"`             // 大会员状态
	Pendant        Pendant      `json:"pendant"`         // 头像框信息
	Nameplate      Nameplate    `json:"nameplate"`       // 勋章信息
	Official       Official     `json:"official"`        // 认证信息
	Birthday       int          `json:"birthday"`        // 生日。时间戳
	IsTourist      int          `json:"is_tourist"`      // 0。**作用尚不明确**
	IsFakeAccount  int          `json:"is_fake_account"` // 0。**作用尚不明确**
	PinPrompting   int          `json:"pin_prompting"`   // 0。**作用尚不明确**
	IsDeleted      int          `json:"is_deleted"`      // 0。**作用尚不明确**
	InRegAudit     int          `json:"in_reg_audit"`
	IsRipUser      bool         `json:"is_rip_user"`
	Profession     MyProfession `json:"profession"` // 专业资质
	Coins          float64      `json:"coins"`      // 硬币数
	Following      int          `json:"following"`  // 粉丝数
	Follower       int          `json:"follower"`   // 粉丝数
}

// GetMyUserSpaceDetail 获取登录用户空间详细信息
func (c *Client) GetMyUserSpaceDetail() (*MyUserSpaceDetail, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/space/myinfo"
	)
	return execute[*MyUserSpaceDetail](c, method, url, nil)
}

type CheckNickNameParam struct {
	Nickname string `json:"nickName"` // 目标昵称。最长为16字符
}

// CheckNickName 检查昵称是否可注册
func (c *Client) CheckNickName(param CheckNickNameParam) error {
	const (
		method = resty.MethodGet
		url    = "https://passport.bilibili.com/web/generic/check/nickname"
	)
	_, err := execute[any](c, method, url, param)
	return err
}

type JoinOldFansParam struct {
	Aid      string `json:"aid,omitempty" request:"query,omitempty"`      // 空串
	UpMid    string `json:"up_mid"`                                       // UP主UID
	Source   string `json:"source,omitempty" request:"query,omitempty"`   // "4"
	Scene    string `json:"scene,omitempty" request:"query,omitempty"`    // "105"
	Platform string `json:"platform,omitempty" request:"query,omitempty"` // "web"
	MobiApp  string `json:"mobi_app,omitempty" request:"query,omitempty"` // "pc"
}

type JoinOldFansResult struct {
	AllowMessage bool   `json:"allow_message"` // true
	InputText    string `json:"input_text"`    // UP主加油！看好你噢
	InputTitle   string `json:"input_title"`   // 感谢你对UP主的特别支持，“老粉”可期！私信留言鼓励下TA吧
}

// JoinOldFans 加入老粉计划
func (c *Client) JoinOldFans(param JoinOldFansParam) (*JoinOldFansResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v1/contract/add_contract"
	)
	return execute[*JoinOldFansResult](c, method, url, param, fillCsrf(c))
}

type FansSendMessageParam struct {
	Aid     string `json:"aid,omitempty" request:"query,omitempty"`    // 空串
	UpMid   string `json:"up_mid"`                                     // UP主UID
	Source  string `json:"source,omitempty" request:"query,omitempty"` // "4"
	Scene   string `json:"scene,omitempty" request:"query,omitempty"`  // "105"
	Content string `json:"content"`                                    // 留言内容
}

type FansSendMessageResult struct {
	SuccessToast string `json:"success_toast"` // "提交成功，UP主已收到留言~"
}

// FansSendMessage 老粉计划发送留言
func (c *Client) FansSendMessage(param FansSendMessageParam) (*FansSendMessageResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v1/contract/add_message"
	)
	return execute[*FansSendMessageResult](c, method, url, param, fillCsrf(c))
}

type BatchGetUserCardsParam struct {
	Uids []int `json:"uids"` // 目标用户的UID列表
}

type BatchGetUserCardsResult struct {
	Mid     int    `json:"mid"`     // mid
	Name    string `json:"name"`    // 昵称
	Face    string `json:"face"`    // 头像链接
	Sign    string `json:"sign"`    // 签名
	Rank    int    `json:"rank"`    // 用户权限等级
	Level   int    `json:"level"`   // 当前等级。0-6 级
	Silence int    `json:"silence"` // 封禁状态。0：正常。1：被封
}

// BatchGetUserCards 获取多用户详细信息
func (c *Client) BatchGetUserCards(param BatchGetUserCardsParam) ([]*BatchGetUserCardsResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.vc.bilibili.com/account/v1/user/cards"
	)
	return execute[[]*BatchGetUserCardsResult](c, method, url, param)
}

type GetUserFollowersParam struct {
	Vmid int `json:"vmid"`                                   // 目标用户 mid
	Ps   int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为 50
	Pn   int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为 1。仅可查看前 1000 名粉丝
}

type GetUserFollowersResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
	Total     int            `json:"total"`      // 粉丝总数
}

// GetUserFollowers 查询用户粉丝明细（需要登录）
func (c *Client) GetUserFollowers(param GetUserFollowersParam) (*GetUserFollowersResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/followers"
	)
	return execute[*GetUserFollowersResult](c, method, url, param)
}

type GetUserFollowingsParam struct {
	Vmid      int    `json:"vmid"`                                           // 目标用户 mid
	OrderType string `json:"order_type,omitempty" request:"query,omitempty"` // 排序方式。当目标用户为自己时有效。按照关注顺序排列：留空。按照最常访问排列：attention
	Ps        int    `json:"ps,omitempty" request:"query,omitempty"`         // 每页项数。默认为 50
	Pn        int    `json:"pn,omitempty" request:"query,omitempty"`         // 页码。默认为 1。其他用户仅可查看前 100 个
}

type GetUserFollowingsResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
	Total     int            `json:"total"`      // 关注总数
}

// GetUserFollowings 查询用户关注明细（需要登录）
func (c *Client) GetUserFollowings(param GetUserFollowingsParam) (*GetUserFollowingsResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/followings"
	)
	return execute[*GetUserFollowingsResult](c, method, url, param)
}

type GetUserFollowings2Param struct {
	Vmid  int    `json:"vmid"`                                      // 目标用户 mid
	Order string `json:"order,omitempty" request:"query,omitempty"` // 排序方式。按照降序排列：desc。按照升序排列：asc。默认降序排列
	Ps    int    `json:"ps,omitempty" request:"query,omitempty"`    // 每页项数。默认为 50
	Pn    int    `json:"pn,omitempty" request:"query,omitempty"`    // 页码。默认为 1。仅可查看前 5 页
}

type UserFollowingsDetail2 struct {
	Mid            int            `json:"mid"`             // 用户 mid
	Attribute      int            `json:"attribute"`       // 关注属性。0：未关注。2：已关注。6：已互粉
	Mtime          int            `json:"mtime"`           // 关注对方时间。时间戳。互关后刷新
	Tag            []int          `json:"tag"`             // 分组 id
	Special        int            `json:"special"`         // 特别关注标志。0：否。1：是
	Uname          string         `json:"uname"`           // 用户昵称
	Face           string         `json:"face"`            // 用户头像 url
	Sign           string         `json:"sign"`            // 用户签名
	OfficialVerify OfficialVerify `json:"official_verify"` // 认证信息
	Vip            Vip            `json:"vip"`             // 会员信息
	Live           int            `json:"live"`            // 是否直播。0：未直播。1：直播中
}

type GetUserFollowings2Result struct {
	List      []UserFollowingsDetail2 `json:"list"`       // 明细列表
	ReVersion any                     `json:"re_version"` // （？）（可能是number，可能是string）
	Total     int                     `json:"total"`      // 关注总数
}

// GetUserFollowings2 查询用户关注明细2
//
// 仅可查看前 5 页，对于已设置可见性隐私关注列表的用户，则返回的List为nil，Total为0
func (c *Client) GetUserFollowings2(param GetUserFollowings2Param) (*GetUserFollowings2Result, error) {
	const (
		method = resty.MethodGet
		url    = "https://app.biliapi.net/x/v2/relation/followings"
	)
	return execute[*GetUserFollowings2Result](c, method, url, param)
}

type GetUserFollowings3Param struct {
	Vmid int `json:"vmid"`                                   // 目标用户mid
	Ps   int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为20
	Pn   int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为1
}

type UserFollowingsDetail3 struct {
	Mid       string `json:"mid"`       // 用户mid
	Attribute int    `json:"attribute"` // 关注属性。0：未关注。2：已关注。6：已互粉
	Uname     string `json:"uname"`     // 用户昵称
	Face      string `json:"face"`      // 用户头像url
}

type GetUserFollowings3Result struct {
	List []UserFollowingsDetail3 `json:"list"` // 明细列表
}

// GetUserFollowings3 查询用户关注明细3
//
// 对于设置了可见性隐私关注列表的用户会返回空列表
func (c *Client) GetUserFollowings3(param GetUserFollowings3Param) (*GetUserFollowings3Result, error) {
	const (
		method = resty.MethodGet
		url    = "https://line3-h5-mobile-api.biligame.com/game/center/h5/user/relationship/following_list"
	)
	return execute[*GetUserFollowings3Result](c, method, url, param)
}

type SearchUserFollowingsParam struct {
	Vmid string `json:"vmid"`                                     // 目标用户 mid
	Name string `json:"name,omitempty" request:"query,omitempty"` // 搜索关键词
	Ps   int    `json:"ps,omitempty" request:"query,omitempty"`   // 每页项数。默认为 50
	Pn   int    `json:"pn,omitempty" request:"query,omitempty"`   // 页码。默认为 1
}

type SearchUserFollowingsResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
	Total     int            `json:"total"`      // 关注总数
}

// SearchUserFollowings 搜索关注明细
func (c *Client) SearchUserFollowings(param SearchUserFollowingsParam) (*SearchUserFollowingsResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/followings/search"
	)
	return execute[*SearchUserFollowingsResult](c, method, url, param)
}

type GetSameFollowingsParam struct {
	Vmid int `json:"vmid"`                                   // 目标用户 mid
	Ps   int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为 50
	Pn   int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为 1
}

type GetSameFollowingsResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
	Total     int            `json:"total"`      // 关注总数
}

// GetSameFollowings 查询共同关注明细
func (c *Client) GetSameFollowings(param GetSameFollowingsParam) (*GetSameFollowingsResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/same/followings"
	)
	return execute[*GetSameFollowingsResult](c, method, url, param)
}

type GetWhispersResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
}

// GetWhispers 查询悄悄关注明细
func (c *Client) GetWhispers() (*GetWhispersResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/whispers"
	)
	return execute[*GetWhispersResult](c, method, url, nil)
}

type GetFriendsResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
}

// GetFriends 查询互相关注明细
func (c *Client) GetFriends() (*GetFriendsResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/friends"
	)
	return execute[*GetFriendsResult](c, method, url, nil)
}

type GetBlacksParam struct {
	Ps int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为 50
	Pn int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为 1
}

type GetBlacksResult struct {
	List      []RelationUser `json:"list"`       // 明细列表
	ReVersion any            `json:"re_version"` // （？）（可能是number，可能是string）
}

// GetBlacks 查询黑名单明细
func (c *Client) GetBlacks(param GetBlacksParam) (*GetBlacksResult, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/blacks"
	)
	return execute[*GetBlacksResult](c, method, url, param)
}

type ModifyRelationAct int

const (
	ModifyRelationActFollow     ModifyRelationAct = iota + 1 // 关注，无法对已注销或不存在的用户进行此操作
	ModifyRelationActUnfollow                                // 取关
	ModifyRelationActWhisper                                 // 悄悄关注，现已下线，使用本操作代码请求接口会提示“请求错误”
	ModifyRelationActUnwhisper                               // 取消悄悄关注
	ModifyRelationActBlack                                   // 拉黑
	ModifyRelationActUnblack                                 // 取消拉黑
	ModifyRelationActUnfollower                              // 踢出粉丝
)

type ModifyRelationParam struct {
	Fid   int               `json:"fid"`    // 目标用户mid
	Act   ModifyRelationAct `json:"act"`    // 操作代码
	ReSrc int               `json:"re_src"` // 关注来源代码。空间：11。视频：14。文章：115。活动页面：222
}

// ModifyRelation 操作用户关系
func (c *Client) ModifyRelation(param ModifyRelationParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/relation/modify"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type BatchModifyRelationParam struct {
	Fids  []int             `json:"fids"`   // 目标用户 mid 列表
	Act   ModifyRelationAct `json:"act"`    // 操作代码。仅可为 1 或 5，故只能进行批量关注和拉黑
	ReSrc int               `json:"re_src"` // 关注来源代码。同上
}

type BatchModifyRelationResult struct {
	FailedFids []int `json:"failed_fids"` // 操作失败的 mid 列表
}

// BatchModifyRelation 批量操作用户关系
func (c *Client) BatchModifyRelation(param BatchModifyRelationParam) (*BatchModifyRelationResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/relation/batch/modify"
	)
	return execute[*BatchModifyRelationResult](c, method, url, param, fillCsrf(c))
}

type GetUserRelationParam struct {
	Fid int `json:"fid"` // 目标用户 mid
}

type RelationDetail struct {
	Mid       int   `json:"mid"`       // 目标用户 mid
	Attribute int   `json:"attribute"` // 关系属性。0：未关注。2：已关注。6：已互粉。128：已拉黑
	Mtime     int   `json:"mtime"`     // 关注对方时间。时间戳。未关注为 0
	Tag       []int `json:"tag"`       // 分组 id
	Special   int   `json:"special"`   // 特别关注标志。0：否。1：是
}

// GetUserRelation 查询用户与自己关系（仅关注）
func (c *Client) GetUserRelation(param GetUserRelationParam) (*RelationDetail, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation"
	)
	return execute[*RelationDetail](c, method, url, param)
}

type GetUserRelation2Param struct {
	Mid int `json:"mid"` // 目标用户mid
}

type GetUserRelation2Result struct {
	Relation   RelationDetail `json:"relation"`    // 目标用户对于当前用户的关系
	BeRelation RelationDetail `json:"be_relation"` // 当前用户对于目标用户的关系
}

// GetUserRelation2 查询用户与自己关系（互相关系）
func (c *Client) GetUserRelation2(param GetUserRelation2Param) (*GetUserRelation2Result, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/space/wbi/acc/relation"
	)
	return execute[*GetUserRelation2Result](c, method, url, param, fillWbiHandler(c.wbi, c.GetCookies()))
}

type BatchGetUserRelationParam struct {
	Fids []int `json:"fids"` // 目标用户 mid
}

// BatchGetUserRelation 批量查询用户与自己关系
func (c *Client) BatchGetUserRelation(param BatchGetUserRelationParam) (map[int]*RelationDetail, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/relations"
	)
	return execute[map[int]*RelationDetail](c, method, url, param)
}

type RelationTag struct {
	Tagid int    `json:"tagid"` // 分组 id。-10：特别关注。0：默认分组
	Name  string `json:"name"`  // 分组名称
	Count int    `json:"count"` // 分组成员数
	Tip   string `json:"tip"`   // 提示信息
}

// GetRelationTags 查询关注分组列表
func (c *Client) GetRelationTags() ([]RelationTag, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/relation/tags"
	)
	return execute[[]RelationTag](c, method, url, nil)
}
