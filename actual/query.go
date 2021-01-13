package actual

import (
	"errors"
	"github.com/curltech/go-colla-biz/actual/entity"
	specentity "github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/util/reflect"
	"strings"
)

func Get(schemaName string, id uint64) *Role {
	role := getCacheRole(schemaName, id)
	if role != nil {
		return role
	} else {
		role := Load(schemaName, id)
		if role != nil {
			setCacheRole(role)
		}

		return role
	}
}

func (this *Role) Get(id uint64) *Role {
	if this.Id == id {
		return this
	}
	for _, roles := range this.Roles {
		for _, role := range roles {
			if role.Id == id {
				return role
			} else {
				r := role.Get(id)
				if r != nil {
					return r
				}
			}
		}
	}

	return nil
}

func (this *Role) getFixedValue(kind string) (interface{}, error) {
	alias := this.getAlias(kind)
	if alias != "" {
		if this.FixedActual != nil {
			return reflect.GetValue(this.FixedActual, alias)
		} else {
			return nil, errors.New("NilFixedActual")
		}
	} else {
		return nil, errors.New("NilAlias")
	}
}

func (this *Role) GetPropertyValue(kind string) (interface{}, error) {
	value, err := this.getFixedValue(kind)
	if err != nil {
		for _, p := range this.Properties {
			if p.Contain(kind) {
				return p.GetValue(kind)
			}
		}
	} else {
		return value, nil
	}

	return nil, errors.New("NoKind")
}

func (this *Role) GetValue(path string) (interface{}, error) {
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
		return role.GetPropertyValue(kind)
	}

	return nil, errors.New("")
}

func (this *Role) Find(path string) (interface{}, string, bool) {
	pp := NewPositionPath(path)
	if pp.startWhere() == StartPosition_Root {
		return this.getTopRole().findExact(pp)
	} else if pp.startWhere() == StartPosition_Current {
		return this.findExact(pp)
	} else {
		rs := make([]*Role, 0)
		return this.findFuzzy(pp, rs)
	}
	return nil, pp.SpecType, false
}

/**
查找子对象，kind,position=-1都不为空，返回单个对象；
position>-1，返回对象数组
kind为空，返回全体子对象数组
*/
func (this *Role) findChildren(specType string, kind string, position int) (interface{}, bool) {
	if specType == specentity.SpecType_Role {
		if kind != "" {
			rs, ok := this.Roles[kind]
			if ok {
				if position <= -1 {
					return rs, len(rs) > 0
				} else if position >= len(rs) {
					return nil, false
				} else {
					r := rs[position]

					return r, r != nil
				}
			}
		} else {
			rs := make([]*Role, 0)
			for _, vs := range this.Roles {
				for _, r := range vs {
					rs = append(rs, r)
				}
			}

			return rs, len(rs) > 0
		}
	} else if specType == specentity.SpecType_Action {
		if kind != "" {
			a, _ := this.ActionResults[kind]

			return a, a != nil
		} else {
			as := make([]*entity.ActionResult, 0)
			for _, a := range this.ActionResults {
				as = append(as, a)
			}

			return as, len(as) > 0
		}
	} else if specType == specentity.SpecType_Property {
		ps := make([]*entity.Property, 0)
		for _, p := range this.Properties {
			if kind != "" {
				if p.Contain(kind) {
					return p, p != nil
				}
			} else {
				ps = append(ps, p)
			}
		}

		return ps, len(ps) > 0
	} else if specType == specentity.SpecType_Connection {
		cs := make([]*entity.Connection, 0)
		for _, c := range this.Connections {
			cs = append(cs, c)
		}

		return cs, len(cs) > 0
	}

	return nil, false
}

/**
 * 严格根据路径从当前位置一步一步地寻找，发现问题就返回null
 *
 * @param pp
 * @return
 */
func (this *Role) findExact(pp *PositionPath) (interface{}, string, bool) {
	node := pp.getCurrent()
	if node == nil {
		return this, pp.SpecType, true
	}
	position := node.Position
	kind := node.Kind
	var o interface{}
	var ok = false
	var specType string
	if kind == "" {
		return this, pp.SpecType, true
	} else if !pp.hasNext() {
		if o == nil && kind != "" {
			if pp.SpecType == specentity.SpecType_Property {
				o, err := reflect.GetValue(this.Role, kind)
				if err == nil {
					return o, pp.SpecType, true
				}

				alias := this.getAlias(kind)
				if alias != "" && this.FixedActual != nil { // 当前的固定属性对象
					v, err := reflect.GetValue(this.FixedActual, alias)
					if err == nil {
						return v, pp.SpecType, true
					}
				}
			}
			o, ok = this.findChildren(pp.SpecType, kind, position)
			if ok {
				return o, pp.SpecType, true
			} else {
				return nil, pp.SpecType, false
			}
		}
	} else {
		if o == nil && kind != "" {
			o, ok = this.findChildren(specentity.SpecType_Role, kind, position)
		}
		if o != nil {
			r, ok := o.(*Role)
			if ok {
				o, specType, ok = r.findExact(pp.next())
			}
		}
	}

	return o, specType, ok
}

/**
 * 从当前位置一步一步地寻找，找不到从下级节点再找， 尽量找到为止，发现问题或没有找到就返回null
 *
 * @param pp
 * @return
 */
func (this *Role) findFuzzy(pp *PositionPath, rs []*Role) (interface{}, string, bool) {
	node := pp.getCurrent()
	position := node.Position
	kind := node.Kind
	var o interface{}
	var ok = false
	var specType string
	if !pp.hasNext() { // 最后的路径节点
		if o == nil && kind != "" {
			alias := this.getAlias(kind)
			if alias != "" { // 当前的固定属性对象
				v, err := reflect.GetValue(this.FixedActual, alias)
				if err == nil {
					return v, pp.SpecType, true
				}
			}
			o, ok = this.findChildren(pp.SpecType, kind, position)
			if o != nil {
				return o, pp.SpecType, ok
			}
		}
		if o == nil && pp.getSize() == 1 { // 当前角色下没找到，继续在下属的角色中寻找
			o, ok = this.findChildren(pp.SpecType, "", -1)
			rs := o.([]*Role)
			if rs != nil && len(rs) > 0 {
				for _, r := range rs {
					o, specType, ok = r.findFuzzy(pp, rs)
					if o != nil {
						return o, specType, ok
					}
				}
			}
		}
	} else { // 带路径的路径
		if o == nil && kind != "" {
			o, ok = this.findChildren(specentity.SpecType_Role, kind, position)
		}
		if o == nil {
			o, ok = this.findChildren(specentity.SpecType_Role, "", -1)
			rs := o.([]*Role)
			if rs != nil && len(rs) > 0 {
				for _, r := range rs {
					o, specType, ok = r.findFuzzy(pp, rs)
					if o != nil {
						return o, specType, ok
					}
				}
			}
		} else if o != nil {
			r, ok := o.(*Role)
			if ok {
				o, specType, ok = r.findFuzzy(pp.next(), rs)
				if o != nil {
					return o, specType, ok
				}
			}

		}
	}

	return o, specType, ok
}
