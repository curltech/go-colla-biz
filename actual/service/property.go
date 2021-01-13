package service

import (
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type PropertyService struct {
	service.OrmBaseService
}

var propertyService = &PropertyService{}

func GetPropertyService() *PropertyService {
	return propertyService
}

func (this *PropertyService) GetSeqName() string {
	return seqname
}

func (this *PropertyService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Property{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *PropertyService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Property, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *PropertyService) ListByParentId(schemaName string, parentIds []uint64) []*entity.Property {
	idss := Split(parentIds)
	properties := make([]*entity.Property, 0)
	for _, ids := range idss {
		ps := make([]*entity.Property, 0)
		conds := fmt.Sprintf("ParentId in (%s)", service.PlaceQuestionMark(len(ids)))
		params := make([]interface{}, len(ids))
		for i, id := range ids {
			params[i] = id
		}
		err := this.Find(&ps, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
		if err == nil && len(ps) > 0 {
			properties = append(properties, ps...)
		}
	}
	if len(properties) > 0 {
		return properties
	}

	return nil
}

func init() {
	service.GetSession().Sync(new(entity.Property))

	propertyService.OrmBaseService.GetSeqName = propertyService.GetSeqName
	propertyService.OrmBaseService.FactNewEntity = propertyService.NewEntity
	propertyService.OrmBaseService.FactNewEntities = propertyService.NewEntities
	container.RegistService("atlProperty", propertyService)
}
