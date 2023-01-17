package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
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

// GetDynamicLikeList 获取动态点赞列表
func GetDynamicLikeList(dynamicId, offset int) (*DynamicLikeList, error) {
	return std.GetDynamicLikeList(dynamicId, offset)
}
func (c *Client) GetDynamicLikeList(dynamicId, offset int) (*DynamicLikeList, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"dynamic_id": strconv.Itoa(dynamicId),
		"offset":     strconv.Itoa(offset),
	}).Get("https://api.vc.bilibili.com/dynamic_like/v1/dynamic_like/spec_item_likes")
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
