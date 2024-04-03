package bilibili

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

type GetVideoDetailInfoParam struct {
	Aid  int    `json:"aid,omitempty" request:"query,omitempty"`  // 稿件avid。avid与bvid任选一个
	Bvid string `json:"bvid,omitempty" request:"query,omitempty"` // 稿件bvid。avid与bvid任选一个
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

type CardVip struct {
	Type               int    `json:"type"`                 // 会员类型。0：无。1：月大会员。2：年度及以上大会员
	Status             int    `json:"status"`               // 会员状态。0：无。1：有
	DueDate            int    `json:"due_date"`             // 会员过期时间。Unix时间戳(毫秒)
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

type Card struct {
	Mid            string         `json:"mid"`              // 用户mid
	Name           string         `json:"name"`             // 用户昵称
	Approve        bool           `json:"approve"`          // false。作用尚不明确
	Sex            string         `json:"sex"`              // 用户性别。男 女 保密
	Rank           string         `json:"rank"`             // 10000。作用尚不明确
	Face           string         `json:"face"`             // 用户头像链接
	FaceNft        int            `json:"face_nft"`         // 是否为 nft 头像。0不是nft头像。1是 nft 头像
	Displayrank    string         `json:"DisplayRank"`      // 0。作用尚不明确
	Regtime        int            `json:"regtime"`          // 0。作用尚不明确
	Spacesta       int            `json:"spacesta"`         // 0。作用尚不明确
	Birthday       string         `json:"birthday"`         // 空。作用尚不明确
	Place          string         `json:"place"`            // 空。作用尚不明确
	Description    string         `json:"description"`      // 空。作用尚不明确
	Article        int            `json:"article"`          // 0。作用尚不明确
	Attentions     []any          `json:"attentions"`       // 空。作用尚不明确
	Fans           int            `json:"fans"`             // 粉丝数
	Friend         int            `json:"friend"`           // 关注数
	Attention      int            `json:"attention"`        // 关注数
	Sign           string         `json:"sign"`             // 签名
	LevelInfo      LevelInfo      `json:"level_info"`       // 等级
	Pendant        Pendant        `json:"pendant"`          // 挂件
	Nameplate      Nameplate      `json:"nameplate"`        // 勋章
	Official       Official       `json:"Official"`         // 认证信息
	OfficialVerify OfficialVerify `json:"official_verify"`  // 认证信息2
	Vip            CardVip        `json:"vip"`              // 大会员状态
	IsSeniorMember int            `json:"is_senior_member"` // 是否为硬核会员。0：否。1：是
}

type CardSpace struct {
	SImg string `json:"s_img"` // 主页头图url 小图
	LImg string `json:"l_img"` // 主页头图url 正常
}

type VideoDetailInfoCard struct {
	Card         Card      `json:"card"`          // UP主名片信息
	Space        CardSpace `json:"space"`         // 主页头图
	Following    bool      `json:"following"`     // 是否关注此用户。true：已关注。false：未关注。需要登录(Cookie) 。未登录为false
	ArchiveCount int       `json:"archive_count"` // 用户稿件数
	ArticleCount int       `json:"article_count"` // 用户专栏数
	Follower     int       `json:"follower"`      // 粉丝数
	LikeNum      int       `json:"like_num"`      // UP主获赞次数
}

type VideoDetailInfo struct {
	View      VideoInfo           `json:"View"`       // 视频基本信息
	Card      VideoDetailInfoCard `json:"Card"`       // 视频UP主信息
	Tags      []VideoTag          `json:"Tags"`       // 视频TAG信息
	Reply     CommentsHotReply    `json:"Reply"`      // 视频热评信息
	Related   []VideoInfo         `json:"Related"`    // 推荐视频信息
	Spec      any                 `json:"Spec"`       // ？。作用尚不明确
	HotShare  any                 `json:"hot_share"`  // ？。作用尚不明确
	Elec      any                 `json:"elec"`       // ？。作用尚不明确
	Recommend any                 `json:"recommend"`  // ？。作用尚不明确
	ViewAddit any                 `json:"view_addit"` // ？。作用尚不明确
}

// GetVideoDetailInfo 获取视频超详细信息
func (c *Client) GetVideoDetailInfo(param GetVideoDetailInfoParam) (*VideoDetailInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/view/detail"
	)
	return execute[*VideoDetailInfo](c, method, url, param)
}

type GetVideoRecommendListParam struct {
	Aid  int    `json:"aid,omitempty" request:"query,omitempty"`  // 稿件avid。avid与bvid任选一个
	Bvid string `json:"bvid,omitempty" request:"query,omitempty"` // 稿件bvid。avid与bvid任选一个
}

// GetVideoRecommendList 获取单视频推荐列表
func (c *Client) GetVideoRecommendList(param GetVideoRecommendListParam) ([]VideoInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/archive/related"
	)
	return execute[[]VideoInfo](c, method, url, param)
}

