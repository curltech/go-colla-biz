package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type Share struct {
	entity.BaseEntity `xorm:"extends"`
	/**
	 * 编号
	 */
	ShareId string `xorm:"varchar(255)" json:"shareId,omitempty"`
	/**
	 * 英文名称
	 */
	Kind string `xorm:"varchar(255)" json:"kind,omitempty"`
	/**
	 * 名称
	 */
	Name string `xorm:"varchar(255)" json:"name,omitempty"`
	// 有效日期
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`
}

func (Share) TableName() string {
	return "stk_share"
}

func (Share) KeyName() string {
	return "ShareId"
}

func (Share) IdName() string {
	return entity.FieldName_Id
}
