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
	Ast    *Ast
	Dump   string
	Source string
	Name   string
	Err    error
}

func FromPackages(pkgNames ...string) []Result {
	var pkgs []*packages.Package
	var err error
	if pkgs, err = loadPackages(pkgNames...); err != nil {
		log.Error(err)
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
	var pkgAsts []Result
	for _, p := range pkgs {
		if errMsg := reduce(p.Errors, accumulator, ""); errMsg != "" {
			r := Result{Name: p.Name, Err: eris.New(errMsg)}
			pkgAsts = append(pkgAsts, r)
			continue
		}
		r := Generate(p.Name, p.Fset, p.Syntax)
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
	return Generate(fname, fset, node)
}

func FromSourceCodes(fname string, codes ...string) []Result {
	var r []Result
	for _, code := range codes {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, fname, code, parser.ParseComments)
		if err != nil {
			r = append(r, Result{Name: fname, Err: err})
		}
		r = append(r, Generate(fname, fset, node))
	}
	return r
}

func Generate(name string, fset *token.FileSet, node any) Result {
	var err error
	var pkgAstBuffer bytes.Buffer
	if pkgAstBuffer, err = dumpAst(fset, node); err != nil {
		return Result{Err: eris.Wrapf(err, "unable to dump AST")}
	}

	var source bytes.Buffer
	if err := printer.Fprint(&source, fset, node); err != nil {
		return Result{Err: eris.Wrapf(err, "unable to collate source")}
	}
	pkgAst, err := BuildAst("", node)
	return Result{
		Ast:    pkgAst,
		Dump:   pkgAstBuffer.String(),
		Source: source.String(),
		Name:   name,
		Err:    err,
	}
}

type accumulatorFunc func(string, packages.Error, int) string

func reduce(collection []packages.Error, accumulator accumulatorFunc, initial string) string {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}
