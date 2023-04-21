package lib

import (
	"fmt"
	"reflect"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func ConvertToJSValue(value interface{}) app.Value {
	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Invalid:
		return app.Null()
	case reflect.Bool:
		return app.ValueOf(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return app.ValueOf(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return app.ValueOf(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return app.ValueOf(rv.Float())
	case reflect.String:
		return app.ValueOf(rv.String())
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		jsArray := app.Window().Get("Array").New()
		for i := 0; i < length; i++ {
			jsArray.SetIndex(i, ConvertToJSValue(rv.Index(i).Interface()))
		}
		return jsArray
	case reflect.Map:
		jsObject := app.Window().Get("Object").New()
		for _, key := range rv.MapKeys() {
			jsObject.Set(key.String(), ConvertToJSValue(rv.MapIndex(key).Interface()))
		}
		return jsObject
	case reflect.Struct:
		jsObject := app.Window().Get("Object").New()
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Type().Field(i)
			if field.PkgPath != "" { // Skip unexported fields
				continue
			}
			jsObject.Set(field.Name, ConvertToJSValue(rv.Field(i).Interface()))
		}
		return jsObject
	default:
		panic(fmt.Sprintf("unsupported type: %T", value))
	}
}
