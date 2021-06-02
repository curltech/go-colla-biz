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
type ShareService struct {
	service.OrmBaseService
}

var shareService = &ShareService{}

func GetShareService() *ShareService {
	return shareService
}

var seqname = "seq_stock"

func (this *ShareService) GetSeqName() string {
	return seqname
}

func (this *ShareService) NewEntity(data []byte) (interface{}, error) {
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

func (this *ShareService) NewEntities(data []byte) (interface{}, error) {
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

func init() {
	service.GetSession().Sync(new(entity.Share))
	shareService.OrmBaseService.GetSeqName = shareService.GetSeqName
	shareService.OrmBaseService.FactNewEntity = shareService.NewEntity
	shareService.OrmBaseService.FactNewEntities = shareService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("stkShare", shareService)
}
