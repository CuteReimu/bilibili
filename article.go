package bilibili

import "github.com/go-resty/resty/v2"

type GetArticlesInfoParam struct {
	Id int `json:"id"` // 文集rlid
}

type Articles struct {
	Id            int    `json:"id"`             // 文集rlid
	Mid           int    `json:"mid"`            // 文集作者mid
	Name          string `json:"name"`           // 文集名称
	ImageUrl      string `json:"image_url"`      // 文集封面图片url
	UpdateTime    int    `json:"update_time"`    // 文集更新时间。时间戳
	Ctime         int    `json:"ctime"`          // 文集创建时间。时间戳
	PublishTime   int    `json:"publish_time"`   // 文集发布时间。时间戳
	Summary       string `json:"summary"`        // 文集简介
	Words         int    `json:"words"`          // 文集字数
	Read          int    `json:"read"`           // 文集阅读量
	ArticlesCount int    `json:"articles_count"` // 文集内文章数量
	State         int    `json:"state"`          // 1或3。作用尚不明确
	Reason        string `json:"reason"`         // 空。作用尚不明确
	ApplyTime     string `json:"apply_time"`     // 空。作用尚不明确
	CheckTime     string `json:"check_time"`     // 空。作用尚不明确
}

type ArticleStats struct {
	View     int `json:"view"`     // 阅读数
	Favorite int `json:"favorite"` // 收藏数
	Like     int `json:"like"`     // 点赞数
	Dislike  int `json:"dislike"`  // 点踩数
	Reply    int `json:"reply"`    // 评论数
	Share    int `json:"share"`    // 分享数
	Coin     int `json:"coin"`     // 投币数
	Dynamic  int `json:"dynamic"`  // 动态转发数
}

type Category struct {
	Id       int    `json:"id"`        // 分类id
	ParentId int    `json:"parent_id"` // 父级分类id
	Name     string `json:"name"`      // 分类名称
}

type Tag struct {
	Tid  int    `json:"tid"`  // 标签id
	Name string `json:"name"` // 标签名称
}

type Media struct {
	Score    int    `json:"score"`     // 0
	MediaId  int    `json:"media_id"`  // 0
	Title    string `json:"title"`     // 空串
	Cover    string `json:"cover"`     // 空串
	Area     string `json:"area"`      // 空串
	TypeId   int    `json:"type_id"`   // 0
	TypeName string `json:"type_name"` // 空串
	Spoiler  int    `json:"spoiler"`   // 0
}

type Article struct {
	Id              int          `json:"id"`         // 专栏文章id
	Category        Category     `json:"category"`   // 分类
	Categories      []Category   `json:"categories"` // 分类
	Title           string       `json:"title"`      // 标题
	Summary         string       `json:"summary"`    // 摘要
	BannerUrl       string       `json:"banner_url"` // 封面图
	TemplateId      int          `json:"template_id"`
	State           int          `json:"state"`
	Author          *Author      `json:"author"` // UP主信息
	Reprint         int          `json:"reprint"`
	ImageUrls       []string     `json:"image_urls"`
	PublishTime     int          `json:"publish_time"` // 发布时间戳。单位：秒
	Ctime           int          `json:"ctime"`        // 提交时间戳。单位：秒
	Stats           ArticleStats `json:"stats"`        // 专栏文章数据统计
	Tags            []Tag        `json:"tags"`         // 标签
	Words           int          `json:"words"`
	Dynamic         string       `json:"dynamic"` // 粉丝动态文案
	OriginImageUrls []string     `json:"origin_image_urls"`
	IsLike          bool         `json:"is_like"`
	Media           *Media       `json:"media"`
	ApplyTime       string       `json:"apply_time"` // 空串
	CheckTime       string       `json:"check_time"` // 空串
	Original        int          `json:"original"`
	ActId           int          `json:"act_id"`
	CoverAvid       int          `json:"cover_avid"`
	Type            int          `json:"type"`
	LikeState       int          `json:"like_state"` // 是否点赞。0：未点赞。1：已点赞。需要登录(Cookie) 。未登录为0
}

type Author struct {
	Mid            int            `json:"mid"`             // 作者mid
	Name           string         `json:"name"`            // 作者昵称
	Face           string         `json:"face"`            // 作者头像url
	OfficialVerify OfficialVerify `json:"official_verify"` // 作者认证信息
	Nameplate      Nameplate      `json:"nameplate"`       // 作者勋章
	Vip            Vip            `json:"vip"`             // 作者大会员状态
}

type ArticlesInfo struct {
	List      Articles  `json:"list"`      // 文集概览
	Articles  []Article `json:"articles"`  // 文集内的文章列表
	Author    Author    `json:"author"`    // 文集作者信息
	Last      Article   `json:"last"`      // -。作用尚不明确。结构与data.articles[]中相似
	Attention bool      `json:"attention"` // 是否关注文集作者。false：未关注。true：已关注。需要登录(Cookie) 。未登录为false
}

// GetArticlesInfo 获取文集基本信息
func (c *Client) GetArticlesInfo(param GetArticlesInfoParam) (*ArticlesInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/article/list/web/articles"
	)
	return execute[*ArticlesInfo](c, method, url, param)
}

