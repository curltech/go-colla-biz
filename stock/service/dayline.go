package service

import (
	"github.com/curltech/go-colla-biz/stock/entity"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
	"io/ioutil"
	"strings"
)

/**
同步表结构，服务继承基本服务的方法
*/
type DayDataService struct {
	service.OrmBaseService
}

var dayDataService = &DayDataService{}

func GetDayDataService() *DayDataService {
	return dayDataService
}

func (this *DayDataService) GetSeqName() string {
	return seqname
}

func (this *DayDataService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Share{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *DayDataService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.DayData, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

/**
读目录下的数据
*/
func (this *DayDataService) parse(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".day")
		if hasSuffix {
			//shareId := strings.TrimSuffix(filename, ".day")
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			for i := 0; i < len(content); i = i + 32 {
				dayDate := string(content[i])
				logger.Sugar.Infof(dayDate)
			}
		}
	}
	return nil
}

func init() {
	service.GetSession().Sync(new(entity.DayData))
	dayDataService.OrmBaseService.GetSeqName = dayDataService.GetSeqName
	dayDataService.OrmBaseService.FactNewEntity = dayDataService.NewEntity
	dayDataService.OrmBaseService.FactNewEntities = dayDataService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("stkDayData", dayDataService)
}
