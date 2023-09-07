package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

var pkgLoaderDir string

func ensureGoModule() {
	var err error
	pkgLoaderDir, err = os.MkdirTemp(os.TempDir(), "go*")
	if err != nil {
		log.Fatal(eris.Wrapf(err, "unable to create package loader directory"))
	}

	CreateWrite(filepath.Join(pkgLoaderDir, "go.mod"), "module main\n\ngo 1.16\n")
}

func loadPackages(pkgNames ...string) (pkgs []*packages.Package, err error) {
	for _, name := range pkgNames {
		log.Infof("Fetching package %s", name)
		// FIXME: Disable path lookup; use absolute path
		_, stderr, err := runCmd("go", "get", "-x", name)

		if err != nil {
			errMsg := fmt.Sprintf("cannot fetch package %v: %v\n", name, err)
			return nil, fmt.Errorf("%s\n===Error Output===\n%v", errMsg, stderr.String())
		}
	}

	return packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes,
		Dir:  pkgLoaderDir,
	}, pkgNames...)
}

//lint:ignore U1000 removePackages will be removed uf not needed after ast.FromPackages is refactored
func removePackages(purgeAll bool, pkgNames ...string) error {
	// nolint: U1000
	if purgeAll {
		log.Info("Purging all build cache")
		_, _, err := runCmd("go", "clean", "-x", "-cache")
		return err
	}

	for _, name := range pkgNames {
		log.Infof("Removing package %s", name)
		_, stderr, err := runCmd("go", "clean", "-x", "-r", "-i", name)

		if err != nil {
			errMsg := fmt.Sprintf("failed removing package %v: %v\n", name, err)
			return fmt.Errorf("%s\n===Error Output===\n%v", errMsg, stderr.String())
		}
	}
	return nil
}

func runCmd(name string, args ...string) (stdout, stderr bytes.Buffer, err error) {
	ensureGoModule()
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = pkgLoaderDir

	log.Infof("Executing command %s\n", cmd.String())
	err = cmd.Run()
	log.Debug("===StdOut: " + cmd.String() + "===\n" + stdout.String())
	log.Debug("===StdErr: " + cmd.String() + "===\n" + stderr.String())
	return stdout, stderr, err
}

func dumpAst(fset *token.FileSet, x any) (pkgAstBuffer bytes.Buffer, err error) {
	err = ast.Fprint(&pkgAstBuffer, fset, x, func(string, reflect.Value) bool {
		return true
	})
	return
}

func DumpEnv() {
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}
}

func SetEnv(k, v string) {
	if err := os.Setenv(k, v); err != nil {
		log.Fatal(err)
	}
}

func GetEnvDefault(k, v string) (d string) {
	var ok bool
	if d, ok = os.LookupEnv(k); !ok {
		d = v
	}
	return
}

func CreateWrite(fname, content string) {
	_ = os.MkdirAll(filepath.Dir(fname), 0755)
	fd, err := os.Create(fname)
	defer func() {
		err := fd.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fd.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenRead(fname string) string {
	bb, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	return string(bb)
}
