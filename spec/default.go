package spec

import (
	"fmt"
	"github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-biz/spec/service"
	"github.com/curltech/go-colla-core/cache"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	baseservice "github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/huandu/xstrings"
	"sort"
	"sync"
	"time"
)

type MetaDefinition struct {
	RoleSpecs       []*entity.RoleSpec       `json:",omitempty"`
	AttributeSpecs  []*entity.AttributeSpec  `json:",omitempty"`
	ActionSpecs     []*entity.ActionSpec     `json:",omitempty"`
	FixedRoleSpecs  []*entity.FixedRoleSpec  `json:",omitempty"`
	ConnectionSpecs []*entity.ConnectionSpec `json:",omitempty"`
}

type RoleSpec struct {
	entity.RoleSpec
	RoleSpecs       map[uint64]*RoleSpec              `json:",omitempty"`
	AttributeSpecs  map[uint64]*entity.AttributeSpec  `json:",omitempty"`
	ActionSpecs     map[uint64]*entity.ActionSpec     `json:",omitempty"`
	FixedRoleSpec   *entity.FixedRoleSpec             `json:",omitempty"`
	FixedAttributes map[string]string                 `json:"-,omitempty"`
	ConnectionSpecs map[string]*entity.ConnectionSpec `json:",omitempty"`
	EffectiveDate   *time.Time
	Parent          *RoleSpec `json:"-,omitempty"`
	Path            string    `json:"-,omitempty"`
}

/**
这个角色定义专用于与前端的交互
*/
type ModelSpec struct {
	entity.RoleSpec
	Children        []*ModelSpec                        `json:"children,omitempty"`
	AttributeSpecs  []*entity.AttributeSpec             `json:",omitempty"`
	ActionSpecs     []*entity.ActionSpec                `json:",omitempty"`
	FixedRoleSpec   *entity.FixedRoleSpec               `json:",omitempty"`
	ConnectionSpecs map[string][]*entity.ConnectionSpec `json:",omitempty"`
	EffectiveDate   *time.Time
}

type RoleSpecs []*RoleSpec

func (as RoleSpecs) Len() int { return len(as) }
func (as RoleSpecs) Less(i, j int) bool {
	return as[i].SpecId < as[j].SpecId
}
func (as RoleSpecs) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

type AttributeSpecs []*entity.AttributeSpec

func (as AttributeSpecs) Len() int { return len(as) }
func (as AttributeSpecs) Less(i, j int) bool {
	return as[i].SpecId < as[j].SpecId
}
func (as AttributeSpecs) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func (this *RoleSpec) GetChildConnectionSpec(subSpecId uint64) *entity.ConnectionSpec {
	key := fmt.Sprintf("%v:%v", this.SpecId, subSpecId)
	v, ok := this.ConnectionSpecs[key]
	if ok {
		return v
	}

	return nil
}

func (this *RoleSpec) GetChildRoleSpec(kind string) *RoleSpec {
	for _, roleSpec := range this.RoleSpecs {
		if kind == roleSpec.Kind {
			return roleSpec
		}
	}

	return nil
}

func (this *RoleSpec) GetChildAttributeSpec(kind string) *entity.AttributeSpec {
	for _, attributeSpec := range this.AttributeSpecs {
		if kind == attributeSpec.Kind {
			return attributeSpec
		}
	}

	return nil
}

func (this *RoleSpec) GetChildActionSpec(kind string) *entity.ActionSpec {
	for _, actionSpec := range this.ActionSpecs {
		if kind == actionSpec.Kind {
			return actionSpec
		}
	}

	return nil
}

func (this *RoleSpec) GetSortedAttributeSpecs(dynamic bool) []*entity.AttributeSpec {
	var as = make(AttributeSpecs, 0)
	for _, v := range this.AttributeSpecs {
		if dynamic && this.FixedAttributes != nil {
			alias, ok := this.FixedAttributes[v.Kind]
			if ok {
				if alias != "" {
					continue
				}
			}
		}
		as = append(as, v)
	}
	sort.Sort(as)

	return as
}

