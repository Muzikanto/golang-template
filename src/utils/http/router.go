package http

import (
	"fmt"
	"go-backend-template/src/utils/errors"
	"go-backend-template/src/utils/request"
	"time"

	"github.com/gin-gonic/gin"
)

func InitRouter(server *Server, path string) *Router {
	routerGroup := server.Engine.Group(path)

	router := &Router{
		Server:      server,
		RouterGroup: routerGroup,
	}

	router.init()

	return router
}

type Router struct {
	*Server
	RouterGroup *gin.RouterGroup
}

func (r *Router) init() {
	r.RouterGroup.Use(r.trace())
	r.RouterGroup.Use(r.recover())
	r.RouterGroup.Use(r.logger())

	r.Engine.NoRoute(r.methodNotFound)
}

//func (r *Router) login(c *gin.Context) {
//	var loginUserDto auth.LoginUserDto
//
//	if err := bindBody(&loginUserDto, c); err != nil {
//		errorResponse(err, nil, r.config.DetailedError()).reply(c)
//		return
//	}
//
//	user, err := r.authService.Login(c, loginUserDto)
//	if err != nil {
//		errorResponse(err, nil, r.config.DetailedError()).reply(c)
//		return
//	}
//
//	OkResponse(user).reply(c)
//}

func (r *Router) methodNotFound(c *gin.Context) {
	err := errors.New(errors.NotFoundError, "method not found")
	ErrorResponse(err, nil, true).Reply(c)
}

func (r *Router) recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		response := InternalErrorResponse(nil)
		c.AbortWithStatusJSON(response.Status, response)
	})
}

func (r *Router) trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("Trace-Id")
		if traceId == "" {
			traceId, _ = r.Crypto.GenerateUUID()
		}

		SetTraceId(c, traceId)
	}
}

func (r *Router) logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var parsedReqInfo request.RequestInfo

		reqInfo, exists := param.Keys[reqInfoKey]
		if exists {
			parsedReqInfo = reqInfo.(request.RequestInfo)
		}

		return fmt.Sprintf("%s - [HTTP] TraceId: %s; UserId: %d; Method: %s; Path: %s; Status: %d, Latency: %s;\n\n",
			param.TimeStamp.Format(time.RFC1123),
			parsedReqInfo.TraceId,
			parsedReqInfo.UserId,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}
