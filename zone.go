package bilibili

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"github.com/pkg/errors"
	"strconv"
)

// ZoneInfo 结构体用来表示CSV文件中的数据, 包含名称、代码、主分区tid、子分区tid,概述和备注等信息
type ZoneInfo struct {
	Name      string //中文名称
	Code      string //代号即英文名
	MasterTid int    //主分区tid
	Tid       int    //子分区tid
	Overview  string //概述,简介
}

//go:embed zone.csv
var embeddedCSV []byte

// readCSV 从文件中读取CSV数据并转换为ZoneInfo切片,内部的一个工具函数。
func readCSV() ([]ZoneInfo, error) {
	// 打开文件

	csvReader := bytes.NewReader(embeddedCSV)

	// 创建CSV读取器
	reader := csv.NewReader(csvReader)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var zoneInfos []ZoneInfo
	// 遍历每一行, 将每一行转换为ZoneInfo对象
	for _, record := range records[1:] { // 跳过标题行
		masterTid, err := strconv.Atoi(record[2]) // 将字符串转换为整数
		if err != nil {
			return nil, errors.WithStack(err)
		}

		tid, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		info := ZoneInfo{
			Name:      record[0],
			Code:      record[1],
			MasterTid: masterTid,
			Tid:       tid,
			Overview:  record[4],
		}
		zoneInfos = append(zoneInfos, info)
	}

	return zoneInfos, nil
}

// GetAllZoneInfos 获取所有ZoneInfo对象
func GetAllZoneInfos() ([]ZoneInfo, error) {
	// 读取CSV文件
	zoneInfos, err := readCSV()
	if err != nil {
		return nil, err
	}

	return zoneInfos, nil
}

// GetDescription 获取ZoneInfo对象的描述信息,描述信息分为四个部分。当前分区，主分区，描述和备注。
func (info ZoneInfo) GetDescription() string {
	var description string
	var masterInfo, _ = GetZoneInfoByTid(info.MasterTid)

	description = "【分区】" + info.Name
	description += "\n【主分区】" + masterInfo.Name

	if info.Overview != "" {
		description += "\n【描述】" + info.Overview
		description += info.Overview
	}
	//备注写到主分区的备注中了
	if masterInfo.Overview != "" {
		description += "\n【备注】" + masterInfo.Overview
	}
	return description
}

// GetZoneInfoByTid 根据名称获取ZoneInfo对象
func GetZoneInfoByTid(tid int) (ZoneInfo, error) {
	// 读取CSV文件
	zoneInfos, err := readCSV()
	if err != nil {
		return ZoneInfo{}, errors.WithStack(err)
	}

	// 遍历ZoneInfo切片, 查找匹配名称的ZoneInfo对象
	for _, info := range zoneInfos {
		if info.Tid == tid {
			return info, nil
		}
	}

	// 如果没有找到匹配的ZoneInfo对象, 返回错误
	return ZoneInfo{}, errors.Errorf("ZoneInfo not found")
}