func (this *RoleSpec) GetSortedRoleSpecs() []*RoleSpec {
	var rs = make(RoleSpecs, 0)
	for _, v := range this.RoleSpecs {
		rs = append(rs, v)
	}
	sort.Sort(rs)

	return rs
}

func (this *RoleSpec) GetFixedService() baseservice.BaseService {
	fixedRoleSpec := this.FixedRoleSpec
	if fixedRoleSpec == nil {
		return nil
	}
	fixedServiceName := fixedRoleSpec.FixedServiceName
	if fixedServiceName == "" {
		return nil
	}
	v := container.GetService(fixedServiceName)
	if v != nil {
		fixedService, ok := v.(baseservice.BaseService)
		if ok {
			return fixedService
		}
	}

	return nil
}

func (this *RoleSpec) GetRoleSpecMap() map[uint64]*RoleSpec {
	roleSpecMap := make(map[uint64]*RoleSpec, 0)
	roleSpecMap[this.RoleSpec.SpecId] = this
	if this.RoleSpecs != nil && len(this.RoleSpecs) > 0 {
		for _, v := range this.RoleSpecs {
			roleSpecMap[v.SpecId] = v
			rMap := v.GetRoleSpecMap()
			if len(rMap) > 0 {
				for k, v := range rMap {
					roleSpecMap[k] = v
				}
			}
		}
	}

	return roleSpecMap
}

var metaDefinition = &MetaDefinition{}
var MemCache = cache.NewMemCache("spec", 1000, 1000)

func GetMetaDefinition() *MetaDefinition {
	return metaDefinition
}

func (md *MetaDefinition) GetRoleSpec(specId uint64, effectiveDate *time.Time) (*RoleSpec, *ModelSpec) {
	return md.getRoleSpec(nil, specId, effectiveDate)
}

func isEffectiveSpec(startDate *time.Time, endDate *time.Time, effectiveDate *time.Time) bool {
	if effectiveDate == nil {
		return true
	}
	currentTime := time.Now()
	if startDate == nil {
		startDate = &currentTime
	}
	if endDate == nil {
		endDate = &currentTime
	}
	return effectiveDate.After(*startDate) && effectiveDate.Before(*endDate)
}

func (md *MetaDefinition) getRoleSpecEntity(specId uint64, effectiveDate *time.Time) *entity.RoleSpec {
	key := fmt.Sprintf("RoleSpec:SpecId:%v", specId)
	v, ok := MemCache.Get(key)
	if ok {
		roleSpec := v.(*entity.RoleSpec)
		effective := isEffectiveSpec(roleSpec.StartDate, roleSpec.EndDate, effectiveDate)
		if effective {
			return roleSpec
		}
	}

	return nil
}

func (md *MetaDefinition) getActionSpecEntity(specId uint64, effectiveDate *time.Time) *entity.ActionSpec {
	key := fmt.Sprintf("ActionSpec:SpecId:%v", specId)
	v, ok := MemCache.Get(key)
	if ok {
		actionSpec := v.(*entity.ActionSpec)
		effective := isEffectiveSpec(actionSpec.StartDate, actionSpec.EndDate, effectiveDate)
		if effective {
			return actionSpec
		}
	}

	return nil
}

func (md *MetaDefinition) GetAttributeSpecEntity(specId uint64, effectiveDate *time.Time) *entity.AttributeSpec {
	key := fmt.Sprintf("AttributeSpec:SpecId:%v", specId)
	v, ok := MemCache.Get(key)
	if ok {
		attributeSpec := v.(*entity.AttributeSpec)
		effective := isEffectiveSpec(attributeSpec.StartDate, attributeSpec.EndDate, effectiveDate)
		if effective {
			return attributeSpec
		}
	}

	return nil
}

func (md *MetaDefinition) getChildrenConnectionSpec(specType string, specId uint64, effectiveDate *time.Time) []*entity.ConnectionSpec {
	key := fmt.Sprintf("ConnectionSpec:SpecType:%v:ParentSpecId:%v", specType, specId)
	css, ok := MemCache.Get(key)
	conns := make([]*entity.ConnectionSpec, 0)
	if ok {
		for _, cs := range css.(map[uint64][]*entity.ConnectionSpec) {
			for _, c := range cs {
				effective := isEffectiveSpec(c.StartDate, c.EndDate, effectiveDate)
				if effective {
					conns = append(conns, c)
					break
				}
			}
		}
	}

	return conns
}

