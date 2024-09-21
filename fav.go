package bilibili

import (
	"github.com/go-resty/resty/v2"
)

type AddFavourFolderParam struct {
	Title   string `json:"title"`                                       // 收藏夹标题
	Intro   string `json:"intro,omitempty" request:"query,omitempty"`   // 收藏夹简介。默认为空
	Privacy int    `json:"privacy,omitempty" request:"query,omitempty"` // 是否公开。默认为公开。0：公开。1：私密
	Cover   string `json:"cover,omitempty" request:"query,omitempty"`   // 封面图url。封面会被审核
}

// AddFavourFolder 新建收藏夹
func (c *Client) AddFavourFolder(param AddFavourFolderParam) (*FavourFolderInfo, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/folder/add"
	)
	return execute[*FavourFolderInfo](c, method, url, param, fillCsrf(c))
}

type EditFavourFolderParam struct {
	MediaId int    `json:"media_id"`                                    // 目标收藏夹mdid
	Title   string `json:"title"`                                       // 修改收藏夹标题
	Intro   string `json:"intro,omitempty" request:"query,omitempty"`   // 修改收藏夹简介
	Privacy int    `json:"privacy,omitempty" request:"query,omitempty"` // 是否公开。默认为公开。。0：公开。1：私密
	Cover   string `json:"cover,omitempty" request:"query,omitempty"`   // 封面图url。封面会被审核
}

// EditFavourFolder 修改收藏夹
func (c *Client) EditFavourFolder(param EditFavourFolderParam) (*FavourFolderInfo, error) {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/folder/edit"
	)
	return execute[*FavourFolderInfo](c, method, url, param, fillCsrf(c))
}

type DeleteFavourFolderParam struct {
	MediaIds []int `json:"media_ids"` // 目标收藏夹mdid列表
}

// DeleteFavourFolder 删除收藏夹
func (c *Client) DeleteFavourFolder(param DeleteFavourFolderParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/folder/del"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type MoveFavourResourcesParam struct {
	SrcMediaId int      `json:"src_media_id"`                                 // 源收藏夹id
	TarMediaId int      `json:"tar_media_id"`                                 // 目标收藏夹id
	Mid        int      `json:"mid"`                                          // 当前用户mid
	Resources  []string `json:"resources"`                                    // 目标内容id列表。格式：{内容id}:{内容类型}。类型：2：视频稿件。12：音频。21：视频合集。内容id：。视频稿件：视频稿件avid。音频：音频auid。视频合集：视频合集id
	Platform   string   `json:"platform,omitempty" request:"query,omitempty"` // 平台标识。可为web
}

