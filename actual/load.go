package actual

import (
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/actual/service"
	"github.com/curltech/go-colla-biz/spec"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	baseservice "github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/mohae/deepcopy"
	"sync"
)

/**
加载子角色，用于懒加载，loadNum为-1的情况
*/
func (this *Role) Load(kind string) []*Role {
	var specId uint64
	for _, roleSpec := range this.RoleSpec.RoleSpecs {
		if roleSpec.Kind == kind {
			specId = roleSpec.SpecId
			break
		}
	}
	if specId > 0 {
		roles, ids := service.GetRoleService().FindSpecOffspring(this.SchemaName, this.Id, specId)
		loader := newRoleLoader(this.SchemaName, 0, this, roles, ids)
		loader.load()
	}

	return this.Roles[kind]
}

/**
从数据库装载实例树
*/
func Load(schemaName string, id uint64) *Role {
	//获取所有的角色实体和角色编号数组
	var ids = make([]uint64, 1)
	ids[0] = id
	roles, ids := service.GetRoleService().FindOffspring(schemaName, ids, 0)

	loader := newRoleLoader(schemaName, id, nil, roles, ids)
	loader.load()
	role := loader.parentRole

	if role != nil {
		setCacheRole(role)
	}

	return role
}

type roleLoader struct {
	/**
	父亲节点，不为空表示局部加载
	*/
	parentRole *Role
	schemaName string
	/**
	顶级节点编号，不为空表示全加载
	*/
	topId uint64

	roles []*entity.Role
	ids   []uint64

	//角色的键值映射
	roleIdMap map[uint64]*Role

	//各元素的父键值映射数组
	propertyMap     map[uint64][]*entity.Property
	connectionMap   map[uint64][]*entity.Connection
	actionResultMap map[uint64][]*entity.ActionResult
	roleMap         map[uint64][]*Role
	//静态服务名与角色数组的对应关系
	fixedSpecMap map[uint64][]*Role
	//节点的角色定义关系
	roleSpecMap map[uint64]*spec.RoleSpec
}

func newRoleLoader(schemaName string, topId uint64, parentRole *Role, roles []*entity.Role, ids []uint64) *roleLoader {
	loader := &roleLoader{}
	loader.ids = ids
	loader.roles = roles

	//有父亲角色，局部加载
	if parentRole != nil {
		loader.parentRole = parentRole
		loader.schemaName = parentRole.SchemaName
	} else { //无父亲角色，全新加载
		loader.schemaName = schemaName
		loader.topId = topId
	}

	loader.propertyMap = make(map[uint64][]*entity.Property, 0)
	loader.connectionMap = make(map[uint64][]*entity.Connection, 0)
	loader.actionResultMap = make(map[uint64][]*entity.ActionResult, 0)
	loader.roleMap = make(map[uint64][]*Role, 0)
	loader.fixedSpecMap = make(map[uint64][]*Role, 0)
	loader.roleIdMap = make(map[uint64]*Role, 0)

	return loader
}

/**
包裹角色实体，设置初始化值和角色定义
*/
func packRole(r *entity.Role) *Role {
	role := &Role{}
	role.Role = *r

	role.UpdateDirtyFlag(baseentity.EntityState_None)
	role.UpdateState(baseentity.EntityState_None)

	return role
}

func (this *roleLoader) load() {
	this.loadElement()
	this.loadRole()
	this.loadRoleSpec()
	this.mountRole()
	this.mountFixedActual()
}

func (this *roleLoader) loadElement() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		//获取所有的属性实体，并按照父节点编号分组
		properties := service.GetPropertyService().ListByParentId(this.schemaName, this.ids)
		for _, p := range properties {
			p.DirtyFlag = baseentity.EntityState_None
			this.groupProperty(p)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//获取所有的连接实体，并按照父节点编号分组
		connections := service.GetConnectionService().ListByParentId(this.schemaName, this.ids)
		for _, c := range connections {
			c.DirtyFlag = baseentity.EntityState_None
			this.groupConnection(c)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		//获取所有的行为实体，并按照父节点编号分组
		actionResults := service.GetActionResultService().ListByParentId(this.schemaName, this.ids)
		for _, a := range actionResults {
			a.DirtyFlag = baseentity.EntityState_None
			this.groupActionResult(a)
		}
	}()
	wg.Wait()
}