type GetVideoInfoParam struct {
	Aid  int    `json:"aid,omitempty" request:"query,omitempty"`  // 稿件avid。avid与bvid任选一个
	Bvid string `json:"bvid,omitempty" request:"query,omitempty"` // 稿件bvid。avid与bvid任选一个
}

type DescV2 struct {
	RawText string `json:"raw_text"` // 简介内容。type=1时显示原文。type=2时显示'@'+raw_text+' '并链接至biz_id的主页
	Type    int    `json:"type"`     // 类型。1：普通，2：@他人
	BizId   int    `json:"biz_id"`   // 被@用户的mid。=0，当type=1
}

type VideoRights struct {
	Bp            int `json:"bp"`              // 是否允许承包
	Elec          int `json:"elec"`            // 是否支持充电
	Download      int `json:"download"`        // 是否允许下载
	Movie         int `json:"movie"`           // 是否电影
	Pay           int `json:"pay"`             // 是否PGC付费
	Hd5           int `json:"hd5"`             // 是否有高码率
	NoReprint     int `json:"no_reprint"`      // 是否显示“禁止转载”标志
	Autoplay      int `json:"autoplay"`        // 是否自动播放
	UgcPay        int `json:"ugc_pay"`         // 是否UGC付费
	IsCooperation int `json:"is_cooperation"`  // 是否为联合投稿
	UgcPayPreview int `json:"ugc_pay_preview"` // 0。作用尚不明确
	NoBackground  int `json:"no_background"`   // 0。作用尚不明确
	CleanMode     int `json:"clean_mode"`      // 0。作用尚不明确
	IsSteinGate   int `json:"is_stein_gate"`   // 是否为互动视频
	Is360         int `json:"is_360"`          // 是否为全景视频
	NoShare       int `json:"no_share"`        // 0。作用尚不明确
	ArcPay        int `json:"arc_pay"`         // 0。作用尚不明确
	FreeWatch     int `json:"free_watch"`      // 0。作用尚不明确
}

type Owner struct {
	Mid  int    `json:"mid"`  // UP主mid
	Name string `json:"name"` // UP主昵称
	Face string `json:"face"` // UP主头像
}

type VideoStat struct {
	Aid        int    `json:"aid"`        // 稿件avid
	View       int    `json:"view"`       // 播放数
	Danmaku    int    `json:"danmaku"`    // 弹幕数
	Reply      int    `json:"reply"`      // 评论数
	Favorite   int    `json:"favorite"`   // 收藏数
	Coin       int    `json:"coin"`       // 投币数
	Share      int    `json:"share"`      // 分享数
	NowRank    int    `json:"now_rank"`   // 当前排名
	HisRank    int    `json:"his_rank"`   // 历史最高排行
	Like       int    `json:"like"`       // 获赞数
	Dislike    int    `json:"dislike"`    // 点踩数。恒为0
	Evaluation string `json:"evaluation"` // 视频评分
	Vt         int    `json:"vt"`         // 作用尚不明确。恒为0
}

type VideoSubtitleAuthor struct {
	Mid           int    `json:"mid"`             // 字幕上传者mid
	Name          string `json:"name"`            // 字幕上传者昵称
	Sex           string `json:"sex"`             // 字幕上传者性别。男 女 保密
	Face          string `json:"face"`            // 字幕上传者头像url
	Sign          string `json:"sign"`            // 字幕上传者签名
	Rank          int    `json:"rank"`            // 10000。作用尚不明确
	Birthday      int    `json:"birthday"`        // 0。作用尚不明确
	IsFakeAccount int    `json:"is_fake_account"` // 0。作用尚不明确
	IsDeleted     int    `json:"is_deleted"`      // 0。作用尚不明确
}

