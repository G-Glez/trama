package errorhandler

import (
	"errors"
	"log/slog"
	"net/http"
	"trama/internal/presentation/api/auth"
	"trama/internal/presentation/api/validation"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Message: validation.GetMessage(err)}
}

func GlobalErrorHandler(err error) (int, ErrorResponse) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return http.StatusBadRequest, NewErrorResponse(err)
	}

	switch {
	case errors.Is(err, auth.ErrUserExists):
		return http.StatusConflict, NewErrorResponse(err)
	case errors.Is(err, auth.ErrUserNotFound),
		errors.Is(err, auth.ErrInvalidPassword):
		return http.StatusUnauthorized, NewErrorResponse(err)
	default:
		slog.Error("Internal server error", "error", err)
		return http.StatusInternalServerError, NewErrorResponse(err)
	}
}
