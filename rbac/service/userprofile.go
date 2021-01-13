package service

import (
	"github.com/curltech/go-colla-biz/rbac/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type UserProfileService struct {
	service.OrmBaseService
}

var userProfileService = &UserProfileService{}

func GetUserProfileService() *UserProfileService {
	return userProfileService
}

func (this *UserProfileService) GetSeqName() string {
	return seqname
}

func (this *UserProfileService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.UserProfile{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *UserProfileService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.UserProfile, 0)
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
	service.GetSession().Sync(new(entity.UserProfile))

	userProfileService.OrmBaseService.GetSeqName = userProfileService.GetSeqName
	userProfileService.OrmBaseService.FactNewEntity = userProfileService.NewEntity
	userProfileService.OrmBaseService.FactNewEntities = userProfileService.NewEntities
	container.RegistService("userProfile", userProfileService)
}
