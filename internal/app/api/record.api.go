package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var RecordSet = wire.NewSet(wire.Struct(new(RecordAPI), "*"))

type RecordAPI struct {
	RecordSrv *service.RecordSrv
}

func (a *RecordAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RecordQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.RecordSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *RecordAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RecordSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
