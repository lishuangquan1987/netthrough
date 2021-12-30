package controllers

import (
	"fmt"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"netthrough.server/business"
	"netthrough.server/models"
)

//client register
func Register(ctx *gin.Context) {
	var request models.RegisterRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(200, models.RegisterResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}
	task := &business.Task{
		ClientId:         request.ClientId,
		ServerListenPort: request.ServerListenPort,
	}
	business.AddTask(task)
	err := task.Start()
	if err != nil {
		ctx.JSON(200, models.RegisterResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.RegisterResponse{
		ClientId:  request.ClientId,
		IsSuccess: true,
	})
}

//client check status to see if has request from outside
func StatusCheck(ctx *gin.Context) {
	var request models.StatusCheckRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(200, models.StatusCheckResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}

	task := From(business.Tasks).WhereT(func(t *business.Task) bool {
		return t.ClientId == request.ClientId
	}).First().(*business.Task)

	if task == nil {
		ctx.JSON(200, models.StatusCheckResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("server has no client named:%s", request.ClientId),
		})
		return
	}
	sessionIds := make([]string, 0)
	//判断WriteBuffer中是否有值
	for guid, ch := range task.WriteBuffer {
		if len(ch) > 0 {
			sessionIds = append(sessionIds, guid)
		}
	}
	ctx.JSON(200, models.StatusCheckResponse{
		ClientId:  request.ClientId,
		IsSuccess: true,
		HasData:   len(sessionIds) > 0,
		SessionId: sessionIds,
	})
}

//client read request from outside
func ReadData(ctx *gin.Context) {
	var request models.ReadDataRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(200, models.ReadDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}

	task := From(business.Tasks).WhereT(func(t *business.Task) bool {
		return t.ClientId == request.ClientId
	}).First().(*business.Task)
	if task == nil {
		ctx.JSON(200, models.ReadDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("server has no client named:%s", request.ClientId),
		})
		return
	}

	ch, ok := task.WriteBuffer[request.SessionId]
	if !ok {
		ctx.JSON(200, models.ReadDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("client in server has no session id:%s", request.SessionId),
		})
		return
	}

	if len(ch) == 0 {
		ctx.JSON(200, models.ReadDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: true,
			HasData:   false,
			ErrMsg:    fmt.Sprintf("session id:%s has no data yet", request.SessionId),
		})
		return
	}

	//读取全部数据，发送到客户端
	data := make([]byte, 0)
	for i := 0; i < len(ch); i++ {
		data = append(data, <-ch...)
	}
	ctx.JSON(200, models.ReadDataResponse{
		ClientId:  request.ClientId,
		SessionId: request.SessionId,
		IsSuccess: true,
		HasData:   true,
		Data:      data,
	})
}

//client write data to outside
func WriteData(ctx *gin.Context) {
	var request models.WriteDataRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(200, models.ReadDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}

	task := From(business.Tasks).WhereT(func(t *business.Task) bool {
		return t.ClientId == request.ClientId
	}).First().(*business.Task)
	if task == nil {
		ctx.JSON(200, models.WriteDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("server has no client named:%s", request.ClientId),
		})
		return
	}

	ch, ok := task.ReadBuffer[request.SessionId]
	if !ok {
		ctx.JSON(200, models.WriteDataResponse{
			ClientId:  request.ClientId,
			SessionId: request.SessionId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("client in server has no session id:%s", request.SessionId),
		})
		return
	}
	//写入值,让task自动读取
	ch <- request.Data
	ctx.JSON(200, models.WriteDataResponse{
		ClientId:  request.ClientId,
		SessionId: request.SessionId,
		IsSuccess: true,
	})
}

func UnRegister(ctx *gin.Context) {
	var request models.UnRegisterRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(200, models.UnRegisterResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    err.Error(),
		})
		return
	}

	task := From(business.Tasks).WhereT(func(t *business.Task) bool {
		return t.ClientId == request.ClientId
	}).First().(*business.Task)
	if task == nil {
		ctx.JSON(200, models.WriteDataResponse{
			ClientId:  request.ClientId,
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("server has no client named:%s", request.ClientId),
		})
		return
	}
	if err := task.Stop(); err != nil {
		ctx.JSON(200, models.WriteDataResponse{
			ClientId:  request.ClientId,
			IsSuccess: true,
		})
	}

}
