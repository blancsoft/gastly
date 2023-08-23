//go:build js && wasm

package js

import (
	"fmt"
	"syscall/js"
)

func recoverPanics(cb js.Value) {
	jsError := global.Get("Error")
	warn := global.Get("console").Get("warn")
	if cb.Type() != js.TypeFunction {
		warn.Invoke("Callback argument to recoverPanics should be a JS function.")
		cb = warn
	}
	if r := recover(); r != nil {
		err := jsError.New(fmt.Sprintf("%+v", r))
		cb.Invoke(err)
	}
}
