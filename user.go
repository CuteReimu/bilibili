package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type OrderType string

const (
	OrderPubDate OrderType = "pubdate"
	OrderClick   OrderType = "click"
	OrderStow    OrderType = "stow"
)

type Video struct {
	Aid          int    `json:"aid"`            // 稿件avid
	Author       string `json:"author"`         // 视频UP主，不一定为目标用户（合作视频）
	Bvid         string `json:"bvid"`           // 稿件bvid
	Comment      int    `json:"comment"`        // 视频评论数
	Copyright    string `json:"copyright"`      // 空，作用尚不明确
	Created      int64  `json:"created"`        // 投稿时间戳
	Description  string `json:"description"`    // 视频简介
	HideClick    bool   `json:"hide_click"`     // 固定值false，作用尚不明确
	IsPay        int    `json:"is_pay"`         // 固定值0，作用尚不明确
	IsUnionVideo int    `json:"is_union_video"` // 是否为合作视频，0：否，1：是
	Length       string `json:"length"`         // 视频长度，MM:SS
	Mid          int    `json:"mid"`            // 视频UP主mid，不一定为目标用户（合作视频）
	Pic          string `json:"pic"`            // 视频封面
	Play         int    `json:"play"`           // 视频播放次数
	Review       int    `json:"review"`         // 固定值0，作用尚不明确
	Subtitle     string `json:"subtitle"`       // 固定值空，作用尚不明确
	Title        string `json:"title"`          // 视频标题
	Typeid       int    `json:"typeid"`         // 视频分区tid
	VideoReview  int    `json:"video_review"`   // 视频弹幕数
}

type GetUserVideosResult struct {
	List struct { // 列表信息
		Tlist map[int]struct { // 投稿视频分区索引
			Count int    `json:"count"` // 投稿至该分区的视频数
			Name  string `json:"name"`  // 该分区名称
			Tid   int    `json:"tid"`   // 该分区tid
		} `json:"tlist"`
		Vlist []Video `json:"vlist"` // 投稿视频列表
	} `json:"list"`
	Page struct { // 页面信息
		Count int `json:"count"` // 总计稿件数
		Pn    int `json:"pn"`    // 当前页码
		Ps    int `json:"ps"`    // 每页项数
	} `json:"page"`
	EpisodicButton struct { // “播放全部“按钮
		Text string `json:"text"` // 按钮文字
		Uri  string `json:"uri"`  // 全部播放页url
	} `json:"episodic_button"`
}

type UserCardResult struct {
	Card struct {
		Mid         string        `json:"mid"`
		Name        string        `json:"name"`
		Approve     bool          `json:"approve"`
		Sex         string        `json:"sex"`
		Rank        string        `json:"rank"`
		Face        string        `json:"face"`
		FaceNft     int           `json:"face_nft"`
		FaceNftType int           `json:"face_nft_type"`
		DisplayRank string        `json:"DisplayRank"`
		Regtime     int           `json:"regtime"`
		Spacesta    int           `json:"spacesta"`
		Birthday    string        `json:"birthday"`
		Place       string        `json:"place"`
		Description string        `json:"description"`
		Article     int           `json:"article"`
		Attentions  []interface{} `json:"attentions"`
		Fans        int           `json:"fans"`
		Friend      int           `json:"friend"`
		Attention   int           `json:"attention"`
		Sign        string        `json:"sign"`
		LevelInfo   struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Pendant struct {
			Pid               int    `json:"pid"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
			NPid              int    `json:"n_pid"`
		} `json:"pendant"`
		Nameplate struct {
			Nid        int    `json:"nid"`
			Name       string `json:"name"`
			Image      string `json:"image"`
			ImageSmall string `json:"image_small"`
			Level      string `json:"level"`
			Condition  string `json:"condition"`
		} `json:"nameplate"`
		Official struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"Official"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"official_verify"`
		Vip struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path                  string `json:"path"`
				Text                  string `json:"text"`
				LabelTheme            string `json:"label_theme"`
				TextColor             string `json:"text_color"`
				BgStyle               int    `json:"bg_style"`
				BgColor               string `json:"bg_color"`
				BorderColor           string `json:"border_color"`
				UseImgLabel           bool   `json:"use_img_label"`
				ImgLabelURIHans       string `json:"img_label_uri_hans"`
				ImgLabelURIHant       string `json:"img_label_uri_hant"`
				ImgLabelURIHansStatic string `json:"img_label_uri_hans_static"`
				ImgLabelURIHantStatic string `json:"img_label_uri_hant_static"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptURL string `json:"avatar_subscript_url"`
			TvVipStatus        int    `json:"tv_vip_status"`
			TvVipPayType       int    `json:"tv_vip_pay_type"`
			TvDueDate          int    `json:"tv_due_date"`
			AvatarIcon         struct {
				IconType     int `json:"icon_type"`
				IconResource struct {
				} `json:"icon_resource"`
			} `json:"avatar_icon"`
			VipType   int `json:"vipType"`
			VipStatus int `json:"vipStatus"`
		} `json:"vip"`
		IsSeniorMember int `json:"is_senior_member"`
	} `json:"card"`
	Space *struct {
		SImg string `json:"s_img,omitempty"`
		LImg string `json:"l_img,omitempty"`
	} `json:"space,omitempty"`
	Following    bool `json:"following"`
	ArchiveCount int  `json:"archive_count"`
	ArticleCount int  `json:"article_count"`
	Follower     int  `json:"follower"`
	LikeNum      int  `json:"like_num"`
}

// GetUserVideos 获取用户投稿视频明细
func (c *Client) GetUserVideos(mid int, order OrderType, tid int, keyword string, pn int, ps int) (*GetUserVideosResult, error) {
	postData := map[string]string{
		"mid": strconv.Itoa(mid),
		"pn":  strconv.Itoa(pn),
		"ps":  strconv.Itoa(ps),
	}
	if len(order) > 0 {
		postData["order"] = string(order)
	}
	if tid != 0 {
		postData["tid"] = strconv.Itoa(tid)
	}
	if len(keyword) > 0 {
		postData["keyword"] = keyword
	}
	resp, err := c.resty.R().SetQueryParams(postData).Get("https://api.bilibili.com/x/space/wbi/arc/search")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取用户视频")
	if err != nil {
		return nil, err
	}
	var ret *GetUserVideosResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetUserCard 获取用户用户名片 免登录
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/user/info.md#%E7%94%A8%E6%88%B7%E5%90%8D%E7%89%87%E4%BF%A1%E6%81%AF
func (c *Client) GetUserCard(mid int, photo bool) (*UserCardResult, error) {
	r := c.resty.R().SetQueryParam("mid", strconv.Itoa(mid)).SetQueryParam("photo", strconv.FormatBool(photo))
	resp, err := r.Get("https://api.bilibili.com/x/web-interface/card")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取用户名片")
	if err != nil {
		return nil, err
	}
	var ret *UserCardResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
