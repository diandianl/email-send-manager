package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/contextx"
	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var SendBatchSet = wire.NewSet(wire.Struct(new(SendBatchAPI), "*"))

type SendBatchAPI struct {
	SendBatchSrv *service.SendBatchSrv
}

func (a *SendBatchAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.SendBatchQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.SendBatchSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *SendBatchAPI) Current(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.SendBatchSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *SendBatchAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.SendBatch
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = contextx.FromUserID(ctx)
	result, err := a.SendBatchSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *SendBatchAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.SendBatchSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
