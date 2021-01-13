package businessentity

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-biz/actual"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/spec"
	entity2 "github.com/curltech/go-colla-biz/spec/entity"
	entity3 "github.com/curltech/go-colla-core/entity"
	baseerror "github.com/curltech/go-colla-core/error"
	"github.com/curltech/go-colla-core/util/convert"
	"time"
)

type BusinessEntity struct {
	*actual.Role
}

func (this *BusinessEntity) GetFloatProperty(path string) float64 {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(float64)
	}

	return 0
}

func (this *BusinessEntity) GetBoolProperty(path string) bool {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(bool)
	}

	return false
}

func (this *BusinessEntity) GetUint64Property(path string) uint64 {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(uint64)
	}

	return 0
}

func (this *BusinessEntity) GetInt64Property(path string) int64 {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(int64)
	}

	return 0
}

func (this *BusinessEntity) GetUintProperty(path string) uint {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(uint)
	}

	return 0
}

func (this *BusinessEntity) GetIntProperty(path string) int {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(int)
	}

	return 0
}

func (this *BusinessEntity) GetStringProperty(path string) string {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(string)
	}

	return ""
}

func (this *BusinessEntity) GetTimeProperty(path string) *time.Time {
	v, err := this.Role.GetValue(path)
	if v != nil && err == nil {
		return v.(*time.Time)
	}

	return nil
}

func (this *BusinessEntity) GetActual(isState bool) map[string]interface{} {
	actual := this.Role.GetActual(isState)

	return actual
}

func (this *BusinessEntity) Update(actual map[string]interface{}, isState bool, isCheckActualId bool) map[uint64]map[string]error {
	errs := this.Role.SetActual(actual, isState, isCheckActualId)

	return errs
}

/**
在BusinessEntity未被取出或者装载的情况下的操作
*/
type BusinessEntityService struct {
}

var businessEntityService = &BusinessEntityService{}

func GetBusinessEntityService() *BusinessEntityService {
	return businessEntityService
}

func (this *BusinessEntityService) Create(schemaName string, specId uint64, effectiveDate *time.Time) (*BusinessEntity, error) {
	roleSpec, _ := spec.GetMetaDefinition().GetRoleSpec(specId, effectiveDate)
	if roleSpec != nil {
		role := actual.Create(schemaName, roleSpec)
		if role != nil {
			return &BusinessEntity{Role: role}, nil
		}
	} else {
		return nil, errors.New(baseerror.Error_NoSpecId)
	}
	return nil, nil
}

func (this *BusinessEntityService) Load(schemaName string, id uint64) *BusinessEntity {
	role := actual.Load(schemaName, id)
	if role != nil {
		return &BusinessEntity{Role: role}
	}
	return nil
}

func (this *BusinessEntityService) Get(schemaName string, id uint64) *BusinessEntity {
	role := actual.Get(schemaName, id)
	if role != nil {
		return &BusinessEntity{Role: role}
	}
	return nil
}

/**
装载实例树，逐条删除，速度较慢
*/
func (this *BusinessEntityService) Delete(schemaName string, id uint64) (int64, error) {
	ids := make([]uint64, 1)
	ids[0] = id

	return actual.Delete(schemaName, ids)
}

func (this *BusinessEntityService) GetActual(schemaName string, id uint64, isState bool) map[string]interface{} {
	role := actual.Get(schemaName, id)
	if role != nil {
		be := &BusinessEntity{Role: role}
		return be.GetActual(isState)
	}
	return nil
}

func (this *BusinessEntityService) Save(schemaName string, id uint64) (*BusinessEntity, error) {
	be := this.Get(schemaName, id)
	if be != nil {
		_, err := be.Role.SaveAll()
		if err == nil {
			return be, nil
		} else {
			return be, err
		}
	}
	return nil, errors.New(baseerror.Error_NotFound)
}

