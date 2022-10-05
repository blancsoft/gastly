package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var indent = []byte("  ")

type JsonPointer []string

func (j *JsonPointer) String() string {
	return "/" + strings.Join(*j, "/")
}

func (j *JsonPointer) Push(v string, o ...string) {
	o = append([]string{v}, o...)
	*j = append(*j, o...)
}

func (j *JsonPointer) Pop() string {
	if len(*j) > 0 {
		n := len(*j) - 1
		v := (*j)[n:]
		*j = (*j)[:n]
		if len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

type AstWriter struct {
	jsonOutput  io.Writer
	pointerMap  map[any]string
	fset        *token.FileSet
	indentLevel int
	last        byte
	currentPath JsonPointer
	minify      bool
}

func (a *AstWriter) Write(data []byte) (n int, err error) {
	var m int
	for _, b := range data {
		if a.last == '\n' {
			for j := a.indentLevel; j > 0; j-- {
				_, err = a.jsonOutput.Write(indent)
				if err != nil {
					return
				}
			}
		}
		a.last = b
	}
	if len(data) > n {
		m, err = a.jsonOutput.Write(data[n:])
		n += m
	}
	return
}

func (a *AstWriter) walk(v reflect.Value, isKey bool, pos map[string]token.Position) {
	if !ast.NotNilFilter("", v) {
		a.printf("null")
		return
	}

	switch v.Kind() {
	case reflect.Interface:
		//z, b := v.Interface().(ast.Node)
		//fmt.Printf("%T, %t\n", z, b)

		if n, ok := v.Interface().(ast.Node); ok {
			loc := map[string]token.Position{
				"start": a.fset.Position(n.Pos()),
				"end":   a.fset.Position(n.End()),
			}
			a.walk(v.Elem(), isKey, loc)
		} else {
			a.walk(v.Elem(), isKey, pos)
		}
		//a.walk(v.Elem(), isKey, nil)

	case reflect.Map:
		a.printf("{\n")
		a.indentLevel++

		a.printf("%q: %q,\n", "@type", v.Type())
		a.printf("%q: %d", "@len", v.Len())
		if pos != nil {
			a.printf(",\n")
			a.printf("%q: ", "@pos")
			a.walk(reflect.ValueOf(pos), false, nil)
		}

		for _, key := range v.MapKeys() {
			a.printf(",\n")
			a.walk(key, true, nil)
			a.printf(": ")
			a.walk(v.MapIndex(key), false, nil)
			a.currentPath.Pop()
		}

		a.indentLevel--
		a.printf("\n")
		a.printf("}")

	case reflect.Pointer:
		// Ast may contain recursive pointers
		ptr := v.Interface()
		if ptrpath, exists := a.pointerMap[ptr]; exists {
			a.printf("{")
			a.printf("\n")
			a.indentLevel++
			a.printf(`"@type": %q,`, "RecursivePtr")
			a.printf("\n")
			a.printf(`"@path": %q`, ptrpath)
			a.printf("\n")
			a.indentLevel--
			a.printf("}")
		} else {
			a.pointerMap[ptr] = a.currentPath.String()
			a.walk(v.Elem(), isKey, pos)
		}

	case reflect.Array, reflect.Slice:
		a.printf("[")
		if v.Len() > 0 {
			a.indentLevel++
			for i, n := 0, v.Len(); i < n; i++ {
				a.currentPath.Push(strconv.Itoa(i))
				a.printf("\n")
				a.walk(v.Index(i), false, nil)
				if i != n-1 {
					a.printf(",")
				}
				a.currentPath.Pop()
			}
			a.indentLevel--
			a.printf("\n")
		}
		a.printf("]")

	case reflect.Struct:
		t := v.Type()
		a.printf("{\n")
		a.indentLevel++
		a.printf("%q: %q", "@type", v.Type())
		if pos != nil {
			a.printf(",\n")
			a.printf("%q: ", "@pos")
			a.walk(reflect.ValueOf(pos), false, nil)
		}

		//if n, ok := v.Interface().(ast.Node); ok {
		//	start := a.fset.Position(n.Pos())
		//	end := a.fset.Position(n.End())
		//	a.printf(",")
		//	a.printf("\n")
		//	a.printf("%q: ", "@pos")
		//	v := map[string]any{"start": start, "end": end}
		//	a.walk(reflect.ValueOf(v), false, false)
		//}

		for i, n := 0, t.NumField(); i < n; i++ {
			// exclude non-exported fields because their
			// values cannot be accessed via reflection
			if name := t.Field(i).Name; ast.IsExported(name) {
				a.printf(",\n")

				a.currentPath.Push(name)
				a.printf("%q: ", name)

				value := v.Field(i)
				a.walk(value, false, nil)
				a.currentPath.Pop()
			}
		}

		a.indentLevel--
		a.printf("\n")
		a.printf("}")

	default:
		v := v.Interface()
		switch v := v.(type) {
		case string:
			// print strings in quotes
			a.printf("%q", v)
			return
		case int, token.Pos:
			// position values can be printed nicely if we have a file set
			if a.fset != nil {
				a.printf("%d", v)
				//a.printf("%q", a.fset.Position(v))
				//a.walk(reflect.ValueOf(a.fset.Position(v)), isKey, false)
			} else {
				a.printf("%d", v)
			}
			a.printf("%d", v)
			return
		}
		// default
		a.printf(`"%v"`, v)
		if isKey {
			a.currentPath.Push(fmt.Sprintf(`"%v"`, v))
		}
	}
}

func (a *AstWriter) printf(format string, z ...any) {
	_, err := fmt.Fprintf(a, format, z...)
	if err != nil {
		log.Error(err)
	}
}

func ParseAst(fset *token.FileSet, v any) (r map[string]any, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var o bytes.Buffer
	a := AstWriter{
		jsonOutput: &o,
		pointerMap: make(map[any]string),
		fset:       fset,
	}
	a.walk(reflect.ValueOf(v), false, nil)

	os.WriteFile("writer.json", o.Bytes(), 0755)
	err = json.Unmarshal(o.Bytes(), &r)
	return
}
