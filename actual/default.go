package actual

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/spec"
	specentity "github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/cache"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/reflect"
	"sync"
)

type Role struct {
	entity.Role `json:",omitempty"`
	RoleSpec    *spec.RoleSpec `json:"-,omitempty"`
	// 以固定属性存在的直接下属属性，在非关系数据库（bolt）的时候所有的字段可以动态地全部存放，Properties将为空
	FixedActual interface{} `json:",omitempty"`
	// 所有下属属性的前端交互的脏标志
	PropertyStates map[string]string `json:"-,omitempty"`
	// 所有的直接下属角色实例，kind为键值的数组
	Roles       map[string][]*Role            `json:",omitempty"`
	Connections map[uint64]*entity.Connection `json:",omitempty"`
	// 所有的直接下属属性元素实例
	Properties []*entity.Property `json:",omitempty"`
	// 所有的直接下属行为结果实例
	ActionResults map[string]*entity.ActionResult `json:",omitempty"`
	/**
	删除的role，保存数据库结束后清除
	*/
	DeleteRoles map[uint64]*Role `json:",omitempty"`
	/**
	删除的角色的编号，用于通知前端，返回前端后清除
	*/
	DeleteActuals []uint64   `json:",omitempty"`
	ParentRole    *Role      `json:"-,omitempty"`
	mux           sync.Mutex `json:"-,omitempty"`
}

var MemCache = cache.NewMemCache("actual", 1, 10)

func (this *Role) PutRole(role *Role) error {
	actualId := role.Id
	parentId := role.ParentId
	if this.Id != actualId && this.Id == parentId {
		role.ParentRole = this
		kind := role.Kind
		if this.Roles == nil {
			this.Roles = make(map[string][]*Role, 0)
		}
		value, ok := this.Roles[kind]
		if !ok {
			value = make([]*Role, 0)
		}
		this.computePath(role, len(value))
		value = append(value, role)
		this.Roles[kind] = value
		role.TopId = this.TopId
	} else {
		logger.Errorf("Role actualId:%v is correct", actualId)

		return errors.New("")
	}

	return nil
}

func (this *Role) computePath(role *Role, pos int) {
	if this.Path == "" {
		role.Path = role.Kind
	} else {
		role.Path = this.Path + "." + role.Kind
	}
	if pos > 0 {
		role.Path = fmt.Sprintf(role.Path+"[%v]", pos)
	}
}

func (this *Role) PutFixedActual(fixedActual interface{}) error {
	v, err := reflect.GetValue(fixedActual, baseentity.FieldName_ParentId)
	if err != nil {
		return errors.New("")
	}
	parentId := v.(uint64)
	if this.Id == parentId {
		reflect.SetValue(fixedActual, "SpecId", this.SpecId)
		reflect.SetValue(fixedActual, "Kind", this.Kind)
		reflect.SetValue(fixedActual, baseentity.FieldName_SchemaName, this.SchemaName)
		reflect.SetValue(fixedActual, "TopId", this.TopId)
		this.FixedActual = fixedActual
	} else {
		logger.Errorf("fixedActual parentId:%v is correct", parentId)

		return errors.New("")
	}

	return nil
}

func (this *Role) PutConnection(conn *entity.Connection) error {
	parentId := conn.ParentId
	actualId := conn.ActualId
	if this.Id != actualId && this.Id == parentId {
		conn.TopId = this.TopId
		conn.SchemaName = this.SchemaName
		if this.Connections == nil {
			this.Connections = make(map[uint64]*entity.Connection, 0)
		}
		conn, ok := this.Connections[actualId]
		if ok {
			logger.Errorf("Repeat Connection")

			return errors.New("Repeat Connection")
		}
		this.Connections[actualId] = conn
	} else {
		logger.Errorf("Connection actualId:%v is correct", actualId)

		return errors.New("")
	}

	return nil
}

func (this *Role) PutActionResult(actionResult *entity.ActionResult) error {
	actualId := actionResult.Id
	parentId := actionResult.ParentId
	kind := actionResult.Kind
	if this.Id != actualId && this.Id == parentId {
		if this.ActionResults == nil {
			this.ActionResults = make(map[string]*entity.ActionResult, 0)
		}
		_, ok := this.ActionResults[kind]
		if ok {
			logger.Errorf("kind:" + kind + ";Same ActionResultEO will be replaced!Please check it")

			return errors.New("Repeat ActionResult")
		}
		this.ActionResults[kind] = actionResult
		actionResult.TopId = this.TopId
	}

	return errors.New("ErrorActualId")
}

