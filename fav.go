package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// AddFavourFolder 新建收藏夹
//
// title：收藏夹标题，必填。intro：收藏夹简介，非必填。
// privacy：是否为私密收藏夹。cover：封面图url。
func AddFavourFolder(title, intro string, privacy bool, cover string) (*FavourFolderInfo, error) {
	return std.AddFavourFolder(title, intro, privacy, cover)
}
func (c *Client) AddFavourFolder(title, intro string, privacy bool, cover string) (*FavourFolderInfo, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"title": title,
		"csrf":  biliJct,
	})
	if len(intro) > 0 {
		r.SetQueryParam("intro", intro)
	}
	if privacy {
		r.SetQueryParam("privacy", "1")
	}
	if len(cover) > 0 {
		r.SetQueryParam("cover", cover)
	}
	resp, err := r.Post("https://api.bilibili.com/x/v3/fav/folder/add")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "新建收藏夹")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *FavourFolderInfo
	err = json.Unmarshal(data, &ret)
	return ret, err
}

// EditFavourFolder 修改收藏夹
//
// media_id：目标收藏夹mdid，必填。
// title：收藏夹标题，必填。intro：收藏夹简介，非必填。
// privacy：是否为私密收藏夹。cover：封面图url。
func EditFavourFolder(mediaId int, title, intro string, privacy bool, cover string) (*FavourFolderInfo, error) {
	return std.EditFavourFolder(mediaId, title, intro, privacy, cover)
}
func (c *Client) EditFavourFolder(mediaId int, title, intro string, privacy bool, cover string) (*FavourFolderInfo, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	r := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"title":    title,
		"csrf":     biliJct,
	})
	if len(intro) > 0 {
		r.SetQueryParam("intro", intro)
	}
	if privacy {
		r.SetQueryParam("privacy", "1")
	}
	if len(cover) > 0 {
		r.SetQueryParam("cover", cover)
	}
	resp, err := r.Post("https://api.bilibili.com/x/v3/fav/folder/edit")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "修改收藏夹")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *FavourFolderInfo
	err = json.Unmarshal(data, &ret)
	return ret, err
}

// DeleteFavourFolder 删除收藏夹
//
// media_ids：目标收藏夹mdid列表，必填。
func DeleteFavourFolder(mediaIds []int) error {
	return std.DeleteFavourFolder(mediaIds)
}
func (c *Client) DeleteFavourFolder(mediaIds []int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	mediaIdsStr := make([]string, 0, len(mediaIds))
	for _, mediaId := range mediaIds {
		mediaIdsStr = append(mediaIdsStr, strconv.Itoa(mediaId))
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
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
func CopyFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	return std.CopyFavourResources(srcMediaId, tarMediaId, mid, resources, platform)
}
func (c *Client) CopyFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	if len(platform) == 0 {
		platform = "web"
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"src_media_id": strconv.Itoa(srcMediaId),
		"tar_media_id": strconv.Itoa(tarMediaId),
		"mid":          strconv.Itoa(mid),
		"resources":    strings.Join(resourcesStr, ","),
		"platform":     platform,
		"csrf":         biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/resource/copy")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量复制收藏内容")
	return err
}

// MoveFavourResources 批量移动收藏内容
func MoveFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	return std.MoveFavourResources(srcMediaId, tarMediaId, mid, resources, platform)
}
func (c *Client) MoveFavourResources(srcMediaId, tarMediaId, mid int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	if len(platform) == 0 {
		platform = "web"
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"src_media_id": strconv.Itoa(srcMediaId),
		"tar_media_id": strconv.Itoa(tarMediaId),
		"mid":          strconv.Itoa(mid),
		"resources":    strings.Join(resourcesStr, ","),
		"platform":     platform,
		"csrf":         biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/resource/move")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量移动收藏内容")
	return err
}

// DeleteFavourResources 批量删除收藏内容
func DeleteFavourResources(mediaId int, resources []Resource, platform string) error {
	return std.DeleteFavourResources(mediaId, resources, platform)
}
func (c *Client) DeleteFavourResources(mediaId int, resources []Resource, platform string) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resourcesStr := make([]string, 0, len(resources))
	for _, resource := range resources {
		resourcesStr = append(resourcesStr, resource.String())
	}
	if len(platform) == 0 {
		platform = "web"
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"media_id":  strconv.Itoa(mediaId),
		"resources": strings.Join(resourcesStr, ","),
		"platform":  platform,
		"csrf":      biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/resource/batch-del")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "批量删除收藏内容")
	return err
}

// CleanFavourResources 清空所有失效收藏内容
func CleanFavourResources(mediaId int) error {
	return std.CleanFavourResources(mediaId)
}
func (c *Client) CleanFavourResources(mediaId int) error {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"csrf":     biliJct,
	}).Post("https://api.bilibili.com/x/v3/fav/resource/clean")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = getRespData(resp, "清空所有失效收藏内容")
	return err
}

type FavourFolderInfo struct {
	Id    int      `json:"id"`    // 收藏夹mlid（完整id），收藏夹原始id+创建者mid尾号2位
	Fid   int      `json:"fid"`   // 收藏夹原始id
	Mid   int      `json:"mid"`   // 创建者mid
	Attr  int      `json:"attr"`  // 属性位（？）
	Title string   `json:"title"` // 收藏夹标题
	Cover string   `json:"cover"` // 	收藏夹封面图片url
	Upper struct { // 创建者信息
		Mid       int    `json:"mid"`        // 创建者mid
		Name      string `json:"name"`       // 创建者昵称
		Face      string `json:"face"`       // 创建者头像url
		Followed  bool   `json:"followed"`   // 是否已关注创建者
		VipType   int    `json:"vip_type"`   // 会员类别，0：无，1：月大会员，2：年度及以上大会员
		VipStatue int    `json:"vip_statue"` // 0：无，1：有
	} `json:"upper"`
	CoverType int      `json:"cover_type"` // 封面图类别（？）
	CntInfo   struct { // 收藏夹状态数
		Collect int `json:"collect"`  // 收藏数
		Play    int `json:"play"`     // 播放数
		ThumbUp int `json:"thumb_up"` // 点赞数
		Share   int `json:"share"`    // 分享数
	} `json:"cnt_info"`
	Type       int    `json:"type"`        // 类型（？）
	Intro      string `json:"intro"`       // 备注
	Ctime      int    `json:"ctime"`       // 创建时间戳
	Mtime      int    `json:"mtime"`       // 收藏时间戳
	State      int    `json:"state"`       // 状态（？）
	FavState   int    `json:"fav_state"`   // 收藏夹收藏状态，已收藏：1，未收藏：0
	LikeState  int    `json:"like_state"`  // 点赞状态，已点赞：1，未点赞：0
	MediaCount int    `json:"media_count"` // 收藏夹内容数量
}

// GetFavourFolderInfo 获取收藏夹元数据
func GetFavourFolderInfo(mediaId int) (*FavourFolderInfo, error) {
	return std.GetFavourFolderInfo(mediaId)
}
func (c *Client) GetFavourFolderInfo(mediaId int) (*FavourFolderInfo, error) {
	biliJct := c.getCookie("bili_jct")
	if len(biliJct) == 0 {
		return nil, errors.New("B站登录过期")
	}
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").SetQueryParams(map[string]string{
		"media_id": strconv.Itoa(mediaId),
		"csrf":     biliJct,
	}).Get("https://api.bilibili.com/x/v3/fav/folder/info")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取收藏夹元数据")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret *FavourFolderInfo
	err = json.Unmarshal(data, &ret)
	return ret, err
}
