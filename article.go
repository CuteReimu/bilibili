package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strconv"
)

type Article struct {
	Id          int      `json:"id"`           // 专栏cvid
	Title       string   `json:"title"`        // 文章标题
	State       int      `json:"state"`        // 固定值0，作用尚不明确
	PublishTime int      `json:"publish_time"` // 发布时间戳（秒）
	Words       int      `json:"words"`        // 文章字数
	ImageUrls   []string `json:"image_urls"`   // 文章封面
	Category    struct { // 文章标签
		Id       int    `json:"id"`
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	} `json:"category"`
	Categories []struct { // 文章标签列表
		Id       int    `json:"id"`
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	} `json:"categories"`
	Summary string   `json:"summary"` // 文章摘要
	Stats   struct { // 文章状态数信息
		View     int `json:"view"`     // 阅读数
		Favorite int `json:"favorite"` // 收藏数
		Like     int `json:"like"`     // 点赞数
		Dislike  int `json:"dislike"`  // 点踩数
		Reply    int `json:"reply"`    // 评论数
		Share    int `json:"share"`    // 分享数
		Coin     int `json:"coin"`     // 投币数
		Dynamic  int `json:"dynamic"`  // 动态转发数
	} `json:"stats"`
	LikeState int `json:"like_state"` // 是否点赞
}

type ArticlesInfo struct {
	List struct { // 文集概览
		Id            int    `json:"id"`             // 文集rlid
		Mid           int    `json:"mid"`            // 文集作者mid
		Name          string `json:"name"`           // 文集名称
		ImageUrl      string `json:"image_url"`      // 文集封面图片url
		UpdateTime    int    `json:"update_time"`    // 文集更新时间戳
		Ctime         int    `json:"ctime"`          // 文集创建时间戳
		PublishTime   int    `json:"publish_time"`   // 文集发布时间戳
		Summary       string `json:"summary"`        // 文集简介
		Words         int    `json:"words"`          // 文集字数
		Read          int    `json:"read"`           // 文集阅读量
		ArticlesCount int    `json:"articles_count"` // 1或3，作用尚不明确
		State         int    `json:"state"`          // 空，作用尚不明确
		Reason        string `json:"reason"`         // 空，作用尚不明确
		ApplyTime     string `json:"apply_time"`     // 空，作用尚不明确
		CheckTime     string `json:"check_time"`     // 空，作用尚不明确
	} `json:"list"`
	Articles []Article `json:"articles"` // 文集内的文章列表
	Author   struct {  // 文集作者信息
		Mid            int            `json:"mid"`  // 作者mid
		Name           string         `json:"name"` // 作者昵称
		Face           string         `json:"face"` // 作者头像url
		Pendant        Pendant        `json:"pendant"`
		OfficialVerify OfficialVerify `json:"official_verify"` // 作者认证信息
		Nameplate      NamePlate      `json:"nameplate"`
		Vip            Vip            `json:"vip"`
	} `json:"author"`
	Last      Article `json:"last"`      // 作用尚不明确
	Attention bool    `json:"attention"` // 是否关注文集作者
}

