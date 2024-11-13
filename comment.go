package bilibili

import "github.com/go-resty/resty/v2"

type GetCommentsDetailParam struct {
	AccessKey string `json:"access_key,omitempty" request:"query,omitempty"` // APP 登录 Token
	Type      int    `json:"type"`                                           // 评论区类型代码，见 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/comment/readme.md
	Oid       int    `json:"oid"`                                            // 目标评论区 id
	Sort      int    `json:"sort,omitempty" request:"query,omitempty"`       // 排序方式。默认为0。0：按时间。1：按点赞数。2：按回复数
	Nohot     int    `json:"nohot,omitempty" request:"query,omitempty"`      // 是否不显示热评。默认为0。1：不显示。0：显示
	Ps        int    `json:"ps,omitempty" request:"query,omitempty"`         // 每页项数。默认为20。定义域：1-20
	Pn        int    `json:"pn,omitempty" request:"query,omitempty"`         // 页码。默认为1
}

type CommentsPage struct {
	Num    int `json:"num"`    // 当前页码
	Size   int `json:"size"`   // 每页项数
	Count  int `json:"count"`  // 根评论条数
	Acount int `json:"acount"` // 总计评论条数
}

type CommentsConfig struct {
	Showadmin  int  `json:"showadmin"`    // 是否显示管理置顶
	Showentry  int  `json:"showentry"`    // (?)
	Showfloor  int  `json:"showfloor"`    // 是否显示楼层号
	Showtopic  int  `json:"showtopic"`    // 是否显示话题
	ShowUpFlag bool `json:"show_up_flag"` // 是否显示“UP 觉得很赞”标志
	ReadOnly   bool `json:"read_only"`    // 是否只读评论区
	ShowDelLog bool `json:"show_del_log"` // 是否显示删除记录
}

type Picture struct {
	ImgSrc    string  `json:"img_src"`    // 图片地址
	ImgWidth  int     `json:"img_width"`  // 图片宽度
	ImgHeight int     `json:"img_height"` // 图片高度
	ImgSize   float64 `json:"img_size"`   // 图片大小。单位KB
}

type CommentContent struct {
	Message  string    `json:"message"`  // 评论内容。**重要**
	Plat     int       `json:"plat"`     // 评论发送端。1：web端。2：安卓客户端。3：ios 客户端。4：wp 客户端
	Device   string    `json:"device"`   // 评论发送平台设备
	Members  []Member  `json:"members"`  // at 到的用户信息
	Emote    any       `json:"emote"`    // 需要渲染的表情转义。评论内容无表情则无此项
	JumpUrl  any       `json:"jump_url"` // 需要高亮的超链转义
	MaxLine  int       `json:"max_line"` // 6。收起最大行数
	Pictures []Picture `json:"pictures"` // 评论图片数组
}

type Folder struct {
	HasFolded bool   `json:"has_folded"` // 是否有被折叠的二级评论
	IsFolded  bool   `json:"is_folded"`  // 评论是否被折叠
	Rule      string `json:"rule"`       // 相关规则页面 url
}

type UpAction struct {
	Like  bool `json:"like"`  // 是否UP主觉得很赞。false：否。true：是
	Reply bool `json:"reply"` // 是否被UP主回复。false：否。true：是
}

type CardLabel struct {
	Rpid             int    `json:"rpid"`              // 评论 rpid
	TextContent      string `json:"text_content"`      // 标签文本。已知有妙评
	TextColorDay     string `json:"text_color_day"`    // 日间文本颜色。十六进制颜色值，下同
	TextColorNight   string `json:"text_color_night"`  // 夜间文本颜色
	LabelColorDay    string `json:"label_color_day"`   // 日间标签颜色
	LabelColorNight  string `json:"label_color_night"` // 夜间标签颜色
	Image            string `json:"image"`             // 作用不明
	Type             int    `json:"type"`              // 1。作用不明
	Background       string `json:"background"`        // 背景图片 url
	BackgroundWidth  int    `json:"background_width"`  // 背景图片宽度
	BackgroundHeight int    `json:"background_height"` // 背景图片高度
	JumpUrl          string `json:"jump_url"`          // 跳转链接
	Effect           int    `json:"effect"`            // 0。作用不明，可能用于控制动画，下同
	EffectStartTime  int    `json:"effect_start_time"` // 0
}

