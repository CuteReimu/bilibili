package bilibili

type Comment struct { // 评论条目对象
	Rpid      int64    `json:"rpid"`      // 评论 rpid
	Oid       int      `json:"oid"`       // 评论区对象 id
	Type      int      `json:"type"`      // 评论区类型代码
	Mid       int      `json:"mid"`       // 发送者 mid
	Root      int      `json:"root"`      // 根评论 rpid，若为一级评论则为 0，大于一级评论则为根评论 id
	Parent    int      `json:"parent"`    // 回复父评论 rpid，若为一级评论则为 0，大于一级评论则为根评论 rpid，大于二级评论为上一级评论 rpid
	Dialog    int      `json:"dialog"`    // 回复对方 rpid，若为一级评论则为 0，大于一级评论则为根评论 rpid，大于二级评论为上一级评论 rpid
	Count     int      `json:"count"`     // 二级评论条数
	Rcount    int      `json:"rcount"`    // 回复评论条数
	Floor     int      `json:"floor"`     // 评论楼层号
	State     int      `json:"state"`     // 作用尚不明确
	Fansgrade int      `json:"fansgrade"` // 是否具有粉丝标签，0：无，1：有
	Attr      int      `json:"attr"`      // 作用尚不明确
	Ctime     int      `json:"ctime"`     // 评论发送时间戳
	Like      int      `json:"like"`      // 评论获赞数
	Action    int      `json:"action"`    // 当前用户操作状态，需要登录(Cookie 或 APP)，否则恒为 0，0：无，1：已点赞，2：已点踩
	Content   struct { // 评论信息
		Message string `json:"message"`
		Plat    int    `json:"plat"`
		Device  string `json:"device"`
	} `json:"content"`
	Replies    interface{} `json:"replies"`     // 评论回复条目预览
	Assist     int         `json:"assist"`      // 作用尚不明确
	ShowFollow bool        `json:"show_follow"` // 作用尚不明确
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