func (md *MetaDefinition) getRoleSpec(parent *RoleSpec, specId uint64, effectiveDate *time.Time) (*RoleSpec, *ModelSpec) {
	var roleSpec *RoleSpec
	var modelSpec *ModelSpec
	spec := md.getRoleSpecEntity(specId, effectiveDate)
	if spec == nil {
		return nil, nil
	}
	roleSpec = &RoleSpec{RoleSpec: *spec}
	roleSpec.AttributeSpecs = make(map[uint64]*entity.AttributeSpec, 0)
	roleSpec.ActionSpecs = make(map[uint64]*entity.ActionSpec, 0)
	roleSpec.RoleSpecs = make(map[uint64]*RoleSpec, 0)
	roleSpec.ConnectionSpecs = make(map[string]*entity.ConnectionSpec, 0)
	modelSpec = &ModelSpec{RoleSpec: *spec}
	modelSpec.AttributeSpecs = make([]*entity.AttributeSpec, 0)
	modelSpec.ActionSpecs = make([]*entity.ActionSpec, 0)
	modelSpec.Children = make([]*ModelSpec, 0)
	modelSpec.ConnectionSpecs = make(map[string][]*entity.ConnectionSpec, 0)
	if parent != nil {
		roleSpec.Parent = parent
		parentPath := parent.Path
		if parentPath != "" {
			roleSpec.Path = parentPath + "." + spec.Kind
		} else {
			roleSpec.Path = spec.Kind
		}
	} else {
		roleSpec.Path = ""
	}
	roleSpec.EffectiveDate = effectiveDate
	modelSpec.EffectiveDate = effectiveDate

	cons := md.getChildrenConnectionSpec(entity.SpecType_Role, specId, effectiveDate)
	if cons != nil {
		for _, con := range cons {
			subSpec := md.getRoleSpecEntity(con.SubSpecId, effectiveDate)
			if subSpec == nil {
				logger.Sugar.Errorf("")
			} else {
				key := fmt.Sprintf("%v:%v", roleSpec.SpecId, con.SubSpecId)
				roleSpec.ConnectionSpecs[key] = con
				childRoleSpec, childModelSpec := md.getRoleSpec(roleSpec, subSpec.SpecId, effectiveDate)
				roleSpec.RoleSpecs[subSpec.SpecId] = childRoleSpec

				connectionSpecs, ok := modelSpec.ConnectionSpecs[entity.SpecType_Role]
				if !ok {
					connectionSpecs = make([]*entity.ConnectionSpec, 0)
				}
				connectionSpecs = append(connectionSpecs, con)
				modelSpec.ConnectionSpecs[entity.SpecType_Role] = connectionSpecs
				modelSpec.Children = append(modelSpec.Children, childModelSpec)
			}
		}
	}

	cons = md.getChildrenConnectionSpec(entity.SpecType_Action, specId, effectiveDate)
	if cons != nil {
		for _, con := range cons {
			subSpec := md.getActionSpecEntity(con.SubSpecId, effectiveDate)
			if subSpec == nil {
				logger.Sugar.Errorf("")
			} else {
				key := fmt.Sprintf("%v:%v", roleSpec.SpecId, con.SubSpecId)
				roleSpec.ConnectionSpecs[key] = con
				roleSpec.ActionSpecs[subSpec.SpecId] = subSpec

				connectionSpecs, ok := modelSpec.ConnectionSpecs[entity.SpecType_Action]
				if !ok {
					connectionSpecs = make([]*entity.ConnectionSpec, 0)
				}
				connectionSpecs = append(connectionSpecs, con)
				modelSpec.ConnectionSpecs[entity.SpecType_Action] = connectionSpecs
				modelSpec.ActionSpecs = append(modelSpec.ActionSpecs, subSpec)
			}
		}
	}

	/**
	角色有静态部分
	*/
	fixedSpecId := spec.FixedSpecId
	if fixedSpecId > 0 {
		fixedRoleSpec := md.GetFixedRoleSpec(fixedSpecId)
		roleSpec.FixedRoleSpec = fixedRoleSpec
		modelSpec.FixedRoleSpec = fixedRoleSpec
	}

	cons = md.getChildrenConnectionSpec(entity.SpecType_Property, specId, effectiveDate)
	if cons != nil {
		for _, con := range cons {
			subSpec := md.GetAttributeSpecEntity(con.SubSpecId, effectiveDate)
			if subSpec == nil {
				logger.Sugar.Errorf("")
			} else {
				key := fmt.Sprintf("%v:%v", roleSpec.SpecId, con.SubSpecId)
				roleSpec.ConnectionSpecs[key] = con
				roleSpec.AttributeSpecs[subSpec.SpecId] = subSpec

				connectionSpecs, ok := modelSpec.ConnectionSpecs[entity.SpecType_Property]
				if !ok {
					connectionSpecs = make([]*entity.ConnectionSpec, 0)
				}
				connectionSpecs = append(connectionSpecs, con)
				modelSpec.ConnectionSpecs[entity.SpecType_Property] = connectionSpecs
				modelSpec.AttributeSpecs = append(modelSpec.AttributeSpecs, subSpec)

				alias := md.GetFixedAttributesAlias(fixedSpecId, subSpec)
				if alias != "" {
					if roleSpec.FixedAttributes == nil {
						roleSpec.FixedAttributes = make(map[string]string)
					}
					roleSpec.FixedAttributes[subSpec.Kind] = alias
				}
			}
		}
	}

	return roleSpec, modelSpec
}

