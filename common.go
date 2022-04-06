package easyValidator

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// 递归的解析并校验http request中的参数
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
			if defaultVal := valType.Tag.Get(DefaultTag); len(defaultVal) != 0 {
				param = append(param, defaultVal)
			} else {
				continue
			}
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

func bindContext(ctx context.Context, Type reflect.Type, Value reflect.Value) error{
	var err error
	var ctxVal interface{}
	for i := 0; i < Value.NumField(); i++ {
		valType, fieldVal := Type.Field(i), Value.Field(i)
		if valType.Type.Kind() == reflect.Struct {
			if err = bindContext(ctx, valType.Type, fieldVal); err != nil {
				return err
			}
		}
		if valType.Type.Kind() == reflect.Ptr && valType.Type.Elem().Kind() == reflect.Struct {
			if err = bindContext(ctx, valType.Type.Elem(), fieldVal.Elem()); err != nil {
				return err
			}
		}
		if ctxVal = ctx.Value(Type.Field(i).Tag.Get(FormTag)); ctxVal != nil {
				ctxType, structType := reflect.TypeOf(ctxVal), valType.Type
				if ctxType == structType {
					fieldVal.Set(reflect.ValueOf(ctxVal))
				} else if isUintKind(ctxType.Kind(), structType.Kind()) {
					fieldVal.SetUint(reflect.ValueOf(ctxVal).Uint())
				} else if isIntKind(ctxType.Kind(), structType.Kind()) {
					fieldVal.SetInt(reflect.ValueOf(ctxVal).Int())
				} else {
					return fmt.Errorf("context解析出现类型错误,结构体中类型为[%v],context中对应值的类型为[%v]",
						structType.Kind().String(), ctxType.Kind().String())
				}
		} else {
			defaultVal := Type.Field(i).Tag.Get(DefaultTag)
			if len(defaultVal) == 0 {
				continue
			}
			if isUintKind(valType.Type.Kind()) {
				origin, _ := strconv.ParseUint(defaultVal, 10, 64)
				fieldVal.SetUint(origin)
			} else if isIntKind(valType.Type.Kind()) {
				origin, _ := strconv.ParseInt(defaultVal, 10, 64)
				fieldVal.SetInt(origin)
			} else if valType.Type.Kind() == reflect.String {
				fieldVal.SetString(defaultVal)
			} else {
				return nil
			}
		}
		if err = check(valType, fieldVal); err != nil {
			return err
		}
	}
	return nil
}

func isUintKind(valType...reflect.Kind) bool {
	uintNumberKind := map[reflect.Kind]bool {
		reflect.Uint: true,
		reflect.Uint16: true,
		reflect.Uint32: true,
		reflect.Uint64: true,
		reflect.Uint8: true,
	}
	for _, v := range valType {
		if _, ok := uintNumberKind[v]; !ok {
			return false
		}
	}
	return true
}

func isIntKind(valType...reflect.Kind) bool {
	intNumberKind := map[reflect.Kind]bool {
		reflect.Int64: true,
		reflect.Int32: true,
		reflect.Int16: true,
		reflect.Int8: true,
		reflect.Int: true,
	}
	for _, v := range valType {
		if _, ok := intNumberKind[v]; !ok {
			return false
		}
	}
	return true
}