func (this *roleLoader) loadRole() {
	if this.parentRole != nil {
		this.groupRole(this.parentRole)
		this.roleIdMap[this.parentRole.Id] = this.parentRole
	}
	//包裹角色实体，并按照父节点编号分组，第一次遍历角色
	for _, r := range this.roles {
		role := packRole(r)
		this.groupRole(role)
		//设置顶级节点作为返回
		if this.parentRole == nil && this.topId == role.Id {
			this.parentRole = role
		}
		this.roleIdMap[r.Id] = role
		this.groupFixedSpec(role)
	}
}

func (this *roleLoader) loadRoleSpec() {
	if this.parentRole != nil {
		if this.parentRole.RoleSpec == nil {
			roleSpec, _ := spec.GetMetaDefinition().GetRoleSpec(this.parentRole.SpecId, this.parentRole.EffectiveDate)
			this.parentRole.RoleSpec = roleSpec
			this.parentRole.Kind = roleSpec.Kind
		}
		this.roleSpecMap = this.parentRole.RoleSpec.GetRoleSpecMap()
	}
}

/**
第二次遍历角色，处理角色树，角色定义
*/
func (this *roleLoader) mountRole() error {
	var err error
	for k, roles := range this.roleMap {
		parentRole, ok := this.roleIdMap[k]
		for _, role := range roles {
			roleSpec := this.roleSpecMap[role.SpecId]
			if roleSpec != nil && role.RoleSpec == nil {
				role.RoleSpec = roleSpec
				role.Kind = roleSpec.Kind
			}
			this.mountElement(role)
			if ok && parentRole != nil {
				e := parentRole.PutRole(role)
				if e != nil {
					err = e
				}
			}
		}
	}

	return err
}

func (this *roleLoader) mountElement(role *Role) {
	//放置属性，遍历属性
	ps := this.propertyMap[role.Id]
	if ps != nil && len(ps) > 0 {
		for _, p := range ps {
			role.PutProperty(p)
		}
	}
	//检查是否属性定义发生变化
	role.UpdateProperty()
	//放置行为
	as := this.actionResultMap[role.Id]
	if as != nil && len(as) > 0 {
		for _, a := range as {
			role.PutActionResult(a)
		}
	}
	//放置连接
	cs := this.connectionMap[role.Id]
	if cs != nil && len(cs) > 0 {
		for _, c := range cs {
			role.PutConnection(c)
		}
	}
}

/**
分组静态服务定义，把相同的静态服务名的角色汇聚
*/
func (this *roleLoader) groupFixedSpec(role *Role) {
	if role.FixedSpecId > 0 {
		fs, ok := this.fixedSpecMap[role.FixedSpecId]
		if !ok {
			fs = make([]*Role, 0)
		}
		fs = append(fs, role)
		this.fixedSpecMap[role.FixedSpecId] = fs
	}
}

