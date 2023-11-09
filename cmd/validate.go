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
	Short:         "Validate a target YAML against validation rules",
	Long:          `This command validates a target YAML file against a set of validation rules defined in a validations file.`,
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
