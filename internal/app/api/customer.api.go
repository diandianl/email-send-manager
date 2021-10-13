package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var CustomerSet = wire.NewSet(wire.Struct(new(CustomerAPI), "*"))

type CustomerAPI struct {
	CustomerSrv *service.CustomerSrv
}

func (a *CustomerAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.CustomerQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.CustomerSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *CustomerAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.CustomerSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *CustomerAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Customer
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.CustomerSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *CustomerAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Customer
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.CustomerSrv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *CustomerAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.CustomerSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *CustomerAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.CustomerSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *CustomerAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.CustomerSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 0)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *CustomerAPI) Import(c *gin.Context) {
	ctx := c.Request.Context()
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	err = a.CustomerSrv.Import(ctx, file)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
