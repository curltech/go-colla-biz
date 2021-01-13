package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type Role struct {
	entity.StatusEntity `xorm:"extends"`
	RoleId              string     `xorm:"varchar(255)" json:"roleId,omitempty"`
	Kind                string     `xorm:"varchar(255)" json:"kind,omitempty"`
	Name                string     `xorm:"varchar(255)" json:"name,omitempty"`
	Label               string     `xorm:"varchar(255)" json:"label,omitempty"`
	ParentId            string     `xorm:"varchar(255)" json:"parentId,omitempty"`
	Description         string     `xorm:"varchar(255)" json:"description,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	DefaultPage         string     `xorm:"varchar(255)" json:"defaultPage,omitempty"`
	OwnedStructureId    string     `xorm:"varchar(255)" json:"ownedStructureId,omitempty"`
	OwnedStructureName  string     `xorm:"varchar(255)" json:"ownedStructureName,omitempty"`
	OwnedStructurePath  string     `xorm:"varchar(255)" json:"ownedStructurePath,omitempty"`
}

func (Role) TableName() string {
	return "rbac_role"
}

func (Role) KeyName() string {
	return "RoleId"
}

func (Role) IdName() string {
	return entity.FieldName_Id
}
