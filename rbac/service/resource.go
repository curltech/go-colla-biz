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
type ResourceService struct {
	service.OrmBaseService
}

var resourceService = &ResourceService{}

func GetResourceService() *ResourceService {
	return resourceService
}

func (this *ResourceService) GetSeqName() string {
	return seqname
}

func (this *ResourceService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Resource{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ResourceService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Resource, 0)
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
	service.GetSession().Sync(new(entity.Resource))

	resourceService.OrmBaseService.GetSeqName = resourceService.GetSeqName
	resourceService.OrmBaseService.FactNewEntity = resourceService.NewEntity
	resourceService.OrmBaseService.FactNewEntities = resourceService.NewEntities
	container.RegistService("rsource", resourceService)
}