type VideoSubtitle struct {
	Id          int                 `json:"id"`           // 字幕id
	Lan         string              `json:"lan"`          // 字幕语言
	LanDoc      string              `json:"lan_doc"`      // 字幕语言名称
	IsLock      bool                `json:"is_lock"`      // 是否锁定
	AuthorMid   int                 `json:"author_mid"`   // 字幕上传者mid
	SubtitleUrl string              `json:"subtitle_url"` // json格式字幕文件url
	Author      VideoSubtitleAuthor `json:"author"`       // 字幕上传者信息
}

type VideoSubtitles struct {
	AllowSubmit bool            `json:"allow_submit"` // 是否允许提交字幕
	List        []VideoSubtitle `json:"list"`         // 字幕列表
}

type StaffVip struct {
	Type      int `json:"type"`       // 成员会员类型。0：无。1：月会员。2：年会员
	Status    int `json:"status"`     // 会员状态。0：无。1：有
	ThemeType int `json:"theme_type"` // 0
}

type Official struct {
	Role  int    `json:"role"`  // 成员认证级别
	Title string `json:"title"` // 成员认证名。无为空
	Desc  string `json:"desc"`  // 成员认证备注。无为空
	Type  int    `json:"type"`  // 成员认证类型。-1：无。0：有
}

type Staff struct {
	Mid        int      `json:"mid"`      // 成员mid
	Title      string   `json:"title"`    // 成员名称
	Name       string   `json:"name"`     // 成员昵称
	Face       string   `json:"face"`     // 成员头像url
	Vip        StaffVip `json:"vip"`      // 成员大会员状态
	Official   Official `json:"official"` // 成员认证信息
	Follower   int      `json:"follower"` // 成员粉丝数
	LabelStyle int      `json:"label_style"`
}

type UserGarb struct {
	UrlImageAniCut string `json:"url_image_ani_cut"` // 某url？
}

type Honor struct {
	Aid                int    `json:"aid"`  // 当前稿件aid
	Type               int    `json:"type"` // 1：入站必刷收录。2：第?期每周必看。3：全站排行榜最高第?名。4：热门
	Desc               string `json:"desc"` // 描述
	WeeklyRecommendNum int    `json:"weekly_recommend_num"`
}

type HonorReply struct {
	Honor []Honor `json:"honor"`
}

type ArgueInfo struct {
	ArgueLink string `json:"argue_link"` // 作用尚不明确
	ArgueMsg  string `json:"argue_msg"`  // 警告/争议提示信息
	ArgueType int    `json:"argue_type"` // 作用尚不明确
}

type VideoInfo struct {
	Bvid               string        `json:"bvid"`         // 稿件bvid
	Aid                int           `json:"aid"`          // 稿件avid
	Videos             int           `json:"videos"`       // 稿件分P总数。默认为1
	Tid                int           `json:"tid"`          // 分区tid
	Tname              string        `json:"tname"`        // 子分区名称
	Copyright          int           `json:"copyright"`    // 视频类型。1：原创。2：转载
	Pic                string        `json:"pic"`          // 稿件封面图片url
	Title              string        `json:"title"`        // 稿件标题
	Pubdate            int           `json:"pubdate"`      // 稿件发布时间。秒级时间戳
	Ctime              int           `json:"ctime"`        // 用户投稿时间。秒级时间戳
	Desc               string        `json:"desc"`         // 视频简介
	DescV2             []DescV2      `json:"desc_v2"`      // 新版视频简介
	State              int           `json:"state"`        // 视频状态。详情见[属性数据文档](attribute_data.md#state字段值(稿件状态))
	Duration           int           `json:"duration"`     // 稿件总时长(所有分P)。单位为秒
	Forward            int           `json:"forward"`      // 撞车视频跳转avid。仅撞车视频存在此字段
	MissionId          int           `json:"mission_id"`   // 稿件参与的活动id
	RedirectUrl        string        `json:"redirect_url"` // 重定向url。仅番剧或影视视频存在此字段。用于番剧&影视的av/bv->ep
	Rights             VideoRights   `json:"rights"`       // 视频属性标志
	Owner              Owner         `json:"owner"`        // 视频UP主信息
	Stat               VideoStat     `json:"stat"`         // 视频状态数
	Dynamic            string        `json:"dynamic"`      // 视频同步发布的的动态的文字内容
	Cid                int           `json:"cid"`          // 视频1P cid
	Dimension          Dimension     `json:"dimension"`    // 视频1P分辨率
	Premiere           any           `json:"premiere"`     // null
	TeenageMode        int           `json:"teenage_mode"`
	IsChargeableSeason bool          `json:"is_chargeable_season"`
	IsStory            bool          `json:"is_story"`
	NoCache            bool          `json:"no_cache"` // 作用尚不明确
	Pages              []VideoPage   `json:"pages"`    // 视频分P列表
	Subtitle           VideoSubtitle `json:"subtitle"` // 视频CC字幕信息
	Staff              []Staff       `json:"staff"`    // 合作成员列表。非合作视频无此项
	IsSeasonDisplay    bool          `json:"is_season_display"`
	UserGarb           UserGarb      `json:"user_garb"` // 用户装扮信息
	HonorReply         HonorReply    `json:"honor_reply"`
	LikeIcon           string        `json:"like_icon"`
	ArgueInfo          ArgueInfo     `json:"argue_info"` // 争议/警告信息
}

