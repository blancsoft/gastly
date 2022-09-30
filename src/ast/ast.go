package ast

import (
	"bytes"
	"fmt"
	//"go/parser"
	//"go/token"

	//"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	Ast  *Ast
	Dump string
	Err  error
}

func Generate(pkgNames ...string) (map[string]Result, error) {
	var pkgs []*packages.Package
	var err error
	if pkgs, err = loadPackages(pkgNames...); err != nil {
		return map[string]Result{}, err
	}
	defer func() {
		err := removePackages(false, pkgNames...)
		if err != nil {
			log.Error(err)
		}
	}()

	accumulator := func(x string, y packages.Error, i int) string {
		return fmt.Sprintf("%s\n========\n%s\n", x, y.Msg)
	}
	pkgAstMap := map[string]Result{}
	for _, p := range pkgs {
		var err error
		if errMsg := lo.Reduce(p.Errors, accumulator, ""); errMsg != "" {
			log.Error(errMsg)
		}

		var pkgAstBuffer bytes.Buffer
		if pkgAstBuffer, err = dumpAst(p.Fset, p.Syntax); err != nil {
			return map[string]Result{}, eris.Wrapf(err, "unable to dump AST")
		}

		pkgAst, err := BuildAst("", p.Syntax)
		pkgAstMap[p.ID] = Result{
			Ast:  pkgAst,
			Dump: pkgAstBuffer.String(),
			Err:  err,
		}
	}
	return pkgAstMap, nil
}
