package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

const (
	UseType_Share     string = "Share"
	UseType_Exclusive string = "Exclusive"
	UseType_HalfShare string = "HalfShare"
	UseType_External  string = "External"
)

/**
 * 静态定义对象是独占还是共享，S共享，M独占，H半独占共享，E是外部表，缺省是独占
 * 独占意味着静态表的数据就是保单数据的一部分，随着保单数据CRUD，另一份保单实例不能访问本保单实例的数据
 * 共享意味着静态表的数据不是保单数据的一部分，只是一个引用，如果已经存在则不用增加，需要和保单一起RU，
 * 可以和另一份保单实例共享静态数据，因此不能随便删除，典型的例子是CIF的party
 * 半共享模式意味着静态表实例是随着保单实例一起创建的，也可以修改，可以单独操作，不属于保单数据的一部分，因此不能随便删除，CRU
 * 外部静态表表示这是一份外部管理的数据，只是在role中纪录了一个外键连接，缺省不被自动装载，但是需要的时候可以手工装载，用于处理数据量很大的场景，
 * 不能在装载保单的时候一起把静态表的数据装载，比如团险的被保人清单
 *
 * @author liu
 */
type FixedRoleSpec struct {
	entity.StatusEntity `xorm:"extends"`
	SpecId              uint64    `json:"specId,omitempty"`
	Kind                string    `xorm:"varchar(255)" json:"kind,omitempty"`
	Name                string    `xorm:"varchar(255)" json:"name,omitempty"`
	Description         string    `xorm:"varchar(255)" json:"description,omitempty"`
	Version             uint64    `json:"version,omitempty"`
	StartDate           time.Time `json:"startDate,omitempty"`
	EndDate             time.Time `json:"endDate,omitempty"`
	/**
	 * 静态实体的全类名
	 *
	 * @param fixedName
	 */
	FixedName        string `xorm:"varchar(255)" json:"fixedName,omitempty"`
	FixedServiceName string `xorm:"varchar(255)" json:"fixedServiceName,omitempty"`
}

func (FixedRoleSpec) TableName() string {
	return "spec_fixedrole"
}

func (FixedRoleSpec) KeyName() string {
	return "SpecId"
}

func (FixedRoleSpec) IdName() string {
	return entity.FieldName_Id
}
