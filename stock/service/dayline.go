package service

import (
	"github.com/curltech/go-colla-biz/stock/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type DayDataService struct {
	service.OrmBaseService
}

var dayDataService = &DayDataService{}

func GetDayDataService() *DayDataService {
	return dayDataService
}

func (this *DayDataService) GetSeqName() string {
	return seqname
}

func (this *DayDataService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Share{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *DayDataService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Share, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *DayDataService) parse(dir string) error {

	return nil
}

func init() {
	service.GetSession().Sync(new(entity.DayData))
	dayDataService.OrmBaseService.GetSeqName = dayDataService.GetSeqName
	dayDataService.OrmBaseService.FactNewEntity = dayDataService.NewEntity
	dayDataService.OrmBaseService.FactNewEntities = dayDataService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("stkDayData", dayDataService)
}
