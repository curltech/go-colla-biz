package service

import (
	"github.com/curltech/go-colla-biz/rbac/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type RoleService struct {
	service.OrmBaseService
}

var roleService = &RoleService{}

func GetRoleService() *RoleService {
	return roleService
}

func (this *RoleService) GetSeqName() string {
	return seqname
}

func (this *RoleService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Role{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *RoleService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Role, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func init() {
	service.GetSession().Sync(new(entity.Role))
	roleService.OrmBaseService.GetSeqName = roleService.GetSeqName
	roleService.OrmBaseService.FactNewEntity = roleService.NewEntity
	roleService.OrmBaseService.FactNewEntities = roleService.NewEntities
	container.RegistService("role", roleService)
}
