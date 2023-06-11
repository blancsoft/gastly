//go:build js && wasm

package ast

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"
)

func ConvertToJSValue(value interface{}) js.Value {
	var rv reflect.Value
	if v, ok := value.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(value)
	}

	switch rv.Kind() {
	case reflect.Invalid:
		return js.Null()
	case reflect.Bool:
		return js.ValueOf(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return js.ValueOf(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return js.ValueOf(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return js.ValueOf(rv.Float())
	case reflect.String:
		return js.ValueOf(rv.String())
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		jsArray := js.Global().Get("Array").New()
		for i := 0; i < length; i++ {
			v := ConvertToJSValue(rv.Index(i).Interface())
			jsArray.SetIndex(i, v)
		}
		return jsArray
	case reflect.Map:
		jsObject := js.Global().Get("Object").New()
		for _, key := range rv.MapKeys() {
			jsObject.Set(key.String(), ConvertToJSValue(rv.MapIndex(key).Interface()))
		}
		return jsObject
	case reflect.Struct:
		jsObject := js.Global().Get("Object").New()
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Type().Field(i)
			if field.PkgPath != "" { // Skip unexported fields
				continue
			}

			// Use JSON tag if specified, otherwise use field name
			fieldName := field.Name
			if jsonTag, ok := field.Tag.Lookup("json"); ok && jsonTag != "" {
				fieldName = strings.SplitN(jsonTag, ",", 2)[0]
			}
			jsObject.Set(fieldName, ConvertToJSValue(rv.Field(i).Interface()))
		}
		return jsObject

	case reflect.Complex64, reflect.Complex128:
		c := rv.Complex()
		jsObject := js.Global().Get("Object").New()
		jsObject.Set("real", real(c))
		jsObject.Set("imag", imag(c))
		return jsObject
	case reflect.Interface:
		fallthrough
	case reflect.Pointer:
		if rv.IsNil() {
			return js.Null()
		}
		return ConvertToJSValue(rv.Elem().Interface())

	default:
		// TODO: Handle implementations for chan, func. Ignore unsafe pointers
		panic(fmt.Sprintf("unsupported type: %T", value))
	}
}
