package httpErrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ApiErr interface {
	Status() int
	Error() string
}

type ApiError struct {
	ErrStatus int         `json:"code,omitempty"`
	ErrError  string      `json:"message,omitempty"`
	ErrCauses interface{} `json:"-"`
}

var (
	BadRequest          = errors.New("Bad Request")
	ContentType         = errors.New("Content type must be `application/json`")
	CannotMarshal       = errors.New("Could not be marshalled")
	NotFound            = errors.New("Not Found")
	BadQueryParams      = errors.New("Bad Query Params")
	InternalServerError = errors.New("Internal Server Error")
	MissingFields       = errors.New("Missing fields")
	ExistsObjectIDError = errors.New("Object with given id already exists")
	RequestTimeoutError = errors.New("Request Timeout")
	CannotBindGivenData = errors.New("Could not bind given data")
	ValidationError     = errors.New("Validation failed for given payload")
	ForbiddenError      = errors.New("Forbidden")
)

func (a ApiError) Status() int {
	return a.ErrStatus
}

func (a ApiError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", a.ErrStatus, a.ErrError, a.ErrCauses)
}

// NewApiError : creates a new ApiError with given input
func NewApiError(code int, err string, causes interface{}) ApiError {
	return ApiError{
		ErrStatus: code,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// ParseErrors : parses error to a specific structure (ApiError)
func ParseErrors(err error) ApiErr {
	switch {
	case strings.Contains(err.Error(), "json: unsupported"):
		return NewApiError(http.StatusBadRequest, CannotMarshal.Error(), err)
	case strings.Contains(err.Error(), "not found"):
		return NewApiError(http.StatusNotFound, NotFound.Error(), err)
	case strings.Contains(err.Error(), "is required"):
		return NewApiError(http.StatusBadRequest, MissingFields.Error(), err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewApiError(http.StatusBadRequest, BadRequest.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewApiError(http.StatusRequestTimeout, RequestTimeoutError.Error(), err)
	case errors.Is(err, CannotBindGivenData):
		return NewApiError(http.StatusBadRequest, CannotBindGivenData.Error(), err)
	case strings.Contains(err.Error(), "validation"):
		return NewApiError(http.StatusBadRequest, ValidationError.Error(), err)
	case strings.Contains(err.Error(), "not allowed"):
		return NewApiError(http.StatusForbidden, ForbiddenError.Error(), err)
	default:
		if apiErr, ok := err.(ApiErr); ok {
			return apiErr
		}
		return NewInternalServerError(err)
	}
}

// parseSqlErrors : if given error is an sql error, parse it explicitly
func parseSqlErrors(err error) ApiErr {
	if strings.Contains(err.Error(), "23505") {
		return NewApiError(http.StatusBadRequest, ExistsObjectIDError.Error(), err)
	}
	return NewApiError(http.StatusBadRequest, BadRequest.Error(), err)
}

// NewInternalServerError : if given error is an internal error, create it explicitly
func NewInternalServerError(causes interface{}) ApiErr {
	result := ApiError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}
