package easyValidator

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func bindHTTP(form url.Values, Type reflect.Type, Value reflect.Value) error {
	var err error
	for i := 0; i < Value.NumField(); i++ {
		valType, fieldVal := Type.Field(i), Value.Field(i)
		if   valType.Type.Kind() == reflect.Struct {
			if err = bindHTTP(form, valType.Type, fieldVal); err != nil {
				return err
			}
		}
		if valType.Type.Kind() == reflect.Ptr && valType.Type.Elem().Kind() == reflect.Struct {
			if err = bindHTTP(form, valType.Type.Elem(), fieldVal.Elem()); err != nil {
				return err
			}
		}

		param := form[valType.Tag.Get(FormTag)]
		if len(param) == 0 {
			continue
		}
		switch fieldVal.Kind().String() {
		case "int8","int16","int32","int64","rune","int":
			origin, _ := strconv.ParseInt(param[0], 10, 64)
			fieldVal.SetInt(origin)
		case "uint8","byte","uint16","uint32","uint64","uint":
			origin, _ := strconv.ParseUint(param[0], 10, 64)
			fieldVal.SetUint(origin)
		case "string":
			fieldVal.SetString(param[0])
		case "bool":
			origin := strings.ToUpper(param[0])
			if origin == "TRUE" {
				fieldVal.SetBool(true)
			} else {
				fieldVal.SetBool(false)
			}
		case "float32","float64":
			origin, _ := strconv.ParseFloat(param[0], 64)
			fieldVal.SetFloat(origin)
		default:
			return errors.New("不支持的解析类型:" + fieldVal.Kind().String())
		}
		if err = check(valType, fieldVal); err != nil {
			return err
		}

	}
	return nil
}