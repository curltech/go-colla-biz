package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type Group struct {
	entity.StatusEntity `xorm:"extends"`
	GroupId             string     `xorm:"varchar(255)" json:"groupId,omitempty"`
	UserId              string     `xorm:"varchar(255)" json:"userId,omitempty"`
	LoginName           string     `xorm:"varchar(255)" json:"loginName,omitempty"`
	RoleId              string     `xorm:"varchar(255)" json:"roleId,omitempty"`
	RoleName            string     `xorm:"varchar(255)" json:"roleName,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
}

func (Group) TableName() string {
	return "rbac_group"
}

func (Group) KeyName() string {
	return "GroupId"
}

func (Group) IdName() string {
	return entity.FieldName_Id
}
