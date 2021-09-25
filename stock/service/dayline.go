package service

import (
	"bytes"
	"encoding/binary"
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
func (this *DayDataService) ParsePath(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for _, file := range files {
		filename := file.Name()
		hasSuffix := strings.HasSuffix(filename, ".day")
		if hasSuffix {
			shareId := strings.TrimSuffix(filename, ".day")
			logger.Sugar.Infof("shareId:", shareId)
			content, err := ioutil.ReadFile(dirname + "/" + filename)
			if err != nil {
				return err
			}
			this.ParseByte(shareId, content)
		}
	}
	return nil
}

func (this *DayDataService) ParseByte(shareId string, content []byte) {
	for i := 0; i < len(content); i = i + 32 {
		dayData := entity.DayData{}
		dayData.ShareId = shareId
		dayData.DayDate = bytesToInt(content[i : i+4])
		dayData.OpeningPrice = bytesToInt(content[i+4 : i+8])
		dayData.CeilingPrice = bytesToInt(content[i+8 : i+12])
		dayData.FloorPrice = bytesToInt(content[i+12 : i+16])
		dayData.ClosingPrice = bytesToInt(content[i+16 : i+20])
		dayData.TurnVolume = bytesToFloat(content[i+20 : i+24])
		dayData.Volume = bytesToInt(content[i+24 : i+28])
		logger.Sugar.Infof("DayData:%s", dayData)
		this.Insert(dayData)
	}
}

func bytesToInt(bys []byte) int32 {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

func bytesToFloat(bys []byte) float32 {
	bytebuff := bytes.NewBuffer(bys)
	var data float32
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

func init() {
	service.GetSession().Sync(new(entity.DayData))
	dayDataService.OrmBaseService.GetSeqName = dayDataService.GetSeqName
	dayDataService.OrmBaseService.FactNewEntity = dayDataService.NewEntity
	dayDataService.OrmBaseService.FactNewEntities = dayDataService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("stkDayData", dayDataService)
}
