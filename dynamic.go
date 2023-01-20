package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"strconv"
	"strings"
)

type SearchDynamicAtResult struct {
	Groups []struct { // 内容分组
		GroupType int        `json:"group_type"` // 2：我的关注，4：其他
		GroupName string     `json:"group_name"` // 分组名字
		Items     []struct { // 用户信息
			Uid                int    `json:"uid"`   // 用户id
			Uname              string `json:"uname"` // 用户昵称
			Face               string `json:"face"`  // 用户头像url
			Fans               int    `json:"fans"`  // 用户粉丝数
			OfficialVerifyType int    `json:"official_verify_type"`
		} `json:"items"`
	} `json:"groups"`
	Gt int `json:"_gt_"` // 固定值0
}

// SearchDynamicAt 根据关键字搜索用户(at别人时的填充列表)
func SearchDynamicAt(uid int, keyword string) (*SearchDynamicAtResult, error) {
	return std.SearchDynamicAt(uid, keyword)
}
func (c *Client) SearchDynamicAt(uid int, keyword string) (*SearchDynamicAtResult, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"uid":     strconv.Itoa(uid),
		"keyword": keyword,
	}).Get("https://api.vc.bilibili.com/dynamic_mix/v1/dynamic_mix/at_search")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "根据关键字搜索用户")
	if err != nil {
		return nil, err
	}
	var ret *SearchDynamicAtResult
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type DynamicRepostDetail struct {
	HasMore int `json:"has_more"` // 是否还有下一页
	Total   int `json:"total"`    // 总计包含
	Items   []struct {
		Desc struct {
			Uid         int   `json:"uid"`
			Type        int   `json:"type"`
			Rid         int64 `json:"rid"`
			Acl         int   `json:"acl"`
			View        int   `json:"view"`
			Repost      int   `json:"repost"`
			Like        int   `json:"like"`
			IsLiked     int   `json:"is_liked"`
			DynamicId   int64 `json:"dynamic_id"`
			Timestamp   int   `json:"timestamp"`
			PreDyId     int64 `json:"pre_dy_id"`
			OrigDyId    int64 `json:"orig_dy_id"`
			OrigType    int   `json:"orig_type"`
			UserProfile struct {
				Info struct {
					Uid     int    `json:"uid"`
					Uname   string `json:"uname"`
					Face    string `json:"face"`
					FaceNft int    `json:"face_nft"`
				} `json:"info"`
				Card struct {
					OfficialVerify struct {
						Type int    `json:"type"`
						Desc string `json:"desc"`
					} `json:"official_verify"`
				} `json:"card"`
				Vip struct {
					VipType    int   `json:"vipType"`
					VipDueDate int64 `json:"vipDueDate"`
					VipStatus  int   `json:"vipStatus"`
					ThemeType  int   `json:"themeType"`
					Label      struct {
						Path        string `json:"path"`
						Text        string `json:"text"`
						LabelTheme  string `json:"label_theme"`
						TextColor   string `json:"text_color"`
						BgStyle     int    `json:"bg_style"`
						BgColor     string `json:"bg_color"`
						BorderColor string `json:"border_color"`
					} `json:"label"`
					AvatarSubscript    int    `json:"avatar_subscript"`
					NicknameColor      string `json:"nickname_color"`
					Role               int    `json:"role"`
					AvatarSubscriptUrl string `json:"avatar_subscript_url"`
				} `json:"vip"`
				Pendant struct {
					Pid               int    `json:"pid"`
					Name              string `json:"name"`
					Image             string `json:"image"`
					Expire            int    `json:"expire"`
					ImageEnhance      string `json:"image_enhance"`
					ImageEnhanceFrame string `json:"image_enhance_frame"`
				} `json:"pendant"`
				Rank      string `json:"rank"`
				Sign      string `json:"sign"`
				LevelInfo struct {
					CurrentLevel int `json:"current_level"`
				} `json:"level_info"`
			} `json:"user_profile"`
			UidType      int    `json:"uid_type"`
			Stype        int    `json:"stype"`
			RType        int    `json:"r_type"`
			InnerId      int    `json:"inner_id"`
			Status       int    `json:"status"`
			DynamicIdStr string `json:"dynamic_id_str"`
			PreDyIdStr   string `json:"pre_dy_id_str"`
			OrigDyIdStr  string `json:"orig_dy_id_str"`
			RidStr       string `json:"rid_str"`
			Origin       struct {
				Uid          int    `json:"uid"`
				Type         int    `json:"type"`
				Rid          int    `json:"rid"`
				Acl          int    `json:"acl"`
				View         int    `json:"view"`
				Repost       int    `json:"repost"`
				Like         int    `json:"like"`
				DynamicId    int64  `json:"dynamic_id"`
				Timestamp    int    `json:"timestamp"`
				PreDyId      int    `json:"pre_dy_id"`
				OrigDyId     int    `json:"orig_dy_id"`
				UidType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerId      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIdStr string `json:"dynamic_id_str"`
				PreDyIdStr   string `json:"pre_dy_id_str"`
				OrigDyIdStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
			} `json:"origin"`
			Previous struct {
				Uid          int    `json:"uid"`
				Type         int    `json:"type"`
				Rid          int64  `json:"rid"`
				Acl          int    `json:"acl"`
				View         int    `json:"view"`
				Repost       int    `json:"repost"`
				Like         int    `json:"like"`
				DynamicId    int64  `json:"dynamic_id"`
				Timestamp    int    `json:"timestamp"`
				PreDyId      int64  `json:"pre_dy_id"`
				OrigDyId     int64  `json:"orig_dy_id"`
				UidType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerId      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIdStr string `json:"dynamic_id_str"`
				PreDyIdStr   string `json:"pre_dy_id_str"`
				OrigDyIdStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
			} `json:"previous"`
		} `json:"desc"`
		Card       string `json:"card"`
		ExtendJson string `json:"extend_json"`
		Display    struct {
			Origin struct {
				EmojiInfo struct {
					EmojiDetails []struct {
						EmojiName string `json:"emoji_name"`
						Id        int    `json:"id"`
						PackageId int    `json:"package_id"`
						State     int    `json:"state"`
						Type      int    `json:"type"`
						Attr      int    `json:"attr"`
						Text      string `json:"text"`
						Url       string `json:"url"`
						Meta      struct {
							Size int `json:"size"`
						} `json:"meta"`
						Mtime int `json:"mtime"`
					} `json:"emoji_details"`
				} `json:"emoji_info"`
				Relation struct {
					Status     int `json:"status"`
					IsFollow   int `json:"is_follow"`
					IsFollowed int `json:"is_followed"`
				} `json:"relation"`
			} `json:"origin"`
			Relation struct {
				Status     int `json:"status"`
				IsFollow   int `json:"is_follow"`
				IsFollowed int `json:"is_followed"`
			} `json:"relation"`
		} `json:"display"`
	} `json:"items"`
	Gt int `json:"_gt_"` // 固定值0
}

