package ast

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/packages"

	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	Name   string            `json:"name,omitempty"`
	Ast    string            `json:"ast,omitempty"`
	Dump   string            `json:"dump,omitempty"`
	Source map[string]string `json:"source,omitempty"`
	Err    error             `json:"err,omitempty"`
	ErrMsg string            `json:"errMsg,omitempty"`
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

		t, d, err := generate(p.Fset, p.Syntax)
		r := Result{
			Name:   p.ID,
			Ast:    t.String(),
			Dump:   d.String(),
			Source: srcs,
			Err:    err,
			ErrMsg: err.Error(),
		}
		pkgAsts = append(pkgAsts, r)
	}
	return pkgAsts
}

func FromSourceCode(fname string, code string) Result {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fname, code, parser.ParseComments)
	if err != nil {
		return Result{Name: fname, Err: err, ErrMsg: err.Error()}
	}

	var source bytes.Buffer
	if err := printer.Fprint(&source, fset, node); err != nil {
		err = eris.Wrapf(err, "unable to collate source")
		return Result{Name: fname, Err: err, ErrMsg: err.Error()}
	}

	t, d, err := generate(fset, node)
	if err != nil {
		return Result{Name: fname, Err: err, ErrMsg: err.Error()}
	}

	return Result{
		Name:   fname,
		Ast:    t.String(),
		Dump:   d.String(),
		Source: map[string]string{fname: source.String()},
	}
}

func generate(fset *token.FileSet, node any) (*bytes.Buffer, bytes.Buffer, error) {
	var err error
	var pkgAstBuffer bytes.Buffer
	if pkgAstBuffer, err = dumpAst(fset, node); err != nil {
		return nil, bytes.Buffer{}, eris.Wrapf(err, "unable to dump AST")
	}

	pkgAst, err := ParseAst(fset, node)
	if err != nil {
		return nil, bytes.Buffer{}, eris.Wrapf(err, "unable to build AST")
	}
	return &pkgAst, pkgAstBuffer, nil
}

type accumulatorFunc func(string, packages.Error, int) string

func reduce(collection []packages.Error, accumulator accumulatorFunc, initial string) string {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}
