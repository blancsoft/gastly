package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/chumaumenze/wago/src/ast"
	log "github.com/sirupsen/logrus"
)

func Execute() {
	pkgNames := []string{
		"github.com/Davincible/goinsta",
		//"github.com/Davincible/goinsta/cmds",
		"github.com/Davincible/goinsta/utilities",
		"github.com/spf13/cobra",
		"github.com/gdexlab/go-render/render",
	}
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
