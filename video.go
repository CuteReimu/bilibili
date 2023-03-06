package bilibili

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
)

var regBv = regexp.MustCompile(`(?i)bv([\dA-Za-z]{10})`)

// GetBvidByShortUrl 通过视频短链接获取bvid
func GetBvidByShortUrl(shortUrl string) (string, error) {
	return std.GetBvidByShortUrl(shortUrl)
}
func (c *Client) GetBvidByShortUrl(shortUrl string) (string, error) {
	resp, err := c.resty().SetRedirectPolicy(resty.NoRedirectPolicy()).R().Get(shortUrl)
	if resp == nil {
		return "", errors.WithStack(err)
	}
	if resp.StatusCode() != 302 {
		return "", errors.Errorf("通过短链接获取视频详细信息失败，status code: %d", resp.StatusCode())
	}
	url := resp.Header().Get("Location")
	ret := regBv.FindString(url)
	if len(ret) == 0 {
		return "", errors.New("无法解析链接：" + url)
	}
	return ret, nil
}

// OfficialInfo 成员认证信息
type OfficialInfo struct {
	Role  int    `json:"role"`  // 成员认证级别，0：无，1 2 7：个人认证，3 4 5 6：机构认证
	Title string `json:"title"` // 成员认证名
	Desc  string `json:"desc"`  // 成员认证备注
	Type  int    `json:"type"`  // 成员认证类型，-1：无，0：有
}

type VideoDimension struct {
	Width  int `json:"width"`  // 当前分P 宽度
	Height int `json:"height"` // 当前分P 高度
	Rotate int `json:"rotate"` // 是否将宽高对换，0：正常，1：对换
}

