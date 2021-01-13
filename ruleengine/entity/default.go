package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

type RuleDefinition struct {
	entity.StatusEntity `xorm:"extends"`
	Kind                string     `xorm:"varchar(255)" json:",omitempty"`
	Name                string     `xorm:"varchar(255)" json:",omitempty"`
	PackageName         string     `xorm:"varchar(255)" json:",omitempty"`
	ExecuteType         string     `xorm:"varchar(255)" json:",omitempty"`
	ContentId           string     `xorm:"varchar(255)" json:",omitempty"`
	RemoteAddress       string     `xorm:"varchar(255)" json:",omitempty"`
	RawContent          string     `xorm:"varchar(32000)" json:",omitempty"`
	Version             string     `xorm:"varchar(255)" json:",omitempty"`
	StartDate           *time.Time `json:",omitempty"`
	EndDate             *time.Time `json:",omitempty"`
}