type ReplyControl struct {
	SubReplyEntryText string `json:"sub_reply_entry_text"` // 回复提示。共 xx 条回复
	SubReplyTitleText string `json:"sub_reply_title_text"` // 回复提示。相关回复共有 xx 条
	TimeDesc          string `json:"time_desc"`            // 时间提示。xx 天/小时 前发布
	Location          string `json:"location"`             // IP属地。IP属地：xx。评论者发送评论时的IP地址属地。仅对2022-07-25 11:00及以后发布的评论有效。需要登录
}

type Comment struct {
	Rpid         int            `json:"rpid"`          // 评论 rpid
	Oid          int            `json:"oid"`           // 评论区对象 id
	Type         int            `json:"type"`          // 评论区类型代码。
	Mid          int            `json:"mid"`           // 发送者 mid
	Root         int            `json:"root"`          // 根评论 rpid。若为一级评论则为 0。大于一级评论则为根评论 id
	Parent       int            `json:"parent"`        // 回复父评论 rpid。若为一级评论则为 0。若为二级评论则为根评论 rpid。大于二级评论为上一级评 论 rpid
	Dialog       int            `json:"dialog"`        // 回复对方 rpid。若为一级评论则为 0。若为二级评论则为该评论 rpid。大于二级评论为上一级评论 rpid
	Count        int            `json:"count"`         // 二级评论条数
	Rcount       int            `json:"rcount"`        // 回复评论条数
	Floor        int            `json:"floor"`         // 评论楼层号。**注：若不支持楼层则无此项**
	State        int            `json:"state"`         // (?)
	Fansgrade    int            `json:"fansgrade"`     // 是否具有粉丝标签。0：无。1：有
	Attr         int            `json:"attr"`          // 某属性位？
	Ctime        int            `json:"ctime"`         // 评论发送时间。时间戳
	RpidStr      string         `json:"rpid_str"`      // 评论rpid。字串格式
	RootStr      string         `json:"root_str"`      // 根评论rpid。字串格式
	ParentStr    string         `json:"parent_str"`    // 回复父评论rpid。字串格式
	Like         int            `json:"like"`          // 评论获赞数
	Action       int            `json:"action"`        // 当前用户操作状态。需要登录(Cookie 或 APP) 。否则恒为 0。0：无。1：已点赞。2：已点踩
	Member       Member         `json:"member"`        // 评论发送者信息
	Content      CommentContent `json:"content"`       // 评论信息
	Replies      []*Comment     `json:"replies"`       // 评论回复条目预览。**仅嵌套一层**。否则为 null
	Assist       int            `json:"assist"`        // (?)
	Folder       Folder         `json:"folder"`        // 折叠信息
	UpAction     UpAction       `json:"up_action"`     // 评论 UP 主操作信息
	ShowFollow   bool           `json:"show_follow"`   // (?)
	Invisible    bool           `json:"invisible"`     // 评论是否被隐藏
	CardLabel    []CardLabel    `json:"card_label"`    // 右上角卡片标签信息
	ReplyControl ReplyControl   `json:"reply_control"` // 评论提示文案信息
}

type CommentsControl struct {
	InputDisable          bool   `json:"input_disable"`            // 是否禁止新增评论。用户涉及合约争议，锁定该用户所有稿件、动态的评论区，不允许新增评论，root_input_text和child_input_text值为“当前评论区不可新增评论”
	RootInputText         string `json:"root_input_text"`          // 评论框文字
	ChildInputText        string `json:"child_input_text"`         // 评论框文字
	BgText                string `json:"bg_text"`                  // 空评论区文字
	WebSelection          bool   `json:"web_selection"`            // 评论是否筛选后可见。false：无需筛选。true：需要筛选
	AnswerGuideText       string `json:"answer_guide_text"`        // 答题页面链接文字
	AnswerGuideIconUrl    string `json:"answer_guide_icon_url"`    // 答题页面图标 url
	AnswerGuideIosUrl     string `json:"answer_guide_ios_url"`     // 答题页面 ios url
	AnswerGuideAndroidUrl string `json:"answer_guide_android_url"` // 答题页面安卓 url
}