// GetVideoInfo 获取视频详细信息
func (c *Client) GetVideoInfo(param GetVideoInfoParam) (*VideoInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/web-interface/view"
	)
	return execute[*VideoInfo](c, method, url, param)
}

// GetVideoDescByAvid 通过Avid获取视频简介
func (c *Client) GetVideoDescByAvid(avid int) (string, error) {
	resp, err := c.resty.R().
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/archive/desc")
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return "", errors.Errorf("获取视频简介失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return "", errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return "", formatError("获取视频简介", code, res.Get("message").String())
	}
	return res.Get("data").String(), errors.WithStack(err)
}

// GetVideoDescByBvid 通过Bvid获取视频简介
func (c *Client) GetVideoDescByBvid(bvid string) (string, error) {
	resp, err := c.resty.R().
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/archive/desc")
	if err != nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != 200 {
		return "", errors.Errorf("获取视频简介失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return "", errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return "", formatError("获取视频简介", code, res.Get("message").String())
	}
	return res.Get("data").String(), errors.WithStack(err)
}

// GetVideoDescByShortUrl 通过短链接获取视频简介
func (c *Client) GetVideoDescByShortUrl(shortUrl string) (string, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return "", errors.New("短链接不是视频链接")
	}
	if err != nil {
		return "", err
	}
	return c.GetVideoDescByBvid(bvid.(string))
}

type Dimension struct {
	Width  int `json:"width"`  // 当前分P 宽度
	Height int `json:"height"` // 当前分P 高度
	Rotate int `json:"rotate"` // 是否将宽高对换。0：正常。1：对换
}

type VideoPage struct {
	Cid       int       `json:"cid"`       // 分P cid
	Page      int       `json:"page"`      // 分P序号。从1开始
	From      string    `json:"from"`      // 视频来源。vupload：普通上传（B站）。hunan：芒果TV。qq：腾讯
	Part      string    `json:"part"`      // 分P标题
	Duration  int       `json:"duration"`  // 分P持续时间。单位为秒
	Vid       string    `json:"vid"`       // 站外视频vid。仅站外视频有效
	Weblink   string    `json:"weblink"`   // 站外视频跳转url。仅站外视频有效
	Dimension Dimension `json:"dimension"` // 当前分P分辨率。部分较老视频无分辨率值
}

// GetVideoPageListByAvid 通过Avid获取视频分P列表(Avid转cid)
func (c *Client) GetVideoPageListByAvid(avid int) ([]*VideoPage, error) {
	resp, err := c.resty.R().
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/player/pagelist")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频分P列表")
	if err != nil {
		return nil, err
	}
	var ret []*VideoPage
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoPageListByBvid 通过Bvid获取视频分P列表(Bvid转cid)
func (c *Client) GetVideoPageListByBvid(bvid string) ([]*VideoPage, error) {
	resp, err := c.resty.R().
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/player/pagelist")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频分P列表")
	if err != nil {
		return nil, err
	}
	var ret []*VideoPage
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoPageListByShortUrl 通过短链接获取视频分P列表
func (c *Client) GetVideoPageListByShortUrl(shortUrl string) ([]*VideoPage, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return nil, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return nil, err
	}
	return c.GetVideoPageListByBvid(bvid.(string))
}

type StatusCount struct {
	View  int `json:"view"`  // 0。作用尚不明确
	Use   int `json:"use"`   // 视频添加TAG数
	Atten int `json:"atten"` // TAG关注
}