type VideoInfo struct {
	Bvid      string     `json:"bvid"`      // 稿件bvid
	Aid       int        `json:"aid"`       // 稿件avid
	Videos    int        `json:"videos"`    // 稿件分P总数，默认为1
	Tid       int        `json:"tid"`       // 分区tid
	Tname     string     `json:"tname"`     // 子分区名称
	Copyright int        `json:"copyright"` // 1：原创，2：转载
	Pic       string     `json:"pic"`       // 稿件封面图片url
	Title     string     `json:"title"`     // 稿件标题
	Pubdate   int64      `json:"pubdate"`   // 稿件发布时间戳
	Ctime     int64      `json:"ctime"`     // 用户投稿时间戳
	Desc      string     `json:"desc"`      // 视频简介
	DescV2    []struct { // 新版视频简介
		RawText string `json:"raw_text"` // 简介内容
		Type    int    `json:"type"`     // 作用尚不明确
		BizId   int    `json:"biz_id"`   // 作用尚不明确
	} `json:"desc_v2"`
	State       int      `json:"state"`        // 视频状态
	Duration    int      `json:"duration"`     // 稿件总时长（所有分P，单位：秒）
	Forward     int      `json:"forward"`      // 撞车视频跳转avid
	MissionId   int      `json:"mission_id"`   // 稿件参与的活动id
	RedirectUrl string   `json:"redirect_url"` // 重定向url，仅番剧或影视视频存在此字段
	Rights      struct { // 视频属性标志，全部都是1表示是，0表示否
		Bp            int `json:"bp"`              // 固定值0，作用尚不明确
		Elec          int `json:"elec"`            // 是否支持充电
		Download      int `json:"download"`        // 是否允许下载
		Movie         int `json:"movie"`           // 是否电影
		Pay           int `json:"pay"`             // 是否PGC付费
		Hd5           int `json:"hd5"`             // 是否有高码率
		NoReprint     int `json:"no_reprint"`      // 是否显示“禁止转载“标志
		Autoplay      int `json:"autoplay"`        // 是否自动播放
		UgcPay        int `json:"ugc_pay"`         // 是否UGC付费
		IsSteinGate   int `json:"is_stein_gate"`   // 是否为互动视频
		IsCooperation int `json:"is_cooperation"`  // 是否为联合投稿
		UgcPayPreview int `json:"ugc_pay_preview"` // 固定值0，作用尚不明确
		NoBackground  int `json:"no_background"`   // 固定值0，作用尚不明确
		CleanMode     int `json:"clean_mode"`
		Is360         int `json:"is_360"`
		NoShare       int `json:"no_share"`
		ArcPay        int `json:"arc_pay"`
		FreeWatch     int `json:"free_watch"`
	}
	Owner struct { // 视频UP主信息
		Mid  int    `json:"mid"`  // UP主mid
		Name string `json:"name"` // UP主昵称
		Face string `json:"face"` // UP主头像url
	} `json:"owner"`
	Stat struct { // 视频状态数
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
		Dislike    int    `json:"dislike"`    // 点踩数，恒为0
		Evaluation string `json:"evaluation"` // 视频评分
		ArgueMsg   string `json:"argue_msg"`  // 警告/争议提示信息
	} `json:"stat"`
	Dynamic            string         `json:"dynamic"`   // 视频同步发布的的动态的文字内容
	Cid                int            `json:"cid"`       // 视频1P cid
	Dimension          VideoDimension `json:"dimension"` // 视频1P分辨率
	TeenageMode        int            `json:"teenage_mode"`
	IsChargeableSeason bool           `json:"is_chargeable_season"`
	NoCache            bool           `json:"no_cache"` // 固定值true，作用尚不明确
	Pages              []struct {     // 视频分P列表，无分P则数组只有1个元素
		Cid       int            `json:"cid"`       // 当前分P cid
		Page      int            `json:"page"`      // 当前分P
		From      string         `json:"from"`      // 视频来源，vupload：普通上传（B站），hunan：芒果TV，qq：腾讯
		Part      string         `json:"part"`      // 当前分P标题
		Duration  int            `json:"duration"`  // 当前分P持续时间（单位：秒）
		Vid       string         `json:"vid"`       // 站外视频vid，仅站外视频有效
		Weblink   string         `json:"weblink"`   // 站外视频跳转url，仅站外视频有效
		Dimension VideoDimension `json:"dimension"` // 当前分P分辨率，部分较老视频无分辨率值
	} `json:"pages"`
	Subtitle struct { // 视频CC字幕信息
		AllowCommit bool       `json:"allow_commit"` // 是否允许提交字幕
		List        []struct { // 字幕列表
			Id          int      `json:"id"`           // 字幕id
			Lan         string   `json:"lan"`          // 字幕语言
			LanDoc      string   `json:"lan_doc"`      // 字幕语言名称
			IsLock      bool     `json:"is_lock"`      // 是否锁定
			AuthorMid   int      `json:"author_mid"`   // 字幕上传者mid
			SubtitleUrl string   `json:"subtitle_url"` // json格式字幕文件url
			Author      struct { // 字幕上传者信息
				Mid           int    `json:"mid"`             // 字幕上传者mid
				Name          string `json:"name"`            // 字幕上传者昵称
				Sex           string `json:"sex"`             // 字幕上传者性别 男 女 保密
				Face          string `json:"face"`            // 字幕上传者头像url
				Sign          string `json:"sign"`            // 字幕上传者签名
				Rank          int    `json:"rank"`            // 固定值10000，作用尚不明确
				Birthday      int    `json:"birthday"`        // 固定值0，作用尚不明确
				IsFakeAccount int    `json:"is_fake_account"` // 固定值0，作用尚不明确
				IsDeleted     int    `json:"is_deleted"`      // 固定值0，作用尚不明确
			} `json:"author"`
		} `json:"list"`
	} `json:"subtitle"`
	IsSeasonDisplay bool       `json:"is_season_display"`
	Staff           []struct { // 合作成员列表，非合作视频无此项
		Mid   int      `json:"mid"`   // 成员mid
		Title string   `json:"title"` // 成员名称
		Name  string   `json:"name"`  // 成员昵称
		Face  string   `json:"face"`  // 成员头像url
		Vip   struct { // 成员大会员状态
			Type      int `json:"type"`       // 成员会员类型，0：无，1：月会员，2：年会员
			Status    int `json:"status"`     // 会员状态，0：无，1：有
			ThemeType int `json:"theme_type"` // 固定值0，作用尚不明确
		} `json:"vip"`
		Official OfficialInfo `json:"official"` // 成员认证信息
		Follower int          `json:"follower"` // 成员粉丝数
	} `json:"staff"`
	UserGarb struct { // 用户装扮信息
		UrlImageAniCut string `json:"url_image_ani_cut"` // 某url，作用尚不明确
	} `json:"user_garb"`
	HonorReply struct {
		Honor []struct {
			Aid                int    `json:"aid"`
			Type               int    `json:"type"`
			Desc               string `json:"desc"`
			WeeklyRecommendNum int    `json:"weekly_recommend_num"`
		} `json:"honor"`
	} `json:"honor_reply"`
}

