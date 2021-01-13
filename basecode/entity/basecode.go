package entity

import (
	entity "github.com/curltech/go-colla-core/entity"
	"time"
)

type BaseCode struct {
	entity.StatusEntity `xorm:"extends"`
	BaseCodeId          string        `xorm:"varchar(16)" json:"baseCodeId,omitempty"`
	BaseCodeType        string        `xorm:"varchar(16)" json:"baseCodeType,omitempty"`
	Kind                string        `xorm:"varchar(16)" json:"kind,omitempty"`
	Name                string        `xorm:"varchar(32)" json:"name,omitempty"`
	Label               string        `xorm:"varchar(32)" json:"label,omitempty"`
	Version             int           `json:"version,omitempty"`
	StartDate           *time.Time    `json:"startDate,omitempty"`
	EndDate             *time.Time    `json:"endDate,omitempty"`
	CodeDetails         []*CodeDetail `xorm:"-" json:"codeDetails,omitempty"`
}

func (BaseCode) TableName() string {
	return "bas_basecode"
}

func (BaseCode) KeyName() string {
	return "BaseCodeId"
}

func (BaseCode) IdName() string {
	return entity.FieldName_Id
}
