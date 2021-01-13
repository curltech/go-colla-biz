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
type ConnectionSpecService struct {
	service.OrmBaseService
}

var connectionSpecService = &ConnectionSpecService{}

func GetConnectionSpecService() *ConnectionSpecService {
	return connectionSpecService
}

func (this *ConnectionSpecService) GetSeqName() string {
	return seqname
}

func (this *ConnectionSpecService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.ConnectionSpec{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ConnectionSpecService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.ConnectionSpec, 0)
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
	service.GetSession().Sync(new(entity.ConnectionSpec))

	connectionSpecService.OrmBaseService.GetSeqName = connectionSpecService.GetSeqName
	connectionSpecService.OrmBaseService.FactNewEntity = connectionSpecService.NewEntity
	connectionSpecService.OrmBaseService.FactNewEntities = connectionSpecService.NewEntities
	container.RegistService("connectionSpec", connectionSpecService)
}
