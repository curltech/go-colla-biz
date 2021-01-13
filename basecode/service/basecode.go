package service

import (
	"errors"
	"github.com/curltech/go-colla-biz/basecode/entity"
	"github.com/curltech/go-colla-core/cache"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/kataras/golog"
	"sync"
)

/**
同步表结构，服务继承基本服务的方法
*/
type BaseCodeService struct {
	service.OrmBaseService
}

var baseCodeService = &BaseCodeService{}

func GetBaseCodeService() *BaseCodeService {
	return baseCodeService
}

var seqname = "seq_basecode"

func (this *BaseCodeService) GetSeqName() string {
	return seqname
}

func (this *BaseCodeService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.BaseCode{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *BaseCodeService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.BaseCode, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

func (this *BaseCodeService) GetBaseCode(baseCodeId string) (*entity.BaseCode, error) {
	key := "BaseCode:" + baseCodeId
	v, ok := MemCache.Get(key)
	if !ok {
		return nil, errors.New("NotFoundBaseCode")
	}
	baseCode, ok := v.(*entity.BaseCode)
	if !ok {
		return nil, errors.New("WrongTypeBaseCode")
	}
	key = "CodeDetail:" + baseCodeId
	v, ok = MemCache.Get(key)
	if !ok {
		return nil, errors.New("NotFoundCodeDetail")
	}
	codeDetails, ok := v.([]*entity.CodeDetail)
	if !ok {
		return nil, errors.New("WrongTypeCodeDetail")
	}
	baseCode.CodeDetails = codeDetails

	return baseCode, nil
}

var MemCache = cache.NewMemCache("basecode", 1000, 1000)

func loadBaseCode() {
	baseCodes := make([]*entity.BaseCode, 0)
	baseCode := &entity.BaseCode{}
	baseCode.Status = baseentity.EntityStatus_Effective
	baseCodeService.Find(&baseCodes, baseCode, "BaseCodeId", 0, 0, "")
	for _, baseCode := range baseCodes {
		key := "BaseCode:" + baseCode.BaseCodeId
		MemCache.SetDefault(key, baseCode)
	}
}

func loadCodeDetail() {
	codeDetails := make([]*entity.CodeDetail, 0)
	codeDetail := &entity.CodeDetail{}
	codeDetail.Status = baseentity.EntityStatus_Effective
	codeDetailService.Find(&codeDetails, codeDetail, "BaseCodeId,SerialId", 0, 0, "")
	for _, codeDetail := range codeDetails {
		key := "CodeDetail:" + codeDetail.BaseCodeId
		var cds []*entity.CodeDetail
		v, ok := MemCache.Get(key)
		if ok {
			cds = v.([]*entity.CodeDetail)
		} else {
			cds = make([]*entity.CodeDetail, 0)
		}
		cds = append(cds, codeDetail)
		MemCache.SetDefault(key, cds)

		if codeDetail.ParentId != "" {
			key = "CodeDetail:" + codeDetail.BaseCodeId + ":" + codeDetail.ParentId
			var cds []*entity.CodeDetail
			v, ok := MemCache.Get(key)
			if ok {
				cds = v.([]*entity.CodeDetail)
			} else {
				cds = make([]*entity.CodeDetail, 0)
			}
			cds = append(cds, codeDetail)
			MemCache.SetDefault(key, cds)
		}
	}
}

func Load() {
	MemCache.Flush()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadBaseCode()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		loadCodeDetail()
	}()
	wg.Wait()

	golog.Infof("BaseCode load completed!")
}

func init() {
	service.GetSession().Sync(new(entity.BaseCode))

	baseCodeService.OrmBaseService.GetSeqName = baseCodeService.GetSeqName
	baseCodeService.OrmBaseService.FactNewEntity = baseCodeService.NewEntity
	baseCodeService.OrmBaseService.FactNewEntities = baseCodeService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("baseCode", baseCodeService)

	Load()
}
