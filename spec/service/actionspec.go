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
type ActionSpecService struct {
	service.OrmBaseService
}

var actionSpecService = &ActionSpecService{}

func GetActionSpecService() *ActionSpecService {
	return actionSpecService
}

func (this *ActionSpecService) GetSeqName() string {
	return seqname
}

func (this *ActionSpecService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.ActionSpec{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ActionSpecService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.ActionSpec, 0)
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
	service.GetSession().Sync(new(entity.ActionSpec))

	actionSpecService.OrmBaseService.GetSeqName = actionSpecService.GetSeqName
	actionSpecService.OrmBaseService.FactNewEntity = actionSpecService.NewEntity
	actionSpecService.OrmBaseService.FactNewEntities = actionSpecService.NewEntities
	container.RegistService("actionSpec", actionSpecService)
}
