// operation.go中封装了一系列用于校验字段的函数
package easyValidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 用于初始化校验分配器并绑定对应的处理函数
func init() {
	defaultOperateFuncMap = make(map[string]func(reflect.Value, ...string) bool, 3)
	defaultOperateFuncMap[PARAM_IT] = IT
	defaultOperateFuncMap[PARAM_GT] = GT
	defaultOperateFuncMap[PARAM_MAX] = MAX
	defaultOperateFuncMap[PARAM_LEN] = LEN
	defaultOperateFuncMap[PARAM_SELECT] = SELECT
	defaultOperateFuncMap[PARAM_GTE] = GTE
	defaultOperateFuncMap[PARAM_ITE] = ITE
}

// 校验规则
// 若字段: age int `form:"age" check:"gt=18,ite=22"`
// 则其对应的ChenckParam为: [[gt, 18], [ite, 22]]
type CheckParam [][]string

// 该函数用于提取出校验标签check的校验规则
func splitCheckTag(str string) (CheckParam, error) {
	if len(str) == 0 {
		return nil, ErrorTag_NoTag
	}
	var param CheckParam = make(CheckParam, 0, 3)
	fields := strings.Split(str, ",")
	for _, v := range fields {
		field := strings.Split(v, "=")
		if len(field) != 2 {
			return nil, ErrorTag_FormatErr
		}
		param = append(param, field)
	}
	return param, nil
}

// 该函数用于对结构体中的某一字段进行校验
func check(valType reflect.StructField, val reflect.Value) error {
	str := valType.Tag.Get(ValidateTag)
	checkParam, err := splitCheckTag(str)
	// 没有check标签代表该字段无需校验
	// 故直接返回nil
	if err == ErrorTag_NoTag {
		return nil
	}
	// 出现标签格式错误
	if err ==  ErrorTag_FormatErr {
		return err
	}
	for _, tag := range checkParam {
		if ok := defaultOperateFuncMap[tag[0]](val, tag[1:]...); !ok{
			err = fmt.Errorf("字段[%v]错误, 值为[%v], 校验规则为[%v]", valType.Name, val, errorTip(str))
			return err
		}
	}
	return nil
}

// 该函数用于人性化的显示出错字段的校验规则
func errorTip(tip string) string {
	if tip == "" {
		return ""
	}
	tip = strings.Replace(tip, ",", ";", -1)
	keys := map[string]string {
		PARAM_GT+"=": "大于",
		PARAM_GTE+"=": "大于等于",
		PARAM_IT+"=": "小于",
		PARAM_ITE+"=": "小于等于",
		PARAM_MAX+"=": "最大长度为",
		PARAM_MIN+"=": "最小长度为",
		PARAM_LEN+"=": "长度等于",
		PARAM_SELECT+"=": "值为下列之一:",
	}
	for k, v := range keys {
		tip = strings.Replace(tip, k, v, -1)
	}
	tip = strings.Replace(tip, " ", ",", -1)
	re, _ := regexp.Compile(`:([\d,]*\d)`)
	tip = re.ReplaceAllString(tip, ":{$1}")
	return tip
}

func checkStruct(Type reflect.Type, Value reflect.Value) error {
	var err error
	for i := 0; i < Value.NumField(); i++ {
		valType, fieldVal := Type.Field(i), Value.Field(i)
		if   valType.Type.Kind() == reflect.Struct {
			if err = checkStruct(valType.Type, fieldVal); err != nil {
				return err
			}
		}
		if valType.Type.Kind() == reflect.Ptr && valType.Type.Elem().Kind() == reflect.Struct {
			if err = checkStruct(valType.Type.Elem(), fieldVal.Elem()); err != nil {
				return err
			}
		}
		if err = check(valType, fieldVal); err != nil {
			return err
		}
	}
	return nil
}

// 校验器分配器
// 用于将校验字段与对应的处理函数绑定
var defaultOperateFuncMap map[string]func(reflect.Value, ...string) bool

func GT(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	switch val.Kind().String() {
	case "float32","float64":
		origin := val.Float()
		if origin <= float64(compare) {
			return false
		}
		return true
	case "uint","uint8","uint16","uint32","uint64":
		origin := val.Uint()
		if origin <= uint64(compare) {
			return false
		}
		return true
	case "int","int8","int16","int32","int64":
		origin := val.Int()
		if origin <= int64(compare) {
			return false
		}
		return true
	default:
		return false
	}
}

func IT(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	switch val.Kind().String() {
	case "float32","float64":
		origin := val.Float()
		if origin >= float64(compare) {
			return false
		}
		return true
	case "uint","uint8","uint16","uint32","uint64":
		origin := val.Uint()
		if origin >= uint64(compare) {
			return false
		}
		return true
	case "int","int8","int16","int32","int64":
		origin := val.Int()
		if origin >= int64(compare) {
			return false
		}
		return true
	default:
		return false
	}
}

func ITE(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	switch val.Kind().String() {
	case "float32","float64":
		origin := val.Float()
		if origin > float64(compare) {
			return false
		}
		return true
	case "uint","uint8","uint16","uint32","uint64":
		origin := val.Uint()
		if origin > uint64(compare) {
			return false
		}
		return true
	case "int","int8","int16","int32","int64":
		origin := val.Int()
		if origin > int64(compare) {
			return false
		}
		return true
	default:
		return false
	}
}

func MAX(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	origin := val.Len()
	if origin > compare {
		return false
	}
	return true
}

func LEN(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	origin := val.Len()
	if origin != compare {
		return false
	}
	return true
}

func SELECT(val reflect.Value, std ...string) bool {
	options := strings.Split(strings.Trim(std[0], " "), " ")
	var str string
	switch val.Kind().String() {
	case "int","int8","int16","int32","int64":
		str = strconv.FormatInt(val.Int(), 10)
	case "uint","uint8","uint16","uint32","uint64":
		str = strconv.FormatUint(val.Uint(), 10)
	case "float32","float64":
		str = strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case "string":
		str = val.String()
	default:
		return false
	}

	for _, v := range options {
		if str == v {
			return true
		}
	}
	return false
}

func GTE(val reflect.Value, std ...string) bool {
	compare, err := strconv.Atoi(std[0])
	if err != nil {
		return false
	}
	switch val.Kind().String() {
	case "float32","float64":
		origin := val.Float()
		if origin < float64(compare) {
			return false
		}
		return true
	case "uint","uint8","uint16","uint32","uint64":
		origin := val.Uint()
		if origin < uint64(compare) {
			return false
		}
		return true
	case "int","int8","int16","int32","int64":
		origin := val.Int()
		if origin < int64(compare) {
			return false
		}
		return true
	default:
		return false
	}
}


