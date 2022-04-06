// Package easyValidator 该包是参数校验器的实现
// 目前使用反射来对用户传入的参数做校验
// TODO: 使用AST代替反射实现
package easyValidator


import (
	"context"
	"errors"
	"net/http"
	"reflect"
)

var (
	ValidatorStruct byte = 0
	ValidatorHttpRequest byte = 1
	ValidatorContext byte = 2
)

func validatorName(index byte) string {
	switch index {
	case ValidatorStruct:
		return "Struct Validator"
	case ValidatorHttpRequest:
		return "HTTP Request Validator"
	case ValidatorContext:
		return "Context Validator"
	default:
		return ""
	}
}

type defaultHttpRequestValidator struct {

}

func NewHttpRequestValidator() *defaultHttpRequestValidator {
	return &defaultHttpRequestValidator{}
}

func (Validate *defaultHttpRequestValidator) ValidatorType() string {
	return validatorName(ValidatorHttpRequest)
}

func (validate *defaultHttpRequestValidator) BindHttpRequest(req *http.Request, val interface{}) error {
	if req == nil {
		return errors.New("http请求为空")
	}
	if val == nil {
		return errors.New("传入结构体为空")
	}
	req.ParseForm()
	Type, Value := reflect.TypeOf(val), reflect.ValueOf(val)
	return bindHTTP(req.Form, Type.Elem(), Value.Elem())
}


type defaultStructValidator struct {

}

func NewStructValidator() *defaultStructValidator {
	return &defaultStructValidator{}
}

func (Validate *defaultStructValidator) ValidatorType() string {
	return validatorName(ValidatorStruct)
}

func (validate *defaultStructValidator) ValidateStruct(val interface{}) error {
	if val == nil {
		return errors.New("传入结构体为空")
	}
	return checkStruct(reflect.TypeOf(val), reflect.ValueOf(val))
}

type defaultContextValidator struct {

}

func NewContextValidator() *defaultContextValidator {
	return &defaultContextValidator{}
}

func (Validate *defaultContextValidator) ValidatorType() string {
	return validatorName(ValidatorStruct)
}

func (validate *defaultContextValidator) BindContext(ctx context.Context, val interface{}) error {
	if ctx == nil {
		return errors.New("上下文Context为空")
	}
	if val == nil {
		return errors.New("传入结构体为空")
	}
	Type, Value := reflect.TypeOf(val), reflect.ValueOf(val)
	return bindContext(ctx, Type.Elem(), Value.Elem())
}