type ShareChannel struct {
	Name         string `json:"name"`          // 分享名称
	Picture      string `json:"picture"`       // 分享图片url
	ShareChannel string `json:"share_channel"` // 分享代号
}

type ArticleInfo struct {
	Like            int            `json:"like"`              // 是否点赞。0：未点赞。1：已点赞。需要登录(Cookie) 。未登录为0
	Attention       bool           `json:"attention"`         // 是否关注文章作者。false：未关注。true：已关注。需要登录(Cookie) 。未登录为false
	Favorite        bool           `json:"favorite"`          // 是否收藏。false：未收藏。true：已收藏。需要登录(Cookie) 。未登录为false
	Coin            int            `json:"coin"`              // 为文章投币数
	Stats           ArticleStats   `json:"stats"`             // 状态数信息
	Title           string         `json:"title"`             // 文章标题
	BannerUrl       string         `json:"banner_url"`        // 文章头图url
	Mid             int            `json:"mid"`               // 文章作者mid
	AuthorName      string         `json:"author_name"`       // 文章作者昵称
	IsAuthor        bool           `json:"is_author"`         // true。作用尚不明确
	ImageUrls       []string       `json:"image_urls"`        // 动态封面
	OriginImageUrls []string       `json:"origin_image_urls"` // 封面图片
	Shareable       bool           `json:"shareable"`         // true。作用尚不明确
	ShowLaterWatch  bool           `json:"show_later_watch"`  // true。作用尚不明确
	ShowSmallWindow bool           `json:"show_small_window"` // true。作用尚不明确
	InList          bool           `json:"in_list"`           // 是否收于文集。false：否。true：是
	Pre             int            `json:"pre"`               // 上一篇文章cvid。无为0
	Next            int            `json:"next"`              // 下一篇文章cvid。无为0
	ShareChannels   []ShareChannel `json:"share_channels"`    // 分享方式列表
	Type            int            `json:"type"`              // 文章类别。0：文章。2：笔记
}

type GetArticleInfoParam struct {
	Id int `json:"id"` // 专栏cvid
}

// GetArticleInfo 获取专栏文章基本信息
func (c *Client) GetArticleInfo(param GetArticleInfoParam) (*ArticleInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/article/viewinfo"
	)
	return execute[*ArticleInfo](c, method, url, param)
}

type LikeArticleParam struct {
	Id   int `json:"id"`   // 文章cvid
	Type int `json:"type"` // 操作方式。1：点赞。2：取消赞
}

// LikeArticle 点赞文章
func (c *Client) LikeArticle(param LikeArticleParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/article/like"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type CoinArticleParam struct {
	Aid      int `json:"aid"`      // 文章cvid
	Upid     int `json:"upid"`     // 文章作者mid
	Multiply int `json:"multiply"` // 投币数量。上限为2
	Avtype   int `json:"avtype"`   // 2。必须为2
}

type CoinArticleResult struct {
	Like bool `json:"like"` // 是否点赞成功。true：成功。false：失败。已赞过则附加点赞失败
}

// CoinArticle 投币文章
func (c *Client) CoinArticle(param CoinArticleParam) (*CoinArticleResult, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/web-interface/coin/add"
	)
	return execute[*CoinArticleResult](c, method, url, param, fillCsrf(c))
}

type FavoritesArticleParam struct {
	Id int `json:"id"` // 文章cvid
}

// FavoritesArticle 收藏文章
func (c *Client) FavoritesArticle(param FavoritesArticleParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/article/favorites/add"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type GetUserArticleListParam struct {
	Mid  int    `json:"mid"`            // 用户uid
	Pn   int    `json:"pn,omitempty"`   // 默认：1
	Ps   int    `json:"ps,omitempty"`   // 默认：30。范围：[1,30]
	Sort string `json:"sort,omitempty"` // publish_time：最新发布。view：最多阅读。fav：最多收藏。默认：publish_time
}

type UserArticleList struct {
	Articles []Article `json:"articles"` // 专栏文章信息列表
	Pn       int       `json:"pn"`       // 本次请求分页页数
	Ps       int       `json:"ps"`       // 本次请求分页大小
	Count    int       `json:"count"`    // 专栏文章总数
}

// GetUserArticleList 获取用户专栏文章列表
func (c *Client) GetUserArticleList(param GetUserArticleListParam) (*UserArticleList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/space/article"
	)
	return execute[*UserArticleList](c, method, url, param, fillCsrf(c))
}

type GetUserArticlesListParam struct {
	Mid      int    `json:"mid"`            // 用户uid
	Sort     int    `json:"sort,omitempty"` // 排序方式。0：最近更新。1：最多阅读
	Jsonp    string `json:"jsonp,omitempty"`
	Callback string `json:"callback,omitempty"`
}

type UserArticlesList struct {
	Lists []ArticlesInfo `json:"lists"` // 文集信息列表
	Total int            `json:"total"` // 文集总数
}

// GetUserArticlesList 获取用户专栏文集列表
func (c *Client) GetUserArticlesList(param GetUserArticlesListParam) (*UserArticlesList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/article/up/lists"
	)
	return execute[*UserArticlesList](c, method, url, param, fillCsrf(c))
}
