package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/chumaumenze/wago/src/ast"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.Flags().Bool("json", false, "Export AST as json")
	generateCmd.Flags().Bool("src", false, "Export AST source code")
	generateCmd.Flags().Bool("file", false, "Read input as file")
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

		var results []ast.Result
		if fileArg {
			for _, filename := range args {
				code := openRead(filename)
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

			base := strings.ReplaceAll(r.Name, "/", "_")
			if exportJson {
				fname := base + ".ast.json"
				bb, err := json.Marshal(r.Ast)
				if err != nil {
					log.Error(err)
				}
				createWrite(fname, string(bb))
			}
			if exportSrc {
				fname := base + ".src.txt"
				createWrite(fname, r.Source)
			}
			fname := base + ".dump.txt"
			createWrite(fname, r.Dump)
		}
	},
}

func createWrite(fname, content string) {
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

func openRead(fname string) string {
	bb, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	return string(bb)
}
