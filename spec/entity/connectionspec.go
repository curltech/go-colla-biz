package entity

import (
	"github.com/curltech/go-colla-core/entity"
	"time"
)

/**
 * The enum Temp password authenticate status.
 */

const (
	//缺省关系，整体和部分的强包含关系，关联的子角色是父角色的一部分，需要一起创建一起加载
	RelationType_Dependency string = "Dependency"
	//整体和部分的强包含关系，关联的子角色是父角色的一部分，需要一起创建但可以不一起加载，在需要的时候加载
	RelationType_Composition string = "Composition"
	//整体和部分的弱包含关系，通过connection对象关联，创建的时候创建，不能一起复制，比如团险中的人员清单，复制的时候复制connection对象
	RelationType_Aggregation string = "Aggregation"
	//对象之间松散的关联关系，比如外部客户信息，关联到CIF，或者联系其他的角色对象，需要通过connection对象关联
	RelationType_Association    string = "Association"
	RelationType_Realization    string = "Realization"
	RelationType_Generalization string = "Generalization"
)

type ConnectionSpec struct {
	entity.StatusEntity `xorm:"extends"`
	SpecType            string     `xorm:"varchar(255)" json:"specType,omitempty"`
	RelationType        string     `xorm:"varchar(255)" json:"relationType,omitempty"`
	ParentSpecId        uint64     `json:"parentSpecId,omitempty"`
	SubSpecId           uint64     `json:"subSpecId,omitempty"`
	Maxmium             int        `json:"maxmium,omitempty"`
	Minmium             int        `json:"minmium,omitempty"`
	Version             uint64     `json:"version,omitempty"`
	StartDate           *time.Time `json:"startDate,omitempty"`
	EndDate             *time.Time `json:"endDate,omitempty"`
	SerialId            uint64     `json:"serialId,omitempty"`
	/**
	 * 缺省的实例创建数
	 */
	BuildNum int `json:"buildNum,omitempty"`
	/**
	 * 只用于role与role之间的连接，表示连接的下的role的装载模式，-1表示不装载，0表示装载，正整数表示缺省装载几行
	 * 本参数用于一个角色的实例太多，不能一次全装入内存的情况，比如团单的人员清单，按日或者按月的消费记录等等
	 * 或者某角色不適合与父親一起裝載，比如是独立的数据源或者远程服务的時候
	 * 或者是本地可以一次性加载，但是作为产品组合，希望单独加载的情况，将影响RoleEO的loadType
	 */
	LoadNum int `json:"loadNum,omitempty"`
}

func (ConnectionSpec) TableName() string {
	return "spec_connection"
}

func (ConnectionSpec) KeyName() string {
	return "SpecId"
}

func (ConnectionSpec) IdName() string {
	return entity.FieldName_Id
}
