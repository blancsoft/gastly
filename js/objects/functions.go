//go:build js && wasm

package objects

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/blancsoft/gastly/ast"
	gastlyJS "github.com/blancsoft/gastly/js"
)

func FromSourceCode(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return fmt.Errorf("FromSourceCode: missing required arguments — expects two arguments")
	}
	fname := args[0].String()
	code := args[1].String()
	result := ast.FromSourceCode(fname, code)
	return gastlyJS.ValueOf(result).Unwrap()
}

func FromPackages(this js.Value, args []js.Value) any {
	var pkgNames []string
	for _, pn := range args {
		pkgNames = append(pkgNames, pn.String())
	}
	result := ast.FromPackages(pkgNames...)
	return gastlyJS.ValueOf(result).Unwrap()
}

func GetRepositoryDetails(this js.Value, args []js.Value) any {
	// TODO: implement me
	panic(errors.New("not implemented"))
}

func FetchRepository(this js.Value, args []js.Value) any {
	// TODO: implement me
	panic(errors.New("not implemented"))
}