// GetVideoInfoByAvid 通过Avid获取视频信息
func GetVideoInfoByAvid(avid int) (*VideoInfo, error) {
	return std.GetVideoInfoByAvid(avid)
}
func (c *Client) GetVideoInfoByAvid(avid int) (*VideoInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/web-interface/view")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频详细信息")
	if err != nil {
		return nil, err
	}
	var ret *VideoInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoInfoByBvid 通过Bvid获取视频信息
func GetVideoInfoByBvid(bvid string) (*VideoInfo, error) {
	return std.GetVideoInfoByBvid(bvid)
}
func (c *Client) GetVideoInfoByBvid(bvid string) (*VideoInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/web-interface/view")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频详细信息")
	if err != nil {
		return nil, err
	}
	var ret *VideoInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoInfoByShortUrl 通过短链接获取视频信息
func GetVideoInfoByShortUrl(shortUrl string) (*VideoInfo, error) {
	return std.GetVideoInfoByShortUrl(shortUrl)
}
func (c *Client) GetVideoInfoByShortUrl(shortUrl string) (*VideoInfo, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoInfoByBvid(bvid)
}

// GetRecommendVideoByAvid 通过Avid获取推荐视频
func GetRecommendVideoByAvid(avid int) ([]*VideoInfo, error) {
	return std.GetRecommendVideoByAvid(avid)
}
func (c *Client) GetRecommendVideoByAvid(avid int) ([]*VideoInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/web-interface/archive/related")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取推荐视频")
	if err != nil {
		return nil, err
	}
	var ret []*VideoInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetRecommendVideoByBvid 通过Bvid获取推荐视频
func GetRecommendVideoByBvid(bvid string) ([]*VideoInfo, error) {
	return std.GetRecommendVideoByBvid(bvid)
}
func (c *Client) GetRecommendVideoByBvid(bvid string) ([]*VideoInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/web-interface/archive/related")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取推荐视频")
	if err != nil {
		return nil, err
	}
	var ret []*VideoInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type OfficialVerify struct {
	Type int    `json:"type"` // 是否认证，-1：无，0：认证
	Desc string `json:"desc"` // 认证信息，无为空
}

type NamePlate struct {
	Nid        int    `json:"nid"`         // 勋章id
	Name       string `json:"name"`        // 勋章名称
	Image      string `json:"image"`       // 挂件图片url 正常
	ImageSmall string `json:"image_small"` // 勋章图片url 小
	Level      string `json:"level"`       // 勋章等级
	Condition  string `json:"condition"`   // 勋章条件
}