/**
返回动态属性或者静态属性
*/
func (md *MetaDefinition) GetChildrenAttributeSpec(specId uint64, effectiveDate *time.Time, dynamic bool) []*entity.AttributeSpec {
	var specs []*entity.AttributeSpec
	parent := md.getRoleSpecEntity(specId, effectiveDate)
	if parent != nil {
		specs = make([]*entity.AttributeSpec, 0)
		cons := md.getChildrenConnectionSpec(entity.SpecType_Property, specId, effectiveDate)

		if len(cons) > 0 {
			for _, con := range cons {
				subSpecId := con.SubSpecId
				spec := md.GetAttributeSpecEntity(subSpecId, effectiveDate)
				fixedSpecId := parent.FixedSpecId
				alias := md.GetFixedAttributesAlias(fixedSpecId, spec)
				if alias != "" {
					if !dynamic {
						specs = append(specs, spec)
						continue
					}
				} else {
					if dynamic {
						specs = append(specs, spec)
					}
				}
			}
		}
	}

	return specs
}

func (md *MetaDefinition) GetFixedRoleSpec(fixedSpecId uint64) *entity.FixedRoleSpec {
	key := fmt.Sprintf("FixedRoleSpec:SpecId:%v", fixedSpecId)
	o, ok := MemCache.Get(key)
	if ok {
		return o.(*entity.FixedRoleSpec)
	}

	return nil
}

//如果返回的alias不为空，说明这是静态属性，字段名就是返回的别名
func (md *MetaDefinition) GetFixedAttributesAlias(fixedSpecId uint64, attributeSpec *entity.AttributeSpec) string {
	if fixedSpecId > 0 {
		fieldMap := md.GetFixedAttributes(fixedSpecId)
		if fieldMap != nil && len(fieldMap) > 0 {
			alias := attributeSpec.Alias
			if alias == "" {
				alias = xstrings.FirstRuneToUpper(attributeSpec.Kind)
			}
			_, ok := fieldMap[alias]
			if ok {
				return alias
			}
		}
	}

	return ""
}

//返回缓存的字段名map，可用于判断字段是否存在
func (md *MetaDefinition) GetFixedAttributes(fixedSpecId uint64) map[string]string {
	key := fmt.Sprintf("FixedRoleSpec:Alias:%v", fixedSpecId)
	o, ok := MemCache.Get(key)
	if ok {
		return o.(map[string]string)
	}

	return nil
}

