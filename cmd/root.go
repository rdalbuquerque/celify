/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "celify",
	Short: "Validate yaml or json files using CEL expressions",
	Long: `celify is a CLI tool designed to interact structured data, supporting both JSON and YAML formats.
	
	Features:
	- Support for both JSON and YAML input formats.
	- Ability to define custom validation rules using CEL
	- Flexibility to provide input data and rules via files or raw strings.
	- Detailed error messages guiding users to the exact validation failure point.
	`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