type Vip struct {
	Type       int   `json:"type"`         // 大会员类型，0：无，1：月度大会员，2：年度及以上大会员
	Status     int   `json:"status"`       // 大会员状态，0：无，1：有
	DueDate    int64 `json:"due_date"`     // 到期时间戳（毫秒）
	VipPayType int   `json:"vip_pay_type"` // 大会员付费类型
	ThemeType  int   `json:"theme_type"`   // 固定值0，作用尚不明确
	Label      struct {
		Path                  string `json:"path"`
		Text                  string `json:"text"`          // 大会员标签上的文字
		LabelTheme            string `json:"label_theme"`   // 大会员标签主题
		TextColor             string `json:"text_color"`    // 大会员文字颜色
		BgStyle               int    `json:"bg_style"`      // 大会员背景样式
		BgColor               string `json:"bg_color"`      // 大会员背景颜色
		BorderColor           string `json:"border_color"`  // 大会员边框颜色
		UseImgLabel           bool   `json:"use_img_label"` // 是否使用图片标签
		ImgLabelUriHans       string `json:"img_label_uri_hans"`
		ImgLabelUriHant       string `json:"img_label_uri_hant"`
		ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"` // 大会员图片标签（简体中文）的url
		ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"` // 大会员图片标签（繁体中文）的url
	} `json:"label"`
	AvatarSubscript    int    `json:"avatar_subscript"` // 作用尚不明确
	NicknameColor      string `json:"nickname_color"`   // 昵称颜色
	Role               int    `json:"role"`
	AvatarSubscriptUrl string `json:"avatar_subscript_url"` // 作用尚不明确
	TvVipStatus        int    `json:"tv_vip_status"`        // TV大会员状态，0：无，1：有
	TvVipPayType       int    `json:"tv_vip_pay_type"`      // TV大会员付费类型
	VipType            int    `json:"vipType"`              // 大会员类型，0：无，1：月度大会员，2：年度及以上大会员
	VipStatus          int    `json:"vipStatus"`            // 大会员状态，0：无，1：有
}

type Pendant struct {
	Pid               int    `json:"pid"`    // 挂件id
	Name              string `json:"name"`   // 挂件名称
	Image             string `json:"image"`  // 挂件图片url
	Expire            int    `json:"expire"` // 固定值0，作用尚不明确
	ImageEnhance      string `json:"image_enhance"`
	ImageEnhanceFrame string `json:"image_enhance_frame"`
}

