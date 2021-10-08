package actual

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/actual/service"
	"github.com/curltech/go-colla-biz/spec"
	specentity "github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/repository"
	baseservice "github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/reflect"
)

func Save(schemaName string, id uint64) (int, error) {
	role := getCacheRole(schemaName, id)
	if role != nil {
		return role.SaveAll()
	}

	return 0, errors.New("NotExistRole")
}

/**
开启事务，保存树
*/
func (this *Role) SaveAll() (int, error) {
	dirtyObject, err := service.GetRoleService().Transaction(func(session repository.DbSession) (interface{}, error) {
		dirtyObject := this.saveRole(session)

		return dirtyObject, nil
	})
	if err != nil {
		return 0, err
	}
	dos, ok := dirtyObject.([]interface{})
	if ok {
		for _, do := range dos {
			dirtyFlag, err := reflect.GetValue(do, "DirtyFlag")
			if err == nil {
				if dirtyFlag == baseentity.EntityState_Deleted {
					role, ok := do.(*Role)
					if ok {
						role.ParentRole.DeleteRoles = nil
					}
				}
				reflect.SetValue(do, "DirtyFlag", baseentity.EntityState_None)
			}
		}

		return len(dos), nil
	}

	return 0, errors.New("NotDirtyObject")
}

/**
角色没有加载，直接删除数据库
*/
func Delete(schemaName string, ids []uint64) (int64, error) {
	roles, ids := service.GetRoleService().FindOffspring(schemaName, ids, -1)
	if roles == nil || len(roles) == 0 {
		return 0, nil
	}
	affected, err := service.GetRoleService().Transaction(func(session repository.DbSession) (interface{}, error) {
		var affected int64
		var fixedSpecMap = make(map[uint64][]*entity.Role, 0)
		for _, role := range roles {
			if role.FixedSpecId > 0 {
				fs, ok := fixedSpecMap[role.FixedSpecId]
				if !ok {
					fs = make([]*entity.Role, 0)
				}
				fs = append(fs, role)
				fixedSpecMap[role.FixedSpecId] = fs
			}
		}
		affected = affected + deleteByParentId(session, ids, &entity.Role{})
		affected = affected + deleteByParentId(session, ids, &entity.Property{})
		affected = affected + deleteByParentId(session, ids, &entity.ActionResult{})
		affected = affected + deleteByParentId(session, ids, &entity.Connection{})
		for fixedSpecId, rs := range fixedSpecMap {
			fixedRoleSpec := spec.GetMetaDefinition().GetFixedRoleSpec(fixedSpecId)
			if fixedRoleSpec == nil {
				continue
			}
			if fixedRoleSpec.FixedServiceName == "" {
				continue
			}
			v := container.GetService(fixedRoleSpec.FixedServiceName)
			if v == nil {
				continue
			}
			fixedService, ok := v.(baseservice.BaseService)
			if !ok {
				continue
			}
			var is = make([]uint64, 0)
			for _, r := range rs {
				is = append(is, r.Id)
			}
			entity, err := fixedService.NewEntity(nil)
			if err == nil {
				affected = affected + deleteByParentId(session, is, entity)
			}
		}

		affected = affected + deleteById(session, ids, &entity.Role{})

		return affected, nil
	})

	return affected.(int64), err
}

func deleteByParentId(session repository.DbSession, is []uint64, entity interface{}) int64 {
	var affected int64
	idss := service.Split(is)
	for _, is := range idss {
		conds := fmt.Sprintf("ParentId in (%s)", baseservice.PlaceQuestionMark(len(is)))
		params := make([]interface{}, len(is))
		for i, id := range is {
			params[i] = id
		}
		c, _ := session.Delete(entity, conds, params...)
		affected = affected + c
	}

	return affected
}

func deleteById(session repository.DbSession, is []uint64, entity interface{}) int64 {
	var affected int64
	idss := service.Split(is)
	for _, is := range idss {
		conds := fmt.Sprintf("Id in (%s)", baseservice.PlaceQuestionMark(len(is)))
		params := make([]interface{}, len(is))
		for i, id := range is {
			params[i] = id
		}
		c, _ := session.Delete(entity, conds, params...)
		affected = affected + c
	}

	return affected
}

/**
角色已经加载，删除角色
*/
func (this *Role) Delete() int64 {
	affected, _ := service.GetRoleService().Transaction(func(session repository.DbSession) (interface{}, error) {
		affected := this.deleteRole(session)

		return affected, nil
	})

	return affected.(int64)
}

func (this *Role) deleteRole(session repository.DbSession) int64 {
	var affected int64
	v, ok := this.findChildren(specentity.SpecType_Property, "", -1)
	if ok {
		properties := v.([]*entity.Property)
		if properties != nil && len(properties) > 0 {
			var ps = make([]interface{}, 1)
			for _, property := range properties {
				ps[0] = property
			}
			c, _ := session.Delete(ps, "")
			affected = affected + c
		}
	}
	v, ok = this.findChildren(specentity.SpecType_Action, "", -1)
	if ok {
		actionResults := v.([]*entity.ActionResult)
		if actionResults != nil && len(actionResults) > 0 {
			var as = make([]interface{}, 1)
			for _, a := range actionResults {
				as[0] = a
			}
			c, _ := session.Delete(as, "")
			affected = affected + c
		}
	}

	//删除关联的静态对象
	if this.FixedActual != nil {
		var fixedActuals = make([]interface{}, 1)
		fixedActuals[0] = this.FixedActual
		c, _ := session.Delete(fixedActuals, "")
		affected = affected + c
	}

	v, ok = this.findChildren(specentity.SpecType_Role, "", -1)
	if ok {
		roles := v.([]*Role)
		if roles != nil && len(roles) > 0 {
			for _, role := range roles {
				affected = affected + role.deleteRole(session)
			}
		}
	}

	var deleteRoles = make([]interface{}, 1)
	deleteRoles[0] = &this.Role
	c, _ := session.Delete(deleteRoles, "")
	affected = affected + c

	return affected
}

