package entity

import (
	entity "github.com/curltech/go-colla-core/entity"
	"time"
)

const (
	MenuType_App      = "App"
	MenuType_Module   = "Module"
	MenuType_Function = "Function"
)

const (
	ClientType_User = "User"
	ClientType_Peer = "Peer"
)

type GuiMenu struct {
	entity.StatusEntity `xorm:"extends"`
	MenuId              string     `xorm:"varchar(16)" json:"menuId,omitempty"`
	MenuType            string     `xorm:"varchar(16)" json:"menuType,omitempty"`
	Kind                string     `xorm:"varchar(16)" json:"kind,omitempty"`
	Name                string     `xorm:"varchar(32)" json:"name,omitempty"`
	Label               string     `xorm:"varchar(32)" json:"label,omitempty"`
	ClientType          string     `xorm:"varchar(32)" json:"clientType,omitempty"`
	Path                string     `xorm:"varchar(32)" json:"path,omitempty"`
	Icon                string     `xorm:"varchar(32)" json:"icon,omitempty"`
	ParentId            string     `xorm:"varchar(16)" json:"parentId,omitempty"`
	SerialId            int        `xorm:"" json:"serialId,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	Children            []*GuiMenu `xorm:"-" json:"children,omitempty"`
}

func (GuiMenu) TableName() string {
	return "gui_menu"
}

func (GuiMenu) KeyName() string {
	return "MenuId"
}

func (GuiMenu) IdName() string {
	return entity.FieldName_Id
}
