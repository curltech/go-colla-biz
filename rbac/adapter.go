package rbac

import (
	"github.com/curltech/go-colla-biz/rbac/entity"
	"github.com/curltech/go-colla-biz/rbac/service"
	"github.com/curltech/go-colla-core/util/reflect"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type Adapter struct {
	isFiltered bool
}

var adapter *Adapter = &Adapter{isFiltered: false}

func GetAdapter() *Adapter {
	return adapter
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	groups := make([]*entity.Group, 0)
	group := entity.Group{}
	group.Status = entity.UserStatus_Enabled
	svc := service.GetGroupService()
	svc.Find(&groups, &group, "", 0, 0, "")
	for _, g := range groups {
		var l = []string{"g",
			g.UserId, g.RoleId}
		line := strings.Join(l, ", ")
		persist.LoadPolicyLine(line, model)
	}

	permissions := make([]*entity.Permission, 0)
	permission := entity.Permission{}
	permission.Status = entity.UserStatus_Enabled
	service.GetPermissionService().Find(&permissions, &permission, "", 0, 0, "")
	for _, p := range permissions {
		var l = []string{"p",
			p.ActorId, p.ResourceId, p.AccessMode}
		line := strings.Join(l, ", ")
		persist.LoadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			permission := entity.Permission{}
			permission.ActorType = ptype
			permission.ActorId = rule[0]
			permission.ResourceId = rule[1]
			permission.AccessMode = rule[2]
			service.GetPermissionService().Insert(&permission)
		}
	}

	for _, ast := range model["g"] {
		for _, rule := range ast.Policy {
			group := entity.Group{}
			group.UserId = rule[0]
			group.RoleId = rule[1]
			service.GetGroupService().Insert(&group)
		}
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	if ptype == "p" {
		permission := entity.Permission{}
		permission.ActorType = ptype
		permission.ActorId = rule[0]
		permission.ResourceId = rule[1]
		permission.AccessMode = rule[2]
		service.GetPermissionService().Insert(&permission)
	} else if ptype == "g" {
		group := entity.Group{}
		group.UserId = rule[0]
		group.RoleId = rule[1]
		service.GetGroupService().Insert(&group)
	}

	return nil
}

// AddPolicies adds multiple policy rule to the storage.
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	if ptype == "p" {
		permissions := make([]*entity.Permission, 0)
		for _, rule := range rules {
			permission := entity.Permission{}
			permission.ActorType = ptype
			permission.ActorId = rule[0]
			permission.ResourceId = rule[1]
			permission.AccessMode = rule[2]
			permissions = append(permissions, &permission)
		}
		service.GetPermissionService().Insert(&permissions)
	} else if ptype == "g" {
		groups := make([]*entity.Group, 0)
		for _, rule := range rules {
			group := entity.Group{}
			group.UserId = rule[0]
			group.RoleId = rule[1]
			groups = append(groups, &group)
		}
		service.GetGroupService().Insert(&groups)
	}

	return nil
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	if ptype == "p" {
		permissions := make([]*entity.Permission, 0)
		permission := entity.Permission{}
		permission.ActorType = ptype
		permission.ActorId = rule[0]
		permission.ResourceId = rule[1]
		permission.AccessMode = rule[2]
		permissions = append(permissions, &permission)
		service.GetPermissionService().Insert(&permissions)
	} else if ptype == "g" {
		groups := make([]*entity.Group, 0)
		group := entity.Group{}
		group.UserId = rule[0]
		group.RoleId = rule[1]
		groups = append(groups, &group)
		service.GetGroupService().Insert(&groups)
	}
	return nil
}

// RemovePolicies removes multiple policy rule from the storage.
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	if ptype == "p" {
		permissions := make([]interface{}, 0)
		for _, rule := range rules {
			permission := entity.Permission{}
			permission.ActorType = ptype
			permission.ActorId = rule[0]
			permission.ResourceId = rule[1]
			permission.AccessMode = rule[2]
			permissions = append(permissions, &permission)
		}
		service.GetPermissionService().Delete(permissions, "")
	} else if ptype == "g" {
		groups := make([]interface{}, 0)
		for _, rule := range rules {
			group := entity.Group{}
			group.UserId = rule[0]
			group.RoleId = rule[1]
			groups = append(groups, &group)
		}
		service.GetGroupService().Delete(groups, "")
	}

	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	if ptype == "p" {
		permissions := make([]interface{}, 0)
		permission := entity.Permission{}
		permission.ActorType = ptype
		idx := fieldIndex + len(fieldValues)
		if fieldIndex <= 0 && idx > 0 {
			permission.ActorId = fieldValues[0-fieldIndex]
		}
		if fieldIndex <= 1 && idx > 1 {
			permission.ResourceId = fieldValues[1-fieldIndex]
		}
		if fieldIndex <= 2 && idx > 2 {
			permission.AccessMode = fieldValues[2-fieldIndex]
		}
		permissions = append(permissions, &permission)
		service.GetPermissionService().Delete(permissions, "")
	} else if ptype == "g" {
		groups := make([]interface{}, 0)
		group := entity.Group{}
		idx := fieldIndex + len(fieldValues)
		if fieldIndex <= 0 && idx > 0 {
			group.UserId = fieldValues[0-fieldIndex]
		}
		if fieldIndex <= 1 && idx > 1 {
			group.RoleId = fieldValues[1-fieldIndex]
		}
		groups = append(groups, &group)
		service.GetGroupService().Delete(groups, "")
	}

	return nil
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	ptype, _ := reflect.GetValue(filter, "PType")

	v0, _ := reflect.GetValue(filter, "V0")
	v1, _ := reflect.GetValue(filter, "V1")
	v2, _ := reflect.GetValue(filter, "V2")

	if ptype == "g" {
		groups := make([]*entity.Group, 0)
		group := entity.Group{}
		group.Status = entity.UserStatus_Enabled
		if v0 != nil {
			group.UserId = v0.(string)
		}
		if v1 != nil {
			group.RoleId = v1.(string)
		}

		service.GetGroupService().Find(groups, group, "", 0, 0, "")
		for _, g := range groups {
			var l = []string{"g",
				g.UserId, g.RoleId}
			line := strings.Join(l, ", ")
			persist.LoadPolicyLine(line, model)
		}
	}

	if ptype == "p" {
		permissions := make([]*entity.Permission, 0)
		permission := entity.Permission{}
		permission.Status = entity.UserStatus_Enabled
		if v0 != nil {
			permission.ActorId = v0.(string)
		}
		if v1 != nil {
			permission.ResourceId = v1.(string)
		}
		if v2 != nil {
			permission.AccessMode = v2.(string)
		}
		service.GetPermissionService().Find(permissions, permission, "", 0, 0, "")
		for _, p := range permissions {
			var l = []string{"p",
				p.ActorId, p.ResourceId, p.AccessMode}
			line := strings.Join(l, ", ")
			persist.LoadPolicyLine(line, model)
		}
	}
	a.isFiltered = true

	return nil
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *Adapter) IsFiltered() bool {
	return a.isFiltered
}