// GetDynamicRepostDetail 获取动态转发列表
func GetDynamicRepostDetail(dynamicId, offset int) (*DynamicRepostDetail, error) {
	return std.GetDynamicRepostDetail(dynamicId, offset)
}
func (c *Client) GetDynamicRepostDetail(dynamicId, offset int) (*DynamicRepostDetail, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"dynamic_id": strconv.Itoa(dynamicId),
		"offset":     strconv.Itoa(offset),
	}).Get("https://api.vc.bilibili.com/dynamic_repost/v1/dynamic_repost/repost_detail")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取动态转发列表")
	if err != nil {
		return nil, err
	}
	var ret *DynamicRepostDetail
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type DynamicLikeList struct {
	ItemLikes []struct { // 点赞信息列表主体
		Uid      int    `json:"uid"`
		Time     int    `json:"time"`
		FaceUrl  string `json:"face_url"`
		Uname    string `json:"uname"`
		UserInfo struct {
			Uid            int    `json:"uid"`
			Uname          string `json:"uname"`
			Face           string `json:"face"`
			Rank           string `json:"rank"`
			OfficialVerify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Vip struct {
				VipType    int   `json:"vipType"`
				VipDueDate int64 `json:"vipDueDate"`
				VipStatus  int   `json:"vipStatus"`
				ThemeType  int   `json:"themeType"`
				Label      struct {
					Path        string `json:"path"`
					Text        string `json:"text"`
					LabelTheme  string `json:"label_theme"`
					TextColor   string `json:"text_color"`
					BgStyle     int    `json:"bg_style"`
					BgColor     string `json:"bg_color"`
					BorderColor string `json:"border_color"`
				} `json:"label"`
				AvatarSubscript    int    `json:"avatar_subscript"`
				NicknameColor      string `json:"nickname_color"`
				Role               int    `json:"role"`
				AvatarSubscriptUrl string `json:"avatar_subscript_url"`
			} `json:"vip"`
			Pendant struct {
				Pid               int    `json:"pid"`
				Name              string `json:"name"`
				Image             string `json:"image"`
				Expire            int    `json:"expire"`
				ImageEnhance      string `json:"image_enhance"`
				ImageEnhanceFrame string `json:"image_enhance_frame"`
			} `json:"pendant"`
			Sign      string `json:"sign"`
			LevelInfo struct {
				CurrentLevel int `json:"current_level"`
			} `json:"level_info"`
		} `json:"user_info"`
		Attend int `json:"attend"`
	} `json:"item_likes"`
	HasMore    int `json:"has_more"`    // 是否还有下一页
	TotalCount int `json:"total_count"` // 总计点赞数
	Gt         int `json:"_gt_"`        // 固定值0
}