type CommentsDetail struct {
	Page        CommentsPage    `json:"page"`         // 页信息
	Config      CommentsConfig  `json:"config"`       // 评论区显示控制
	Replies     []*Comment      `json:"replies"`      // 评论列表
	Hots        []*Comment      `json:"hots"`         // 热评列表
	Upper       *Comment        `json:"upper"`        // 置顶评论
	Top         any             `json:"top"`          // (?)
	Notice      *Notice         `json:"notice"`       // 评论区公告信息
	Vote        int             `json:"vote"`         // 投票评论?
	Blacklist   int             `json:"blacklist"`    // (?)
	Assist      int             `json:"assist"`       // (?)
	Mode        int             `json:"mode"`         // 评论区类型id
	SupportMode []int           `json:"support_mode"` // 评论区支持的类型id
	Folder      Folder          `json:"folder"`       // 折叠相关信息
	LotteryCard any             `json:"lottery_card"` // (?)
	ShowBvid    bool            `json:"show_bvid"`    // 显示bvid?
	Control     CommentsControl `json:"control"`      // 评论区输入属性
}

// GetCommentsDetail 获取评论区明细
func (c *Client) GetCommentsDetail(param GetCommentsDetailParam) (*CommentsDetail, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v2/reply"
	)
	return execute[*CommentsDetail](c, method, url, param)
}

type GetCommentReplyParam struct {
	AccessKey string `json:"access_key,omitempty" request:"query,omitempty"` // APP登录 Token
	Type      int    `json:"type"`                                           // 评论区类型代码，见 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/comment/readme.md
	Oid       int    `json:"oid"`                                            // 目标评论区 id
	Root      int    `json:"root"`                                           // 根回复 rpid
	Ps        int    `json:"ps,omitempty" request:"query,omitempty"`         // 每页项数。默认为20。定义域：1-49 。 但 data_replies 的最大内容数为20,因此设置为49其实也只会有20条回复被返回
	Pn        int    `json:"pn,omitempty" request:"query,omitempty"`         // 页码。默认为1
}

type CommentReply struct {
	Config   CommentsConfig  `json:"config"`    // 评论区显示控制
	Control  CommentsControl `json:"control"`   // 评论区输入属性
	Page     CommentsPage    `json:"page"`      // 页面信息
	Replies  []*Comment      `json:"replies"`   // 评论对话树列表。最大内容数为20
	Root     *Comment        `json:"root"`      // 根评论信息
	ShowBvid bool            `json:"show_bvid"` // 显示 bvid?
	ShowText string          `json:"show_text"` // (?)
	ShowType int             `json:"show_type"` // (?)
	Upper    Upper           `json:"upper"`     // UP主 mid
}

// GetCommentReply 获取指定评论的回复，按照回复顺序排序
func (c *Client) GetCommentReply(param GetCommentReplyParam) (*CommentReply, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v2/reply/reply"
	)
	return execute[*CommentReply](c, method, url, param)
}

type GetCommentsHotReplyParam struct {
	Type int `json:"type"`                                   // 评论区类型代码
	Oid  int `json:"oid"`                                    // 目标评论区 id
	Root int `json:"root"`                                   // 根回复 rpid
	Ps   int `json:"ps,omitempty" request:"query,omitempty"` // 每页项数。默认为20。定义域：1-49
	Pn   int `json:"pn,omitempty" request:"query,omitempty"` // 页码。默认为1
}

type CommentsHotReply struct {
	Page    CommentsPage `json:"page"`    // 页面信息
	Replies []Comment    `json:"replies"` // 热评列表
}

// GetCommentsHotReply 获取评论区热评
func (c *Client) GetCommentsHotReply(param GetCommentsHotReplyParam) (*CommentsHotReply, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v2/reply/hot"
	)
	return execute[*CommentsHotReply](c, method, url, param)
}
