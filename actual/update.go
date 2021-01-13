package actual

import (
	"errors"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/actual/service"
	"github.com/curltech/go-colla-biz/spec"
	specentity "github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/kataras/golog"
	"strings"
)

func (this *Role) CanbeAdd(specId uint64) int {
	roleSpec, ok := this.RoleSpec.RoleSpecs[specId]
	if !ok {
		return 0
	}
	_, max := this.GetMultiplicity(specId)
	kind := roleSpec.Kind
	count := this.GetCount(kind)
	if count < max {
		return max - count
	} else {
		golog.Errorf("Max:%v,count:%v cannot add Speckind:%v", max, count, kind)
		return 0
	}
}

func (this *Role) AddRole(subRoleSpec *spec.RoleSpec, values map[string]interface{}) *Role {
	count := this.CanbeAdd(subRoleSpec.SpecId)
	if count > 0 {
		role := createRole(this.SchemaName, this, subRoleSpec, 0, values)

		return role
	}

	return nil
}

func (this *Role) AddRoles(subRoleSpec *spec.RoleSpec, count int, values []map[string]interface{}) *[]*Role {
	if count > 0 {
		roles := make([]*Role, 0)
		for i := 0; i < count; i++ {
			if i < len(values) {
				value := (values)[i]
				if value != nil {
					role := this.AddRole(subRoleSpec, value)
					roles = append(roles, role)
				}
			}
		}

		return &roles
	}

	return nil
}

/**
 * 检查属性和定义是否匹配，如果定义有增加，则增加新的属性
 * @return
 */
func (this *Role) UpdateProperty() error {
	attributeSpecs := this.RoleSpec.GetSortedAttributeSpecs(true)
	if attributeSpecs != nil && len(attributeSpecs) > 0 {
		last := this.Properties[len(this.Properties)-1]
		realIndex := (last.SerialId)*entity.AttributeSpecNumber + last.CurrentIndex
		gap := len(attributeSpecs) - realIndex
		if gap == 0 {
			return nil
		}
		if gap < 0 {
			golog.Errorf("OverAttributeSpec")

			return errors.New("OverAttributeSpec")
		}
		propertySvc := service.GetPropertyService()
		j := last.CurrentIndex
		var property *entity.Property
		index := 0
		for i := realIndex; i < len(attributeSpecs); i++ {
			attributeSpec := attributeSpecs[i]
			dataType := attributeSpec.DataType
			var defaultValue interface{}
			if attributeSpec.DefaultValue != "" {
				defaultValue = attributeSpec.DefaultValue
				if dataType != "" {
					defaultValue, _ = convert.ToObject(attributeSpec.DefaultValue, dataType)
				}
			}
			if property == nil {
				j = 0
				property = entity.NewProperty()
				property.SerialId = index
				index++
				property.DirtyFlag = baseentity.EntityState_New
				property.State = baseentity.EntityState_New
			}
			if defaultValue != nil {
				property.PutValue(j, defaultValue)
			}
			this.computeProperty(property, attributeSpec, j)
			j++

			// 在PropertyEO全部填满之后或者没有属性需要填充的时候放入角色
			if property != nil {
				if j >= entity.AttributeSpecNumber || i >= len(attributeSpecs) {
					property.CurrentIndex = j
					property.ParentId = this.Id
					property.Id = propertySvc.GetSeq()
					property.SchemaName = this.SchemaName
					property.TopId = this.TopId
					this.PutProperty(property)
					property = nil
				}
			}
		}
	}

	return nil
}

func (this *Role) setFixedValue(kind string, value interface{}) (bool, error) {
	old, err := this.getFixedValue(kind)
	if err == nil {
		if old == value {
			return false, nil
		} else {
			err = reflect.SetValue(this.FixedActual, kind, value)
			if err != nil {
				return false, err
			} else {
				internalFixedActual := this.FixedActual.(specentity.IInternalFixedActual)
				internalFixedActual.UpdateDirtyFlag(baseentity.EntityState_Modified)
				baseEntity := this.FixedActual.(baseentity.IBaseEntity)
				baseEntity.UpdateState(baseentity.EntityState_Modified)

				this.UpdateDirtyFlag(baseentity.EntityState_Modified)
				this.UpdateState(baseentity.EntityState_Modified)
				this.updateState(kind, baseentity.EntityState_Modified)

				return true, nil
			}
		}
	} else {
		return false, err
	}
}

/**
 * @param kind
 * @param value
 * @return
 * @Description: 根据属性kind和位置设置属性值 先设置角色对象的值， 再设置静态属性的值，不成功设置到动态属性的值，
 *               设置成功，更新角色脏标志
 */
func (this *Role) SetPropertyValue(kind string, value interface{}) (bool, error) {
	if kind != "" {
		ok, err := this.setFixedValue(kind, value)
		if err == nil {
			return ok, nil
		} else {
			for _, p := range this.Properties {
				if p.Contain(kind) {
					ok, err = p.SetValue(kind, value)
					if ok {
						this.UpdateDirtyFlag(baseentity.EntityState_Modified)
						this.UpdateState(baseentity.EntityState_Modified)
						p.UpdateDirtyFlag(baseentity.EntityState_Modified)
						p.UpdateState(baseentity.EntityState_Modified)
						this.updateState(kind, baseentity.EntityState_Modified)

						return true, err
					} else {
						return false, err
					}
				}
			}
		}
	}

	return false, errors.New("NoKind")
}

func (this *Role) updateState(kind string, state string) {
	// 现在是新的
	old, ok := this.PropertyStates[kind]
	if ok && baseentity.EntityState_New == old && baseentity.EntityState_Modified == state {
		return
	}

	this.PropertyStates[kind] = state
}

func (this *Role) SetValue(path string, value interface{}) (bool, error) {
	kind := path
	var role = this
	if strings.Contains(path, ".") {
		rolePath := GetRolePath(path)
		v, _, ok := this.Find(rolePath)
		if ok && v != nil {
			role = v.(*Role)
			kind = GetLastKind(path)
		}
	}
	if kind != "" {
		return role.SetPropertyValue(kind, value)
	}

	return false, errors.New("")
}