func (md *MetaDefinition) GetConnection(specType string, parentId uint64, subId uint64, effectiveDate *time.Time) *entity.ConnectionSpec {
	key := fmt.Sprintf("ConnectionSpec:SpecType:%v:ParentSpecId:%v", specType, parentId)
	v, ok := MemCache.Get(key)
	if ok {
		cons := v.(*map[uint64][]entity.ConnectionSpec)
		if cons != nil {
			v, ok = (*cons)[subId]
			if ok {
				cs := v.([]entity.ConnectionSpec)
				for _, c := range cs {
					effective := isEffectiveSpec(c.StartDate, c.EndDate, effectiveDate)
					if effective {
						return &c
					}
				}
			}
		}
	}

	return nil
}

func loadRoleSpec() {
	roleSpec := entity.RoleSpec{}
	roleSpec.Status = baseentity.EntityStatus_Effective
	svc := service.GetRoleSpecService()
	svc.Find(&metaDefinition.RoleSpecs, &roleSpec, "SpecId", 0, 0, "")
}

func cacheRoleSpec(roleSpecs []*entity.RoleSpec) {
	for _, roleSpec := range roleSpecs {
		key := fmt.Sprintf("RoleSpec:SpecId:%v", roleSpec.SpecId)
		_, ok := MemCache.Get(key)
		if ok {
			logger.Sugar.Errorf("RoleSpec specId:%v repeat", roleSpec.SpecId)
		} else {
			MemCache.SetDefault(key, roleSpec)
		}

		key = fmt.Sprintf("RoleSpec:Kind:%v", roleSpec.Kind)
		l, ok := MemCache.Get(key)
		var s []*entity.RoleSpec
		if !ok {
			l = make([]*entity.RoleSpec, 0)
		}
		s = l.([]*entity.RoleSpec)
		s = append(s, roleSpec)
		MemCache.SetDefault(key, s)
	}
}

func loadFixedRoleSpec() {
	fixedroleSpec := entity.FixedRoleSpec{}
	fixedroleSpec.Status = baseentity.EntityStatus_Effective
	svc := service.GetFixedRoleSpecService()
	svc.Find(&metaDefinition.FixedRoleSpecs, &fixedroleSpec, "SpecId", 0, 0, "")
}

func cacheFixedRoleSpec(fixedRoleSpecs []*entity.FixedRoleSpec) {
	for _, fixedRoleSpec := range fixedRoleSpecs {
		key := fmt.Sprintf("FixedRoleSpec:SpecId:%v", fixedRoleSpec.SpecId)
		_, ok := MemCache.Get(key)
		if ok {
			logger.Sugar.Errorf("FixedRoleSpec specId:%v repeat", fixedRoleSpec.SpecId)
		} else {
			MemCache.SetDefault(key, fixedRoleSpec)
		}

		key = fmt.Sprintf("FixedRoleSpec:Alias:%v", fixedRoleSpec.SpecId)
		if fixedRoleSpec.FixedServiceName != "" {
			v := container.GetService(fixedRoleSpec.FixedServiceName)
			if v != nil {
				fixedSvc := v.(baseservice.BaseService)
				fixedActual, err := fixedSvc.NewEntity(nil)
				if err == nil {
					fieldnameMap, err := reflect.GetFieldNames(fixedActual, false)
					if err == nil && len(fieldnameMap) > 0 {
						_, ok = MemCache.Get(key)
						if ok {
							logger.Sugar.Errorf("FixedRoleSpec specId:%v repeat", fixedRoleSpec.SpecId)
						} else {
							MemCache.SetDefault(key, fieldnameMap)
						}
					}
				}
			}
		}
	}
}

func loadAttributeSpec() {
	attributeSpec := entity.AttributeSpec{}
	attributeSpec.Status = baseentity.EntityStatus_Effective
	svc := service.GetAttributeSpecService()
	svc.Find(&metaDefinition.AttributeSpecs, &attributeSpec, "SpecId", 0, 0, "")
}

