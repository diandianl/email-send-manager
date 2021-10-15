package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"email-send-manager/internal/app/api"
)

var _ IRouter = (*Router)(nil)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	CustomerAPI  *api.CustomerAPI
	TemplateAPI  *api.TemplateAPI
	RecordAPI    *api.RecordAPI
	SendBatchAPI *api.SendBatchAPI
	SettingAPI   *api.SettingAPI
} // end

// Register 注册路由
func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// Prefixes 路由前缀列表
func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	v1 := g.Group("/v1")
	{
		gCustomer := v1.Group("customers")
		{
			gCustomer.GET("", a.CustomerAPI.Query)
			gCustomer.GET(":id", a.CustomerAPI.Get)
			gCustomer.POST("", a.CustomerAPI.Create)
			gCustomer.POST("import", a.CustomerAPI.Import)
			gCustomer.PUT(":id", a.CustomerAPI.Update)
			gCustomer.DELETE(":id", a.CustomerAPI.Delete)
			gCustomer.PATCH(":id/enable", a.CustomerAPI.Enable)
			gCustomer.PATCH(":id/disable", a.CustomerAPI.Disable)
		}

		gTemplate := v1.Group("templates")
		{
			gTemplate.GET("", a.TemplateAPI.Query)
			gTemplate.GET(":id", a.TemplateAPI.Get)
			gTemplate.POST("", a.TemplateAPI.Create)
			gTemplate.PUT(":id", a.TemplateAPI.Update)
			gTemplate.DELETE(":id", a.TemplateAPI.Delete)
			gTemplate.PATCH(":id/enable", a.TemplateAPI.Enable)
			gTemplate.PATCH(":id/disable", a.TemplateAPI.Disable)
		}

		gRecord := v1.Group("records")
		{
			gRecord.GET("", a.RecordAPI.Query)
			gRecord.DELETE(":id", a.RecordAPI.Delete)
		}

		gSendBatch := v1.Group("send-batches")
		{
			gSendBatch.GET("current", a.SendBatchAPI.Current)
			gSendBatch.POST("", a.SendBatchAPI.Create)
			gSendBatch.DELETE("", a.SendBatchAPI.Cancel)
		}

		gSetting := v1.Group("settings")
		{
			gSetting.GET(":key", a.SettingAPI.Get)
			gSetting.POST("", a.SettingAPI.Upsert)
		}

	} // v1 end
}
