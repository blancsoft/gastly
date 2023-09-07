//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/chumaumenze/gastly/lib/objects"
	"github.com/chumaumenze/gjs"
)

func main() {
	gastly := js.ValueOf(map[string]any{
		"ping": js.FuncOf(func(this js.Value, args []js.Value) any {
			return gjs.ValueOf("pong!").Into()
		}),
		"FromSourceCode": js.FuncOf(objects.FromSourceCode),
		"FromPackages":   js.FuncOf(objects.FromPackages),
	})
	js.Global().Set("Gastly", gastly)

	for {
		quitChannel := make(chan *struct{})
		func() {
			defer func() {
				if v := recover(); v != nil {
					fmt.Printf("PANIC caught: %s\n", v)
				}
			}()

			println("This is GASTly Renderer")
			// block until interrupt/terminate signal
			if terminate := <-quitChannel; terminate != nil {
				fmt.Println("Goodbye from Gastly!")
			}

			fmt.Println("Recovering from panic...")
			return
		}()
	}
}
