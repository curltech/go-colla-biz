package entity

import (
	baseentity "github.com/curltech/go-colla-core/entity"
)

type ActionSpec struct {
	Specification `xorm:"extends"`
	ExecuteType   string `xorm:"varchar(255)" json:"executeType,omitempty"`
	/**
	 * 在不同的executeType时代表不同含义：规则集和规则流名
	 *
	 * 类名和方法名，spring bean名和方法名，工作流定义编号
	 */
	ActionClass string `xorm:"varchar(255)" json:"actionClass,omitempty"`

	ActionName string `xorm:"varchar(255)" json:"actionName,omitempty"`

	BusinessType string `xorm:"varchar(255)" json:"businessType,omitempty"`

	DataType string `xorm:"varchar(255)" json:"dataType,omitempty"`

	ActionVersion string `xorm:"varchar(255)" json:"ActionVersion,omitempty"`
}

func (ActionSpec) TableName() string {
	return "spec_action"
}

func (ActionSpec) KeyName() string {
	return "SpecId"
}

func (ActionSpec) IdName() string {
	return baseentity.FieldName_Id
}
