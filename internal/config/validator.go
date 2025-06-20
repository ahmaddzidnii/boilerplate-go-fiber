package config

import "github.com/go-playground/validator/v10"

//var Validator *validator.Validate
//
//func InitializeValidator() {
//	validate := validator.New()
//	Validator = validate
//}

func InitValidator() *validator.Validate {
	validate := validator.New()
	return validate
}