// CopyFavourResources 批量复制收藏内容
func (c *Client) CopyFavourResources(param MoveFavourResourcesParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/resource/copy"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

// MoveFavourResources 批量复制收藏内容
func (c *Client) MoveFavourResources(param MoveFavourResourcesParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/resource/move"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type DeleteFavourResourcesParam struct {
	Resources []int  `json:"resources"`                                    // 目标内容id列表。格式：{内容id}:{内容类型}。类型：2：视频稿件。12：音频。21：视频合集。内容id：。视频稿件：视频稿件avid。音频：音频auid。视频合集：视频合集id
	MediaId   int    `json:"media_id"`                                     // 目标收藏夹id
	Platform  string `json:"platform,omitempty" request:"query,omitempty"` // 平台标识。可为web
}

// DeleteFavourResources 批量删除收藏内容
func (c *Client) DeleteFavourResources(param DeleteFavourResourcesParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/resource/batch-del"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type MediaIdParam struct {
	MediaId int `json:"media_id"` // 目标收藏夹id
}

// CleanFavourResources 清空所有失效收藏内容
func (c *Client) CleanFavourResources(param MediaIdParam) error {
	const (
		method = resty.MethodPost
		url    = "https://api.bilibili.com/x/v3/fav/resource/clean"
	)
	_, err := execute[any](c, method, url, param, fillCsrf(c))
	return err
}

type CntInfo struct {
	Collect int `json:"collect"`  // 收藏数
	Play    int `json:"play"`     // 播放数
	ThumbUp int `json:"thumb_up"` // 点赞数
	Share   int `json:"share"`    // 分享数
}

type FavourFolderInfo struct {
	Id         int     `json:"id"`          // 收藏夹mlid（完整id），收藏夹原始id+创建者mid尾号2位
	Fid        int     `json:"fid"`         // 收藏夹原始id
	Mid        int     `json:"mid"`         // 创建者mid
	Attr       int     `json:"attr"`        // 属性位（？）
	Title      string  `json:"title"`       // 收藏夹标题
	Cover      string  `json:"cover"`       // 	收藏夹封面图片url
	Upper      Upper   `json:"upper"`       // 创建者信息
	CoverType  int     `json:"cover_type"`  // 封面图类别（？）
	CntInfo    CntInfo `json:"cnt_info"`    // 收藏夹状态数
	Type       int     `json:"type"`        // 类型（？）
	Intro      string  `json:"intro"`       // 备注
	Ctime      int     `json:"ctime"`       // 创建时间戳
	Mtime      int     `json:"mtime"`       // 收藏时间戳
	State      int     `json:"state"`       // 状态（？）
	FavState   int     `json:"fav_state"`   // 收藏夹收藏状态，已收藏：1，未收藏：0
	LikeState  int     `json:"like_state"`  // 点赞状态，已点赞：1，未点赞：0
	MediaCount int     `json:"media_count"` // 收藏夹内容数量
}

// GetFavourFolderInfo 获取收藏夹元数据
func (c *Client) GetFavourFolderInfo(param MediaIdParam) (*FavourFolderInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/folder/info"
	)
	return execute[*FavourFolderInfo](c, method, url, param)
}

type GetAllFavourFolderInfoParam struct {
	UpMid int `json:"up_mid"`                                   // 目标用户mid
	Type  int `json:"type,omitempty" request:"query,omitempty"` // 目标内容属性。默认为全部。0：全部。2：视频稿件
	Rid   int `json:"rid,omitempty" request:"query,omitempty"`  // 目标内容id。视频稿件：视频稿件avid
}

type AllFavourFolderInfo struct {
	Count int        `json:"count"` // 创建的收藏夹总数
	List  []struct { // 创建的收藏夹列表
		Id         int    `json:"id"`          // 收藏夹mlid（完整id），收藏夹原始id+创建者mid尾号2位
		Fid        int    `json:"fid"`         // 收藏夹原始id
		Mid        int    `json:"mid"`         // 创建者mid
		Attr       int    `json:"attr"`        // 属性位（？）
		Title      string `json:"title"`       // 收藏夹标题
		FavState   int    `json:"fav_state"`   // 目标id是否存在于该收藏夹，存在于该收藏夹：1，不存在于该收藏夹：0
		MediaCount int    `json:"media_count"` // 收藏夹内容数量
	} `json:"list"`
}

// GetAllFavourFolderInfo 获取指定用户创建的所有收藏夹信息
func (c *Client) GetAllFavourFolderInfo(param GetAllFavourFolderInfoParam) (*AllFavourFolderInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/folder/created/list-all"
	)
	return execute[*AllFavourFolderInfo](c, method, url, param)
}

