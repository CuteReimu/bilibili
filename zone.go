package bilibili

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ZoneInfo 结构体用来表示CSV文件中的数据, 包含名称、代码、主分区tid、子分区tid,概述和备注等信息
type ZoneInfo struct {
	Name      string //中文名称
	Code      string //代号即英文名
	MasterTid int    //主分区tid
	Tid       int    //子分区tid
	Overview  string //概述,简介
	Remark    string //备注
}

// readCSV 从文件中读取CSV数据并转换为ZoneInfo切片,内部的一个工具函数。
func readCSV(filename string) ([]ZoneInfo, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()
	
	// 创建CSV读取器
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	
	var zoneInfos []ZoneInfo
	// 遍历每一行, 将每一行转换为ZoneInfo对象
	for _, record := range records[1:] { // 跳过标题行
		masterTid, err := strconv.Atoi(record[2]) // 将字符串转换为整数
		if err != nil {
			return nil, err
		}
		
		tid, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}
		
		info := ZoneInfo{
			Name:      record[0],
			Code:      record[1],
			MasterTid: masterTid,
			Tid:       tid,
			Overview:  record[4],
			Remark:    record[5],
		}
		zoneInfos = append(zoneInfos, info)
	}
	
	return zoneInfos, nil
}

// GetAllZoneInfos 获取所有ZoneInfo对象
func GetAllZoneInfos() ([]ZoneInfo, error) {
	// 读取CSV文件
	zoneInfos, err := readCSV("zone.csv")
	if err != nil {
		return nil, err
	}
	
	return zoneInfos, nil
}

// GetDescription 获取ZoneInfo对象的描述信息
func (info ZoneInfo) GetDescription() string {
	var description string
	description = "【分区】" + info.Name
	if info.Overview != "" {
		description += "\n【描述】" + info.Overview
	}
	if info.Remark != "" {
		description += "\n【备注】" + info.Remark
	}
	return description
}

// GetZoneInfoByTid 根据名称获取ZoneInfo对象
func GetZoneInfoByTid(tid int) (ZoneInfo, error) {
	// 读取CSV文件
	zoneInfos, err := readCSV("zone.csv")
	if err != nil {
		return ZoneInfo{}, err
	}
	
	//fmt.Println(zoneInfos)
	// 遍历ZoneInfo切片, 查找匹配名称的ZoneInfo对象
	for _, info := range zoneInfos {
		if info.Tid == tid {
			return info, nil
		}
	}
	
	// 如果没有找到匹配的ZoneInfo对象, 返回错误
	return ZoneInfo{}, fmt.Errorf("ZoneInfo not found")
}
