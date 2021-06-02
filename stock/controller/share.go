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
type ShareController struct {
	controller.BaseController
}

var shareController *ShareController

func (this *ShareController) ParseJSON(json []byte) (interface{}, error) {
	var entities = make([]*entity.Share, 0)
	err := message.Unmarshal(json, &entities)

	return &entities, err
}

/**
注册bean管理器，注册序列
*/
func init() {
	shareController = &ShareController{
		BaseController: controller.BaseController{
			BaseService: service.GetShareService(),
		},
	}
	shareController.BaseController.ParseJSON = shareController.ParseJSON
	container.RegistController("stkShare", shareController)
}
