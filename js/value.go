//go:build js && wasm

package js

import "syscall/js"

type Value js.Value

func (v Value) Unwrap() js.Value {
	return js.Value(v)
}
