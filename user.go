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
		method = "GET"
		url    = "https://api.bilibili.com/x/space/wbi/arc/search"
	)
	return execute[*UserVideos](c, method, url, param, fillWbiHandler(c.wbi, c.GetCookies()))
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
		method = "GET"
		url    = "https://api.bilibili.com/x/web-interface/card"
	)
	return execute[*UserCard](c, method, url, param)
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
