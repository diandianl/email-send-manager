package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var SendBatchSet = wire.NewSet(wire.Struct(new(SendBatchAPI), "*"))

type SendBatchAPI struct {
	SendBatchSrv *service.SendBatchSrv
}

func (a *SendBatchAPI) Current(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.SendBatchSrv.Current(ctx)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *SendBatchAPI) StartSendBatch(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.SendBatch
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.SendBatchSrv.StartSendBatch(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *SendBatchAPI) Cancel(c *gin.Context) {
	err := a.SendBatchSrv.TerminateCurrent(c.Request.Context())
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
