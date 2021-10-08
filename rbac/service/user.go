package service

import (
	"errors"
	"fmt"
	"github.com/curltech/go-colla-biz/rbac/entity"
	"github.com/curltech/go-colla-core/cache"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/crypto/std"
	entity2 "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

var MemCache = cache.NewMemCache("rbac", 1, 10)

/**
同步表结构，服务继承基本服务的方法
*/
type UserService struct {
	service.OrmBaseService
}

var userService = &UserService{}

func GetUserService() *UserService {
	return userService
}

var seqname = "seq_rbac"

func (this *UserService) GetSeqName() string {
	return seqname
}

func (this *UserService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.User{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *UserService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.User, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *UserService) getCacheKey(key string) string {
	return "user:" + key
}

func (this *UserService) Logout(userName string) {
	user := this.GetUser(userName)
	if user != nil {
		MemCache.Delete(this.getCacheKey(user.UserName))
	}
}

func (this *UserService) Regist(user *entity.User) (*entity.User, error) {
	var err error
	old := &entity.User{}
	old.UserName = user.UserName
	if old.UserName == "" || user.Password == "" {
		return nil, errors.New("NoUserName")
	}
	ok, _ := this.Get(old, false, "", "")
	if ok {
		return nil, errors.New("ExistUserName")
	}
	user.Status = entity2.EntityStatus_Draft
	id := this.GetSeq()
	user.Id = id
	user.UserId = fmt.Sprintf("%v", id)
	if user.PlainPassword == user.ConfirmPassword {
		this.EncryptPassword(user)
		affected, _ := this.Insert(user)
		if affected <= 0 {
			err = errors.New("ErrorInsert")
		}
	} else {
		err = errors.New("UnmatchPassword")
	}

	return user, err
}

func (this *UserService) UpdateStatus(userName string, status string) {
	user := &entity.User{}
	user.UserName = userName
	ok, _ := this.Get(user, false, "", "")
	if ok {
		this.Update(user, []string{"Status"}, "")
	}
}

func (this *UserService) EncryptPassword(user *entity.User) {
	user.Password = std.EncodeBase64(std.Hash(user.UserName+user.PlainPassword, "sha3_256"))
}

func (this *UserService) Auth(userName string, password string) (*entity.User, error) {
	user := &entity.User{}
	user.UserName = userName
	user.PlainPassword = password
	user.Status = entity.UserStatus_Enabled
	this.EncryptPassword(user)
	ok, _ := this.Get(user, false, "", "")
	if !ok {
		logger.Sugar.Errorf("%v auth fail!", userName)

		return nil, errors.New("AuthFail")
	}

	return user, nil
}

func (this *UserService) Login(userName string, password string) (*entity.User, error) {
	this.Logout(userName)
	user, err := this.Auth(userName, password)
	if err != nil {
		return nil, err
	}
	key := this.getCacheKey(user.UserName)
	MemCache.SetDefault(key, user)
	logger.Sugar.Infof("%v successfully login!", userName)

	return user, nil
}

func (this *UserService) GetUser(userName string) *entity.User {
	key := this.getCacheKey(userName)
	v, ok := MemCache.Get(key)
	if ok {
		user := v.(*entity.User)

		return user
	}
	user := &entity.User{}
	user.UserName = userName
	user.Status = entity.UserStatus_Enabled
	ok, _ = this.Get(user, false, "", "")
	if ok {
		return user
	}

	return nil
}

func init() {
	service.GetSession().Sync(new(entity.User))

	userService.OrmBaseService.GetSeqName = userService.GetSeqName
	userService.OrmBaseService.FactNewEntity = userService.NewEntity
	userService.OrmBaseService.FactNewEntities = userService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("user", userService)
}