type FavourInfo struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Intro    string `json:"intro"`
	Page     int    `json:"page"`
	Duration int    `json:"duration"`
	Upper    struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"upper"`
	Attr    int `json:"attr"`
	CntInfo struct {
		Collect int `json:"collect"`
		Play    int `json:"play"`
		Danmaku int `json:"danmaku"`
	} `json:"cnt_info"`
	Link    string `json:"link"`
	Ctime   int    `json:"ctime"`
	Pubtime int    `json:"pubtime"`
	FavTime int    `json:"fav_time"`
	BvId    string `json:"bv_id"`
	Bvid    string `json:"bvid"`
	Season  any    `json:"season"`
	Ugc     struct {
		FirstCid int `json:"first_cid"` // 视频cid
	} `json:"ugc"`
}

type GetFavourInfoParam struct {
	Resources []string `json:"resources"`                                    // 目标内容id列表。格式：{内容id}:{内容类型}。类型：2：视频稿件。12：音频。21：视频合集。内容id：视频稿件：视频稿件avid。音频：音频auid。视频合集：视频合集id。注意：一次最多只能请求100个内容id，超过100个内容id将不放回数据。
	Platform  string   `json:"platform,omitempty" request:"query,omitempty"` // 平台标识。可为web（影响内容列表类型）
}

// GetFavourInfo 获取收藏内容
func (c *Client) GetFavourInfo(param GetFavourInfoParam) ([]FavourInfo, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/resource/infos"
	)
	return execute[[]FavourInfo](c, method, url, param)
}

type GetFavourListParam struct {
	MediaId  int    `json:"media_id"`                                     // 目标收藏夹mlid（完整id）
	Tid      int    `json:"tid,omitempty" request:"query,omitempty"`      // 分区tid。默认为全部分区。0：全部分区
	Keyword  string `json:"keyword,omitempty" request:"query,omitempty"`  // 搜索关键字
	Order    string `json:"order,omitempty" request:"query,omitempty"`    // 排序方式。按收藏时间:mtime。按播放量: view。按投稿时间：pubtime
	Type     int    `json:"type,omitempty" request:"query,omitempty"`     // 查询范围。0：当前收藏夹（对应media_id）。 1：全部收藏夹
	Ps       int    `json:"ps"`                                           // 每页数量。定义域：1-20
	Pn       int    `json:"pn,omitempty" request:"query,omitempty"`       // 页码。默认为1
	Platform string `json:"platform,omitempty" request:"query,omitempty"` // 平台标识。可为web（影响内容列表类型）
}

type FavourList struct {
	Info struct { // 收藏夹元数据
		Id    int      `json:"id"`    // 收藏夹mlid（完整id），收藏夹原始id+创建者mid尾号2位
		Fid   int      `json:"fid"`   // 收藏夹原始id
		Mid   int      `json:"mid"`   // 创建者mid
		Attr  int      `json:"attr"`  // 属性，0：正常，1：失效
		Title string   `json:"title"` // 收藏夹标题
		Cover string   `json:"cover"` // 收藏夹封面图片url
		Upper struct { // 创建者信息
			Mid       int    `json:"mid"`        // 创建者mid
			Name      string `json:"name"`       // 创建者昵称
			Face      string `json:"face"`       // 创建者头像url
			Followed  bool   `json:"followed"`   // 是否已关注创建者
			VipType   int    `json:"vip_type"`   // 会员类别，0：无，1：月大会员，2：年度及以上大会员
			VipStatue int    `json:"vip_statue"` // 会员开通状态，0：无，1：有
		} `json:"upper"`
		CoverType int      `json:"cover_type"` // 封面图类别（？）
		CntInfo   struct { // 收藏夹状态数
			Collect int `json:"collect"`  // 收藏数
			Play    int `json:"play"`     // 播放数
			ThumbUp int `json:"thumb_up"` // 点赞数
			Share   int `json:"share"`    // 分享数
		} `json:"cnt_info"`
		Type       int    `json:"type"`        // 类型（？），一般是11
		Intro      string `json:"intro"`       // 备注
		Ctime      int    `json:"ctime"`       // 创建时间戳
		Mtime      int    `json:"mtime"`       // 收藏时间戳
		State      int    `json:"state"`       // 状态（？），一般为0
		FavState   int    `json:"fav_state"`   // 收藏夹收藏状态，已收藏收藏夹：1，未收藏收藏夹：0
		LikeState  int    `json:"like_state"`  // 点赞状态，已点赞：1，未点赞：0
		MediaCount int    `json:"media_count"` // 收藏夹内容数量
	} `json:"info"`
	Medias []struct { // 收藏夹内容
		Id       int      `json:"id"`       // 内容id，视频稿件：视频稿件avid，音频：音频auid，视频合集：视频合集id
		Type     int      `json:"type"`     // 内容类型，2：视频稿件，12：音频，21：视频合集
		Title    string   `json:"title"`    // 标题
		Cover    string   `json:"cover"`    // 封面url
		Intro    string   `json:"intro"`    // 简介
		Page     int      `json:"page"`     // 视频分P数
		Duration int      `json:"duration"` // 音频/视频时长
		Upper    struct { // UP主信息
			Mid  int    `json:"mid"`  // UP主mid
			Name string `json:"name"` // UP主昵称
			Face string `json:"face"` // UP主头像url
		} `json:"upper"`
		Attr    int      `json:"attr"` // 属性位（？）
		CntInfo struct { // 状态数
			Collect int `json:"collect"` // 收藏数
			Play    int `json:"play"`    // 播放数
			Danmaku int `json:"danmaku"` // 弹幕数
		} `json:"cnt_info"`
		Link    string `json:"link"`     // 跳转uri
		Ctime   int    `json:"ctime"`    // 投稿时间戳
		Pubtime int    `json:"pubtime"`  // 发布时间戳
		FavTime int    `json:"fav_time"` // 收藏时间戳
		BvId    string `json:"bv_id"`    // 视频稿件bvid
		Bvid    string `json:"bvid"`     // 视频稿件bvid
		Ugc     struct {
			FirstCid int `json:"first_cid"` // 视频cid
		} `json:"ugc"`
	} `json:"medias"`
	HasMore bool `json:"has_more"`
}

// GetFavourList 获取收藏夹内容明细列表
func (c *Client) GetFavourList(param GetFavourListParam) (*FavourList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/resource/list"
	)
	return execute[*FavourList](c, method, url, param)
}

type GetFavourIdsParam struct {
	MediaId  int    `json:"media_id"`                                     // 目标收藏夹mlid（完整id）
	Platform string `json:"platform,omitempty" request:"query,omitempty"` // 平台标识。可为web（影响内容列表类型）
}

type FavourId struct {
	Id   int    `json:"id"`    // 内容id，视频稿件：视频稿件avid，音频：音频auid，视频合集：视频合集id
	Type int    `json:"type"`  // 内容类型，2：视频稿件，12：音频，21：视频合集
	BvId string `json:"bv_id"` // 视频稿件bvid
	Bvid string `json:"bvid"`  // 视频稿件bvid
}

// GetFavourIds 获取收藏夹全部内容id
func (c *Client) GetFavourIds(param GetFavourIdsParam) ([]FavourId, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/resource/ids"
	)
	return execute[[]FavourId](c, method, url, param)
}

type SelfFavourList struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	MediaListResponse struct {
		Count int `json:"count"`
		List  []struct {
			ID       int64  `json:"id"`
			Fid      int    `json:"fid"`
			Mid      int    `json:"mid"`
			Attr     int    `json:"attr"`
			AttrDesc string `json:"attr_desc"`
			Title    string `json:"title"`
			Cover    string `json:"cover"`
			Upper    struct {
				Mid  int    `json:"mid"`
				Name string `json:"name"`
				Face string `json:"face"`
			} `json:"upper"`
			CoverType  int         `json:"cover_type"`
			Intro      string      `json:"intro"`
			Ctime      int         `json:"ctime"`
			Mtime      int         `json:"mtime"`
			State      int         `json:"state"`
			FavState   int         `json:"fav_state"`
			MediaCount int         `json:"media_count"`
			ViewCount  int         `json:"view_count"`
			Vt         int         `json:"vt"`
			IsTop      bool        `json:"is_top"`
			RecentFav  interface{} `json:"recent_fav"`
			PlaySwitch int         `json:"play_switch"`
			Type       int         `json:"type"`
			Link       string      `json:"link"`
			Bvid       string      `json:"bvid"`
		} `json:"list"`
		HasMore bool `json:"has_more"`
	} `json:"mediaListResponse"`
	URI string `json:"uri"`
}

// GetSelfFavourList 获取自己的收藏夹列表
func (c *Client) GetSelfFavourList() ([]SelfFavourList, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/v3/fav/folder/list4navigate"
	)
	return execute[[]SelfFavourList](c, method, url, nil)
}
