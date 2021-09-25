package entity

import (
	"github.com/curltech/go-colla-core/entity"
)

/**
文件名即股票代码
每32个字节为一天数据
每4个字节为一个字段，每个字段内低字节在前
00 ~ 03 字节：年月日, 整型
04 ~ 07 字节：开盘价*1000， 整型
08 ~ 11 字节：最高价*1000,  整型
12 ~ 15 字节：最低价*1000,  整型
16 ~ 19 字节：收盘价*1000,  整型
20 ~ 23 字节：成交额（元），float型
24 ~ 27 字节：成交量（手），整型
28 ~ 31 字节：上日收盘*1000, 整型
*/
type DayData struct {
	entity.BaseEntity `xorm:"extends"`
	ShareId           string  `xorm:"varchar(255)" json:"shareId,omitempty"`
	DayDate           int32   `json:"dayDate,omitempty"`
	OpeningPrice      int32   `json:"OpeningPrice,omitempty"`
	CeilingPrice      int32   `json:"CeilingPrice,omitempty"`
	FloorPrice        int32   `json:"FloorPrice,omitempty"`
	ClosingPrice      int32   `json:"ClosingPrice,omitempty"`
	TurnVolume        float32 `json:"TurnVolume,omitempty"`
	Volume            int32   `json:"Volume,omitempty"`
	Previous          int32   `json:"Previous,omitempty"`
}

func (DayData) TableName() string {
	return "stk_daydata"
}

func (DayData) KeyName() string {
	return entity.FieldName_Id
}

func (DayData) IdName() string {
	return entity.FieldName_Id
}