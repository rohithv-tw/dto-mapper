package gjson

import (
	"github.com/tidwall/gjson"
	"reflect"
)

func GetValueBasedOnType(result gjson.Result, fieldKind reflect.Kind) interface{} {
	switch fieldKind {
	case reflect.Bool:
		return result.Bool()
	case reflect.Int:
		return int(result.Int())
	case reflect.Int8:
		return int8(result.Int())
	case reflect.Int16:
		return int16(result.Int())
	case reflect.Int32:
		return int32(result.Int())
	case reflect.Int64:
		return result.Int()
	case reflect.Uint:
		return uint(result.Uint())
	case reflect.Uint8:
		return uint8(result.Uint())
	case reflect.Uint16:
		return uint16(result.Uint())
	case reflect.Uint32:
		return uint32(result.Uint())
	case reflect.Uint64:
		return result.Uint()
	case reflect.Float32:
		return float32(result.Float())
	case reflect.Float64:
		return result.Float()
	default:
		return result.Value()
	}
}
