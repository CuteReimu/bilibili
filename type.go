package bilibili

import (
	"strconv"
)

type Size struct {
	Width  int
	Height int
}

type FormatCtrl struct {
	Location int    `json:"location"` // 从全文第几个字开始
	Type     int    `json:"type"`     // 1：At
	Length   int    `json:"length"`   // 长度总共多少个字
	Data     string `json:"data"`     // 当Type为1时，这里填At的人的Uid
}

type ResourceType int

var (
	ResourceTypeVideo     ResourceType = 2  // 视频稿件
	ResourceTypeAudio                  = 12 // 音频
	ResourceTypeVideoList              = 21 // 视频合集
)

type Resource struct {
	Id   int
	Type ResourceType
}

func (r Resource) String() string {
	return strconv.Itoa(r.Id) + ":" + strconv.Itoa(int(r.Type))
}
