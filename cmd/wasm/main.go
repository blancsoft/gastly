//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/blancsoft/gastly/js/objects"
)

func main() {
	js.Global().Set("Gastly", js.ValueOf(map[string]any{
		"FromSourceCode": js.FuncOf(objects.FromSourceCode),
		"FromPackages":   js.FuncOf(objects.FromPackages),
	}))

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
