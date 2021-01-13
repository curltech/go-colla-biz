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
type AttributeSpecService struct {
	service.OrmBaseService
}

var attributeSpecService = &AttributeSpecService{}

func GetAttributeSpecService() *AttributeSpecService {
	return attributeSpecService
}

func (this *AttributeSpecService) GetSeqName() string {
	return seqname
}

func (this *AttributeSpecService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.AttributeSpec{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *AttributeSpecService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.AttributeSpec, 0)
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
	service.GetSession().Sync(new(entity.AttributeSpec))

	attributeSpecService.OrmBaseService.GetSeqName = attributeSpecService.GetSeqName
	attributeSpecService.OrmBaseService.FactNewEntity = attributeSpecService.NewEntity
	attributeSpecService.OrmBaseService.FactNewEntities = attributeSpecService.NewEntities
	container.RegistService("attributeSpec", attributeSpecService)
}
