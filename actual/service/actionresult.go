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
type ActionResultService struct {
	service.OrmBaseService
}

var actionResultService = &ActionResultService{}

func GetActionResultService() *ActionResultService {
	return actionResultService
}

func (this *ActionResultService) GetSeqName() string {
	return seqname
}

func (this *ActionResultService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.ActionResult{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ActionResultService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.ActionResult, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *ActionResultService) ListByParentId(schemaName string, parentIds []uint64) []*entity.ActionResult {
	idss := Split(parentIds)
	actionResults := make([]*entity.ActionResult, 0)
	for _, ids := range idss {
		as := make([]*entity.ActionResult, 0)
		conds := fmt.Sprintf("ParentId in (%s)", service.PlaceQuestionMark(len(ids)))
		params := make([]interface{}, len(ids))
		for i, id := range ids {
			params[i] = id
		}
		err := this.Find(&as, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
		if err == nil && len(as) > 0 {
			actionResults = append(actionResults, as...)
		}
	}
	if len(actionResults) > 0 {
		return actionResults
	}

	return nil
}

func init() {
	service.GetSession().Sync(new(entity.ActionResult))

	actionResultService.OrmBaseService.GetSeqName = actionResultService.GetSeqName
	actionResultService.OrmBaseService.FactNewEntity = actionResultService.NewEntity
	actionResultService.OrmBaseService.FactNewEntities = actionResultService.NewEntities
	container.RegistService("atlActionResult", actionResultService)
}