// GetArticlesInfo 获取文集基本信息
func GetArticlesInfo(id int) (*ArticlesInfo, error) {
	return std.GetArticlesInfo(id)
}
func (c *Client) GetArticlesInfo(id int) (*ArticlesInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("id", strconv.Itoa(id)).Get("https://api.bilibili.com/x/article/list/web/articles")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取文集基本信息")
	if err != nil {
		return nil, err
	}
	var ret *ArticlesInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type ArticleViewInfo struct {
	Like      int      `json:"like"`      // 是否点赞，0：未点赞，1：已点赞
	Attention bool     `json:"attention"` // 是否关注文章作者
	Favorite  bool     `json:"favorite"`  // 是否收藏
	Coin      int      `json:"coin"`      // 为文章投币数
	Stats     struct { // 状态数信息
		View     int `json:"view"`     // 阅读数
		Favorite int `json:"favorite"` // 收藏数
		Like     int `json:"like"`     // 点赞数
		Dislike  int `json:"dislike"`  // 点踩数
		Reply    int `json:"reply"`    // 评论数
		Share    int `json:"share"`    // 分享数
		Coin     int `json:"coin"`     // 投币数
		Dynamic  int `json:"dynamic"`  // 动态转发数
	} `json:"stats"`
	Title           string     `json:"title"`             // 文章标题
	BannerUrl       string     `json:"banner_url"`        // 文章头图url
	Mid             int        `json:"mid"`               // 文章作者mid
	AuthorName      string     `json:"author_name"`       // 文章作者昵称
	IsAuthor        bool       `json:"is_author"`         // 固定值true，作用尚不明确
	ImageUrls       []string   `json:"image_urls"`        // 动态封面图片url
	OriginImageUrls []string   `json:"origin_image_urls"` // 文章封面图片url
	Shareable       bool       `json:"shareable"`         // 固定值true，作用尚不明确
	ShowLaterWatch  bool       `json:"show_later_watch"`  // 固定值true，作用尚不明确
	ShowSmallWindow bool       `json:"show_small_window"` // 固定值true，作用尚不明确
	InList          bool       `json:"in_list"`           // 是否收于文集
	Pre             int        `json:"pre"`               // 上一篇文章cvid
	Next            int        `json:"next"`              // 下一篇文章cvid
	ShareChannels   []struct { // 分享方式列表
		Name         string `json:"name"`          // 分享名称：QQ，QQ空间，微信，朋友圈，微博
		Picture      string `json:"picture"`       // 分享图片url
		ShareChannel string `json:"share_channel"` // 分享代号：QQ，QZONE，WEIXIN，WEIXIN_MOMENT，SINA
	} `json:"share_channels"`
}

// GetArticleViewInfo 获取专栏文章基本信息
func GetArticleViewInfo(id int) (*ArticleViewInfo, error) {
	return std.GetArticleViewInfo(id)
}
func (c *Client) GetArticleViewInfo(id int) (*ArticleViewInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("id", strconv.Itoa(id)).Get("https://api.bilibili.com/x/article/viewinfo")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取专栏文章基本信息")
	if err != nil {
		return nil, err
	}
	var ret *ArticleViewInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// LikeArticle 点赞文章，like为false表示取消点赞
func LikeArticle(id int, like bool) error {
	return std.LikeArticle(id, like)
}
func (c *Client) LikeArticle(id int, like bool) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	var typeNum string
	if like {
		typeNum = "1"
	} else {
		typeNum = "2"
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"id":   strconv.Itoa(id),
		"type": typeNum,
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/article/like")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "点赞文章")
	return err
}

// CoinArticle 投币文章，id为文章cvid，upid为作者mid，mutiply为投币数量。返回的bool值为是否附加点赞成功，若已赞过则附加点赞失败
func CoinArticle(id, upid, multiply int) (bool, error) {
	return std.CoinArticle(id, upid, multiply)
}
func (c *Client) CoinArticle(id, upid, multiply int) (bool, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return false, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"aid":      strconv.Itoa(id),
		"upid":     strconv.Itoa(upid),
		"multiply": strconv.Itoa(multiply),
		"avtype":   "2",
		"csrf":     biliJct,
	}).Post("https://api.bilibili.com/x/web-interface/coin/add")
	if err != nil {
		return false, errors.WithStack(err)
	}
	data, err := getRespData(resp, "投币文章")
	if err != nil {
		return false, err
	}
	return gjson.GetBytes(data, "like").Bool(), nil
}

// FavourArticle 收藏文章
func FavourArticle(id int) error {
	return std.FavourArticle(id)
}
func (c *Client) FavourArticle(id int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"id":   strconv.Itoa(id),
		"csrf": biliJct,
	}).Post("https://api.bilibili.com/x/article/favorites/add")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "收藏文章")
	return err
}

