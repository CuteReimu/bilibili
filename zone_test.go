package bilibili

import "testing"

func TestZoneInfo(t *testing.T) {
	println("开始测试。")
	defer println("测试结束")

	infos, _ := GetAllZoneInfos()
	println("一共有", len(infos), "个区")
	for _, zoneInfo := range infos {
		println(zoneInfo.GetDescription())
	}

	zoneInfo, err := GetZoneInfoByTid(24)

	if err != nil {
		t.Error(err)
	}
	println(zoneInfo.GetDescription())

}
