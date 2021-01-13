package actual

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	entity2 "github.com/curltech/go-colla-biz/spec/entity"
	baseentity "github.com/curltech/go-colla-core/entity"
	baseerror "github.com/curltech/go-colla-core/error"
	"github.com/curltech/go-colla-core/util/collection"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/huandu/xstrings"
	basereflect "reflect"
)

func (this *Role) GetActual(isState bool) map[string]interface{} {
	actual := make(map[string]interface{}, 0)
	if isState && (this.State == "" || this.State == baseentity.EntityState_None) {
		return actual
	}
	options := collection.MapOptions{}
	roleActuals := collection.StructToMap(this.InternalFixedActual, &options)
	for k, v := range roleActuals {
		actual[xstrings.FirstRuneToLower(k)] = v
	}
	actual[baseentity.JsonFieldName_Path] = this.Path

	if isState && this.DeleteActuals != nil {
		actual["deleteActuals"] = this.DeleteActuals
	}

	if this.FixedActual != nil {
		for kind, alias := range this.RoleSpec.FixedAttributes {
			state, ok := this.PropertyStates[kind]
			if isState && !ok {
				continue
			}
			if isState && ok && (state == "" || state == baseentity.EntityState_None) {
				continue
			}
			value, err := reflect.GetValue(this.FixedActual, alias)
			if err == nil {
				kind = xstrings.FirstRuneToLower(kind)
				actual[kind] = value
			}
		}
	}

	if this.Properties != nil && len(this.Properties) > 0 {
		for _, property := range this.Properties {
			for i := 0; i < property.CurrentIndex; i++ {
				v, err := property.Get(entity.AttributeType_Kind, i)
				if err == nil {
					kind := v.(string)
					state, ok := this.PropertyStates[kind]
					if isState && !ok {
						continue
					}
					if isState && ok && (state == "" || state == baseentity.EntityState_None) {
						continue
					}
					value, err := property.GetValue(kind)
					if err == nil {
						kind = xstrings.FirstRuneToLower(kind)
						actual[kind] = value
					}
				}
			}
		}
	}

	if this.ActionResults != nil && len(this.ActionResults) > 0 {
		for kind, actionResult := range this.ActionResults {
			state, ok := this.PropertyStates[kind]
			if isState && !ok {
				continue
			}
			if isState && ok && (state == "" || state == baseentity.EntityState_None) {
				continue
			}
			kind = xstrings.FirstRuneToLower(kind)
			actual[kind] = actionResult.Value
		}
	}

	if this.Roles != nil && len(this.Roles) > 0 {
		for kind, rs := range this.Roles {
			as := make([]interface{}, 0)
			for _, r := range rs {
				a := r.GetActual(isState)
				if a != nil && len(a) > 0 {
					as = append(as, a)
				}
			}
			kind = xstrings.FirstRuneToLower(kind)
			actual[kind] = as
		}
	}
	this.clearState()

	return actual
}

