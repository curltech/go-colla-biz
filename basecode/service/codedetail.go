package service

import (
	"github.com/curltech/go-colla-biz/basecode/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type CodeDetailService struct {
	service.OrmBaseService
}

var codeDetailService = &CodeDetailService{}

func GetCodeDetailService() *CodeDetailService {
	return codeDetailService
}

func (this *CodeDetailService) GetSeqName() string {
	return seqname
}

func (this *CodeDetailService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.CodeDetail{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *CodeDetailService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.CodeDetail, 0)
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
	service.GetSession().Sync(new(entity.CodeDetail))

	codeDetailService.OrmBaseService.GetSeqName = codeDetailService.GetSeqName
	codeDetailService.OrmBaseService.FactNewEntity = codeDetailService.NewEntity
	codeDetailService.OrmBaseService.FactNewEntities = codeDetailService.NewEntities
	container.RegistService("codeDetail", codeDetailService)
}
