package controller

import (
	"github.com/curltech/go-colla-biz/controller"
	"github.com/curltech/go-colla-biz/stock/entity"
	"github.com/curltech/go-colla-biz/stock/service"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/util/message"
)

/**
控制层代码需要做数据转换，调用服务层的代码，由于数据转换的结构不一致，因此每个实体（外部rest方式访问）的控制层都需要写一遍
*/
type DayDataController struct {
	controller.BaseController
}

var dayDataController *DayDataController

func (this *DayDataController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.DayData, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

/**
注册bean管理器，注册序列
*/
func init() {
	dayDataController = &DayDataController{
		BaseController: controller.BaseController{
			BaseService: service.GetDayDataService(),
		},
	}
	dayDataController.BaseController.ParseJSON = dayDataController.ParseJSON
	container.RegistController("stkDayData", dayDataController)
}
