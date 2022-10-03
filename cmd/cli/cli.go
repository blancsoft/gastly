package main

import (
	"fmt"
	"os"

	"github.com/chumaumenze/wago/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile = ""
var cliVersion = fmt.Sprintf("%s %s (%s) %s", cmd.Version(), cmd.Target(), cmd.CommitDate(), cmd.Commit())
var rootCmd = &cobra.Command{
	Use:     "wago",
	Short:   "Go packages to webassembly",
	Long:    `Go packages to webassembly`,
	Version: cliVersion,
}

func init() {
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		}
	})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Set configuration file")

	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
