package service

import (
	"github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type FixedRoleSpecService struct {
	service.OrmBaseService
}

var fixedRoleSpecService = &FixedRoleSpecService{}

func GetFixedRoleSpecService() *FixedRoleSpecService {
	return fixedRoleSpecService
}

func (this *FixedRoleSpecService) GetSeqName() string {
	return seqname
}

func (this *FixedRoleSpecService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.FixedRoleSpec{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *FixedRoleSpecService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.FixedRoleSpec, 0)
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
	service.GetSession().Sync(new(entity.FixedRoleSpec))

	fixedRoleSpecService.OrmBaseService.GetSeqName = fixedRoleSpecService.GetSeqName
	fixedRoleSpecService.OrmBaseService.FactNewEntity = fixedRoleSpecService.NewEntity
	fixedRoleSpecService.OrmBaseService.FactNewEntities = fixedRoleSpecService.NewEntities
	container.RegistService("fixedRoleSpec", fixedRoleSpecService)
}
