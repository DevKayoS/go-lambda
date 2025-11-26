package errors

import "net/http"

func BadRequest(msg string) ApiError {
	return ApiError{
		Code:       "BAD_REQUEST",
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

func Unathorized(msg string) ApiError {
	return ApiError{
		Code:       "UNAUTHORIZED",
		Message:    msg,
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(msg string) ApiError {
	return ApiError{
		Code:       "FORBIDDEN",
		Message:    msg,
		StatusCode: http.StatusForbidden,
	}
}

func NotFound(msg string) ApiError {
	return ApiError{
		Code:       "NOT_FOUND",
		Message:    msg,
		StatusCode: http.StatusNotFound,
	}
}

func Internal(msg string, err error) ApiError {
	return ApiError{
		Code:       "INTERNAL_ERROR",
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
