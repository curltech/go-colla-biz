package service

import (
	"github.com/curltech/go-colla-biz/rbac/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type GroupService struct {
	service.OrmBaseService
}

var groupService = &GroupService{}

func GetGroupService() *GroupService {
	return groupService
}

func (this *GroupService) GetSeqName() string {
	return seqname
}

func (this *GroupService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Group{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *GroupService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Group, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *GroupService) FindRolesByUserId(userId string) (string, error) {
	groups := make([]*entity.Group, 0)
	condiBean := &entity.Group{UserId: userId}
	err := this.Find(&groups, condiBean, "", 0, 0, "", nil)
	if err != nil {
		logger.Sugar.Errorf("FindByUserId error:%v", err.Error())
		return "", err
	}
	roles := ""
	i := 0
	for _, group := range groups {
		if i > 0 {
			roles += ","
		}
		roles += group.RoleId
		i++
	}
	return roles, nil
}

func init() {
	service.GetSession().Sync(new(entity.Group))

	groupService.OrmBaseService.GetSeqName = groupService.GetSeqName
	groupService.OrmBaseService.FactNewEntity = groupService.NewEntity
	groupService.OrmBaseService.FactNewEntities = groupService.NewEntities
	container.RegistService("group", groupService)
}
