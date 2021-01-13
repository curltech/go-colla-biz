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
type PermissionService struct {
	service.OrmBaseService
}

var permissionService = &PermissionService{}

func GetPermissionService() *PermissionService {
	return permissionService
}

func (this *PermissionService) GetSeqName() string {
	return seqname
}

func (this *PermissionService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Permission{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PermissionService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Permission, 0)
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
	service.GetSession().Sync(new(entity.Permission))

	permissionService.OrmBaseService.GetSeqName = permissionService.GetSeqName
	permissionService.OrmBaseService.FactNewEntity = permissionService.NewEntity
	permissionService.OrmBaseService.FactNewEntities = permissionService.NewEntities
	container.RegistService("permission", permissionService)
}
