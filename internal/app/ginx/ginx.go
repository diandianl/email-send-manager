package ginx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
	"email-send-manager/pkg/logger"
)

// 定义上下文中的键
const (
	prefix           = "gin-admin"
	ReqBodyKey       = prefix + "/req-body"
	ResBodyKey       = prefix + "/res-body"
)

// ParseParamID Parse path id
func ParseParamID(c *gin.Context, key string) uint64 {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ParseQuery 解析Query参数
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ResOK 响应OK
func ResOK(c *gin.Context) {
	ResSuccess(c, schema.StatusResult{Status: schema.OKStatus})
}

// ResPage 响应分页数据
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResSuccess(c, list)
}

// ResSuccess 响应成功
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON 响应JSON数据
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError 响应错误
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()
	var res *errors.ResponseError

	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(status) > 0 {
		res.Status = status[0]
	}

	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}

		if status := res.Status; status >= 400 && status < 500 {
			logger.WithContext(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.WithContext(logger.NewStackContext(ctx, err)).Errorf(err.Error())
		}
	}

	eitem := schema.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.Status, schema.ErrorResult{Error: eitem})
}
