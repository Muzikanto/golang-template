package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/request"
)

type reqInfoKeyType = string

const (
	reqInfoKey reqInfoKeyType = "request-info"
)

func SetTraceId(c *gin.Context, traceId string) {
	info, exists := c.Get(reqInfoKey)
	if exists {
		parsedInfo := info.(request.RequestInfo)
		parsedInfo.TraceId = traceId

		c.Set(reqInfoKey, parsedInfo)

		return
	}

	c.Set(reqInfoKey, request.RequestInfo{TraceId: traceId})
}

func SetUserId(c *gin.Context, userId int64) {
	info, exists := c.Get(reqInfoKey)
	if exists {
		parsedInfo := info.(request.RequestInfo)
		parsedInfo.UserId = userId

		c.Set(reqInfoKey, parsedInfo)

		return
	}

	c.Set(reqInfoKey, request.RequestInfo{UserId: userId})
}

func GetReqInfo(c *gin.Context) request.RequestInfo {
	info, ok := c.Get(reqInfoKey)
	if ok {
		return info.(request.RequestInfo)
	}

	return request.RequestInfo{}
}

func ContextWithReqInfo(c *gin.Context) context.Context {
	info, ok := c.Get(reqInfoKey)
	if ok {
		return request.WithRequestInfo(c, info.(request.RequestInfo))
	}

	return request.WithRequestInfo(c, request.RequestInfo{})
}
