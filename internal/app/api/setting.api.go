package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/ginx"
	"email-send-manager/internal/app/schema"
	"email-send-manager/internal/app/service"
)

var SettingSet = wire.NewSet(wire.Struct(new(SettingAPI), "*"))

type SettingAPI struct {
	SettingSrv *service.SettingSrv
}

func (a *SettingAPI) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.SettingSrv.Get(ctx, c.Param("key"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *SettingAPI) Upsert(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Setting
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.SettingSrv.Upsert(ctx, item)

	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
