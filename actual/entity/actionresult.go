package entity

import (
	"github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
)

type ActionResult struct {
	entity.InternalFixedActual `xorm:"extends"`
	ExecuteType                string `xorm:"varchar(255)" json:"executeType,omitempty"`
	/**
	 * 在不同的executeType时代表不同含义：规则集和规则流名 类名和方法名，spring bean名和方法名，工作流定义编号
	 */
	ActionClass       string `xorm:"varchar(255)" json:"actionClass,omitempty"`
	ActionName        string `xorm:"varchar(255)" json:"actionName,omitempty"`
	Value             string `xorm:"varchar(255)" json:"value,omitempty"`
	DataType          string `xorm:"varchar(255)" json:"dataType,omitempty"`
	ProcessInstanceId string `xorm:"varchar(255)" json:"processInstanceId,omitempty"`
}

func (ActionResult) TableName() string {
	return "atl_actionresult"
}

func (ActionResult) KeyName() string {
	return "ActualId"
}

func (ActionResult) IdName() string {
	return baseentity.FieldName_Id
}
