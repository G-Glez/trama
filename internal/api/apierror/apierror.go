package apierror

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"trama/internal/battlelog"
	"trama/internal/core"
	"trama/internal/user"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *Error) Error() string {
	return e.Message
}

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func BadRequest(msg string) *Error {
	return New(http.StatusBadRequest, msg)
}

func Internal(msg string) *Error {
	return New(http.StatusInternalServerError, msg)
}

func HandleError(ctx *gin.Context, err error) {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		ctx.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
		return
	}

	switch {
	case errors.Is(err, core.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrUserNotFound),
		errors.Is(err, battlelog.ErrMatchNotFound),
		errors.Is(err, battlelog.ErrTeamNotFound),
		errors.Is(err, battlelog.ErrTournamentNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrEmailTaken):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrInvalidPassword):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, battlelog.ErrInvalidPoints),
		errors.Is(err, battlelog.ErrInvalidWinner):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, core.ErrDB):
		slog.Error("database error", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	case errors.Is(err, core.ErrDataCorruption):
		slog.Error("data corruption", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
