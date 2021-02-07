package service

import (
	"errors"
	"github.com/curltech/go-colla-biz/gui/entity"
	"github.com/curltech/go-colla-core/cache"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type GuiMenuService struct {
	service.OrmBaseService
}

var guiMenuService = &GuiMenuService{}

func GetGuiMenuService() *GuiMenuService {
	return guiMenuService
}

var seqname = "seq_gui"

func (this *GuiMenuService) GetSeqName() string {
	return seqname
}

func (this *GuiMenuService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.GuiMenu{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *GuiMenuService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.GuiMenu, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *GuiMenuService) GetGuiMenu(menuId string) (*entity.GuiMenu, error) {
	key := "GuiMenu:" + menuId
	v, ok := MemCache.Get(key)
	if !ok {
		return nil, errors.New("NotFoundGuiMenu")
	}
	guiMenu, ok := v.(*entity.GuiMenu)
	if !ok {
		return nil, errors.New("WrongTypeGuiMenu")
	}
	key = "GuiMenuChildren:" + menuId
	v, ok = MemCache.Get(key)
	if !ok {
		return nil, errors.New("NotFoundGuiMenu")
	}
	guiMenus, ok := v.([]*entity.GuiMenu)
	if !ok {
		return nil, errors.New("WrongTypeGuiMenu")
	}
	guiMenu.Children = guiMenus

	return guiMenu, nil
}

var MemCache = cache.NewMemCache("guimenu", 1000, 1000)

func loadGuiMenu() {
	guiMenus := make([]*entity.GuiMenu, 0)
	guiMenu := &entity.GuiMenu{}
	guiMenu.Status = baseentity.EntityStatus_Effective
	guiMenuService.Find(&guiMenus, guiMenu, "MenuId,SerialId", 0, 0, "")
	for _, guiMenu := range guiMenus {
		key := "GuiMenu:" + guiMenu.MenuId
		MemCache.SetDefault(key, guiMenu)

		key = "GuiMenuChildren:" + guiMenu.ParentId
		var menus []*entity.GuiMenu
		v, ok := MemCache.Get(key)
		if ok {
			menus = v.([]*entity.GuiMenu)
		} else {
			menus = make([]*entity.GuiMenu, 0)
		}
		menus = append(menus, guiMenu)
		MemCache.SetDefault(key, menus)
	}
	logger.Sugar.Infof("GuiMenu load completed!")
}

func Load() {
	MemCache.Flush()
	go loadGuiMenu()
}

func init() {
	service.GetSession().Sync(new(entity.GuiMenu))

	guiMenuService.OrmBaseService.GetSeqName = guiMenuService.GetSeqName
	guiMenuService.OrmBaseService.FactNewEntity = guiMenuService.NewEntity
	guiMenuService.OrmBaseService.FactNewEntities = guiMenuService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("guiMenu", guiMenuService)

	Load()
}