func (this *BusinessEntityService) Version(schemaName string, id uint64) (*BusinessEntity, error) {
	role := actual.Get(schemaName, id)
	if role != nil {
		role = role.Version()
		if role != nil {
			be := &BusinessEntity{Role: role}

			return be, nil
		}
	}
	return nil, errors.New(baseerror.Error_NotFound)
}

func (this *BusinessEntityService) Find(schemaName string, id uint64, path string, isState bool) (interface{}, string, bool) {
	role := actual.Get(schemaName, id)
	if role != nil {
		v, specType, ok := role.Find(path)
		if ok {
			if specType == entity2.SpecType_Role {
				role := v.(*actual.Role)
				be := &BusinessEntity{Role: role}

				return be, specType, ok
			} else {

				return v, specType, ok
			}
		}

	}
	return nil, "", false
}

func (this *BusinessEntityService) AddRole(schemaName string, id uint64, parentId uint64, kind string, values map[string]interface{}) (*BusinessEntity, *actual.Role) {
	be := this.Get(schemaName, id)
	if be != nil {
		parent := be.Role.Get(parentId)
		if parent != nil {
			roleSpec := parent.RoleSpec.GetChildRoleSpec(kind)
			if roleSpec != nil {
				child := parent.AddRole(roleSpec, values)

				return be, child
			}
		}
	}

	return nil, nil
}

func (this *BusinessEntityService) LoadRole(schemaName string, id uint64, parentId uint64, kind string) (*BusinessEntity, []*actual.Role) {
	be := this.Get(schemaName, id)
	if be != nil {
		parent := be.Role.Get(parentId)
		if parent != nil {
			roleSpec := parent.RoleSpec.GetChildRoleSpec(kind)
			if roleSpec != nil {
				child := parent.Load(kind)

				return be, child
			}
		}
	}

	return nil, nil
}

func (this *BusinessEntityService) RemoveRole(schemaName string, id uint64, roleId uint64, path string) (*BusinessEntity, error) {
	be := this.Get(schemaName, id)
	if be != nil {
		var role *actual.Role
		if roleId > 0 {
			role = be.Role.Get(roleId)
		} else if path != "" {
			v, specType, ok := be.Role.Find(path)
			if ok && specType == entity2.SpecType_Role {
				role = v.(*actual.Role)
			}
		}
		if role != nil {
			ok := role.ParentRole.RemoveRole(role)
			if ok {
				return be, nil
			} else {
				return be, errors.New(baseerror.Error_RemoveFail)
			}
		} else {
			return be, errors.New(baseerror.Error_NotExist)
		}
	}

	return nil, errors.New(baseerror.Error_NotFound)
}

func (this *BusinessEntityService) Update(actual map[string]interface{}) (*BusinessEntity, map[uint64]map[string]error) {
	var err error
	var schemaName string
	v, ok := actual[entity3.JsonFieldName_SchemaName]
	if ok {
		schemaName = v.(string)
	}
	var id uint64
	v, ok = actual[entity3.JsonFieldName_Id]
	if ok {
		v, err = convert.ToObject(fmt.Sprintf("%v", v), entity2.DataType_Uint64)
		if err == nil {
			id = v.(uint64)
		}
	}
	if id > 0 {
		be := this.Get(schemaName, id)
		if be != nil {
			errs := be.Update(actual, true, true)

			return be, errs
		}
	}

	return nil, nil
}

