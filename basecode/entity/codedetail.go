package entity

import (
	entity "github.com/curltech/go-colla-core/entity"
	"time"
)

type CodeDetail struct {
	entity.StatusEntity `xorm:"extends"`
	BaseCodeId          string     `xorm:"varchar(16)" json:"baseCodeId,omitempty"`
	CodeDetailId        string     `xorm:"varchar(16)" json:"codeDetailId,omitempty"`
	Label               string     `xorm:"varchar(16)" json:"label,omitempty"`
	Value               string     `xorm:"varchar(32)" json:"value,omitempty"`
	Version             int        `json:"version,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	ParentId            string     `xorm:"varchar(16)" json:"parentId,omitempty"`
	SerialId            int        `json:"serialId,omitempty"`
}

func (CodeDetail) TableName() string {
	return "bas_codedetail"
}

func (CodeDetail) KeyName() string {
	return "CodeDetailId"
}

func (CodeDetail) IdName() string {
	return entity.FieldName_Id
}