type VideoDetailInfo struct {
	View VideoInfo `json:"View"` // 视频基本信息
	Card struct {  // 视频UP主信息
		Card struct { // UP主名片信息
			Mid         string   `json:"mid"`           // 用户mid
			Name        string   `json:"name"`          // 用户昵称
			Approve     bool     `json:"approve"`       // 固定值false，作用尚不明确
			Sex         string   `json:"sex"`           // 用户性别 男 女 保密
			Rank        string   `json:"rank"`          // 固定值"10000"，作用尚不明确
			Face        string   `json:"face"`          // 用户头像链接
			FaceNft     int      `json:"face_nft"`      // 是否为 nft 头像，0：不是nft头像，1：是 nft 头像
			FaceNftType int      `json:"face_nft_type"` // ntf 头像类型
			DisplayRank string   `json:"DisplayRank"`   // 固定值"0"，作用尚不明确
			Regtime     int      `json:"regtime"`       // 固定值0，作用尚不明确
			Spacesta    int      `json:"spacesta"`      // 固定值0，作用尚不明确
			Birthday    string   `json:"birthday"`      // 固定值""，作用尚不明确
			Place       string   `json:"place"`         // 固定值""，作用尚不明确
			Description string   `json:"description"`   // 固定值""，作用尚不明确
			Article     int      `json:"article"`       // 固定值0，作用尚不明确
			Fans        int      `json:"fans"`          // 粉丝数
			Friend      int      `json:"friend"`        // 关注数
			Attention   int      `json:"attention"`     // 关注数
			Sign        string   `json:"sign"`          // 签名
			LevelInfo   struct { // 等级
				CurrentLevel int `json:"current_level"` // 当前等级，0-6级
				CurrentMin   int `json:"current_min"`   // 固定值0，作用尚不明确
				CurrentExp   int `json:"current_exp"`   // 固定值0，作用尚不明确
				NextExp      int `json:"next_exp"`      // 固定值0，作用尚不明确
			} `json:"level_info"`
			Pendant        Pendant        `json:"pendant"`          // 挂件
			Nameplate      NamePlate      `json:"nameplate"`        // 勋章
			Official       OfficialInfo   `json:"Official"`         // 认证信息
			OfficialVerify OfficialVerify `json:"official_verify"`  // 认证信息2
			Vip            Vip            `json:"vip"`              // 大会员状态
			IsSeniorMember int            `json:"is_senior_member"` // 是否为硬核会员，0：否，1：是
		} `json:"card"`
		Space struct { // 主页头图
			SImg string `json:"s_img"` // 主页头图url 小图
			LImg string `json:"l_img"` // 主页头图url 正常
		} `json:"space"`
		Following    bool `json:"following"`     // 是否关注此用户，true：已关注，false：未关注，需要登录(Cookie)，未登录为false
		ArchiveCount int  `json:"archive_count"` // 用户稿件数
		ArticleCount int  `json:"article_count"` // 固定值0，作用尚不明确
		Follower     int  `json:"follower"`      // 粉丝数
		LikeNum      int  `json:"like_num"`      // UP主获赞次数
	} `json:"Card"`
	Tags     []VideoTag  `json:"Tags"`    // 视频TAG信息
	Reply    HotReply    `json:"Reply"`   // 视频热评信息
	Related  []VideoInfo `json:"Related"` // 推荐视频信息
	HotShare struct {
		Show bool `json:"show"` // 固定为false，作用尚不明确
	} `json:"hot_share"`
	ViewAddit struct {
		Field1 bool `json:"63"` // 固定为false，作用尚不明确
		Field2 bool `json:"64"` // 固定为false，作用尚不明确
	} `json:"view_addit"`
}

