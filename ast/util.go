package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"os/exec"
	"reflect"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

func loadPackages(pkgNames ...string) ([]*packages.Package, error) {
	defer func() {
		err := removePackages(false, pkgNames...)
		if err != nil {
			log.Error(err)
		}
	}()
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
	}, pkgNames...)
}

func removePackages(purgeAll bool, pkgNames ...string) error {
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
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Infof("Executing command %s\n", cmd.String())
	err = cmd.Run()
	log.Debug("===StdOut: %s===\n%s", cmd.String(), stdout.String())
	log.Debug("===StdErr: %s===\n%s", cmd.String(), stderr.String())
	return stdout, stderr, err
}

func dumpAst(fset *token.FileSet, x any) (pkgAstBuffer bytes.Buffer, err error) {
	err = ast.Fprint(&pkgAstBuffer, fset, x, func(string, reflect.Value) bool {
		return true
	})
	return
}
