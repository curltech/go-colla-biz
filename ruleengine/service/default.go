package service

import (
	"github.com/curltech/go-colla-biz/ruleengine/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type RuleDefinitionService struct {
	service.OrmBaseService
}

var ruleDefinitionService = &RuleDefinitionService{}

func GetRuleDefinitionService() *RuleDefinitionService {
	return ruleDefinitionService
}

var seqname = "seq_rule"

func (this *RuleDefinitionService) GetSeqName() string {
	return seqname
}

func (this *RuleDefinitionService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.RuleDefinition{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *RuleDefinitionService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.RuleDefinition, 0)
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
	service.GetSession().Sync(new(entity.RuleDefinition))

	ruleDefinitionService.OrmBaseService.GetSeqName = ruleDefinitionService.GetSeqName
	ruleDefinitionService.OrmBaseService.FactNewEntity = ruleDefinitionService.NewEntity
	ruleDefinitionService.OrmBaseService.FactNewEntities = ruleDefinitionService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("ruleDefinition", ruleDefinitionService)
}
