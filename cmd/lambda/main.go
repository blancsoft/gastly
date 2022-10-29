package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"

	"github.com/blancsoft/gastly/ast"
)

var GOROOT = "/opt/go"
var GOPATH string
var PATH string

func init() {
	// Use temporary directory for GOPATH
	var err error
	if GOPATH, err = os.MkdirTemp(os.TempDir(), "go*"); err != nil {
		log.Fatal(err)
	}

	PATH = filepath.Join(GOROOT, "bin") + ":" + ast.GetEnvDefault("PATH", "")
	GOPATH = GOPATH + ":" + ast.GetEnvDefault("GOPATH", "")

	ast.SetEnv("PATH", strings.Trim(PATH, ":"))
	ast.SetEnv("GOPATH", strings.Trim(GOPATH, ":"))
	ast.SetEnv("GOROOT", strings.Trim(GOROOT, ":"))
	ast.SetEnv("GOCACHE", os.TempDir())
	ast.DumpEnv()
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, pkgNames []string) ([]ast.Result, error) {
	return ast.FromPackages(pkgNames...), ctx.Err()
}
