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
type RoleSpecService struct {
	service.OrmBaseService
}

var roleSpecService = &RoleSpecService{}

func GetRoleSpecService() *RoleSpecService {
	return roleSpecService
}

var seqname = "seq_spec"

func (this *RoleSpecService) GetSeqName() string {
	return seqname
}

func (this *RoleSpecService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.RoleSpec{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *RoleSpecService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.RoleSpec, 0)
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
	service.GetSession().Sync(new(entity.RoleSpec))

	roleSpecService.OrmBaseService.GetSeqName = roleSpecService.GetSeqName
	roleSpecService.OrmBaseService.FactNewEntity = roleSpecService.NewEntity
	roleSpecService.OrmBaseService.FactNewEntities = roleSpecService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("roleSpec", roleSpecService)
}