/**
根据映射中的id,parentId,path等字段获取角色和需要修改数据的子角色，并修改数据values映射
*/
func (this *BusinessEntityService) SetValue(data map[string]interface{}) (*BusinessEntity, error) {
	var err error
	var schemaName string
	v, ok := data[entity3.JsonFieldName_SchemaName]
	if ok {
		schemaName = v.(string)
	}
	var id uint64
	v, ok = data[entity3.JsonFieldName_TopId]
	if ok {
		v, err = convert.ToObject(fmt.Sprintf("%v", v), entity2.DataType_Uint64)
		if err == nil {
			id = v.(uint64)
		}
		delete(data, entity3.JsonFieldName_TopId)
	}
	var roleId uint64
	v, ok = data[entity3.JsonFieldName_Id]
	if ok {
		v, err = convert.ToObject(fmt.Sprintf("%v", v), entity2.DataType_Uint64)
		if err == nil {
			roleId = v.(uint64)
		}
		delete(data, entity3.JsonFieldName_Id)
	}
	var path string
	v, ok = data[entity3.JsonFieldName_Path]
	if ok {
		path = v.(string)
		delete(data, entity3.JsonFieldName_Path)
	}
	be := this.Get(schemaName, id)
	if be != nil {
		var values = data
		var role *actual.Role
		if roleId > 0 {
			role = be.Role.Get(roleId)
		}
		if role == nil && path != "" {
			var specType string
			v, specType, ok = be.Role.Find(path)
			if ok && specType == entity2.SpecType_Role {
				role, ok = v.(*actual.Role)
			}
		}
		if role == nil {
			return be, errors.New(baseerror.Error_NotExist)
		} else {
			for kind, value := range values {
				_, e := role.SetPropertyValue(kind, value)
				if e != nil {
					err = e
				}
			}
		}

		return be, err
	}

	return nil, errors.New(baseerror.Error_NotFound)
}

type ActualDifference struct {
	SpecType string
	Path     string
	Kind     string
	State    string
	OldValue interface{}
	NewValue interface{}
}

type RuleMessage struct {
	Kind            string
	ReturnValue     string
	MessageCode     string
	Message         string
	BusinessMessage string
	InParaMeterDesc string
	Parameters      []interface{}
}

type BusinessEntityDifference struct {
	PromptMessages map[string]map[string]*RuleMessage
	/**
	 * @Fields compareResult : 所有的差异记录在一个多键值的映射中，包含角色和属性
	 *         角色的变化包含包含增加，修改，删除的角色（由BusinessEntity包裹）， 在相同角色下的属性修改，变化的结果是一个列表，
	 *         每一行代表一个修改的属性，每个差异代表属性的名称，旧值和新值
	 */
	CompareRoleResult       map[string][]*ActualDifference
	ComparePropertiesResult map[string][]*ActualDifference
}

func (this *BusinessEntityDifference) compareProperties(src *actual.Role, target *actual.Role) bool {
	isSame := true
	srcProperties := make(map[string]interface{})
	if src.Properties != nil {
		for _, p := range src.Properties {
			for i := 0; i < p.CurrentIndex; i++ {
				k, err := p.Get(entity.AttributeType_Kind, i)
				if err == nil {
					kind := k.(string)
					value, err := p.Get(entity.AttributeType_Value, i)
					if err == nil {
						srcProperties[kind] = value
					}
				}
			}
		}
	}

	if target.Properties != nil {
		for _, p := range target.Properties {
			for i := 0; i < p.CurrentIndex; i++ {
				k, err := p.Get(entity.AttributeType_Kind, i)
				if err == nil {
					kind := k.(string)
					targetValue, err := p.Get(entity.AttributeType_Value, i)
					if err == nil {
						srcValue, ok := srcProperties[kind]
						if ok {
							if srcValue == targetValue {

							} else {
								path := src.Path + "." + kind
								diff := &ActualDifference{SpecType: entity2.SpecType_Property, Path: path, Kind: kind, OldValue: srcValue, NewValue: targetValue}
								ads, ok := this.ComparePropertiesResult[path]
								if !ok {
									ads = make([]*ActualDifference, 0)
								}
								ads = append(ads, diff)
								this.ComparePropertiesResult[path] = ads
								isSame = false
							}
						}
					}
				}
			}
		}
	}

	return isSame
}