/**
把属性放入角色中，属性可能只有值，无specId，计算属性也未填充
*/
func (this *Role) PutProperty(property *entity.Property) (err error) {
	actualId := property.Id
	parentId := property.ParentId
	if this.Id != actualId && this.Id == parentId {
		property.SetTopId(this.TopId)
		var attributeSpecs = this.RoleSpec.GetSortedAttributeSpecs(true)
		//遍历动态属性
		for i := 0; i < entity.AttributeSpecNumber; i++ {
			serialId := property.SerialId
			realIndex := serialId*entity.AttributeSpecNumber + i
			if realIndex < len(attributeSpecs) {
				attributeSpec := attributeSpecs[realIndex]
				if attributeSpec != nil {
					v, err := property.Get(entity.AttributeType_Kind, i)
					if err == nil {
						kind := v.(string)
						if kind == "" {
							this.computeProperty(property, attributeSpec, i)
						}
					} else {
						this.computeProperty(property, attributeSpec, i)
					}
				}
			} else { // 真实位置已经超过了属性的个数，不用再取了(attributeSpecs))
				break
			}
		}
		if this.Properties == nil {
			this.Properties = make([]*entity.Property, 0)
		}
		if len(this.Properties) == property.SerialId {
			this.Properties = append(this.Properties, property)
		} else {
			logger.Errorf("ErrorSerialId")
			err = errors.New("ErrorSerialId")
		}
	}

	return err
}

/**
把动态属性转换成静态属性
*/
func (this *Role) transferProperty() error {
	var err error
	fixedService := this.RoleSpec.GetFixedService()
	if fixedService != nil && this.FixedActual == nil {
		fixedActual, err := fixedService.NewEntity(nil)
		if err == nil {
			baseEntity := fixedActual.(baseentity.IBaseEntity)
			baseEntity.UpdateState(baseentity.EntityState_New)
			internalFixedActual := fixedActual.(specentity.IInternalFixedActual)
			internalFixedActual.UpdateDirtyFlag(baseentity.EntityState_New)
			this.PutFixedActual(fixedActual)

			if this.Properties != nil && len(this.Properties) > 0 {
				for _, p := range this.Properties {
					for i := 0; i < p.CurrentIndex; i++ {
						v, err := p.Get(entity.AttributeType_Kind, i)
						if err == nil {
							kind := v.(string)
							alias := this.getAlias(kind)
							if alias != "" {
								v, err = p.GetValue(kind)
								err = reflect.SetValue(fixedActual, alias, v)
							}
						}
					}
					//需要删除吗？
					p.DirtyFlag = baseentity.EntityState_Deleted
				}
			}
		}
	}

	return err
}

/**
填充属性的计算值和脏标志，index代表位置，假设此时可能property的specId无值，serialId有值，value有值
*/
func (this *Role) computeProperty(property *entity.Property, attributeSpec *specentity.AttributeSpec, index int) error {
	//serialId:=property.SerialId
	attributeSpecId := attributeSpec.SpecId
	var err error
	if attributeSpecId > 0 {
		// 检查此位置的specId
		v, err := property.Get(entity.AttributeType_SpecId, index)
		if err == nil {
			specId := v.(uint64)
			// 如果属性此位置的specId已经存在，而且不等于要设置的属性specId，出错
			if specId > 0 {
				if attributeSpecId != specId {
					logger.Errorf("property's currentAttributeSpecId:%v,but attributeSpecId:%v", specId, attributeSpecId)
					err = errors.New("ConflictSpecId")
				}
				attributeSpecId = specId
			}
		}
		property.Set(entity.AttributeType_SpecId, index, attributeSpecId)
		// 检查此位置的kind
		var kind string
		v, err = property.Get(entity.AttributeType_Kind, index)
		if err == nil {
			kind = v.(string)
		}
		if kind == "" {
			kind = attributeSpec.Kind
			// 属性定义的kind不为空，则根据属性定义设置kind，dataType，pattern
			if kind != "" {
				property.Set(entity.AttributeType_Kind, index, kind)
				property.Set(entity.AttributeType_DataType, index, attributeSpec.DataType)
				property.Set(entity.AttributeType_Pattern, index, attributeSpec.Pattern)

				// 计算自己的路径
				path := this.Path + "." + attributeSpec.Kind
				// 设置创建的属性的路径，目前不支持多个同名属性，也没有必要
				property.Set(entity.AttributeType_Path, index, path)

				_, ok := this.PropertyStates[kind]
				if ok {
					logger.Errorf("kind:%v repeat PropertyKind error!", kind)
				}
				state := property.State
				if state == "" {
					state = baseentity.EntityState_None
				}
				if this.PropertyStates == nil {
					this.PropertyStates = make(map[string]string, 0)
				}
				this.updateState(kind, state)
			} else { // 定义的kind为空，出错
				logger.Errorf("AttributeSpecId:%v has no kind!", attributeSpecId)
				err = errors.New("NoKind")
			}
		}
	}

	return err
}

func (this *Role) getAlias(kind string) string {
	if this.RoleSpec.FixedAttributes != nil && this.RoleSpec != nil {
		alias, _ := this.RoleSpec.FixedAttributes[kind]

		return alias
	}

	return ""
}

func (this *Role) getTopRole() *Role {
	if this.ParentRole == nil {
		return this
	}
	return this.ParentRole.getTopRole()
}

func (this *Role) GetCount(kind string) int {
	var count = 0
	v, ok := this.findChildren(specentity.SpecType_Role, kind, -1)
	if ok {
		rs := v.([]*Role)
		if rs != nil {
			count = len(rs)
		}
	}

	return count
}

func (this *Role) GetMultiplicity(specId uint64) (int, int) {
	conSpec := this.RoleSpec.GetChildConnectionSpec(specId)
	if conSpec != nil {
		return conSpec.Minmium, conSpec.Maxmium
	}

	return 0, 0
}
