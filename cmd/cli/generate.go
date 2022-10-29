package main

import (
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/blancsoft/gastly/ast"
)

func init() {
	generateCmd.Flags().Bool("json", false, "Export AST as json")
	generateCmd.Flags().Bool("src", false, "Export AST source code")
	generateCmd.Flags().Bool("file", false, "Read input as file")
	generateCmd.Flags().String("out", "output", "Output directory")
}

var generateCmd = &cobra.Command{
	Use:   "generate pkg...",
	Short: "Generate AST files",
	Long:  `Generate AST files for the specified packages`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exportJson, _ := cmd.Flags().GetBool("json")
		exportSrc, _ := cmd.Flags().GetBool("src")
		fileArg, _ := cmd.Flags().GetBool("file")
		outdirArg, _ := cmd.Flags().GetString("out")

		var results []ast.Result
		if fileArg {
			for _, filename := range args {
				code := ast.OpenRead(filename)
				r := ast.FromSourceCode(filename, code)
				results = append(results, r)
			}
		} else {
			results = ast.FromPackages(args...)
		}

		for _, r := range results {
			if r.Err != nil {
				log.Error(r.Err)
				continue
			}

			base := filepath.Join(outdirArg, strings.ReplaceAll(r.Name, "/", "_"))
			if exportJson {
				fname := base + ".ast.json"
				//a, _ := json.Marshal(r.Ast)
				ast.CreateWrite(fname, r.Ast)
			}
			if exportSrc {
				for k, v := range r.Source {
					b := filepath.Base(k)
					b = strings.TrimSuffix(b, filepath.Ext(b))
					fname := filepath.Join(base, b+".src.go")
					ast.CreateWrite(fname, v)
				}
			}
			fname := base + ".dump.txt"
			ast.CreateWrite(fname, string(r.Dump))
		}
	},
}
