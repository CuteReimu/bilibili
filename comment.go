package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type Comment struct {
	Rpid      int64  `json:"rpid"`       // 评论 rpid
	Oid       int    `json:"oid"`        // 评论区对象 id
	Type      int    `json:"type"`       // 评论区类型代码
	Mid       int    `json:"mid"`        // 发送者 mid
	Root      int    `json:"root"`       // 根评论 rpid，若为一级评论则为 0，大于一级评论则为根评论 id
	Parent    int    `json:"parent"`     // 回复父评论 rpid，若为一级评论则为 0，若为二级评论则为根评论 rpid，大于二级评论为上一级评论 rpid
	Dialog    int    `json:"dialog"`     // 回复对方 rpid，若为一级评论则为 0，若为二级评论则为该评论 rpid，大于二级评论为上一级评论 rpid
	Count     int    `json:"count"`      // 二级评论条数
	Rcount    int    `json:"rcount"`     // 回复评论条数
	Floor     int    `json:"floor"`      // 评论楼层号
	State     int    `json:"state"`      // 作用尚不明确
	Fansgrade int    `json:"fansgrade"`  // 是否具有粉丝标签，0：无，1：有
	Attr      int    `json:"attr"`       // 作用尚不明确
	Ctime     int64  `json:"ctime"`      // 评论发送时间戳
	RpidStr   string `json:"rpid_str"`   // 评论 rpid 字符串格式
	RootStr   string `json:"root_str"`   // 根评论 rpid 字符串格式
	ParentStr string `json:"parent_str"` // 回复父评论 rpid 字符串格式
	Like      int    `json:"like"`       // 评论获赞数
	Action    int    `json:"action"`     // 当前用户操作状态，需要登录(Cookie 或 APP)，否则恒为 0，0：无，1：已点赞，2：已点踩
	Member    struct {
		Mid         int64  `json:"mid"`    // 发送者 mid
		Uname       string `json:"uname"`  // 发送者昵称
		Sex         string `json:"sex"`    // 发送者性别
		Sign        string `json:"sign"`   // 发送者签名
		Avatar      string `json:"avatar"` // 发送者头像 url
		Rank        string `json:"rank"`
		DisplayRank string `json:"DisplayRank"`
		LevelInfo   struct {
			CurrentLevel int `json:"current_level"` // 用户等级
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"` // 发送者等级
		Pendant *struct {
			Pid               int    `json:"pid"`   // 头像框 id
			Name              string `json:"name"`  // 头像框名称
			Image             string `json:"image"` // 头像框图片 url
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"` // 发送者头像框信息
		Nameplate *struct {
			Nid        int    `json:"nid"`         // 勋章 id
			Name       string `json:"name"`        // 勋章名称
			Image      string `json:"image"`       // 挂件图片 url 正常
			ImageSmall string `json:"image_small"` // 勋章图片 url 小
			Level      string `json:"level"`       // 勋章等级
			Condition  string `json:"condition"`   // 勋章条件
		} `json:"nameplate"` // 发送者勋章信息
		OfficialVerify *struct {
			Type int    `json:"type"` // 是否认证，-1：无，0：个人认证，1：机构认证
			Desc string `json:"desc"` // 认证信息，无为空
		} `json:"official_verify"` // 发送者认证信息
		Vip *struct {
			VipType       int    `json:"vipType"`    // 大会员类型，0：无，1：月会员，2：年以上会员
			VipDueDate    int64  `json:"vipDueDate"` // 大会员到期时间，毫秒时间戳
			DueRemark     string `json:"dueRemark"`
			AccessStatus  int    `json:"accessStatus"`
			VipStatus     int    `json:"vipStatus"` // 大会员状态，0：无，1：有
			VipStatusWarn string `json:"vipStatusWarn"`
			ThemeType     int    `json:"theme_type"` // 会员样式 id
			Label         struct {
				Path        string `json:"path"`
				Text        string `json:"text"`        // 会员类型文案
				LabelTheme  string `json:"label_theme"` // 会员类型，vip：大会员，annual_vip：年度大会员，ten_annual_vip：十年大会员，hundred_annual_vip：百年大会员
				TextColor   string `json:"text_color"`
				BgStyle     int    `json:"bg_style"`
				BgColor     string `json:"bg_color"`
				BorderColor string `json:"border_color"`
			} `json:"label"` // 会员铭牌样式
			AvatarSubscript    int    `json:"avatar_subscript"`
			AvatarSubscriptUrl string `json:"avatar_subscript_url"`
			NicknameColor      string `json:"nickname_color"` // 昵称颜色
		} `json:"vip"` // 发送者会员信息
		FansDetail *struct {
			Uid          int    `json:"uid"`        // 用户 mid
			MedalId      int    `json:"medal_id"`   // 粉丝标签 id
			MedalName    string `json:"medal_name"` // 粉丝标签名
			Score        int    `json:"score"`
			Level        int    `json:"level"` // 当前标签等级
			Intimacy     int    `json:"intimacy"`
			MasterStatus int    `json:"master_status"`
			IsReceive    int    `json:"is_receive"`
		} `json:"fans_detail"` // 发送者粉丝标签
		Following   int `json:"following"`   // 是否关注该用户，需要登录(Cookie或APP)，否则恒为0，0：未关注，1：已关注
		IsFollowed  int `json:"is_followed"` // 是否被该用户关注，需要登录(Cookie或APP)，否则恒为0，0：未关注，1：已关注
		UserSailing *struct {
			Pendant *struct {
				Id                int    `json:"id"`    // 头像框 id
				Name              string `json:"name"`  // 头像框名称
				Image             string `json:"image"` // 头像框图片 url
				JumpUrl           string `json:"jump_url"`
				Type              string `json:"type"` // 装扮类型，suit：一般装扮，vip_suit：vip装扮
				ImageEnhance      string `json:"image_enhance"`
				ImageEnhanceFrame string `json:"image_enhance_frame"`
			} `json:"pendant"` // 头像框信息
			Cardbg *struct {
				Id      int    `json:"id"`       // 评论条目装扮 id
				Name    string `json:"name"`     // 评论条目装扮名称
				Image   string `json:"image"`    // 评论条目装扮图片 url
				JumpUrl string `json:"jump_url"` // 评论条目装扮商城页面 url
				Fan     struct {
					IsFan   int    `json:"is_fan"`   // 是否为粉丝专属装扮，0：否，1：是
					Number  int    `json:"number"`   // 粉丝专属编号
					Color   string `json:"color"`    // 数字颜色
					Name    string `json:"name"`     // 装扮名称
					NumDesc string `json:"num_desc"` // 粉丝专属编号，字串格式
				} `json:"fan"` // 粉丝专属信息
				Type string `json:"type"` // 装扮类型，suit：一般装扮，vip_suit：vip装扮
			} `json:"cardbg"` // 评论卡片装扮
			CardbgWithFocus interface{} `json:"cardbg_with_focus"` // null
		} `json:"user_sailing"` // 发送者评论条目装扮信息
		IsContractor bool   `json:"is_contractor"` // 是否为合作用户
		ContractDesc string `json:"contract_desc"` // 合作用户说明
	} `json:"member"`
	Content struct {
		Message string        `json:"message"` // 评论内容
		Plat    int           `json:"plat"`    // 评论发送端，1：web端，2：安卓客户端，3：ios客户端，4：wp客户端
		Device  string        `json:"device"`  // 评论发送平台设备
		Members []interface{} `json:"members"` // at到的用户信息
		Emote   map[string]struct {
			Id        int    `json:"id"`         // 表情 id
			PackageId int    `json:"package_id"` // 表情包 id
			State     int    `json:"state"`
			Type      int    `json:"type"` // 表情类型，1：免费，2：会员专属，3：购买所得，4：颜文字
			Attr      int    `json:"attr"`
			Text      string `json:"text"` // 表情转义符
			Url       string `json:"url"`  // 表情图片 url
			Meta      struct {
				Size  int    `json:"size"`  // 表情尺寸信息，1：小，2：大
				Alias string `json:"alias"` // 简写名，无则无此项
			} `json:"meta"`
			Mtime     int64  `json:"mtime"`      // 表情创建时间，时间戳
			JumpTitle string `json:"jump_title"` // 表情名称
		} `json:"emote"` // 需要渲染的表情转义，评论内容无表情则无此项
		JumpUrl map[string]struct {
			Title          string `json:"title"` // 标题
			State          int    `json:"state"` // 图标 url
			PrefixIcon     string `json:"prefixIcon"`
			AppUrlSchema   string `json:"appUrlSchema"`
			AppName        string `json:"appName"`
			AppPackageName string `json:"appPackageName"`
			ClickReport    string `json:"clickReport"` // 上报 id
		} `json:"jump_url"` // 需要高亮的超链转义
		MaxLine  int `json:"max_line"` // 收起最大行数
		Pictures []struct {
			ImgSrc    string `json:"img_src"`    // 图片地址
			ImgWidth  int    `json:"img_width"`  // 图片宽度
			ImgHeight int    `json:"img_height"` // 图片高度
			ImgSize   int    `json:"img_size"`   // 图片大小，单位KB
		} `json:"pictures"` // 评论图片数组
	} `json:"content"`
	Replies []Comment `json:"replies"` // 评论回复条目预览，仅嵌套一层
	Assist  int       `json:"assist"`  // 作用尚不明确
	Folder  struct {
		HasFolded bool   `json:"has_folded"` // 是否有被折叠的二级评论
		IsFolded  bool   `json:"is_folded"`  // 评论是否被折叠
		Rule      string `json:"rule"`       // 相关规则页面 url
	} `json:"folder"` // 折叠信息
	UpAction struct {
		Like  bool `json:"like"`  // 是否UP主觉得很赞
		Reply bool `json:"reply"` // 是否被UP主回复
	} `json:"up_action"` // 评论 UP 主操作信息
	ShowFollow bool `json:"show_follow"` // 作用尚不明确
	Invisible  bool `json:"invisible"`   // 评论是否被隐藏
	CardLabel  struct {
		Rpid             int64  `json:"rpid"`              // 评论 rpid
		TextContent      string `json:"text_content"`      // 标签文本，已知有"妙评"
		TextColorDay     string `json:"text_color_day"`    // 日间文本颜色，十六进制颜色值
		TextColorNight   string `json:"text_color_night"`  // 夜间文本颜色，十六进制颜色值
		LabelColorDay    string `json:"label_color_day"`   // 日间标签颜色，十六进制颜色值
		LabelColorNight  string `json:"label_color_night"` // 夜间标签颜色，十六进制颜色值
		Image            string `json:"image"`             // 作用不明
		Type             string `json:"type"`              // 1，作用不明
		Background       string `json:"background"`        // 背景图片 url
		BackgroundWidth  int    `json:"background_width"`  // 背景图片宽度
		BackgroundHeight int    `json:"background_height"` // 背景图片高度
		JumpUrl          string `json:"jump_url"`          // 跳转链接
		Effect           int    `json:"effect"`            // 0，作用不明，可能用于控制动画
		EffectStartTime  int    `json:"effect_start_time"` // 0，作用不明，可能用于控制动画
	} `json:"card_label"` // 右上角卡片标签信息
	ReplyControl struct {
		SubReplyEntryText string `json:"sub_reply_entry_text"` // 回复提示，"共 xx 条回复"
		SubReplyTitleText string `json:"sub_reply_title_text"` // 回复提示，"相关回复共有 xx 条"
		TimeDesc          string `json:"time_desc"`            // 时间提示，"xx 天/小时 前发布"
		// 评论者发送评论时的IP地址属地，仅对2022-07-25 11:00及以后发布的评论有
		Location string `json:"location"`
	} `json:"reply_control"`
}

// HotReply 视频热评信息
type HotReply struct {
	Page struct { // 页面信息
		Acount int `json:"acount"` // 总评论数
		Count  int `json:"count"`  // 热评数
		Num    int `json:"num"`    // 当前页码
		Size   int `json:"size"`   // 每页项数
	} `json:"page"`
	Replies []Comment `json:"replies"` // 热评列表
}

// GetVideoComment 获取视频评论，sort：0按时间、1按点赞数、2按回复数
//
// oidType：见 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/comment/readme.md
func GetVideoComment(oidType, oid, sort, root int) (*HotReply, error) {
	return std.GetVideoComment(oidType, oid, sort, root)
}

// GetVideoComment 用于获取视频评论 及二级评论
func (c *Client) GetVideoComment(oidType, oid, sort, root int) (*HotReply, error) {
	// 创建一个包含 type、oid 和 sort 参数的映射
	params := map[string]string{
		"type": strconv.Itoa(oidType),
		"oid":  strconv.Itoa(oid),
		"sort": strconv.Itoa(sort),
	}

	// 设置请求 URL
	url := "https://api.bilibili.com/x/v2/reply"
	if root != 0 {
		// 如果 root 不为 0,则需要访问 /x/v2/reply/reply ,并添加 root 参数
		url = "https://api.bilibili.com/x/v2/reply/reply"
		params["root"] = strconv.Itoa(root)
	}

	// 发送 HTTP GET 请求
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(params).Get(url)
	if err != nil {
		// 返回错误
		return nil, errors.WithStack(err)
	}

	// 处理响应数据
	data, err := getRespData(resp, "获取视频评论")
	if err != nil {
		// 发生错误,返回错误
		return nil, err
	}

	var ret *HotReply
	err = json.Unmarshal(data, &ret)
	// 返回 HotReply
	return ret, errors.WithStack(err)
}
