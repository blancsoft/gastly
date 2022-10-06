package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chumaumenze/wago/src/ast"
)

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, pkgNames []string) ([]ast.Result, error) {
	return ast.FromPackages(pkgNames...), ctx.Err()
}
