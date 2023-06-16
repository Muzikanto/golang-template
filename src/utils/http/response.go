package http

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/utils/errors"
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
		return errors.New(errors.BadRequestError, err.Error())
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

func InternalErrorResponse(data interface{}) *Response {
	status, message := http.StatusInternalServerError, "internal error"

	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