func cacheAttributeSpec(attributeSpecs []*entity.AttributeSpec) {
	for _, a := range attributeSpecs {
		key := fmt.Sprintf("AttributeSpec:SpecId:%v", a.SpecId)
		_, ok := MemCache.Get(key)
		if ok {
			logger.Sugar.Errorf("AttributeSpec specId:%v:%v repeat", a.SpecId, a.Kind)
		} else {
			MemCache.SetDefault(key, a)
		}

		key = "AttributeSpec:Kind:" + a.Kind
		l, ok := MemCache.Get(key)
		var s []*entity.AttributeSpec
		if !ok {
			l = make([]*entity.AttributeSpec, 0)
		}
		s = l.([]*entity.AttributeSpec)
		s = append(s, a)
		MemCache.SetDefault(key, s)
	}
}

func loadActionSpec() {
	actionSpec := entity.ActionSpec{}
	actionSpec.Status = baseentity.EntityStatus_Effective
	svc := service.GetActionSpecService()
	svc.Find(&metaDefinition.ActionSpecs, &actionSpec, "SpecId", 0, 0, "")
}

func cacheActionSpec(actionSpecs []*entity.ActionSpec) {
	for _, a := range actionSpecs {
		key := fmt.Sprintf("ActionSpec:SpecId:%v", a.SpecId)
		_, ok := MemCache.Get(key)
		if ok {
			logger.Sugar.Errorf("ActionSpec specId:%v repeat", a.SpecId)
		} else {
			MemCache.SetDefault(key, a)
		}

		key = "ActionSpec:Kind:" + a.Kind
		l, ok := MemCache.Get(key)
		var s []*entity.ActionSpec
		if !ok {
			l = make([]*entity.ActionSpec, 0)
		}
		s = l.([]*entity.ActionSpec)
		s = append(s, a)
		MemCache.SetDefault(key, s)
	}
}

func loadConnectionSpec() {
	connectionSpec := entity.ConnectionSpec{}
	connectionSpec.Status = baseentity.EntityStatus_Effective
	svc := service.GetConnectionSpecService()
	svc.Find(&metaDefinition.ConnectionSpecs, &connectionSpec, "ParentSpecId,SubSpecId", 0, 0, "")
}

func cacheConnectionSpec(connectionSpecs []*entity.ConnectionSpec) {
	for _, c := range connectionSpecs {
		key := fmt.Sprintf("ConnectionSpec:SpecType:%v:ParentSpecId:%v", c.SpecType, c.ParentSpecId)
		v, ok := MemCache.Get(key)
		var m map[uint64][]*entity.ConnectionSpec
		if ok {
			m = v.(map[uint64][]*entity.ConnectionSpec)
		} else {
			m = make(map[uint64][]*entity.ConnectionSpec, 0)
		}
		var s []*entity.ConnectionSpec
		v, ok = m[c.SubSpecId]
		if ok {
			s = v.([]*entity.ConnectionSpec)
		} else {
			s = make([]*entity.ConnectionSpec, 0)
		}
		s = append(s, c)
		m[c.SubSpecId] = s
		MemCache.SetDefault(key, m)
	}
}

func init() {
	Load()
}

func Load() {
	MemCache.Flush()
	metaDefinition.RoleSpecs = make([]*entity.RoleSpec, 0)
	metaDefinition.FixedRoleSpecs = make([]*entity.FixedRoleSpec, 0)
	metaDefinition.AttributeSpecs = make([]*entity.AttributeSpec, 0)
	metaDefinition.ActionSpecs = make([]*entity.ActionSpec, 0)
	metaDefinition.ConnectionSpecs = make([]*entity.ConnectionSpec, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadRoleSpec()
		cacheRoleSpec(metaDefinition.RoleSpecs)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadActionSpec()
		cacheActionSpec(metaDefinition.ActionSpecs)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadFixedRoleSpec()
		cacheFixedRoleSpec(metaDefinition.FixedRoleSpecs)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadAttributeSpec()
		cacheAttributeSpec(metaDefinition.AttributeSpecs)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadConnectionSpec()
		cacheConnectionSpec(metaDefinition.ConnectionSpecs)
	}()
	wg.Wait()

	logger.Sugar.Infof("MetaDefinition load completed!")
}
