/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "celify",
	Short: "Validate yaml or json files using CEL expressions",
	Long: `celify is a CLI tool designed to validate structured data, supporting both JSON and YAML formats, against a defined set of Common Expression Language (CEL) rules.

	The tool aims to provide developers and operators with a flexible and concise way to enforce custom validation checks on their data, ensuring data consistency, correctness, and adherence to specified rules.
	
	Features:
	- Support for both JSON and YAML input formats.
	- Ability to define custom validation rules using CEL
	- Flexibility to provide input data and rules via files or raw strings.
	- Detailed error messages guiding users to the exact validation failure point.
	
	Examples:
	
	1. Validate a YAML file against a set of rules:
	   $ celify validate --target deployment.yaml --validations validations.yaml

	1. Validate a JSON file against a set of rules:
	   $ celify validate --target tfplan.json --validations validations.yaml
	
	2. Validate remote data:
	   $ celify validate --target "$(curl -s https://example.com/data.json)" --validations validations.yaml
	
	`,
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
