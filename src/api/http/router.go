package http

import (
	"fmt"
	errors2 "go-backend-template/src/internal/base/errors"
	"go-backend-template/src/internal/base/request"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func initRouter(server *Server) {
	router := &router{
		Server: server,
	}

	router.init()
}

type router struct {
	*Server
}

func (r *router) init() {
	r.Engine.Use(r.trace())
	r.Engine.Use(r.recover())
	r.Engine.Use(r.logger())

	// r.Engine.POST("/login", r.login)
	r.Engine.NoRoute(r.methodNotFound)
}

//func (r *router) login(c *gin.Context) {
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

func (r *router) methodNotFound(c *gin.Context) {
	err := errors2.New(errors2.NotFoundError, "method not found")
	errorResponse(err, nil, r.config.DetailedError()).reply(c)
}

func (r *router) recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		response := internalErrorResponse(nil)
		c.AbortWithStatusJSON(response.Status, response)
	})
}

func (r *router) trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("Trace-Id")
		if traceId == "" {
			traceId, _ = r.crypto.GenerateUUID()
		}

		setTraceId(c, traceId)
	}
}

func (r *router) logger() gin.HandlerFunc {
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

func bindBody(payload interface{}, c *gin.Context) error {
	err := c.BindJSON(payload)

	if err != nil {
		return errors2.New(errors2.BadRequestError, err.Error())
	}

	return nil
}

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OkResponse(data interface{}) *response {
	return &response{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}
}

func internalErrorResponse(data interface{}) *response {
	status, message := http.StatusInternalServerError, "internal error"

	return &response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func errorResponse(err error, data interface{}, withDetails bool) *response {
	status, message, details := parseError(err)

	if withDetails && details != "" {
		message = details
	}
	return &response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func (r *response) reply(c *gin.Context) {
	c.JSON(r.Status, r)
}
