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
type ConnectionService struct {
	service.OrmBaseService
}

var connectionService = &ConnectionService{}

func GetConnectionService() *ConnectionService {
	return connectionService
}

func (this *ConnectionService) GetSeqName() string {
	return seqname
}

func (this *ConnectionService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Connection{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *ConnectionService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Connection, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *ConnectionService) ListByParentId(schemaName string, parentIds []uint64) []*entity.Connection {
	idss := Split(parentIds)
	connections := make([]*entity.Connection, 0)
	for _, ids := range idss {
		cs := make([]*entity.Connection, 0)
		conds := fmt.Sprintf("ParentId in (%s)", service.PlaceQuestionMark(len(ids)))
		params := make([]interface{}, len(ids))
		for i, id := range ids {
			params[i] = id
		}
		err := this.Find(&cs, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
		if err == nil && len(cs) > 0 {
			connections = append(connections, cs...)
		}
	}
	if len(connections) > 0 {
		return connections
	}

	return nil
}

func init() {
	service.GetSession().Sync(new(entity.Connection))

	connectionService.OrmBaseService.GetSeqName = connectionService.GetSeqName
	connectionService.OrmBaseService.FactNewEntity = connectionService.NewEntity
	connectionService.OrmBaseService.FactNewEntities = connectionService.NewEntities
	container.RegistService("atlConnection", connectionService)
}