// GetDynamicLikeList 获取动态点赞列表。offset是非必填项
func GetDynamicLikeList(dynamicId, offset int) (*DynamicLikeList, error) {
	return std.GetDynamicLikeList(dynamicId, offset)
}
func (c *Client) GetDynamicLikeList(dynamicId, offset int) (*DynamicLikeList, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("dynamic_id", strconv.Itoa(dynamicId))
	if offset != 0 {
		r = r.SetQueryParam("offset", strconv.Itoa(offset))
	}
	resp, err := r.Get("https://api.vc.bilibili.com/dynamic_like/v1/dynamic_like/spec_item_likes")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取动态点赞列表")
	if err != nil {
		return nil, err
	}
	var ret *DynamicLikeList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type DynamicLiveUserList struct {
	Count int        `json:"count"` // 直播者数量
	Group string     `json:"group"` // 固定值"default"，作用尚不明确
	Items []struct { // 直播者列表
		Uid   int    `json:"uid"`   // 直播者id
		Uname string `json:"uname"` // 直播者昵称
		Face  string `json:"face"`  // 直播者头像
		Link  string `json:"link"`  // 直播链接
		Title string `json:"title"` // 直播标题
	} `json:"items"`
	Gt int `json:"_gt_"` // 固定值0，作用尚不明确
}

// GetDynamicLiveUserList 获取正在直播的已关注者。size是非必填项
func GetDynamicLiveUserList(size int) (*DynamicLiveUserList, error) {
	return std.GetDynamicLiveUserList(size)
}
func (c *Client) GetDynamicLiveUserList(size int) (*DynamicLiveUserList, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded")
	if size != 0 {
		r = r.SetQueryParam("size", strconv.Itoa(size))
	}
	resp, err := r.Get("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_live_users")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取正在直播的已关注者")
	if err != nil {
		return nil, err
	}
	var ret *DynamicLiveUserList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type DynamicUpList struct {
	ButtonStatement string     `json:"button_statement"` // 固定值空，作用尚不明确
	Items           []struct { // 更新者列表
		UserProfile struct {
			Info struct {
				Uid   int    `json:"uid"`
				Uname string `json:"uname"`
				Face  string `json:"face"`
			} `json:"info"`
			Card struct {
				OfficialVerify struct {
					Type int    `json:"type"`
					Desc string `json:"desc"`
				} `json:"official_verify"`
			} `json:"card"`
			Vip struct {
				VipType       int    `json:"vipType"`
				VipDueDate    int64  `json:"vipDueDate"`
				DueRemark     string `json:"dueRemark"`
				AccessStatus  int    `json:"accessStatus"`
				VipStatus     int    `json:"vipStatus"`
				VipStatusWarn string `json:"vipStatusWarn"`
				ThemeType     int    `json:"themeType"`
				Label         struct {
					Path string `json:"path"`
				} `json:"label"`
			} `json:"vip"`
			Pendant struct {
				Pid          int    `json:"pid"`
				Name         string `json:"name"`
				Image        string `json:"image"`
				Expire       int    `json:"expire"`
				ImageEnhance string `json:"image_enhance"`
			} `json:"pendant"`
			Rank      string `json:"rank"`
			Sign      string `json:"sign"`
			LevelInfo struct {
				CurrentLevel int    `json:"current_level"`
				CurrentMin   int    `json:"current_min"`
				CurrentExp   int    `json:"current_exp"`
				NextExp      string `json:"next_exp"`
			} `json:"level_info"`
		} `json:"user_profile"`
		HasUpdate int `json:"has_update"`
	} `json:"items"`
	Gt int `json:"_gt_"` // 固定值0，作用尚不明确
}

// GetDynamicUpList 获取发布新动态的已关注者。size参数，0：不开启青少年模式，1：开启青少年模式
func GetDynamicUpList(size int) (*DynamicUpList, error) {
	return std.GetDynamicUpList(size)
}
func (c *Client) GetDynamicUpList(size int) (*DynamicUpList, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("size", strconv.Itoa(size)).Get("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/w_dyn_uplist")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取发布新动态的已关注者")
	if err != nil {
		return nil, err
	}
	var ret *DynamicUpList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// RemoveDynamic 删除动态
func RemoveDynamic(dynamicId int) error {
	return std.RemoveDynamic(dynamicId)
}
func (c *Client) RemoveDynamic(dynamicId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"dynamic_id": strconv.Itoa(dynamicId),
		"csrf":       biliJct,
	}).Post("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/rm_dynamic")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "删除动态")
	return err
}

