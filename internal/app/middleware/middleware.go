package middleware

import (
	"email-send-manager/internal/app/ginx"
	"email-send-manager/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NoMethodHandler 未找到请求方法的处理函数
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginx.ResError(c, errors.ErrMethodNotAllow)
	}
}

// NoRouteHandler 未找到请求路由的处理函数
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginx.ResError(c, errors.ErrNotFound)
	}
}

func IndexHandler(fs http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.FileFromFS("/", fs)
	}
}

// SkipperFunc 定义中间件跳过函数
type SkipperFunc func(*gin.Context) bool

// AllowPathPrefixSkipper 检查请求路径是否包含指定的前缀，如果包含则跳过
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// AllowPathPrefixNoSkipper 检查请求路径是否包含指定的前缀，如果包含则不跳过
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// SkipHandler 统一处理跳过函数
func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}