func (this *Role) SetActual(actual map[string]interface{}, isState bool, isCheckActualId bool) map[uint64]map[string]error {
	var errs = make(map[uint64]map[string]error, 0)
	var err error
	v, ok := actual[baseentity.JsonFieldName_State]
	state := ""
	if ok {
		state = v.(string)
	}
	if state == baseentity.EntityState_Deleted {
		if this.ParentRole != nil {
			this.ParentRole.RemoveRole(this)
		}

		return errs
	}
	v, ok = actual[baseentity.JsonFieldName_Id]
	var id uint64
	if ok {
		v, err = convert.ToObject(fmt.Sprintf("%v", v), entity2.DataType_Uint64)
		if err == nil {
			id = v.(uint64)
		}
	}
	v, ok = actual[baseentity.JsonFieldName_SpecId]
	var specId uint64
	if ok {
		v, err = convert.ToObject(fmt.Sprintf("%v", v), entity2.DataType_Uint64)
		if err == nil {
			specId = v.(uint64)
		}
	}
	v, ok = actual[baseentity.JsonFieldName_Kind]
	var kind string
	if ok {
		kind = v.(string)
	}
	canModify := true
	if isState && (state == "" || state == baseentity.EntityState_None) {
		canModify = false
	}
	if specId > 0 && specId != this.SpecId {
		canModify = false
		setError(errs, this.Id, kind, errors.New(baseerror.Error_WrongSpecId))
	}
	if isCheckActualId && id > 0 && id != this.Id {
		canModify = false
		setError(errs, this.Id, kind, errors.New(baseerror.Error_WrongId))
	}
	if !canModify {
		this.State = baseentity.EntityState_None
		return errs
	}
	for k, value := range actual {
		v := basereflect.ValueOf(value)
		v = basereflect.Indirect(v)
		typ := v.Type().Kind()
		if typ == basereflect.Struct || typ == basereflect.Map || typ == basereflect.Slice || typ == basereflect.Array {
			roles, _ := this.Roles[k]
			if typ == basereflect.Struct || typ == basereflect.Map {
				if roles == nil || len(roles) == 0 {
					roles = this.addRoles(k, 1)
				}
				if roles != nil && len(roles) > 0 {
					ess := roles[0].SetActual(value.(map[string]interface{}), isState, isCheckActualId)
					mergeError(errs, ess)
				}
			} else if typ == basereflect.Slice || typ == basereflect.Array {
				slice := value.([]map[string]interface{})
				roleMap := make(map[uint64]*Role)
				if isCheckActualId {
					if roles != nil && len(roles) > 0 {
						for _, role := range roles {
							roleMap[role.Id] = role
						}
					}
				}
				for i := 0; i < len(slice); i++ {
					if isCheckActualId {
						var actualId uint64
						v1, ok := slice[i][baseentity.JsonFieldName_Id]
						if ok {
							actualId = v1.(uint64)
						}
						if actualId > 0 {
							role, ok := roleMap[actualId]
							if ok {
								ess := role.SetActual(slice[i], isState, isCheckActualId)
								mergeError(errs, ess)
							} else {
								setError(errs, this.Id, k, errors.New(baseerror.Error_NotExist))
							}
						} else {
							rs := this.addRoles(k, 1)
							ess := rs[0].SetActual(slice[i], isState, isCheckActualId)
							mergeError(errs, ess)
						}
					} else {
						if i < len(roles) {
							ess := roles[i].SetActual(slice[i], isState, isCheckActualId)
							mergeError(errs, ess)
						} else {
							setError(errs, this.Id, k, errors.New("ErrorRolesIndex"))
						}
					}
				}
			}
		} else {
			_, err := this.SetPropertyValue(xstrings.FirstRuneToUpper(k), value)
			if err != nil {
				setError(errs, this.Id, k, err)
			}
		}
	}
	this.clearState()

	return errs
}

func (this *Role) addRoles(kind string, count int) []*Role {
	roles := make([]*Role, 0)
	for _, roleSpec := range this.RoleSpec.RoleSpecs {
		if roleSpec.Kind == kind {
			role := this.AddRole(roleSpec, nil)
			if role != nil {
				roles = append(roles, role)
			}
			return roles
		}
	}

	return nil
}

func (this *Role) clearState() {
	this.DeleteActuals = nil
	this.State = baseentity.EntityState_None
	if this.FixedActual != nil {
		reflect.SetValue(this.FixedActual, baseentity.FieldName_State, baseentity.EntityState_None)
	}
	for kind, _ := range this.PropertyStates {
		this.PropertyStates[kind] = baseentity.EntityState_None
	}
	for _, p := range this.Properties {
		p.State = baseentity.EntityState_None
	}
}

func setError(errs map[uint64]map[string]error, id uint64, kind string, err error) {
	es, ok := errs[id]
	if !ok {
		es = make(map[string]error, 0)
		errs[id] = es
	}
	es[kind] = err
}

func mergeError(errs map[uint64]map[string]error, ess map[uint64]map[string]error) {
	for id, es := range ess {
		for kind, err := range es {
			setError(errs, id, kind, err)
		}
	}
}