// DynamicCard 动态卡片内容。因为 ActivityInfos 、 Desc 、 Display 等字段会随着此动态类型不同发生一定的变化，
// // 无法统一，因此都转换成了 map[string]interface{} ，请自行解析
type DynamicCard struct {
	ActivityInfos map[string]interface{} `json:"activity_infos"` // 该条动态参与的活动
	Card          string                 `json:"card"`           // 动态详细信息
	Desc          map[string]interface{} `json:"desc"`           // 动态相关信息
	Display       map[string]interface{} `json:"display"`        // 动态部分的可操作项
	ExtendJson    string                 `json:"extend_json"`    // 动态扩展项
}

type DynamicDetail struct {
	Card   *DynamicCard `json:"card"` // 动态卡片内容
	Result int          `json:"result"`
	Gt     int          `json:"_gt_"`
}

// GetDynamicDetail 获取特定动态卡片信息
func GetDynamicDetail(dynamicId int) (*DynamicDetail, error) {
	return std.GetDynamicDetail(dynamicId)
}
func (c *Client) GetDynamicDetail(dynamicId int) (*DynamicDetail, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("dynamic_id", strconv.Itoa(dynamicId)).Get("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取特定动态卡片信息")
	if err != nil {
		return nil, err
	}
	var ret *DynamicDetail
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type DynamicPortal struct {
	MyInfo struct { // 个人关注的一些信息
		Dyns      int      `json:"dyns"`      // 个人动态
		Face      string   `json:"face"`      // 头像url
		FaceNft   int      `json:"face_nft"`  // 含义尚不明确
		Follower  int      `json:"follower"`  // 粉丝数量
		Following int      `json:"following"` // 我的关注
		LevelInfo struct { // 本人等级内容
			CurrentExp   int   `json:"current_exp"`
			CurrentLevel int   `json:"current_level"` // 当前等级，0-6级
			CurrentMin   int   `json:"current_min"`
			LevelUp      int64 `json:"level_up"`
			NextExp      int   `json:"next_exp"`
		} `json:"level_info"`
		Mid      int      `json:"mid"`  // 账户mid
		Name     string   `json:"name"` // 账户名称
		Official struct { // 认证信息
			Desc  string `json:"desc"`  // 认证备注
			Role  int    `json:"role"`  // 认证类型，0：无，1 2 7：个人认证，3 4 5 6：机构认证
			Title string `json:"title"` // 认证信息
			Type  int    `json:"type"`  // 是否认证，-1：无，0：认证
		} `json:"official"`
		SpaceBg string   `json:"space_bg"` // 账户个人中心的背景横幅url
		Vip     struct { // vip信息
			AvatarSubscript    int      `json:"avatar_subscript"`     // 是否显示会员图标，0：不显示，1：显示
			AvatarSubscriptUrl string   `json:"avatar_subscript_url"` // 大会员角标地址
			DueDate            int64    `json:"due_date"`             // 会员过期时间，Unix时间戳（毫秒）
			Label              struct { // 会员标签
				BgColor               string `json:"bg_color"`                  // 会员标签背景颜色，颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
				BgStyle               int    `json:"bg_style"`                  // 固定值1，作用尚不明确
				BorderColor           string `json:"border_color"`              // 会员标签边框颜色，未使用
				ImgLabelUriHans       string `json:"img_label_uri_hans"`        // 固定值空
				ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"` // 大会员牌子图片，简体版
				ImgLabelUriHant       string `json:"img_label_uri_hant"`        // 固定值空
				ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"` // 大会员牌子图片，繁体版
				LabelTheme            string `json:"label_theme"`               // 会员标签，vip，annual_vip，ten_annual_vip，hundred_annual_vip，fools_day_hundred_annual_vip
				Path                  string `json:"path"`                      // 固定值空，作用尚不明确
				Text                  string `json:"text"`                      // 会员类型文案，大会员，年度大会员，十年大会员，百年大会员，最强绿鲤鱼
				TextColor             string `json:"text_color"`                // 会员标签文字颜色
				UseImgLabel           bool   `json:"use_img_label"`             // 固定值true
			} `json:"label"`
			NicknameColor string `json:"nickname_color"`  // 会员昵称颜色，颜色码，一般为#FB7299，曾用于愚人节改变大会员配色
			Role          int    `json:"role"`            // 大会员类型，1：月度大会员，3：年度大会员，7：十年大会员，15：百年大会员
			Status        int    `json:"status"`          // 会员状态，0：无，1：有
			ThemeType     int    `json:"theme_type"`      // 固定值0，作用尚不明确
			TvVipPayType  int    `json:"tv_vip_pay_type"` // 电视大会员支付类型
			TvVipStatus   int    `json:"tv_vip_status"`   // 电视大会员状态，0：未开通
			Type          int    `json:"type"`            // 会员类型，0：无，1：月大会员，2：年度及以上大会员
			VipPayType    int    `json:"vip_pay_type"`    // 支付类型，0：未支付，1：已支付
		} `json:"vip"`
	} `json:"my_info"`
	UpList []struct { // 最近更新的up主列表
		Face            string `json:"face"`       // UP主头像
		HasUpdate       bool   `json:"has_update"` // 最近是否有更新
		IsReserveRecall bool   `json:"is_reserve_recall"`
		Mid             int    `json:"mid"`   // UP主mid
		Uname           string `json:"uname"` // UP主昵称
	} `json:"up_list"`
}

// GetDynamicPortal 获取最近更新UP主列表（其实就是获取自己的动态门户）
func GetDynamicPortal() (*DynamicPortal, error) {
	return std.GetDynamicPortal()
}
func (c *Client) GetDynamicPortal() (*DynamicPortal, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").Get("https://api.bilibili.com/x/polymer/web-dynamic/v1/portal")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取最近更新UP主列表")
	if err != nil {
		return nil, err
	}
	var ret *DynamicPortal
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// UploadDynamicBfs 为图片动态上传图片
func UploadDynamicBfs(fileName string, file io.Reader, category string) (url string, size Size, err error) {
	return std.UploadDynamicBfs(fileName, file, category)
}
func (c *Client) UploadDynamicBfs(fileName string, file io.Reader, category string) (url string, size Size, err error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return "", Size{}, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFileReader("file_up", fileName, file).SetQueryParams(map[string]string{
		"category": category,
		"csrf":     biliJct,
	}).Post("https://api.bilibili.com/x/dynamic/feed/draw/upload_bfs")
	if err != nil {
		return "", Size{}, errors.WithStack(err)
	}
	ret, err := getRespData(resp, "为图片动态上传图片")
	if err != nil {
		return "", Size{}, err
	}
	var data struct {
		ImageUrl    string `json:"image_url"`
		ImageWidth  int    `json:"image_width"`
		ImageHeight int    `json:"image_height"`
	}
	err = json.Unmarshal(ret, &data)
	return data.ImageUrl, Size{Width: data.ImageWidth, Height: data.ImageHeight}, err
}

// CreateDynamic 发表纯文本动态
//
// TODO: extension 字段尚不知如何使用，需自行填写。参考 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/dynamic/publish.md
func CreateDynamic(content, extension string, atUids []int, ctrl []*FormatCtrl) (dynamicId int, err error) {
	return std.CreateDynamic(content, extension, atUids, ctrl)
}
func (c *Client) CreateDynamic(content, extension string, atUids []int, ctrl []*FormatCtrl) (dynamicId int, err error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return 0, errors.New("B站登录过期")
	}
	ctrlJson, err := json.Marshal(ctrl)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	atUidsStr := make([]string, 0, len(atUids))
	for _, atUid := range atUids {
		atUidsStr = append(atUidsStr, strconv.Itoa(atUid))
	}
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"dynamic_id":        "0",
		"type":              "4",
		"rid":               "0",
		"content":           content,
		"up_choose_comment": "0",
		"up_close_comment":  "0",
		"at_uids":           strings.Join(atUidsStr, ","),
		"ctrl":              string(ctrlJson),
		"csrf_token":        biliJct,
		"csrf":              biliJct,
	})
	if len(extension) > 0 {
		r = r.SetQueryParam("extension", extension)
	}
	resp, err := r.Post("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/create")
	if err != nil {
		return 0, errors.WithStack(err)
	}
	ret, err := getRespData(resp, "发表纯文本动态")
	if err != nil {
		return 0, err
	}
	return int(gjson.GetBytes(ret, "dynamic_id").Int()), nil
}