/**
保存实例树到数据库，根据脏标志只保存修改的对象，返回值是发生变化的实例对象的数组
*/
func (this *Role) saveRole(session repository.DbSession) []interface{} {
	dirtyObject := make([]interface{}, 0)
	fixedActual := this.FixedActual
	if fixedActual != nil {
		dirtyFlag, err := reflect.GetValue(fixedActual, "DirtyFlag")
		if err != nil {
			dirtyFlag = baseentity.EntityState_None
		}
		if dirtyFlag == baseentity.EntityState_New {
			session.Insert(fixedActual)
		} else if dirtyFlag == baseentity.EntityState_Modified {
			fixedActuals := make([]interface{}, 1)
			fixedActuals[0] = fixedActual
			session.Update(fixedActuals, nil, "")
		} else if dirtyFlag == baseentity.EntityState_Deleted {
			fixedActuals := make([]interface{}, 1)
			fixedActuals[0] = fixedActual
			session.Delete(fixedActuals, "")
		}
		dirtyObject = append(dirtyObject, fixedActual)
	}
	dirtyFlag := this.DirtyFlag
	if dirtyFlag == baseentity.EntityState_New {
		session.Insert(&this.Role)
	} else if dirtyFlag == baseentity.EntityState_Modified {
		roles := make([]interface{}, 1)
		roles[0] = &this.Role
		session.Update(roles, nil, "")
	} else if dirtyFlag == baseentity.EntityState_Deleted {
		roles := make([]interface{}, 1)
		roles[0] = &this.Role
		session.Delete(roles, "")
	}
	dirtyObject = append(dirtyObject, this)

	if this.DeleteRoles != nil && len(this.DeleteRoles) > 0 {
		for _, r := range this.DeleteRoles {
			r.deleteRole(session)
			dirtyObject = append(dirtyObject, r)
		}
	}

	v, ok := this.findChildren(specentity.SpecType_Property, "", -1)
	if ok {
		properties := v.([]*entity.Property)
		if properties != nil && len(properties) > 0 {
			for _, property := range properties {
				dirtyFlag := property.DirtyFlag
				if dirtyFlag == baseentity.EntityState_New {
					session.Insert(property)
				} else if dirtyFlag == baseentity.EntityState_Modified {
					ps := make([]interface{}, 1)
					ps[0] = property
					session.Update(ps, nil, "")
				} else if dirtyFlag == baseentity.EntityState_Deleted {
					ps := make([]interface{}, 1)
					ps[0] = property
					session.Delete(ps, "")
				}
				dirtyObject = append(dirtyObject, property)
			}
		}
	}
	v, ok = this.findChildren(specentity.SpecType_Action, "", -1)
	if ok {
		actionResults := v.([]*entity.ActionResult)
		if actionResults != nil && len(actionResults) > 0 {
			for _, actionResult := range actionResults {
				dirtyFlag := actionResult.DirtyFlag
				if dirtyFlag == baseentity.EntityState_New {
					session.Insert(actionResult)
				} else if dirtyFlag == baseentity.EntityState_Modified {
					as := make([]interface{}, 1)
					as[0] = actionResult
					session.Update(as, nil, "")
				} else if dirtyFlag == baseentity.EntityState_Deleted {
					as := make([]interface{}, 1)
					as[0] = actionResult
					session.Delete(as, "")
				}
				dirtyObject = append(dirtyObject, actionResult)
			}
		}
	}
	v, ok = this.findChildren(specentity.SpecType_Connection, "", -1)
	if ok {
		connections := v.([]*entity.Connection)
		if connections != nil && len(connections) > 0 {
			for _, connection := range connections {
				dirtyFlag := connection.DirtyFlag
				if dirtyFlag == baseentity.EntityState_New {
					session.Insert(connection)
				} else if dirtyFlag == baseentity.EntityState_Modified {
					cs := make([]interface{}, 1)
					cs[0] = connection
					session.Update(cs, nil, "")
				} else if dirtyFlag == baseentity.EntityState_Deleted {
					cs := make([]interface{}, 1)
					cs[0] = connection
					session.Delete(cs, "")
				}
				dirtyObject = append(dirtyObject, connection)
			}
		}
	}
	v, ok = this.findChildren(specentity.SpecType_Role, "", -1)
	if ok {
		roles := v.([]*Role)
		if roles != nil && len(roles) > 0 {
			for _, role := range roles {
				do := role.saveRole(session)
				if len(do) > 0 {
					dirtyObject = append(dirtyObject, do...)
				}
			}
		}
	}

	return dirtyObject
}
