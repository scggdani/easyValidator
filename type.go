package easyValidator

import (
	"context"
	"fmt"
	"net/http"
)

type Validator interface {
	ValidatorType() string
}

// 该接口用于从http.Request中解析出参数用于初始化用户传入的结构体并对字段进行校验
type HttpRequestValidator interface {
	BindHttpRequest(req *http.Request, val interface{}) error
}

// 该接口用于检验结构体字段的规范性
type StructValidator interface {
	ValidateStruct(val interface{}) error
}

// 该接口用于从Context上下文中解析参数并校验
type ContextValidator interface {
	BindContext(ctx *context.Context, val interface{}) error
}

var (
	ErrorTag_NoTag = fmt.Errorf("校验标签[%v]不存在", ValidateTag)
	ErrorTag_FormatErr = fmt.Errorf("校验标签格式错误")
)

// 校验和参数解析需要用到的标签
const (
	ValidateTag = "check"
	FormTag = "form"
	DefaultTag = "default" // 用于定义默认值
)

// 针对数字类型
var (
	PARAM_IT  = "it" // 小于, 针对数字类型
	PARAM_GT = "gt" // 大于, 针对数字类型
	PARAM_ITE = "ite" // 小于等于, 针对数字类型
	PARAM_GTE = "gte" // 大于等于, 针对数字类型
	PARAM_MAX = "max" // 最大长度, 针对字符串、切片
	PARAM_MIN = "min" // 最小长度, 针对字符串、切片
	PARAM_LEN = "len" // 指定长度, 针对字符串、切片
	PARAM_SELECT = "select" // 应该为其中之一, 针对字符串和数值
)