package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

const (
	SpecType_Role       string = "Role"
	SpecType_Action     string = "Action"
	SpecType_Property   string = "Property"
	SpecType_Connection string = "Connection"
)

const (
	FieldName_SpecId = "SpecId"
	FieldName_Kind   = "Kind"
	FieldName_TopId  = "TopId"
)

type Specification struct {
	entity.StatusEntity `xorm:"extends"`
	SpecId              uint64     `json:"specId,omitempty"`
	Kind                string     `xorm:"varchar(255)" json:"kind,omitempty"`
	Name                string     `xorm:"varchar(255)" json:"name,omitempty"`
	Description         string     `xorm:"varchar(255)" json:"description,omitempty"`
	Version             uint64     `json:"version,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	SerialId            uint64     `json:"serialId,omitempty"`
}

/**
 * 自己定义的所有事实的祖先类， 包括静态对象，
 *
 * 外部的静态对象（包括外联和共享）不要继承本类
 */
type InternalFixedActual struct {
	entity.BaseEntity `xorm:"extends"`
	ParentId          uint64   `xorm:"index" json:"parentId,omitempty"`
	SpecId            uint64   `json:"specId,omitempty"`
	Kind              string   `xorm:"-" json:"kind,omitempty"`
	TopId             uint64   `json:"topId,omitempty"`
	Revision          uint64   `json:"revision,omitempty"`
	SchemaName        string   `xorm:"varchar(255)" json:"schemaName,omitempty"`
	DirtyFlag         string   `xorm:"-" json:"dirtyFlag,omitempty"`
	dirtyFlagFields   []string `xorm:"-" json:"-"`
}

type IInternalFixedActual interface {
	UpdateDirtyFlag(dirtyFlag string)
	SetParentId(parentId uint64)
	SetTopId(topId uint64)
}

func (this *InternalFixedActual) SetParentId(parentId uint64) {
	this.ParentId = parentId
}

func (this *InternalFixedActual) SetTopId(topId uint64) {
	this.TopId = topId
}

func (this *InternalFixedActual) UpdateDirtyFlag(dirtyFlag string) {
	// 现在是新的
	if entity.EntityState_New == this.DirtyFlag && entity.EntityState_Modified == dirtyFlag {
		return
	}

	this.DirtyFlag = dirtyFlag
}