// DynamicList 包含置顶及热门的动态列表
//
// TODO: 因为不清楚 attentions 字段（关注列表）的格式，暂未对此字段进行解析
type DynamicList struct {
	Cards         *DynamicCard `json:"cards"` // 动态列表
	FounderUid    int          `json:"founder_uid,omitempty"`
	HasMore       int          `json:"has_more"` // 当前话题是否有额外的动态，0：无额外动态，1：有额外动态
	IsDrawerTopic int          `json:"is_drawer_topic,omitempty"`
	Offset        string       `json:"offset"` // 接下来获取列表时的偏移值，一般为当前获取的话题列表下最后一个动态id
	Gt            int          `json:"_gt_"`   // 固定值0，作用尚不明确
}

// FetchDynamics 获取包含置顶及热门的动态列表，topicId与topicName任选一个
func FetchDynamics(topicId int, topicName string, sortby, offset int) (*DynamicList, error) {
	return std.FetchDynamics(topicId, topicName, sortby, offset)
}
func (c *Client) FetchDynamics(topicId int, topicName string, sortby, offset int) (*DynamicList, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded")
	if topicId != 0 {
		r = r.SetQueryParam("topic_id", strconv.Itoa(topicId))
	} else {
		r = r.SetQueryParam("topic_name", topicName)
	}
	if sortby != 0 {
		r = r.SetQueryParam("sortby", strconv.Itoa(sortby))
	}
	if offset != 0 {
		r = r.SetQueryParam("offset", strconv.Itoa(offset))
	}
	resp, err := r.Get("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/fetch_dynamics")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取动态列表")
	if err != nil {
		return nil, err
	}
	var ret *DynamicList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

// GetTopicHistory 获取历史动态列表，topicId与topicName任选一个
func GetTopicHistory(topicId int, topicName string, offsetDynamicId int) (*DynamicList, error) {
	return std.GetTopicHistory(topicId, topicName, offsetDynamicId)
}
func (c *Client) GetTopicHistory(topicId int, topicName string, offsetDynamicId int) (*DynamicList, error) {
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParam("offset_dynamic_id", strconv.Itoa(offsetDynamicId))
	if topicId != 0 {
		r = r.SetQueryParam("topic_id", strconv.Itoa(topicId))
	} else {
		r = r.SetQueryParam("topic_name", topicName)
	}
	resp, err := r.Get("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_history")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取历史动态列表")
	if err != nil {
		return nil, err
	}
	var ret *DynamicList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
