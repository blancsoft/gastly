package wago

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PkgConfigValue struct {
	// The package path just like you would import it in Go
	Path string `json:"path"`

	// Where this output should be written to.
	// If you specify a folder it will be written to a file within that folder.
	// By default, it is written into the Golang package folder.
	OutputPath string `json:"output_path,omitempty"`

	// Customize the indentation (use \t if you want tabs)
	Indent string `json:"indent,omitempty"`

	// Specify your own custom type translations, useful for custom types, `time.Time` and `null.String`.
	// Be default unrecognized types will be output as `any /* name */`.
	TypeMappings map[string]string `json:"type_mappings,omitempty"`

	// This content will be put at the top of the output Typescript file.
	// You would generally use this to import custom types.
	Frontmatter string `json:"frontmatter,omitempty"`

	// Filenames of Go source files that should not be included in the Typescript output.
	ExcludeFiles []string `json:"exclude_files,omitempty"`

	// Filenames of Go source files that should be included in the Typescript output.
	IncludeFiles []string `json:"include_files,omitempty"`
}

type Config struct {
	Packages []PkgConfigValue `json:"packages"`
}

func init() {
	viper.SetConfigName("wago")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/wago/")
	viper.AddConfigPath("$HOME/.wago/")
	viper.AddConfigPath("./")

	viper.SetEnvPrefix("WAGO")
	viper.AutomaticEnv()

	// TODO: Update the defaults
	viper.SetDefault("exportPath", "./")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found; using default values")
		} else {
			// Config file was found but another error was produced
			log.Fatal(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	viper.WatchConfig()
}
