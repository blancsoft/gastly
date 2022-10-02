package ast

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"

	//"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	Name   string
	Ast    *Ast
	Dump   string
	Source map[string]string
	Err    error
}

func FromPackages(pkgNames ...string) []Result {
	var pkgs []*packages.Package
	var err error
	if pkgs, err = loadPackages(pkgNames...); err != nil {
		log.Error(err)
	}
	accumulator := func(x string, y packages.Error, i int) string {
		return fmt.Sprintf("%s\n========\n%s\n", x, y.Msg)
	}
	var pkgAsts []Result
	for _, p := range pkgs {
		if errMsg := reduce(p.Errors, accumulator, ""); errMsg != "" {
			r := Result{Name: p.ID, Err: eris.New(errMsg)}
			pkgAsts = append(pkgAsts, r)
			continue
		}

		srcs := map[string]string{}
		for i := 0; i < len(p.Syntax); i++ {
			fname := p.GoFiles[i]
			node := p.Syntax[i]

			var source bytes.Buffer
			if err := printer.Fprint(&source, token.NewFileSet(), node); err != nil {
				r := Result{Name: p.ID, Err: eris.Wrapf(err, "unable to collate pkg source")}
				pkgAsts = append(pkgAsts, r)
				continue
			}
			srcs[fname] = source.String()
		}

		t, d, err := Generate(p.Fset, p.Syntax)
		r := Result{
			Name:   p.ID,
			Ast:    t,
			Dump:   d.String(),
			Source: srcs,
			Err:    err,
		}
		pkgAsts = append(pkgAsts, r)
	}
	return pkgAsts
}

func FromSourceCode(fname string, code string) Result {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fname, code, parser.ParseComments)
	if err != nil {
		return Result{Name: fname, Err: err}
	}

	var source bytes.Buffer
	if err := printer.Fprint(&source, fset, node); err != nil {
		return Result{Name: fname, Err: eris.Wrapf(err, "unable to collate source")}
	}

	t, d, err := Generate(fset, node)
	if err != nil {
		return Result{Name: fname, Err: err}
	}

	return Result{
		Name:   fname,
		Ast:    t,
		Dump:   d.String(),
		Source: map[string]string{fname: source.String()},
		Err:    err,
	}
}

func Generate(fset *token.FileSet, node any) (*Ast, bytes.Buffer, error) {
	var err error
	var pkgAstBuffer bytes.Buffer
	if pkgAstBuffer, err = dumpAst(fset, node); err != nil {
		return nil, bytes.Buffer{}, eris.Wrapf(err, "unable to dump AST")
	}

	pkgAst, err := BuildAst("", node)
	if err != nil {
		return nil, bytes.Buffer{}, eris.Wrapf(err, "unable to build AST")
	}
	return pkgAst, pkgAstBuffer, nil
}

type accumulatorFunc func(string, packages.Error, int) string

func reduce(collection []packages.Error, accumulator accumulatorFunc, initial string) string {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}
