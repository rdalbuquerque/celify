// in cmd/validate.go
package cmd

import (
	"celify/pkg/validate"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var target string
var validations string
var expression string

var validateCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "validate",
	Short:         "Validate yaml or json files using CEL expressions",
	Long: `celify is a CLI tool designed to validate structured data, supporting both JSON and YAML formats, against a defined set of Common Expression Language (CEL) rules.

	The tool aims to provide developers and operators with a flexible and concise way to enforce custom validation checks on their data, ensuring data consistency, correctness, and adherence to specified rules.
	
	Examples:
	
	1. Validate a YAML file against a set of rules:
	   $ celify validate --target deployment.yaml --validations validations.yaml

	1. Validate a JSON file against a set of rules:
	   $ celify validate --target tfplan.json --validations validations.yaml
	
	2. Validate remote data:
	   $ celify validate --target "$(curl -s https://example.com/data.json)" --validations validations.yaml
	
	3. Validate a YAML file against a single expression:
	   $ celify validate --target deployment.yaml --expression "object.spec.replicas > 1"
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validations != "" && expression != "" {
			return errors.Errorf("You can only provide either a validations file or a single expression")
		}
		cmd.SilenceUsage = true
		if validations != "" {
			return validate.Validate(validations, target)
		} else {
			return validate.ValidateSingleExpression(expression, target)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you define the flags for the command
	validateCmd.Flags().StringVarP(&target, "target", "t", "", "Path to target file or raw string data")
	validateCmd.Flags().StringVarP(&validations, "validations", "v", "", "Path to the validations YAML file or raw string data - this has to be in correcy yaml format")
	validateCmd.Flags().StringVarP(&expression, "expression", "e", "", "single cel expression to evaluate against the target data")
}
