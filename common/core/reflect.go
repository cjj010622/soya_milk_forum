package core

import (
	"reflect"
)

func ElemValueIfPointer(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	}
	return v
}

// IsDefaultZero 判断对象是否是默认值
func IsDefaultZero(obj any) bool {
	bObj := reflect.ValueOf(obj)
	if bObj.Kind() == reflect.Pointer {
		return bObj.IsNil()
	}

	bObj = ElemValueIfPointer(bObj)

	switch bObj.Kind() {
	case reflect.Int8:
		return obj.(int8) == 0
	case reflect.Int16:
		return obj.(int16) == 0
	case reflect.Int32:
		return obj.(int32) == 0
	case reflect.Int64:
		return obj.(int64) == 0
	case reflect.Int:
		return obj.(int) == 0
	case reflect.String:
		return obj.(string) == ""
	case reflect.Slice:
		return bObj.Len() == 0
	case reflect.Pointer:
		return obj == nil
	}

	return false
}
