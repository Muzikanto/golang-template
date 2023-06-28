package core

import (
	"github.com/gin-gonic/gin"
	error2 "go-backend-template/src/utils/error"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//

func OkResponse(data interface{}) *Response {
	return &Response{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}
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

//

func (r *Response) Reply(c *gin.Context) {
	c.JSON(r.Status, r)
}

func BindBody(payload interface{}, c *gin.Context) error {
	err := c.BindJSON(payload)

	if err != nil {
		return error2.New(error2.BadRequestError, err.Error())
	}

	return nil
}

func BindQuery(payload interface{}, c *gin.Context) error {
	err := c.BindQuery(payload)

	if err != nil {
		return error2.New(error2.BadRequestError, err.Error())
	}

	return nil
}

func Bind(payload interface{}, c *gin.Context) error {
	err := c.Bind(payload)

	if err != nil {
		return error2.New(error2.BadRequestError, err.Error())
	}

	return nil
}

// utils

func convertErrorStatusToHTTP(status error2.Status) int {
	switch status {
	case error2.BadRequestError:
		return http.StatusBadRequest
	case error2.ValidationError:
		return http.StatusBadRequest
	case error2.UnauthorizedError:
		return http.StatusUnauthorized
	case error2.WrongCredentialsError:
		return http.StatusUnauthorized
	case error2.NotFoundError:
		return http.StatusNotFound
	case error2.AlreadyExistsError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func parseError(err error) (status int, message, details string) {
	var baseErr *error2.Error

	if castErr, ok := err.(*error2.Error); ok {
		baseErr = castErr
	}
	if baseErr == nil {
		baseErr = error2.Wrap(err, error2.InternalError, "")
	}

	status = convertErrorStatusToHTTP(baseErr.Status())
	message = baseErr.Error()
	details = baseErr.DetailedError()

	return
}
