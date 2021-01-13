package entity

import (
	baseentity "github.com/curltech/go-colla-core/entity"
)

type RoleSpec struct {
	Specification `xorm:"extends"`
	/**
	 * 对应静态对象的外部名称(kind)
	 */
	ExternalKind string `xorm:"varchar(255)" json:"externalKind,omitempty"`
	/**
	 * 对应内部静态对象的编号(specId)，表示角色定义的静态属性是产品定义的一部分，
	 *
	 * 在获取保单是会一起获取
	 */
	FixedSpecId uint64 `json:"fixedSpecId,omitempty"`
	/**
	 * 对应外挂静态对象的编号(specId)，表示角色的主键会外部关联外部静态对象，
	 * 适用的场景是团单下有很多被保人，比如1000个，不能作为保单的对象全部取出来，这样性能会成问题，
	 * 所以采用外部关联的办法，分页取出外部关联的静态对象，外部静态对象的实例数很多
	 */
	ExternalFixedSpecId uint64 `json:"externalFixedSpecId,omitempty"`
}

func (RoleSpec) TableName() string {
	return "spec_role"
}

func (RoleSpec) KeyName() string {
	return "SpecId"
}

func (RoleSpec) IdName() string {
	return baseentity.FieldName_Id
}
