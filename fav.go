package bilibili

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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

// DeleteFavourFolder 删除收藏夹
//
// media_ids：目标收藏夹mdid列表，必填。
func (c *Client) DeleteFavourFolder(mediaIds []int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	mediaIdsStr := make([]string, 0, len(mediaIds))
	for _, mediaId := range mediaIds {
		mediaIdsStr = append(mediaIdsStr, strconv.Itoa(mediaId))
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"media_ids": strings.Join(mediaIdsStr, ","),
		"csrf":      biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/folder/del")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "删除收藏夹")
	return err
}

// CopyFavourResources 批量复制收藏内容
func (c *Client) CopyFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	r := c.resty.R().SetQueryParams(map[string]string{
		"src_media_id": strconv.Itoa(srcMediaId),
		"tar_media_id": strconv.Itoa(tarMediaId),
		"mid":          strconv.Itoa(mid),
		"resources":    strings.Join(resourcesStr, ","),
		"csrf":         biliJct,
	})
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Post("https://api.bilibili.com/x/v3/fav/resource/copy")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量复制收藏内容")
	return err
}

// MoveFavourResources 批量移动收藏内容
func (c *Client) MoveFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	r := c.resty.R().SetQueryParams(map[string]string{
		"src_media_id": strconv.Itoa(srcMediaId),
		"tar_media_id": strconv.Itoa(tarMediaId),
		"mid":          strconv.Itoa(mid),
		"resources":    strings.Join(resourcesStr, ","),
		"csrf":         biliJct,
	})
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Post("https://api.bilibili.com/x/v3/fav/resource/move")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量移动收藏内容")
	return err
}

// DeleteFavourResources 批量删除收藏内容
func (c *Client) DeleteFavourResources(mediaId int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	r := c.resty.R().SetQueryParams(map[string]string{
		"media_id":  strconv.Itoa(mediaId),
		"resources": strings.Join(resourcesStr, ","),
		"csrf":      biliJct,
	})
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Post("https://api.bilibili.com/x/v3/fav/resource/batch-del")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量删除收藏内容")
	return err
}

// CleanFavourResources 清空所有失效收藏内容
func (c *Client) CleanFavourResources(mediaId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty.R().SetQueryParams(map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"csrf":     biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/resource/clean")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "清空所有失效收藏内容")
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
func (c *Client) GetFavourFolderInfo(mediaId int) (*FavourFolderInfo, error) {
	resp, err := c.resty.R().
		SetQueryParam("media_id", strconv.Itoa(mediaId)).Get("https://api.bilibili.com/x/v3/fav/folder/info")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取收藏夹元数据")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *FavourFolderInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
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

func (c *Client) GetAllFavourFolderInfo(upMid, attrType, rid int) (*AllFavourFolderInfo, error) {
	r := c.resty.R().SetQueryParams(map[string]string{
		"up_mid": strconv.Itoa(upMid),
		"type":   strconv.Itoa(attrType),
	})
	if rid != 0 {
		r = r.SetQueryParam("rid", strconv.Itoa(rid))
	}
	resp, err := r.Get("https://api.bilibili.com/x/v3/fav/folder/created/list-all")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取指定用户创建的所有收藏夹信息")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *AllFavourFolderInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
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
	Link    string      `json:"link"`
	Ctime   int         `json:"ctime"`
	Pubtime int         `json:"pubtime"`
	FavTime int         `json:"fav_time"`
	BvId    string      `json:"bv_id"`
	Bvid    string      `json:"bvid"`
	Season  interface{} `json:"season"`
}

// GetFavourInfo 获取收藏内容

func (c *Client) GetFavourInfo(resources []Resource, platform string) ([]*FavourInfo, error) {
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	r := c.resty.R().SetQueryParam("resources", strings.Join(resourcesStr, ","))
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Get("https://api.bilibili.com/x/v3/fav/resource/infos")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取收藏内容")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret []*FavourInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
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
	} `json:"medias"`
	HasMore bool `json:"has_more"`
}

// GetFavourList 获取收藏夹内容明细列表

func (c *Client) GetFavourList(mediaId, tid int, keyword, order string, searchType, ps, pn int, platform string) (*FavourList, error) {
	if pn == 0 {
		pn = 1
	}
	r := c.resty.R().SetQueryParams(map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"tid":      strconv.Itoa(tid),
		"type":     strconv.Itoa(searchType),
		"ps":       strconv.Itoa(ps),
		"pn":       strconv.Itoa(pn),
	})
	if len(keyword) > 0 {
		r = r.SetQueryParam("keyword", keyword)
	}
	if len(order) > 0 {
		r = r.SetQueryParam("order", order)
	}
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Get("https://api.bilibili.com/x/v3/fav/resource/list")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取收藏夹内容明细列表")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *FavourList
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type FavourId struct {
	Id   int    `json:"id"`    // 内容id，视频稿件：视频稿件avid，音频：音频auid，视频合集：视频合集id
	Type int    `json:"type"`  // 内容类型，2：视频稿件，12：音频，21：视频合集
	BvId string `json:"bv_id"` // 视频稿件bvid
	Bvid string `json:"bvid"`  // 视频稿件bvid
}

// GetFavourIds 获取收藏夹全部内容id

func (c *Client) GetFavourIds(mediaId int, platform string) ([]*FavourId, error) {
	r := c.resty.R().SetQueryParam("media_id", strconv.Itoa(mediaId))
	if len(platform) > 0 {
		r = r.SetQueryParam("platform", platform)
	}
	resp, err := r.Get("https://api.bilibili.com/x/v3/fav/resource/ids")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取收藏夹全部内容id")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret []*FavourId
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
