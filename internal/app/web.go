package app

import (
	"email-send-manager/internal/app/config"
	"email-send-manager/internal/app/middleware"
	"email-send-manager/internal/app/router"
	"email-send-manager/pkg/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

// InitGinEngine 初始化gin引擎
func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	//prefixes := r.Prefixes()

	// Recover
	//app.Use(middleware.RecoveryMiddleware())

	// Trace ID
	//app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Access logger
	//app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Router register
	r.Register(app)

	fe, _ := fs.Sub(static.Static, "assets")

	app.StaticFS("assets", http.FS(fe))

	app.GET("", func(c *gin.Context) {
		c.FileFromFS("/", http.FS(fe))
	})

	app.NoRoute(middleware.IndexHandler(http.FS(fe)))
	app.NoMethod(middleware.NoMethodHandler())

	// Website
	//if dir := config.C.WWW; dir != "" {
	//	app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	//}

	return app
}
