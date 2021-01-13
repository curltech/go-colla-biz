package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type UserProfile struct {
	entity.StatusEntity `xorm:"extends"`
	UserId              string     `xorm:"varchar(255)" json:"userId,omitempty"`
	UserType            string     `xorm:"varchar(255)" json:"userType,omitempty"`
	LastLoginDate       *time.Time `json:"lastLoginDate,omitempty"`
	LoginCount          int64      `json:"loginCount,omitempty"`
	EmailUser           string     `xorm:"varchar(255)" json:"emailUser,omitempty"`
	EmailPassword       string     `xorm:"varchar(255)" json:"emailPassword,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	RegisterIp          string     `xorm:"varchar(255)" json:"registerIp,omitempty"`
	RegisterSource      string     `xorm:"varchar(255)" json:"registerSource,omitempty"`
	DefaultPage         string     `xorm:"varchar(255)" json:"defaultPage,omitempty"`
	TempPassword        string     `xorm:"varchar(255)" json:"tempPassword,omitempty"`
	TempPasswordTime    *time.Time `json:"tempPasswordTime,omitempty"`
	TempPasswordCount   int64      `json:"tempPasswordCount,omitempty"`
	ActiveStatus        string     `xorm:"varchar(255)" json:"activeStatus,omitempty"`
}

func (UserProfile) TableName() string {
	return "rbac_userprofile"
}

func (UserProfile) KeyName() string {
	return "UserId"
}

func (UserProfile) IdName() string {
	return entity.FieldName_Id
}
