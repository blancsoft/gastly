package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/chumaumenze/wago/src/ast"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate pkg...",
	Short: "Generate AST files",
	Long:  `Generate AST files for the specified packages`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, pkgNames []string) {
		exportJson, _ := cmd.Flags().GetBool("json")
		results, err := ast.Generate(pkgNames...)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range pkgNames {
			r := results[v]
			if r.Err != nil {
				log.Error(r.Err)
			}

			base := strings.ReplaceAll(v, "/", "_")
			if exportJson {
				afd, err := os.Create(base + ".ast.json")
				if err != nil {
					log.Error(err)
				}
				bb, err := json.Marshal(r.Ast)
				if err != nil {
					log.Error(err)
				}
				_, err = afd.WriteString(string(bb))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				dfd, err := os.Create(base + ".dump.txt")
				if err != nil {
					log.Error(err)
				}
				_, err = dfd.WriteString(r.Dump)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

func init() {
	generateCmd.Flags().Bool("json", false, "Export AST as json")
}