// VideoTag 视频TAG信息
type VideoTag struct {
	TagId        int         `json:"tag_id"`        // tag_id
	TagName      string      `json:"tag_name"`      // TAG名称
	Cover        string      `json:"cover"`         // TAG图片url
	HeadCover    string      `json:"head_cover"`    // TAG页面头图url
	Content      string      `json:"content"`       // TAG介绍
	ShortContent string      `json:"short_content"` // TAG简介
	Type         int         `json:"type"`          // ？？？
	State        int         `json:"state"`         // 0
	Ctime        int         `json:"ctime"`         // 创建时间。时间戳
	Count        StatusCount `json:"count"`         // 状态数
	IsAtten      int         `json:"is_atten"`      // 是否关注。0：未关注。1：已关注。需要登录(Cookie) 。未登录为0
	Likes        int         `json:"likes"`         // 0。作用尚不明确
	Hates        int         `json:"hates"`         // 0。作用尚不明确
	Attribute    int         `json:"attribute"`     // 0。作用尚不明确
	Liked        int         `json:"liked"`         // 是否已经点赞。0：未点赞。1：已点赞。需要登录(Cookie) 。未登录为0
	Hated        int         `json:"hated"`         // 是否已经点踩。0：未点踩。1：已点踩。需要登录(Cookie) 。未登录为0
	ExtraAttr    int         `json:"extra_attr"`    // ? ? ?
}

// GetVideoTagsByAvid 通过Avid获取视频TAG
func (c *Client) GetVideoTagsByAvid(avid int) ([]*VideoTag, error) {
	resp, err := c.resty.R().
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/tag/archive/tags")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频TAG")
	if err != nil {
		return nil, err
	}
	var ret []*VideoTag
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoTagsByBvid 通过Bvid获取视频TAG
func (c *Client) GetVideoTagsByBvid(bvid string) ([]*VideoTag, error) {
	resp, err := c.resty.R().
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/tag/archive/tags")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频TAG")
	if err != nil {
		return nil, err
	}
	var ret []*VideoTag
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoTagsByShortUrl 通过短链接获取视频TAG
func (c *Client) GetVideoTagsByShortUrl(shortUrl string) ([]*VideoTag, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return nil, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return nil, err
	}
	return c.GetVideoTagsByBvid(bvid.(string))
}

// LikeVideoTag 点赞视频TAG，重复访问为取消
func (c *Client) LikeVideoTag(avid, tagId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid":    strconv.Itoa(avid),
		"tag_id": strconv.Itoa(tagId),
		"csrf":   biliJct,
	}).Post("https://api.bilibili.com/x/tag/archive/like2")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "点赞视频TAG")
	return err
}

// HateVideoTag 点踩视频TAG，重复访问为取消
func (c *Client) HateVideoTag(avid, tagId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid":    strconv.Itoa(avid),
		"tag_id": strconv.Itoa(tagId),
		"csrf":   biliJct,
	}).Post("https://api.bilibili.com/x/tag/archive/hate2")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "点踩视频TAG")
	return err
}

// LikeVideoByAvid 通过Avid点赞视频，like为false表示取消点赞
func (c *Client) LikeVideoByAvid(avid int, like bool) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	var likeNum string
	if like {
		likeNum = "1"
	} else {
		likeNum = "2"
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid":  strconv.Itoa(avid),
		"like": likeNum,
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/archive/like")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "点赞视频")
	return err
}

// LikeVideoByBvid 通过Bvid点赞视频，like为false表示取消点赞
func (c *Client) LikeVideoByBvid(bvid string, like bool) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	var likeNum string
	if like {
		likeNum = "1"
	} else {
		likeNum = "2"
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"bvid": bvid,
		"like": likeNum,
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/archive/like")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "点赞视频")
	return err
}

// LikeVideoByShortUrl 通过短链接点赞视频，like为false表示取消点赞
func (c *Client) LikeVideoByShortUrl(shortUrl string, like bool) error {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return errors.New("短链接不是视频链接")
	}
	if err != nil {
		return err
	}
	return c.LikeVideoByBvid(bvid.(string), like)
}

// CoinVideoByAvid 通过Avid投币视频，multiply为投币数量，上限为2，like为是否附加点赞。返回是否附加点赞成功
func (c *Client) CoinVideoByAvid(avid int, multiply int, like bool) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	var likeNum string
	if like {
		likeNum = "1"
	} else {
		likeNum = "0"
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid":         strconv.Itoa(avid),
		"select_like": likeNum,
		"multiply":    strconv.Itoa(multiply),
		"csrf":        biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/coin/add")
	if err != nil {
		return false, errors.WithStack(err)
	}
	data, err := getRespData(resp, "投币视频")
	if err != nil {
		return false, err
	}
	return gjson.GetBytes(data, "like").Bool(), nil
}

