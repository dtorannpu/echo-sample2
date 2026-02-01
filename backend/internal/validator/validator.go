package validator

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			e := make(map[string]string)
			for _, err := range validationErrors {
				e[err.Field()] = err.Tag()
			}
			return echo.NewHTTPError(http.StatusBadRequest, e)
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
