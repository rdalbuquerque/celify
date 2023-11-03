// in cmd/validate.go
package cmd

import (
	"celify/pkg/validate"
	"os"

	"github.com/fatih/color"
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
		var result bool
		var err error
		if validations != "" && expression != "" {
			return errors.Errorf("You can only provide either a validations file or a single expression")
		}
		cmd.SilenceUsage = true
		if validations != "" {
			result, err = validate.Validate(validations, target)
			if err != nil {
				return errors.Errorf(color.New(color.FgRed).Add(color.Bold).Sprintf("Error during validation: %v", err))
			}
		} else {
			result, err = validate.ValidateSingleExpression(expression, target)
			if err != nil {
				return errors.Errorf(color.New(color.FgHiRed).Add(color.Bold).Sprintf("Error during validation: %v", err))
			}
		}
		if !result {
			color.New(color.FgRed).Add(color.Bold).Println("validation failed")
			os.Exit(1)
		}
		color.New(color.FgGreen).Add(color.Bold).Println("validation passed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you define the flags for the command
	validateCmd.Flags().StringVarP(&target, "target", "t", "", "Path to target file or raw string data")
	validateCmd.Flags().StringVarP(&validations, "validations", "v", "", "Path to the validations YAML file or raw string data - this has to be in correcy yaml format")
	validateCmd.Flags().StringVarP(&expression, "expression", "e", "", "single cel expression to evaluate against the target data")
}
