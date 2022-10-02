//go:build (js && ecmascript) || (js && wasm)

package main

//go:generate gopherjs build -o "../../dist/gastly.js" --minify --no_cache --verbose
//go:generate go build -x -o "../../dist/gastly.wasm"

import (
	"github.com/chumaumenze/wago/src/ast"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	println("This is GASTly Renderer")

	js.Global.Set("Gastly", map[string]any{
		"FromSourceCode":  ast.FromSourceCode,
		"FromSourceCodes": ast.FromSourceCodes,
		"FromPackages":    ast.FromPackages,
		"Generate":        ast.Generate,
	})
}
