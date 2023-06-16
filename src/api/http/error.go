package http

import (
	errors2 "go-backend-template/src/internal/base/errors"
	"net/http"
)

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
