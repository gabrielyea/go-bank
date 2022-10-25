package handlers

import (
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		res := util.IsSupported(currency)
		return res
	}
	return false
}