type UserArticleList struct {
	Articles []struct { // 专栏文章信息列表
		Id       int      `json:"id"` // 专栏文章id
		Category struct { // 分类
			Id       int    `json:"id"`        // 	分类id
			ParentId int    `json:"parent_id"` // 父级分类id
			Name     string `json:"name"`      // 分类名称
		} `json:"category"`
		Categories []struct { // 分类
			Id       int    `json:"id"`        // 分类id
			ParentId int    `json:"parent_id"` // 父级分类id
			Name     string `json:"name"`      // 分类名称
		} `json:"categories"`
		Title      string   `json:"title"`      // 标题
		Summary    string   `json:"summary"`    // 摘要
		BannerUrl  string   `json:"banner_url"` // 封面图
		TemplateId int      `json:"template_id"`
		State      int      `json:"state"`
		Author     struct { // UP主信息
			Mid     int      `json:"mid"`  // 用户uid
			Name    string   `json:"name"` // 用户名
			Face    string   `json:"face"` // 头像
			Pendant struct { // 头像框信息
				Pid    int    `json:"pid"`    // 头像框id
				Name   string `json:"name"`   // 头像框名称
				Image  string `json:"image"`  // 头像框图片url
				Expire int    `json:"expire"` // 过期时间
			} `json:"pendant"`
			OfficialVerify struct { // 账号认证信息
				Type int    `json:"type"` // 是否认证，-1：无，0：个人认证，1：机构认证
				Desc string `json:"desc"` // 认证备注
			} `json:"official_verify"`
			Nameplate struct { // 成就勋章信息
				Nid        int    `json:"nid"`         // 勋章id
				Name       string `json:"name"`        // 勋章名称
				Image      string `json:"image"`       // 勋章图标
				ImageSmall string `json:"image_small"` // 勋章图标（小）
				Level      string `json:"level"`       // 勋章等级
				Condition  string `json:"condition"`   // 获取条件
			} `json:"nameplate"`
			Vip struct { // 大会员信息
				Type       int      `json:"type"`         // 大会员类型，0：无，1：月大会员，2：年度及以上大会员
				Status     int      `json:"status"`       // 大会员状态，0：无，1：有
				DueDate    int      `json:"due_date"`     // 大会员过期时间时间戳，单位：毫秒
				VipPayType int      `json:"vip_pay_type"` // 支付类型
				ThemeType  int      `json:"theme_type"`   // 固定值0
				Label      struct { // 大会员标签
					Path       string `json:"path"`        // 空串
					Text       string `json:"text"`        // 会员类型文案，大会员，年度大会员，十年大会员，百年大会员，最强绿鲤鱼
					LabelTheme string `json:"label_theme"` // 会员标签，vip，annual_vip，ten_annual_vip，hundred_annual_vip，fools_day_hundred_annual_vip
				} `json:"label"`
				AvatarSubscript int    `json:"avatar_subscript"` // 是否显示大会员图标，0：不显示，1：显示
				NicknameColor   string `json:"nickname_color"`   // 大会员昵称颜色
			} `json:"vip"`
		} `json:"author"`
		Reprint     int      `json:"reprint"`
		ImageUrls   []string `json:"image_urls"`
		PublishTime int      `json:"publish_time"` // 发布时间戳，单位：秒
		Ctime       int      `json:"ctime"`        // 提交时间戳，单位：秒
		Stats       struct { // 专栏文章数据统计
			View     int `json:"view"`     // 浏览数
			Favorite int `json:"favorite"` // 收藏数
			Like     int `json:"like"`     // 点赞数
			Dislike  int `json:"dislike"`  // 点踩数
			Reply    int `json:"reply"`    // 回复数
			Share    int `json:"share"`    // 转发数
			Coin     int `json:"coin"`     // 投币数
			Dynamic  int `json:"dynamic"`
		} `json:"stats"`
		Words           int      `json:"words"`
		Dynamic         string   `json:"dynamic"` // 粉丝动态文案
		OriginImageUrls []string `json:"origin_image_urls"`
		IsLike          bool     `json:"is_like"`
		Media           struct {
			Score    int    `json:"score"`
			MediaId  int    `json:"media_id"`
			Title    string `json:"title"`
			Cover    string `json:"cover"`
			Area     string `json:"area"`
			TypeId   int    `json:"type_id"`
			TypeName string `json:"type_name"`
			Spoiler  int    `json:"spoiler"`
		} `json:"media"`
		ApplyTime string     `json:"apply_time"`
		CheckTime string     `json:"check_time"`
		Original  int        `json:"original"`
		ActId     int        `json:"act_id"`
		CoverAvid int        `json:"cover_avid"`
		Type      int        `json:"type"`
		Tags      []struct { // 标签
			Tid  int    `json:"tid"`  // 标签id
			Name string `json:"name"` // 标签名称
		} `json:"tags,omitempty"`
	} `json:"articles"`
	Pn    int `json:"pn"`    // 本次请求分页页数
	Ps    int `json:"ps"`    // 本次请求分页大小
	Count int `json:"count"` // 专栏文章总数
}

// GetUserArticleList 获取用户专栏文章列表。sort可选值"publish_time"，"view"，"fav"，不填默认"publish_time"。
func GetUserArticleList(mid, pn, ps int, sort string) (*UserArticleList, error) {
	return std.GetUserArticleList(mid, pn, ps, sort)
}
func (c *Client) GetUserArticleList(mid, pn, ps int, sort string) (*UserArticleList, error) {
	if len(sort) == 0 {
		sort = "publish_time"
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"mid":  strconv.Itoa(mid),
		"pn":   strconv.Itoa(pn),
		"ps":   strconv.Itoa(ps),
		"sort": sort,
	}).Get("https://api.bilibili.com/x/space/article")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取用户专栏文章列表")
	if err != nil {
		return nil, err
	}
	var ret *UserArticleList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type UserArticlesList struct {
	Lists []struct { // 文集信息列表
		Id            int    `json:"id"`           // 文集id
		Mid           int    `json:"mid"`          // 作者uid
		Name          string `json:"name"`         // 文集名称
		ImageUrl      string `json:"image_url"`    // 封面
		UpdateTime    int    `json:"update_time"`  // 最后更新时间戳，单位：秒
		Ctime         int    `json:"ctime"`        // 创建时间戳，单位：秒
		PublishTime   int    `json:"publish_time"` // 发布时间戳，单位：秒
		Summary       string `json:"summary"`
		Words         int    `json:"words"`          // 总字数
		Read          int    `json:"read"`           // 阅读量
		ArticlesCount int    `json:"articles_count"` // 包含文章数
		State         int    `json:"state"`          // 固定值1
		Reason        string `json:"reason"`
		ApplyTime     string `json:"apply_time"`
		CheckTime     string `json:"check_time"`
	} `json:"lists"`
	Total int `json:"total"` // 文集总数
}

// GetUserArticlesList 获取用户专栏文集列表。sort可选值，0：最近更新，1：最多阅读。
func GetUserArticlesList(mid, sort int) (*UserArticlesList, error) {
	return std.GetUserArticlesList(mid, sort)
}
func (c *Client) GetUserArticlesList(mid, sort int) (*UserArticlesList, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"mid":  strconv.Itoa(mid),
		"sort": strconv.Itoa(sort),
	}).Get("https://api.bilibili.com/x/article/up/lists")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取用户专栏文集列表")
	if err != nil {
		return nil, err
	}
	var ret *UserArticlesList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
