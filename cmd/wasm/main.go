//go:build js && wasm

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"syscall/js"

	"github.com/blancsoft/gastly/ast"
)

func main() {
	println("This is GASTly Renderer")

	gastly := js.ValueOf(map[string]any{
		"FromSourceCode": js.FuncOf(func(this js.Value, args []js.Value) any {
			if len(args) < 2 {
				return fmt.Errorf("FromSourceCode: missing required arguments — expects two arguments")
			}
			fname := args[0].String()
			code := args[1].String()
			result := ast.FromSourceCode(fname, code)
			return ast.ConvertToJSValue(result)
		}),
		"FromPackages": js.FuncOf(func(this js.Value, args []js.Value) any {
			var pkgNames []string
			for _, pn := range args {
				pkgNames = append(pkgNames, pn.String())
			}
			result := ast.FromPackages(pkgNames...)
			return ast.ConvertToJSValue(result)
		}),
	})
	js.Global().Set("Gastly", gastly)

	count := 0
	for {
		fmt.Println("Count: ", count)
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
		func() {
			defer func() {
				if v := recover(); v != nil {
					fmt.Printf("PANIC caught: %s\n", v)
				}
			}()

			fmt.Println("This is Gastly!")

			// block until interrupt/terminate signal
			if terminate := <-quitChannel; terminate != nil {
				fmt.Println("Goodbye from Gastly!")
			}

			fmt.Println("Recovering from panic...")
			return
		}()
		count++
	}
}
