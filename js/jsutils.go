//go:build js && wasm

package js

import (
	"errors"
	"reflect"
	"strings"
	"syscall/js"
)

var (
	global  = js.Global()
	array   = global.Get("Array")
	object  = global.Get("Object")
	promise = global.Get("Promise")
)

func ValueOf(v any) Value {
	switch v.(type) {
	case nil, js.Value:
		return Value(js.ValueOf(v))
	default:
		v := reflect.ValueOf(v)
		return Value(valueOf(v))
	}
}

func valueOfSlice(v reflect.Value) js.Value {
	length := v.Len()
	jsArray := array.New()
	for i := 0; i < length; i++ {
		v := valueOf(v.Index(i))
		jsArray.SetIndex(i, v)
	}
	return jsArray
}

func valueOfMap(v reflect.Value) js.Value {
	jsObject := object.New()
	for _, key := range v.MapKeys() {
		jsObject.Set(key.String(), valueOf(v.MapIndex(key)))
	}
	return jsObject
}

func valueOfStruct(v reflect.Value) js.Value {
	jsObject := object.New()
	for i := 0; i < v.NumField(); i++ {
		// Ignore unexported fields
		if field := v.Field(i); field.CanInterface() {
			sf := v.Type().Field(i)

			// Use JSON tag if specified, otherwise use field name
			name := sf.Name
			if jsonTag := sf.Tag.Get("json"); jsonTag != "" {
				name = strings.SplitN(jsonTag, ",", 2)[0]
			}

			jsObject.Set(name, valueOf(v.Field(i)))
		}

	}
	return jsObject
}

func valueOfComplex(v reflect.Value) js.Value {
	c := v.Complex()
	jsObject := object.New()
	jsObject.Set("real", real(c))
	jsObject.Set("imag", imag(c))
	return jsObject
}

func valueOfPointer(v reflect.Value) js.Value {
	if v.IsNil() {
		return js.Null()
	}
	return valueOf(v.Elem())
}

func valueOfFunc(v reflect.Value) (jsFunc js.Func) {
	// TODO: Implement this function
	//	t := v.Type()
	//	wrapper := reflect.MakeFunc(t, func(args []reflect.Value) []reflect.Value {
	//		// Convert arguments to js.Value
	//		jsArgs := make([]js.Value, len(args))
	//		for i, arg := range args {
	//			jsArgs[i] = valueOf(arg)
	//		}
	//
	//		// Call the original function
	//		result := v.Call(jsArgs)
	//
	//		// Convert the result to []any and then to js.Value
	//		jsResult := make([]any, len(result))
	//		for i, res := range result {
	//			jsResult[i] = res.Interface()
	//		}
	//		return []reflect.Value{reflect.ValueOf(jsResult)}
	//	})
	//
	//	// convert []jsvalues to []reflectvalues
	//	//    ensure len([]jsvalues == []reflectvalues)
	//	//    factor in variadic funcs args
	//	//    ensure each jsvalue types match func arg and ret types
	//	// pass []reflectvalues to fn.call
	//	// handle unexpected errors/panics
	//	// collect return values
	//	// convert ret []reflectvalues to []jsvalues
	//
	//	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	//		// Convert js.Value arguments to []reflect.Value
	//		reflectArgs := make([]reflect.Value, len(args))
	//		for i, arg := range args {
	//			reflectArgs[i] = convertToReflectValue(arg)
	//		}
	//
	//		// Call the wrapper function with the reflect arguments
	//		reflectResult := wrapper.Call(reflectArgs)
	//
	//		// Extract the []interface{} result and convert it to js.Value
	//		jsResult := reflectResult[0].Interface().([]interface{})
	//		return convertToJSValueResult(jsResult)
	//	})
	//
	return jsFunc
}

func valueOf(v reflect.Value) js.Value {
	switch v.Kind() {
	case reflect.Invalid:
		return js.Undefined()
	case reflect.Slice, reflect.Array:
		return valueOfSlice(v)
	case reflect.Map:
		return valueOfMap(v)
	case reflect.Struct:
		return valueOfStruct(v)
	case reflect.Complex64, reflect.Complex128:
		return valueOfComplex(v)
	case reflect.Interface, reflect.Pointer:
		return valueOfPointer(v)
	case reflect.Func:
		//return valueOfFunc(v)
		panic(errors.New("not implemented"))
	default:
		return js.ValueOf(v.Interface())
	}
}

type PromiseHandlerFunc func(resolve, reject js.Value) any

func Promisify(jsFunc PromiseHandlerFunc) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]
		defer recoverPanics(reject)

		jsFunc(resolve, reject)

		return nil
	})

	return promise.New(handler)
}
