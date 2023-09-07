//go:build js && wasm
// +build js,wasm

package objects

import (
	"fmt"
	"syscall/js"

	"github.com/chumaumenze/gastly/lib/ast"
	"github.com/chumaumenze/gjs"
)

func FromSourceCode(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return fmt.Errorf("FromSourceCode: missing required arguments — expects two arguments")
	}
	fname := args[0].String()
	code := args[1].String()
	result := ast.FromSourceCode(fname, code)
	return gjs.ValueOf(result).Into()
}

func FromPackages(this js.Value, args []js.Value) any {
	var pkgNames []string
	for _, pn := range args {
		pkgNames = append(pkgNames, pn.String())
	}
	result := ast.FromPackages(pkgNames...)
	return gjs.ValueOf(result).Into()
}
