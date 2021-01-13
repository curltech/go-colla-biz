package actual

import (
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-biz/actual/service"
	"github.com/curltech/go-colla-biz/spec"
	specentity "github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	baseservice "github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/reflect"
)

func newRole(roleSpec *spec.RoleSpec) *Role {
	role := &Role{}

	role.RoleSpec = roleSpec
	role.FixedSpecId = roleSpec.FixedSpecId
	role.Kind = roleSpec.Kind
	role.SpecId = roleSpec.SpecId

	role.DirtyFlag = baseentity.EntityState_New
	role.State = baseentity.EntityState_New
	role.EffectiveDate = roleSpec.EffectiveDate

	roleSvc := service.GetRoleService()
	role.Id = roleSvc.GetSeq()

	return role
}

func Create(schemaName string, roleSpec *spec.RoleSpec) *Role {
	role := createRole(schemaName, nil, roleSpec, 0, nil)
	if role != nil {
		setCacheRole(role)
	}

	return role
}

func createRole(schemaName string, parent *Role, roleSpec *spec.RoleSpec, position int, values map[string]interface{}) *Role {
	if parent != nil && schemaName == "" {
		schemaName = parent.SchemaName
	}
	//创建新的角色
	role := newRole(roleSpec)
	role.SchemaName = schemaName
	role.FirstId = role.Id

	// 最顶级节点
	var loadNum int = 0
	specId := roleSpec.SpecId
	//顶级节点
	if parent == nil {
		if schemaName == "" {
			// 未来加入分表算法
			schemaName = "" //schemaNameComponent.getSchemaName(actual.getActualId());
		}
		role.ParentId = 0
		role.LoadNum = loadNum
		role.TopId = role.Id
	} else {
		//自己不是顶级节点，则顶级节点的状态改变
		parent.UpdateState(baseentity.EntityState_Modified)
		//获取与上级节点之间的连接定义
		conSpec := parent.RoleSpec.GetChildConnectionSpec(specId)
		if conSpec != nil {
			loadNum := conSpec.LoadNum
			if loadNum != 0 {
				role.LoadNum = loadNum
			}
		}
		role.ParentId = parent.Id
		parent.PutRole(role)
	}
	//创建关联的静态对象
	fixedRoleSpec := roleSpec.FixedRoleSpec
	if fixedRoleSpec != nil {
		serviceName := fixedRoleSpec.FixedServiceName
		if serviceName != "" {
			fixedSvc := container.GetService(serviceName)
			if fixedSvc != nil {
				fixedService := fixedSvc.(baseservice.BaseService)
				fixedActual, err := fixedService.NewEntity(nil)
				if err == nil {
					baseEntity := fixedActual.(baseentity.IBaseEntity)
					baseEntity.SetId(fixedService.GetSeq())
					baseEntity.UpdateState(baseentity.EntityState_New)

					internalFixedActual := fixedActual.(specentity.IInternalFixedActual)
					internalFixedActual.SetParentId(role.Id)
					internalFixedActual.UpdateDirtyFlag(baseentity.EntityState_New)
					role.PutFixedActual(fixedActual)
				}
			}
		}
	}

	propertySvc := service.GetPropertyService()
	// 创建角色的所有子属性和行为属性
	attributeSpecs := roleSpec.GetSortedAttributeSpecs(false)
	if attributeSpecs != nil && len(attributeSpecs) > 0 {
		// 遍历每个要创建的属性
		// 全部属性计数器
		i := 0
		// PropertyEO的属性计数器，
		// 当等于AttributeSpecNumber时创建新的PropertyEO,重新计数
		j := 0
		var property *entity.Property
		// 动态属性计算器
		index := 0
		for _, attributeSpec := range attributeSpecs {
			dataType := attributeSpec.DataType
			var defaultValue interface{}
			if attributeSpec.DefaultValue != "" {
				defaultValue = attributeSpec.DefaultValue
				if dataType != "" {
					defaultValue, _ = convert.ToObject(attributeSpec.DefaultValue, dataType)
				}
				if values != nil {
					v, ok := (values)[attributeSpec.Kind]
					if ok {
						defaultValue = v
					}
				}
			}
			// 如果属性是动态属性
			alias := role.getAlias(attributeSpec.Kind)
			if alias == "" {
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
				role.computeProperty(property, attributeSpec, j)
				j++
			} else {
				if defaultValue != nil {
					reflect.SetValue(role.FixedActual, alias, defaultValue)
				}
			}
			i++
			// 在PropertyEO全部填满之后或者没有属性需要填充的时候放入角色
			if property != nil {
				if j >= entity.AttributeSpecNumber || i >= len(attributeSpecs) {
					property.CurrentIndex = j
					property.ParentId = role.Id
					property.Id = propertySvc.GetSeq()
					property.SchemaName = schemaName
					role.PutProperty(property)
					property = nil
				}
			}
		}
	}

	actionSvc := service.GetActionResultService()
	actionSpecs := roleSpec.ActionSpecs
	if actionSpecs != nil && len(actionSpecs) > 0 {
		for _, actionSpec := range actionSpecs {
			conSpec := roleSpec.GetChildConnectionSpec(actionSpec.SpecId)
			if conSpec != nil {
				count := 0
				if conSpec != nil {
					count = conSpec.BuildNum
				}
				for i := 0; i < count; i++ {
					actionResult := &entity.ActionResult{}
					actionResult.Id = actionSvc.GetSeq()
					actionResult.ParentId = role.Id
					actionResult.SpecId = actionSpec.SpecId
					actionResult.Kind = actionSpec.Kind
					actionResult.SchemaName = schemaName
					actionResult.DirtyFlag = baseentity.EntityState_New
					actionResult.State = baseentity.EntityState_New
				}
			}
		}
	}

	// 创建子角色
	roleSpecs := roleSpec.GetSortedRoleSpecs()
	if roleSpecs != nil && len(roleSpecs) > 0 {
		// 遍历子角色定义
		for _, subRoleSpec := range roleSpecs {
			conSpec := roleSpec.GetChildConnectionSpec(subRoleSpec.SpecId)
			if conSpec != nil {
				relationType := conSpec.RelationType
				buildNum := conSpec.BuildNum
				min := conSpec.Minmium
				//要创建角色的数量
				count := 0
				if buildNum >= min {
					count = buildNum
				} else {
					count = min
				}
				if relationType == specentity.RelationType_Association {
					//不创建角色对象，只建立关联对象
					count = 0
				}

				var value map[string]interface{}
				if values != nil {
					v, ok := values[role.Kind]
					if ok {
						value = v.(map[string]interface{})
					}
				}
				for i := 0; i < count; i++ {
					//创建子角色
					subRole := createRole(schemaName, role, subRoleSpec, i, value)
					//创建对应的连接对象
					if relationType == specentity.RelationType_Aggregation || relationType == specentity.RelationType_Association {
						var connection = &entity.Connection{}
						connection.ParentId = role.Id
						connection.DirtyFlag = baseentity.EntityState_New
						connection.State = baseentity.EntityState_New
						connection.ActualId = subRole.Id
						role.PutConnection(connection)
					}
				}
			}
		}
	}

	return role
}
