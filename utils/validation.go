package utils

import (
	"github.com/casmeyu/micro-user/structs"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error, errors *[]*structs.IError) {
	for _, err := range err.(validator.ValidationErrors) {
		var el structs.IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		*errors = append(*errors, &el)
	}
}
