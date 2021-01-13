package entity

import (
	"github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
	"time"
)

type Role struct {
	entity.InternalFixedActual `xorm:"extends"`
	/**
	 * 静态对象的类型编号，对应
	 */
	FixedSpecId uint64 `json:"fixedSpecId,omitempty"`
	/**
	 * 连接静态对象的外键，记录静态对象的某字段的值，如果是静态對象的主键，那么角色对应的是单条记录， 如果不是，可能对应多条静态对象记录
	 */
	FixedActualId uint64 `json:"fixedActualId,omitempty"`
	/**
	 * 装载配置参数，-1表示不装载，0表示正常装载，不装载表示模型之间的一种弱关系，父亲角色中不放入子角色的对象，不一起保存，也不一起加载
	 */
	LoadNum int `json:"loadNum,omitempty"`
	// 角色定义的有效日期，表示角色事实遵循的角色定义的有效日期
	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`
	Version       int        `json:"version,omitempty"`
	// 角色的开始有效日期
	StartDate *time.Time `json:"startDate,omitempty"`
	// 角色的有效结束日期，表示角色事实被后一个版本替换掉的日期
	EndDate *time.Time `json:"endDate,omitempty"`
	// 角色的位置，指多个角色实例的时候的次序，不代表在数组中的实际位置
	//Position int `json:"position,omitempty"`
	// 外部引用编号
	ReferenceId string `xorm:"varchar(255)" json:"referenceId,omitempty"`
	FirstId     uint64 `json:"firstId,omitempty"`
	PreviousId  uint64 `json:"previousId,omitempty"`
	Path        string `xorm:"-" json:"path,omitempty"`
}

func (Role) TableName() string {
	return "atl_role"
}

func (Role) KeyName() string {
	return "ActualId"
}

func (Role) IdName() string {
	return baseentity.FieldName_Id
}