// CoinVideoByBvid 通过Bvid投币视频，multiply为投币数量，上限为2，like为是否附加点赞。返回是否附加点赞成功
func (c *Client) CoinVideoByBvid(bvid string, multiply int, like bool) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	var likeNum string
	if like {
		likeNum = "1"
	} else {
		likeNum = "0"
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"bvid":        bvid,
		"select_like": likeNum,
		"multiply":    strconv.Itoa(multiply),
		"csrf":        biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/coin/add")
	if err != nil {
		return false, errors.WithStack(err)
	}
	data, err := getRespData(resp, "投币视频")
	if err != nil {
		return false, err
	}
	return gjson.GetBytes(data, "like").Bool(), nil
}

// CoinVideoByShortUrl 通过短链接投币视频，multiply为投币数量，上限为2，like为是否附加点赞。返回是否附加点赞成功
func (c *Client) CoinVideoByShortUrl(shortUrl string, multiply int, like bool) (bool, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return false, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return false, err
	}
	return c.CoinVideoByBvid(bvid.(string), multiply, like)
}

// FavourVideoByAvid 通过Avid收藏视频，addMediaIds和delMediaIds为要增加/删除的收藏列表，非必填。返回是否为未关注用户收藏
func (c *Client) FavourVideoByAvid(avid int, addMediaIds, delMediaIds []int) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	addMediaIdStr := make([]string, 0, len(addMediaIds))
	for _, id := range addMediaIds {
		addMediaIdStr = append(addMediaIdStr, strconv.Itoa(id))
	}
	delMediaIdStr := make([]string, 0, len(delMediaIds))
	for _, id := range delMediaIds {
		delMediaIdStr = append(delMediaIdStr, strconv.Itoa(id))
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"rid":           strconv.Itoa(avid),
		"type":          "2",
		"add_media_ids": strings.Join(addMediaIdStr, ","),
		"del_media_ids": strings.Join(delMediaIdStr, ","),
		"csrf":          biliJct,
	}).Post("https://api.bilibili.com/medialist/gateway/coll/resource/deal")
	if err != nil {
		return false, errors.WithStack(err)
	}
	data, err := getRespData(resp, "收藏视频")
	if err != nil {
		return false, err
	}
	return gjson.GetBytes(data, "prompt").Bool(), nil
}

// FavourVideoByBvid 通过Bvid收藏视频，addMediaIds和delMediaIds为要增加/删除的收藏列表，非必填。返回是否为未关注用户收藏
func (c *Client) FavourVideoByBvid(bvid string, addMediaIds, delMediaIds []int) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	addMediaIdStr := make([]string, 0, len(addMediaIds))
	for _, id := range addMediaIds {
		addMediaIdStr = append(addMediaIdStr, strconv.Itoa(id))
	}
	delMediaIdStr := make([]string, 0, len(delMediaIds))
	for _, id := range delMediaIds {
		delMediaIdStr = append(delMediaIdStr, strconv.Itoa(id))
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"rid":           bvid,
		"type":          "2",
		"add_media_ids": strings.Join(addMediaIdStr, ","),
		"del_media_ids": strings.Join(delMediaIdStr, ","),
		"csrf":          biliJct,
	}).Post("https://api.bilibili.com/medialist/gateway/coll/resource/deal")
	if err != nil {
		return false, errors.WithStack(err)
	}
	data, err := getRespData(resp, "收藏视频")
	if err != nil {
		return false, err
	}
	return gjson.GetBytes(data, "prompt").Bool(), nil
}

// FavourVideoByShortUrl 通过短链接收藏视频，addMediaIds和delMediaIds为要增加/删除的收藏列表，非必填。返回是否为未关注用户收藏
func (c *Client) FavourVideoByShortUrl(shortUrl string, addMediaIds, delMediaIds []int) (bool, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return false, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return false, err
	}
	return c.FavourVideoByBvid(bvid.(string), addMediaIds, delMediaIds)
}

type LikeCoinFavourResult struct {
	Like     bool `json:"like"`     // 是否点赞成功
	Coin     bool `json:"coin"`     // 是否投币成功
	Fav      bool `json:"fav"`      // 是否收藏成功
	Multiply int  `json:"multiply"` // 投币枚数
}

