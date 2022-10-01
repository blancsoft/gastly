package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile = ""
var cliVersion = fmt.Sprintf("%s %s (%s) %s", Version(), Target(), CommitDate(), Commit())
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
