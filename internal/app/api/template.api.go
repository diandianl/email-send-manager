package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var TemplateSet = wire.NewSet(wire.Struct(new(TemplateAPI), "*"))

type TemplateAPI struct {
	TemplateSrv *service.TemplateSrv
}

func (a *TemplateAPI) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.TemplateQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.TemplateSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *TemplateAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.TemplateSrv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *TemplateAPI) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Template
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.TemplateSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *TemplateAPI) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Template
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.TemplateSrv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *TemplateAPI) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TemplateSrv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *TemplateAPI) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TemplateSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *TemplateAPI) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TemplateSrv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 0)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