// LikeCoinFavourVideoByAvid 通过Avid一键三连视频
func (c *Client) LikeCoinFavourVideoByAvid(avid int) (*LikeCoinFavourResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid":  strconv.Itoa(avid),
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/archive/like/triple")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "一键三连视频")
	if err != nil {
		return nil, err
	}
	var ret *LikeCoinFavourResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// LikeCoinFavourVideoByBvid 通过Bvid一键三连视频
func (c *Client) LikeCoinFavourVideoByBvid(bvid string) (*LikeCoinFavourResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"bvid": bvid,
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/archive/like/triple")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "一键三连视频")
	if err != nil {
		return nil, err
	}
	var ret *LikeCoinFavourResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// LikeCoinFavourVideoByShortUrl 通过短链接一键三连视频
func (c *Client) LikeCoinFavourVideoByShortUrl(shortUrl string) (*LikeCoinFavourResult, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return nil, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return nil, err
	}
	return c.LikeCoinFavourVideoByBvid(bvid.(string))
}

type VideoOnlineInfo struct {
	Total      string   `json:"total"` // 所有终端总计人数，例如“10万+”
	Count      string   `json:"count"` // web端实时在线人数
	ShowSwitch struct { // 数据显示控制
		Total bool `json:"total"` // 是否展示所有终端总计人数
		Count bool `json:"count"` // 是否展示web端实时在线人数
	} `json:"show_switch"`
}

// GetVideoOnlineInfoByAvid 通过Avid获取视频在线人数
func (c *Client) GetVideoOnlineInfoByAvid(avid, cid int) (*VideoOnlineInfo, error) {
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"aid": strconv.Itoa(avid),
		"cid": strconv.Itoa(cid),
	}).Get("https://api.bilibili.com/x/player/online/total")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频在线人数")
	if err != nil {
		return nil, err
	}
	var ret *VideoOnlineInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoOnlineInfoByBvid 通过Bvid获取视频在线人数
func (c *Client) GetVideoOnlineInfoByBvid(bvid string, cid int) (*VideoOnlineInfo, error) {
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"bvid": bvid,
		"cid":  strconv.Itoa(cid),
	}).Get("https://api.bilibili.com/x/player/online/total")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频在线人数")
	if err != nil {
		return nil, err
	}
	var ret *VideoOnlineInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoOnlineInfoByShortUrl 通过短链接获取视频在线人数
func (c *Client) GetVideoOnlineInfoByShortUrl(shortUrl string, cid int) (*VideoOnlineInfo, error) {
	typ, bvid, err := c.UnwrapShortUrl(shortUrl)
	if typ != "bvid" {
		return nil, errors.New("短链接不是视频链接")
	}
	if err != nil {
		return nil, err
	}
	return c.GetVideoOnlineInfoByBvid(bvid.(string), cid)
}

type VideoPbPInfo struct {
	StepSec int      `json:"step_sec"` // 采样间隔时间（单位为秒，由视频时长决定）
	Tagstr  string   `json:"tagstr"`   // 作用尚不明确
	Events  struct { // 数据本体
		Default []float64 `json:"default"` // 顶点值列表（顶点个数由视频时长和采样时间决定）
	} `json:"events"`
	Debug string `json:"debug"` // 调试信息（json字串）
}

