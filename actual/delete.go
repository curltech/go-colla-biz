package actual

import (
	"github.com/curltech/go-colla-core/entity"
	"github.com/kataras/golog"
)

func (this *Role) CanbeRemove(specId uint64) int {
	roleSpec, ok := this.RoleSpec.RoleSpecs[specId]
	if !ok {
		return 0
	}
	min, _ := this.GetMultiplicity(specId)
	kind := roleSpec.Kind
	count := this.GetCount(kind)
	if count > min {
		return count - min
	} else {
		golog.Errorf("Min:%v,count:%v cannot remove SpecId:%v", min, count, specId)
	}

	return 0
}

func (this *Role) RemoveRole(role *Role) bool {
	if this.CanbeRemove(role.SpecId) > 0 {
		rs, ok := this.Roles[role.Kind]
		if ok {
			if len(rs) > 0 {
				if this.DeleteRoles == nil {
					this.DeleteRoles = make(map[uint64]*Role, 0)
				}
				if this.DeleteActuals == nil {
					this.DeleteActuals = make([]uint64, 0)
				}
				needComputePos := false
				pos := -1
				for k, r := range rs {
					if needComputePos {
						this.computePath(r, k-1)
					}
					if r.Id == role.Id {
						role.DirtyFlag = entity.EntityState_Deleted
						this.DeleteRoles[role.Id] = role
						this.DeleteActuals = append(this.DeleteActuals, role.Id)
						pos = k
						needComputePos = true
					}
				}
				if pos > -1 {
					rs = append((rs)[:pos], (rs)[pos+1:]...)
					if len(rs) == 0 {
						this.Roles[role.Kind] = nil
					} else {
						this.Roles[role.Kind] = rs
					}
				}
				this.UpdateState(entity.EntityState_Modified)

				return true
			}
		}
	} else {
		golog.Errorf("connection count less min parent:%v;%v removed role:%v;%v", this.Id, this.Kind, role.Id, role.Kind)
	}

	return false
}

func (this *Role) getDeleteRoles() []*Role {
	roles := make([]*Role, 0)
	for _, role := range this.DeleteRoles {
		roles = append(roles, role)
	}

	return roles
}
