package controller

import (
	"curltech.io/camsi/camsi-node/p2p/chain/handler"
	"curltech.io/camsi/camsi-node/p2p/chain/service"
	"curltech.io/camsi/camsi-node/p2p/msg"
	"errors"
	"fmt"
	"github.com/curltech/go-colla-core/config"
	"github.com/curltech/go-colla-core/container"
	"github.com/curltech/go-colla-core/crypto"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/debug"
	"github.com/curltech/go-colla-core/util/message"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"mime/multipart"
	"strings"
)

var TemplateParams = make(map[string]interface{})

func RegistTemplateParam(key string, param interface{}) error {
	_, ok := TemplateParams[key]
	if !ok {
		TemplateParams[key] = param

		return nil
	}

	return errors.New("Exist")
}

func init() {
	RegistTemplateParam("Host", config.ServerParams.Name)
}

func HTMLController(ctx iris.Context) {
	name := ctx.Params().GetString("name")
	name = strings.ReplaceAll(name, ".html", "")
	fn := debug.Trace("render view:" + name)
	defer fn()
	ctx.View(name, TemplateParams)
}

func MainController(ctx iris.Context) {
	serviceName := ctx.Params().Get("serviceName")
	methodName := ctx.Params().Get("methodName")
	logger.Sugar.Infof("MainController call %v.%v", serviceName, methodName)
	args := make([]interface{}, 1)
	args[0] = ctx
	msg := fmt.Sprintf("call servicename:%v methodName:%v", serviceName, methodName)
	fn := debug.Trace(msg)
	defer fn()
	controller := container.GetController(serviceName)
	if controller == nil {
		panic("NoController")
	}
	reflect.Call(controller, methodName, args)
}

/**
接收Receive的p2p chain协议请求，消息是ChainMessage的格式
*/
func ReceiveController(ctx iris.Context) {
	chainMessage := &msg.ChainMessage{}
	err := ctx.ReadJSON(chainMessage)
	if err != nil {
		ctx.JSON(err.Error())
	} else {
		chainMessage.LocalConnectAddress = ctx.RemoteAddr()
		response, err := service.Receive(chainMessage)
		handler.SetResponse(chainMessage, response)
		if err != nil {
			ctx.JSON(err.Error())
		} else {
			ctx.JSON(response)
		}
	}
}

////////////////////

/**
接收Receive的p2p chain协议请求，消息是PCChainMessage的格式
*/
func ReceivePCController(ctx iris.Context) {
	chainMessage := &msg.PCChainMessage{}
	err := ctx.ReadJSON(chainMessage)
	if err != nil {
		ctx.JSON(err.Error())
	} else {
		securityContext := &crypto.SecurityContext{}
		err = message.TextUnmarshal(chainMessage.SecurityContextString, securityContext)
		if err != nil {
			ctx.JSON(err.Error())
		} else {
			chainMessage.SecurityContext = securityContext
			response, err := service.ReceivePC(chainMessage)
			if err != nil {
				ctx.JSON(err.Error())
			} else {
				if response != nil {
					chainMessage, err = handler.EncryptPC(response)
					if err != nil {
						ctx.JSON(err.Error())
					}
					ctx.JSON(chainMessage)
				}
			}
		}
	}
}

type UploadParam struct {
	ServiceName string
}

const postMaxSize = 256 * iris.MB

func UploadController(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(postMaxSize)
	var serviceName = ctx.PostValue("serviceName")
	var methodName = ctx.PostValue("methodName")
	if serviceName == "" || methodName == "" {
		ctx.StopWithJSON(iris.StatusInternalServerError, "NoService")

		return
	}
	logger.Sugar.Infof("UploadController call %v.%v", serviceName, methodName)

	err := ctx.Request().ParseMultipartForm(postMaxSize)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	}
	form := ctx.Request().MultipartForm
	if form == nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, "BlankForm")

		return
	}
	if form.File == nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, "NilFormFile")

		return
	}
	var files = make([]multipart.File, 0)
	for _, heads := range form.File {
		for _, head := range heads {
			logger.Sugar.Infof("file:%v", head.Filename)
			file, err := head.Open()
			if err != nil {
				ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

				return
			}
			files = append(files, file)
			/**
			下面是读取数据到[]byte
			*/
			//buf, err := ioutil.ReadAll(file)
			//if err != nil {
			//	ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())
			//
			//	return
			//}
			defer file.Close()
		}
	}

	args := make([]interface{}, 1)
	args[0] = files
	msg := fmt.Sprintf("call servicename:%v methodName:%v", serviceName, methodName)
	fn := debug.Trace(msg)
	defer fn()
	svc := container.GetService(serviceName)
	if svc == nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, "NoService")

		return
	}
	result, err := reflect.Call(svc, methodName, args)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, err.Error())

		return
	} else {
		ctx.JSON(result)
	}
}

func DownloadController(ctx iris.Context) {
	params := make(map[string]interface{}, 0)
	ctx.ReadJSON(&params)

	serviceName := params["serviceName"]
	methodName := params["methodName"]
	destName := params["destName"]

	logger.Sugar.Infof("DownloadController call %v.%v", serviceName, methodName)
	ctx.ContentType("")
	ctx.ResponseWriter().Header().Set(context.ContentDispositionHeaderKey, "attachment;filename="+destName.(string))
	args := make([]interface{}, 2)
	args[0] = params["condiBean"]
	msg := fmt.Sprintf("call servicename:%v methodName:%v", serviceName, methodName)
	fn := debug.Trace(msg)
	defer fn()
	svc := container.GetService(serviceName.(string))
	if svc == nil {
		panic("NoService")
	}
	result, err := reflect.Call(svc, methodName.(string), args)
	if err != nil {
		ctx.JSON(err.Error())
	} else {
		ctx.ResponseWriter().Write(result[0].([]byte))
		ctx.ResponseWriter().Flush()
	}
}

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

func wrapValidationErrors(errs validator.ValidationErrors) []validationError {
	validationErrors := make([]validationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, validationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}

	return validationErrors
}
