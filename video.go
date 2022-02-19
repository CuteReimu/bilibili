package bilibili

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

type VideoDimension struct {
	Width  int `json:"width"`  // 当前分P 宽度
	Height int `json:"height"` // 当前分P 高度
	Rotate int `json:"rotate"` // 是否将宽高对换，0：正常，1：对换
}

type VideoInfo struct {
	Code    int      `json:"code"`    // 返回值，0表示成功
	Message string   `json:"message"` // 错误信息
	TTL     int      `json:"ttl"`     // 固定值1，作用尚不明确
	Data    struct { // 信息本体
		Bvid      string     `json:"bvid"`      // 稿件bvid
		Aid       int        `json:"aid"`       // 稿件avid
		Videos    int        `json:"videos"`    // 稿件分P总数，默认为1
		Tid       int        `json:"tid"`       // 分区tid
		Tname     string     `json:"tname"`     // 子分区名称
		Copyright int        `json:"copyright"` // 1：原创，2：转载
		Pic       string     `json:"pic"`       // 稿件封面图片url
		Title     string     `json:"title"`     // 稿件标题
		Pubdate   int        `json:"pubdate"`   // 稿件发布时间戳
		Ctime     int        `json:"ctime"`     // 用户投稿时间戳
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
		Dynamic   string         `json:"dynamic"`   // 视频同步发布的的动态的文字内容
		Cid       int            `json:"cid"`       // 视频1P cid
		Dimension VideoDimension `json:"dimension"` // 视频1P分辨率
		NoCache   bool           `json:"no_cache"`  // 固定值true，作用尚不明确
		Pages     []struct {     // 视频分P列表，无分P则数组只有1个元素
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
		Staff []struct { // 合作成员列表，非合作视频无此项
			Mid   int      `json:"mid"`   // 成员mid
			Title string   `json:"title"` // 成员名称
			Name  string   `json:"name"`  // 成员昵称
			Face  string   `json:"face"`  // 成员头像url
			Vip   struct { // 成员大会员状态
				Type      int `json:"type"`       // 成员会员类型，0：无，1：月会员，2：年会员
				Status    int `json:"status"`     // 会员状态，0：无，1：有
				ThemeType int `json:"theme_type"` // 固定值0，作用尚不明确
			} `json:"vip"`
			Official struct { // 成员认证信息
				Role  int    `json:"role"`  // 成员认证级别，0：无，1 2 7：个人认证，3 4 5 6：机构认证
				Title string `json:"title"` // 成员认证名
				Desc  string `json:"desc"`  // 成员认证备注
				Type  int    `json:"type"`  // 成员认证类型，-1：无，0：有
			} `json:"official"`
			Follower int `json:"follower"` // 成员粉丝数
		} `json:"staff"`
		UserGarb struct { // 用户装扮信息
			UrlImageAniCut string `json:"url_image_ani_cut"` // 某url，作用尚不明确
		} `json:"user_garb"`
	} `json:"data"`
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
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取视频详细信息失败，status code: %d", resp.StatusCode())
	}
	var ret *VideoInfo
	err = json.Unmarshal(resp.Body(), &ret)
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
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("获取视频详细信息失败，status code: %d", resp.StatusCode())
	}
	var ret *VideoInfo
	err = json.Unmarshal(resp.Body(), &ret)
	return ret, errors.WithStack(err)
}

var regBv = regexp.MustCompile("(?i)bv([0-9A-Za-z]{10})")

// GetVideoInfoByShortUrl 通过短链接获取视频信息
func GetVideoInfoByShortUrl(shortUrl string) (*VideoInfo, error) {
	return std.GetVideoInfoByShortUrl(shortUrl)
}
func (c *Client) GetVideoInfoByShortUrl(shortUrl string) (*VideoInfo, error) {
	resp, err := c.resty().SetRedirectPolicy(resty.NoRedirectPolicy()).R().Get(shortUrl)
	if resp == nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 302 {
		return nil, errors.Errorf("获取视频详细信息失败，status code: %d", resp.StatusCode())
	}
	url := resp.Header().Get("Location")
	ret := regBv.FindAllStringSubmatch(url, 1)
	if len(ret) != 1 {
		return nil, errors.New("通过短链接获取视频信息失败：" + url)
	}
	return GetVideoInfoByBvid(ret[0][0])
}