/**
按照汇聚的静态服务定义装载静态对象，并放置到角色中
*/
func (this *roleLoader) mountFixedActual() {
	var wg sync.WaitGroup
	for _, roles := range this.fixedSpecMap {
		ids := make([]uint64, 0)
		roleMap := make(map[uint64]*Role)
		for _, role := range roles {
			ids = append(ids, role.Id)
			roleMap[role.Id] = role
		}
		fixedService := roles[0].RoleSpec.GetFixedService()
		if fixedService == nil {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			fixedEntities, err := fixedService.NewEntities(nil)
			if err != nil {
				return
			}
			idss := service.Split(ids)
			for _, is := range idss {
				conds := fmt.Sprintf("ParentId in (%s)", baseservice.PlaceQuestionMark(len(is)))
				params := make([]interface{}, len(is))
				for i, id := range is {
					params[i] = id
				}
				err := fixedService.Find(fixedEntities, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
				if err != nil {
					continue
				}
				entities := reflect.ToArray(fixedEntities)
				if entities != nil && len(entities) > 0 {
					for _, fixedActual := range entities {
						v, err := reflect.GetValue(fixedActual, baseentity.FieldName_ParentId)
						if err == nil && v != nil {
							parentId := v.(uint64)
							role, ok := roleMap[parentId]
							if ok && role != nil {
								reflect.SetValue(fixedActual, "DirtyFlag", baseentity.EntityState_None)
								role.PutFixedActual(fixedActual)
								delete(roleMap, parentId)
							}
						}
					}
				}
			}
			//有些角色没有找到静态实例，要么数据有问题，要么修改了静态定义（需要将动态属性转成静态属性）
			if len(roleMap) > 0 {
				for _, role := range roleMap {
					role.transferProperty()
				}
			}
		}()
		wg.Wait()
	}
}

/**
分组角色，按照父节点编号和kind汇聚
*/
func (this *roleLoader) groupRole(role *Role) {
	rs, ok := this.roleMap[role.ParentId]
	if !ok {
		rs = make([]*Role, 0)
	}
	rs = append(rs, role)
	this.roleMap[role.ParentId] = rs
}

/**
分组属性，按照父节点编号汇聚
*/
func (this *roleLoader) groupProperty(property *entity.Property) {
	ps, ok := this.propertyMap[property.ParentId]
	if !ok {
		ps = make([]*entity.Property, 0)
	}
	ps = append(ps, property)
	this.propertyMap[property.ParentId] = ps
}

/**
分组行为，按照父节点编号汇聚
*/
func (this *roleLoader) groupActionResult(actionResult *entity.ActionResult) {
	as, ok := this.actionResultMap[actionResult.ParentId]
	if !ok {
		as = make([]*entity.ActionResult, 0)
	}
	as = append(as, actionResult)
	this.actionResultMap[actionResult.ParentId] = as
}

/**
分组连接，按照父节点编号汇聚
*/
func (this *roleLoader) groupConnection(connection *entity.Connection) {
	cs, ok := this.connectionMap[connection.ParentId]
	if !ok {
		cs = make([]*entity.Connection, 0)
	}
	cs = append(cs, connection)
	this.connectionMap[connection.ParentId] = cs
}

/**
克隆角色树，所有的数据保留，编号更新
*/
func (role *Role) Version() *Role {
	cloneRole := &Role{}
	cloneRole.Role = deepcopy.Copy(role.Role).(entity.Role)
	cloneRole.Id = service.GetRoleService().GetSeq()
	cloneRole.DirtyFlag = baseentity.EntityState_New
	cloneRole.RoleSpec = role.RoleSpec

	if role.FixedActual != nil {
		cloneFixedActual := deepcopy.Copy(role.FixedActual)
		roleSpec := role.RoleSpec
		if roleSpec != nil {
			fixedRoleSpec := roleSpec.FixedRoleSpec
			if fixedRoleSpec != nil {
				if fixedRoleSpec.FixedServiceName != "" {
					v := container.GetService(fixedRoleSpec.FixedServiceName)
					if v != nil {
						svc := v.(baseservice.BaseService)
						id := svc.GetSeq()
						reflect.SetValue(cloneFixedActual, baseentity.FieldName_Id, id)
					}
				}
			}
		}
		reflect.SetValue(cloneFixedActual, baseentity.FieldName_ParentId, cloneRole.Id)
		reflect.SetValue(cloneFixedActual, "DirtyFlag", baseentity.EntityState_New)
		cloneRole.FixedActual = cloneFixedActual
	}

	if role.Properties != nil && len(role.Properties) > 0 {
		for _, property := range role.Properties {
			cloneProperety := deepcopy.Copy(property).(*entity.Property)
			cloneProperety.Id = service.GetPropertyService().GetSeq()
			cloneProperety.ParentId = cloneRole.Id
			cloneProperety.DirtyFlag = baseentity.EntityState_New
			cloneRole.PutProperty(cloneProperety)
		}
	}

	if role.ActionResults != nil && len(role.ActionResults) > 0 {
		for _, actionResult := range role.ActionResults {
			cloneActionResult := deepcopy.Copy(actionResult).(*entity.ActionResult)
			cloneActionResult.Id = service.GetActionResultService().GetSeq()
			cloneActionResult.ParentId = cloneRole.Id
			cloneActionResult.DirtyFlag = baseentity.EntityState_New
			cloneRole.PutActionResult(cloneActionResult)
		}
	}

	if role.Connections != nil && len(role.Connections) > 0 {
		for _, connection := range role.Connections {
			cloneConnection := deepcopy.Copy(connection).(*entity.Connection)
			cloneConnection.Id = service.GetConnectionService().GetSeq()
			cloneConnection.ParentId = cloneRole.Id
			cloneConnection.DirtyFlag = baseentity.EntityState_New
			cloneRole.PutConnection(cloneConnection)
		}
	}

	if role.Roles != nil && len(role.Roles) > 0 {
		for _, rs := range role.Roles {
			for _, r := range rs {
				cloneR := r.Version()
				cloneR.ParentId = cloneRole.Id
				cloneR.ParentRole = cloneRole
				cloneRole.PutRole(cloneR)
			}
		}
	}

	if cloneRole != nil {
		if cloneRole.ParentId == 0 {
			setCacheRole(cloneRole)
		}
	}

	return cloneRole
}
