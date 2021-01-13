package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type Permission struct {
	entity.StatusEntity `xorm:"extends"`
	ActorId             string     `xorm:"varchar(255)" json:"actorId,omitempty"`
	ActorType           string     `xorm:"varchar(255)" json:"actorType,omitempty"`
	ActorKind           string     `xorm:"varchar(255)" json:"actorKind,omitempty"`
	ResourceId          string     `xorm:"varchar(255)" json:"resourceId,omitempty"`
	ResourceType        string     `xorm:"varchar(255)" json:"resourceType,omitempty"`
	ResourceKind        string     `xorm:"varchar(255)" json:"resourceKind,omitempty"`
	ResourceName        string     `xorm:"varchar(255)" json:"resourceName,omitempty"`
	ResourcePath        string     `xorm:"varchar(255)" json:"resourcePath,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	AccessMode          string     `xorm:"varchar(255)" json:"accessMode,omitempty"`
	Priority            int64      `json:"priority,omitempty"`
	Path                string     `xorm:"varchar(255)" json:"path,omitempty"`
	Value               string     `xorm:"varchar(255)" json:"value,omitempty"`
}

func (Permission) TableName() string {
	return "rbac_permission"
}

func (Permission) KeyName() string {
	return entity.FieldName_Id
}

func (Permission) IdName() string {
	return entity.FieldName_Id
}