func (this *BusinessEntityDifference) recordRole(src *actual.Role, target *actual.Role, state string) {
	var diff *ActualDifference
	var path string
	var kind string
	if src != nil {
		path = src.Path
		kind = src.Kind
	} else if target != nil {
		path = target.Path
		kind = target.Kind
	}

	diff = &ActualDifference{SpecType: entity2.SpecType_Role, Path: path, Kind: kind, OldValue: src, NewValue: target}
	diff.State = state
	ads := this.ComparePropertiesResult[path]
	if ads == nil {
		ads = make([]*ActualDifference, 0)
	}
	ads = append(ads, diff)

	this.CompareRoleResult[path] = ads
}

func (this *BusinessEntityDifference) CompareRole(src *actual.Role, target *actual.Role) {
	if src != nil && target != nil { // 比较对象不能为空
		if src == target { // 如果是两个相同的对象，不用比较
			return
		} else {
			// 直接比较角色的属性是否相同
			isSame := this.compareProperties(src, target)
			if !isSame { // 属性不同，记录角色是改变的
				this.recordRole(src, target, entity3.EntityState_Modified)
			}
		}
		// 获取源的所有子角色，子角色是一个角色列表，建立新的映射作比较
		srcChildren := make(map[string][]*actual.Role)
		if src.Roles != nil && len(src.Roles) > 0 {
			for kind, children := range src.Roles {
				newChildren := make([]*actual.Role, 0)
				for _, child := range children {
					newChildren = append(newChildren, child)
				}
				srcChildren[kind] = newChildren
			}
		}
		if srcChildren != nil && len(srcChildren) > 0 { // 源角色有子角色
			if target.Roles != nil && len(target.Roles) > 0 { // 目标角色也有子角色
				for kind, targetChildrenList := range target.Roles { // 遍历目标角色的子角色
					// 获取同名的源角色
					srcChildrenList := srcChildren[kind]
					if srcChildrenList != nil { // 源角色存在，目标为角色数组
						if len(targetChildrenList) > 0 { // 目标为角色数组有角色
							if len(srcChildrenList) > 0 {
								for i := 0; i < len(targetChildrenList); i++ {
									targetChildRole := targetChildrenList[i]
									if len(srcChildrenList) > 0 {
										sRole := srcChildrenList[0]
										srcChildrenList = srcChildrenList[1:]
										this.CompareRole(sRole, targetChildRole)
									} else {
										this.CompareRole(nil, targetChildRole)
									}
								}
							} else { //源无角色
								for _, tRole := range targetChildrenList {
									this.CompareRole(nil, tRole)
								}
							}
						} else { // 目标为角色数组无角色
							delete(srcChildren, kind)
							if len(srcChildrenList) > 0 {
								for i := 0; i < len(srcChildrenList); i++ {
									child := srcChildrenList[i]
									if child != nil {
										this.CompareRole(child, nil)
									}
								}
							}
						}
					} else { // 源角色不存在，目标为角色数组
						for _, tRole := range targetChildrenList {
							this.CompareRole(nil, tRole)
						}
					}
				}
				// 最后检查源子角色数组中是否还有遗留的角色，如果有，是删除角色
				if len(srcChildren) > 0 {
					for _, chldrenRole := range srcChildren {
						for _, role := range chldrenRole {
							this.CompareRole(role, nil)
						}
					}
				}
			} else { // 目标角色无子角色
				if srcChildren != nil && len(srcChildren) > 0 {
					for _, chldrenRole := range srcChildren {
						for _, role := range chldrenRole {
							this.CompareRole(role, nil)
						}
					}
				}
			}
		} else { // // 源角色无子角色
			if target.Roles != nil && len(target.Roles) > 0 {
				for _, chldrenRole := range target.Roles {
					for _, role := range chldrenRole {
						this.CompareRole(nil, role)
					}
				}
			}
		}
	} else if src != nil && target == nil {
		this.recordRole(src, target, entity3.EntityState_Deleted)
	} else if src == nil && target != nil {
		this.recordRole(src, target, entity3.EntityState_New)
	}
}
