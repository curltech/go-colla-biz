package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type Resource struct {
	entity.StatusEntity `xorm:"extends"`
	ResourceId          string     `xorm:"varchar(255)" json:"resourceId,omitempty"`
	ResourceType        string     `xorm:"varchar(255)" json:"resourceTypeResourceType,omitempty"`
	Name                string     `xorm:"varchar(255)" json:"name,omitempty"`
	Path                string     `xorm:"varchar(255)" json:"path,omitempty"`
	ParentId            string     `xorm:"varchar(255)" json:"parentId,omitempty"`
	Value               string     `xorm:"varchar(255)" json:"value,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	AccessMode          string     `xorm:"varchar(255)" json:"accessMode,omitempty"`
	Priority            int64      `json:"priority,omitempty"`
	OwnedStructureId    string     `xorm:"varchar(255)" json:"ownedStructureId,omitempty"`
	OwnedStructureName  string     `xorm:"varchar(255)" json:"ownedStructureName,omitempty"`
	OwnedStructurePath  string     `xorm:"varchar(255)" json:"ownedStructurePath,omitempty"`
}

func (Resource) TableName() string {
	return "rbac_resource"
}

func (Resource) KeyName() string {
	return "ResourceId"
}

func (Resource) IdName() string {
	return entity.FieldName_Id
}