// GetVideoPbPInfo 获取视频弹幕趋势顶点列表（高能进度条）
func (c *Client) GetVideoPbPInfo(cid int) (*VideoPbPInfo, error) {
	resp, err := c.resty.R().
		SetQueryParam("cid", strconv.Itoa(cid)).Get("https://api.bilibili.com/pbp/data")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频弹幕趋势顶点列表")
	if err != nil {
		return nil, err
	}
	var ret *VideoPbPInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type VideoStatusNumber struct {
	Aid        int         `json:"aid"`        // 稿件avid
	Bvid       string      `json:"bvid"`       // 稿件bvid
	View       interface{} `json:"view"`       // 播放次数（有值则为一个int，如果被屏蔽了则为字符串"--"）
	Danmaku    int         `json:"danmaku"`    // 弹幕条数
	Reply      int         `json:"reply"`      // 评论条数
	Favorite   int         `json:"favorite"`   // 收藏人数
	Coin       int         `json:"coin"`       // 投币枚数
	Share      int         `json:"share"`      // 分享次数
	Like       int         `json:"like"`       // 获赞次数
	NowRank    int         `json:"now_rank"`   // 固定值0，作用尚不明确
	HisRank    int         `json:"his_rank"`   // 历史最高排行
	Dislike    int         `json:"dislike"`    // 固定值0，作用尚不明确
	NoReprint  int         `json:"no_reprint"` // 禁止转载标志，0：无，1：禁止
	Copyright  int         `json:"copyright"`  // 版权标志，1：自制，2：转载
	ArgueMsg   string      `json:"argue_msg"`  // 警告信息
	Evaluation string      `json:"evaluation"` // 视频评分
}

// GetVideoStatusNumberByAvid 通过Avid获取视频状态数视频
func (c *Client) GetVideoStatusNumberByAvid(avid int) (*VideoStatusNumber, error) {
	resp, err := c.resty.R().
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/web-interface/archive/stat")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频状态数视频")
	if err != nil {
		return nil, err
	}
	var ret *VideoStatusNumber
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoStatusNumberByBvid 通过Bvid获取视频状态数
func (c *Client) GetVideoStatusNumberByBvid(bvid string) (*VideoStatusNumber, error) {
	resp, err := c.resty.R().
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/web-interface/archive/stat")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频状态数")
	if err != nil {
		return nil, err
	}
	var ret *VideoStatusNumber
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetTopRecommendVideo 获取首页视频推荐列表，freshType相关性（默认为3），ps单页返回的记录条数（默认为8）
func (c *Client) GetTopRecommendVideo(freshType, ps int) ([]*VideoInfo, error) {
	request := c.resty.R().SetQueryParam("version", "1")
	if freshType != 0 {
		request.SetQueryParam("fresh_type", strconv.Itoa(freshType))
	}
	if ps != 0 {
		request.SetQueryParam("ps", strconv.Itoa(ps))
	}
	resp, err := request.Get("https://api.bilibili.com/x/web-interface/index/top/rcmd")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取首页视频推荐列表")
	if err != nil {
		return nil, err
	}

	type MixedData struct {
		Item []*VideoInfo `json:"item"`
	}

	var mixedRet MixedData
	err = json.Unmarshal(data, &mixedRet)
	return mixedRet.Item, errors.WithStack(err)
}

type ArchivesList struct {
	Aids     []int `json:"aids"`
	Archives []struct {
		Aid              int    `json:"aid"`
		Bvid             string `json:"bvid"`
		Ctime            int    `json:"ctime"`
		Duration         int    `json:"duration"`
		EnableVt         bool   `json:"enable_vt"`
		InteractiveVideo bool   `json:"interactive_video"`
		Pic              string `json:"pic"`
		PlaybackPosition int    `json:"playback_position"`
		Pubdate          int    `json:"pubdate"`
		Stat             struct {
			View int `json:"view"`
			Vt   int `json:"vt"`
		} `json:"stat"`
		State     int    `json:"state"`
		Title     string `json:"title"`
		UgcPay    int    `json:"ugc_pay"`
		VtDisplay string `json:"vt_display"`
	} `json:"archives"`
	Meta struct {
		Category    int    `json:"category"`
		Cover       string `json:"cover"`
		Description string `json:"description"`
		Mid         int    `json:"mid"`
		Name        string `json:"name"`
		Ptime       int    `json:"ptime"`
		SeasonID    int    `json:"season_id"`
		Total       int    `json:"total"`
	} `json:"meta"`
	Page struct {
		PageNum  int `json:"page_num"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	} `json:"page"`
}

// GetArchivesList 获取视频合集信息 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/collection.md#%E8%8E%B7%E5%8F%96%E8%A7%86%E9%A2%91%E5%90%88%E9%9B%86%E4%BF%A1%E6%81%AF
func (c *Client) GetArchivesList(mid int, sid int, pn int, ps int, sortReverse bool) (*ArchivesList, error) {
	postData := map[string]string{
		"mid":          strconv.Itoa(mid),
		"page_num":     strconv.Itoa(pn),
		"page_size":    strconv.Itoa(ps),
		"season_id":    strconv.Itoa(sid),
		"sort_reverse": strconv.FormatBool(sortReverse),
	}
	resp, err := c.resty.R().SetQueryParams(postData).Get("https://api.bilibili.com/x/polymer/web-space/seasons_archives_list")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var ret *ArchivesList
	data, err := getRespData(resp, "获取合集信息")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