// GetVideoDetailInfoByAvid 通过Avid获取视频超详细信息
func GetVideoDetailInfoByAvid(avid int) (*VideoDetailInfo, error) {
	return std.GetVideoDetailInfoByAvid(avid)
}
func (c *Client) GetVideoDetailInfoByAvid(avid int) (*VideoDetailInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("aid", strconv.Itoa(avid)).Get("https://api.bilibili.com/x/web-interface/view/detail")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频超详细信息")
	if err != nil {
		return nil, err
	}
	var ret *VideoDetailInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoDetailInfoByBvid 通过Bvid获取视频超详细信息
func GetVideoDetailInfoByBvid(bvid string) (*VideoDetailInfo, error) {
	return std.GetVideoDetailInfoByBvid(bvid)
}
func (c *Client) GetVideoDetailInfoByBvid(bvid string) (*VideoDetailInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("bvid", bvid).Get("https://api.bilibili.com/x/web-interface/view/detail")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取视频超详细信息")
	if err != nil {
		return nil, err
	}
	var ret *VideoDetailInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetVideoDetailInfoByShortUrl 通过短链接获取视频超详细信息
func GetVideoDetailInfoByShortUrl(shortUrl string) (*VideoDetailInfo, error) {
	return std.GetVideoDetailInfoByShortUrl(shortUrl)
}
func (c *Client) GetVideoDetailInfoByShortUrl(shortUrl string) (*VideoDetailInfo, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoDetailInfoByBvid(bvid)
}

// GetVideoDescByAvid 通过Avid获取视频简介
func GetVideoDescByAvid(avid int) (string, error) {
	return std.GetVideoDescByAvid(avid)
}
func (c *Client) GetVideoDescByAvid(avid int) (string, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoDescByBvid(bvid string) (string, error) {
	return std.GetVideoDescByBvid(bvid)
}
func (c *Client) GetVideoDescByBvid(bvid string) (string, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoDescByShortUrl(shortUrl string) (string, error) {
	return std.GetVideoDescByShortUrl(shortUrl)
}
func (c *Client) GetVideoDescByShortUrl(shortUrl string) (string, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return "", err
	}
	return GetVideoDescByBvid(bvid)
}

type VideoPage struct {
	Cid       int    `json:"cid"`
	Page      int    `json:"page"`
	From      string `json:"from"`
	Part      string `json:"part"`
	Duration  int    `json:"duration"`
	Vid       string `json:"vid"`
	Weblink   string `json:"weblink"`
	Dimension struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		Rotate int `json:"rotate"`
	} `json:"dimension"`
}

// GetVideoPageListByAvid 通过Avid获取视频分P列表(Avid转cid)
func GetVideoPageListByAvid(avid int) ([]*VideoPage, error) {
	return std.GetVideoPageListByAvid(avid)
}
func (c *Client) GetVideoPageListByAvid(avid int) ([]*VideoPage, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoPageListByBvid(bvid string) ([]*VideoPage, error) {
	return std.GetVideoPageListByBvid(bvid)
}
func (c *Client) GetVideoPageListByBvid(bvid string) ([]*VideoPage, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoPageListByShortUrl(shortUrl string) ([]*VideoPage, error) {
	return std.GetVideoPageListByShortUrl(shortUrl)
}
func (c *Client) GetVideoPageListByShortUrl(shortUrl string) ([]*VideoPage, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoPageListByBvid(bvid)
}

// VideoTag 视频TAG信息
type VideoTag []struct {
	TagId        int      `json:"tag_id"`        // tag_id
	TagName      string   `json:"tag_name"`      // TAG名称
	Cover        string   `json:"cover"`         // TAG图片url
	HeadCover    string   `json:"head_cover"`    // TAG页面头图url
	Content      string   `json:"content"`       // TAG介绍
	ShortContent string   `json:"short_content"` // TAG简介
	Type         int      `json:"type"`          // 作用尚不明确
	State        int      `json:"state"`         // 固定值0，作用尚不明确
	Ctime        int      `json:"ctime"`         // 创建时间戳
	Count        struct { // 状态数
		View  int `json:"view"`  // 固定值0，作用尚不明确
		Use   int `json:"use"`   // 视频添加TAG数
		Atten int `json:"atten"` // TAG关注
	} `json:"count"`
	IsAtten         int    `json:"is_atten"`   // 是否关注，0：未关注，1：已关注，需要登录(Cookie)，未登录为0
	Likes           int    `json:"likes"`      // 固定值0，作用尚不明确
	Hates           int    `json:"hates"`      // 固定值0，作用尚不明确
	Attribute       int    `json:"attribute"`  // 固定值0，作用尚不明确
	Liked           int    `json:"liked"`      // 是否已经点赞，0：未点赞，1：已点赞，需要登录(Cookie)，未登录为0
	Hated           int    `json:"hated"`      // 是否已经点踩，0：未点踩，1：已点踩，需要登录(Cookie)，未登录为0
	ExtraAttr       int    `json:"extra_attr"` // 作用尚不明确
	MusicId         string `json:"music_id"`
	TagType         string `json:"tag_type"`
	IsActivity      bool   `json:"is_activity"`
	Color           string `json:"color"`
	Alpha           int    `json:"alpha"`
	IsSeason        bool   `json:"is_season"`
	SubscribedCount int    `json:"subscribed_count"`
	ArchiveCount    string `json:"archive_count"`
	FeaturedCount   int    `json:"featured_count"`
	JumpUrl         string `json:"jump_url"`
}

// GetVideoTagsByAvid 通过Avid获取视频TAG
func GetVideoTagsByAvid(avid int) ([]*VideoTag, error) {
	return std.GetVideoTagsByAvid(avid)
}
func (c *Client) GetVideoTagsByAvid(avid int) ([]*VideoTag, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoTagsByBvid(bvid string) ([]*VideoTag, error) {
	return std.GetVideoTagsByBvid(bvid)
}
func (c *Client) GetVideoTagsByBvid(bvid string) ([]*VideoTag, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoTagsByShortUrl(shortUrl string) ([]*VideoTag, error) {
	return std.GetVideoTagsByShortUrl(shortUrl)
}
func (c *Client) GetVideoTagsByShortUrl(shortUrl string) ([]*VideoTag, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoTagsByBvid(bvid)
}

// LikeVideoTag 点赞视频TAG，重复访问为取消
func LikeVideoTag(avid, tagId int) error {
	return std.LikeVideoTag(avid, tagId)
}
func (c *Client) LikeVideoTag(avid, tagId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func HateVideoTag(avid, tagId int) error {
	return std.HateVideoTag(avid, tagId)
}
func (c *Client) HateVideoTag(avid, tagId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func LikeVideoByAvid(avid int, like bool) error {
	return std.LikeVideoByAvid(avid, like)
}
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
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func LikeVideoByBvid(bvid string, like bool) error {
	return std.LikeVideoByBvid(bvid, like)
}
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
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func LikeVideoByShortUrl(shortUrl string, like bool) error {
	return std.LikeVideoByShortUrl(shortUrl, like)
}
func (c *Client) LikeVideoByShortUrl(shortUrl string, like bool) error {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return err
	}
	return LikeVideoByBvid(bvid, like)
}

// CoinVideoByAvid 通过Avid投币视频，multiply为投币数量，上限为2，like为是否附加点赞。返回是否附加点赞成功
func CoinVideoByAvid(avid int, multiply int, like bool) (bool, error) {
	return std.CoinVideoByAvid(avid, multiply, like)
}
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
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func CoinVideoByBvid(bvid string, multiply int, like bool) (bool, error) {
	return std.CoinVideoByBvid(bvid, multiply, like)
}
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
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func CoinVideoByShortUrl(shortUrl string, multiply int, like bool) (bool, error) {
	return std.CoinVideoByShortUrl(shortUrl, multiply, like)
}
func (c *Client) CoinVideoByShortUrl(shortUrl string, multiply int, like bool) (bool, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return false, err
	}
	return CoinVideoByBvid(bvid, multiply, like)
}

// FavourVideoByAvid 通过Avid收藏视频，addMediaIds和delMediaIds为要增加/删除的收藏列表，非必填。返回是否为未关注用户收藏
func FavourVideoByAvid(avid int, addMediaIds, delMediaIds []int) (bool, error) {
	return std.FavourVideoByAvid(avid, addMediaIds, delMediaIds)
}
func (c *Client) FavourVideoByAvid(avid int, addMediaIds, delMediaIds []int) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	var addMediaIdStr, delMediaIdStr []string
	for _, id := range addMediaIds {
		addMediaIdStr = append(addMediaIdStr, strconv.Itoa(id))
	}
	for _, id := range delMediaIds {
		delMediaIdStr = append(delMediaIdStr, strconv.Itoa(id))
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func FavourVideoByBvid(bvid string, addMediaIds, delMediaIds []int) (bool, error) {
	return std.FavourVideoByBvid(bvid, addMediaIds, delMediaIds)
}
func (c *Client) FavourVideoByBvid(bvid string, addMediaIds, delMediaIds []int) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	var addMediaIdStr, delMediaIdStr []string
	for _, id := range addMediaIds {
		addMediaIdStr = append(addMediaIdStr, strconv.Itoa(id))
	}
	for _, id := range delMediaIds {
		delMediaIdStr = append(delMediaIdStr, strconv.Itoa(id))
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func FavourVideoByShortUrl(shortUrl string, addMediaIds, delMediaIds []int) (bool, error) {
	return std.FavourVideoByShortUrl(shortUrl, addMediaIds, delMediaIds)
}
func (c *Client) FavourVideoByShortUrl(shortUrl string, addMediaIds, delMediaIds []int) (bool, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return false, err
	}
	return FavourVideoByBvid(bvid, addMediaIds, delMediaIds)
}

type LikeCoinFavourResult struct {
	Like     bool `json:"like"`     // 是否点赞成功
	Coin     bool `json:"coin"`     // 是否投币成功
	Fav      bool `json:"fav"`      // 是否收藏成功
	Multiply int  `json:"multiply"` // 投币枚数
}

// LikeCoinFavourVideoByAvid 通过Avid一键三连视频
func LikeCoinFavourVideoByAvid(avid int) (*LikeCoinFavourResult, error) {
	return std.LikeCoinFavourVideoByAvid(avid)
}
func (c *Client) LikeCoinFavourVideoByAvid(avid int) (*LikeCoinFavourResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func LikeCoinFavourVideoByBvid(bvid string) (*LikeCoinFavourResult, error) {
	return std.LikeCoinFavourVideoByBvid(bvid)
}
func (c *Client) LikeCoinFavourVideoByBvid(bvid string) (*LikeCoinFavourResult, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func LikeCoinFavourVideoByShortUrl(shortUrl string) (*LikeCoinFavourResult, error) {
	return std.LikeCoinFavourVideoByShortUrl(shortUrl)
}
func (c *Client) LikeCoinFavourVideoByShortUrl(shortUrl string) (*LikeCoinFavourResult, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return LikeCoinFavourVideoByBvid(bvid)
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
func GetVideoOnlineInfoByAvid(avid, cid int) (*VideoOnlineInfo, error) {
	return std.GetVideoOnlineInfoByAvid(avid, cid)
}
func (c *Client) GetVideoOnlineInfoByAvid(avid, cid int) (*VideoOnlineInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func GetVideoOnlineInfoByBvid(bvid string, cid int) (*VideoOnlineInfo, error) {
	return std.GetVideoOnlineInfoByBvid(bvid, cid)
}
func (c *Client) GetVideoOnlineInfoByBvid(bvid string, cid int) (*VideoOnlineInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func GetVideoOnlineInfoByShortUrl(shortUrl string, cid int) (*VideoOnlineInfo, error) {
	return std.GetVideoOnlineInfoByShortUrl(shortUrl, cid)
}
func (c *Client) GetVideoOnlineInfoByShortUrl(shortUrl string, cid int) (*VideoOnlineInfo, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoOnlineInfoByBvid(bvid, cid)
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
func GetVideoPbPInfo(cid int) (*VideoPbPInfo, error) {
	return std.GetVideoPbPInfo(cid)
}
func (c *Client) GetVideoPbPInfo(cid int) (*VideoPbPInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoStatusNumberByAvid(avid int) (*VideoStatusNumber, error) {
	return std.GetVideoStatusNumberByAvid(avid)
}
func (c *Client) GetVideoStatusNumberByAvid(avid int) (*VideoStatusNumber, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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
func GetVideoStatusNumberByBvid(bvid string) (*VideoStatusNumber, error) {
	return std.GetVideoStatusNumberByBvid(bvid)
}
func (c *Client) GetVideoStatusNumberByBvid(bvid string) (*VideoStatusNumber, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
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

// GetVideoStatusNumberByShortUrl 通过短链接获取视频状态数
func GetVideoStatusNumberByShortUrl(shortUrl string) (*VideoStatusNumber, error) {
	return std.GetVideoStatusNumberByShortUrl(shortUrl)
}
func (c *Client) GetVideoStatusNumberByShortUrl(shortUrl string) (*VideoStatusNumber, error) {
	bvid, err := c.GetBvidByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	return GetVideoStatusNumberByBvid(bvid)
}
