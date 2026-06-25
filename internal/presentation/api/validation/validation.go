package validation

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetMessage(err error) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}

	if len(validationErrors) == 0 {
		return err.Error()
	}

	msgs := make([]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		var msg string
		switch fieldErr.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", fieldErr.Field())
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", fieldErr.Field(), fieldErr.Param())
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", fieldErr.Field(), fieldErr.Param())
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
		default:
			slog.Warn("Unknown validation error", "field", fieldErr.Field(), "tag", fieldErr.Tag(), "param", fieldErr.Param())
			msg = fieldErr.Error()
		}
		msgs = append(msgs, msg)
	}

	return strings.Join(msgs, "; ")
}
