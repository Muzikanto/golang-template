package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	errors2 "go-backend-template/src/internal/base/errors"
	"go-backend-template/src/internal/base/request"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OkResponse(data interface{}) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}
}

func (r *Response) Reply(c *gin.Context) {
	c.JSON(r.Status, r)
}

func BindBody(payload interface{}, c *gin.Context) error {
	err := c.BindJSON(payload)

	if err != nil {
		return errors2.New(errors2.BadRequestError, err.Error())
	}

	return nil
}

func ErrorResponse(err error, data interface{}, withDetails bool) *Response {
	status, message, details := parseError(err)

	if withDetails && details != "" {
		message = details
	}
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func parseError(err error) (status int, message, details string) {
	var baseErr *errors2.Error

	if castErr, ok := err.(*errors2.Error); ok {
		baseErr = castErr
	}
	if baseErr == nil {
		baseErr = errors2.Wrap(err, errors2.InternalError, "")
	}

	status = convertErrorStatusToHTTP(baseErr.Status())
	message = baseErr.Error()
	details = baseErr.DetailedError()

	return
}

func convertErrorStatusToHTTP(status errors2.Status) int {
	switch status {
	case errors2.BadRequestError:
		return http.StatusBadRequest
	case errors2.ValidationError:
		return http.StatusBadRequest
	case errors2.UnauthorizedError:
		return http.StatusUnauthorized
	case errors2.WrongCredentialsError:
		return http.StatusUnauthorized
	case errors2.NotFoundError:
		return http.StatusNotFound
	case errors2.AlreadyExistsError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

//

type reqInfoKeyType = string

const (
	reqInfoKey reqInfoKeyType = "request-info"
)

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
