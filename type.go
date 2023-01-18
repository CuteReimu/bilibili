package bilibili

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
