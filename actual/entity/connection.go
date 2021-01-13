package entity

import (
	"github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
	"time"
)

//表示两个role之间的多对多连接
type Connection struct {
	entity.InternalFixedActual `xorm:"extends"`
	ActualId                   uint64     `json:"actualId,omitempty"`
	EffectiveDate              *time.Time `json:"effectiveDate,omitempty"`
	Version                    uint64     `json:"version,omitempty"`
	// 角色的开始有效日期
	StartDate *time.Time `json:"startDate,omitempty"`
	// 角色的有效结束日期，表示角色事实被后一个版本替换掉的日期
	EndDate *time.Time `json:"endDate,omitempty"`
}

func (Connection) TableName() string {
	return "atl_connection"
}

func (Connection) KeyName() string {
	return "ActualId"
}

func (Connection) IdName() string {
	return baseentity.FieldName_Id
}
