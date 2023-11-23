package validate

import (
	"fmt"
	"os"

	"celify/pkg/evaluator"
	"celify/pkg/helpers"
	"celify/pkg/printer"

	"celify/pkg/models"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func ValidateSingleExpression(expression, targetInput string, supressObjects bool) error {
	targetData, err := readTarget(targetInput)
	if err != nil {
		return errors.Errorf("Error reading target: %v", err)
	}

	eval, err := evaluator.NewEvaluator(targetData)
	if err != nil {
		return errors.Errorf("Error creating evaluator: %v", err)
	}
	result := eval.EvaluateRule(models.ValidationRule{Expression: expression})
	printer := printer.NewPrinter(eval)
	printer.PrintResults([]models.EvaluationResult{result}, supressObjects)
	return getErrors([]models.EvaluationResult{result})
}

func Validate(validationInput, targetInput string, supressObjects bool) error {
	// Load validation rules
	var validations models.ValidationConfig
	_, err := unmarshalData(validationInput, &validations)
	if err != nil {
		return errors.Errorf("Error reading validations: %v", err)
	}

	// Load target YAML data
	targetData, err := readTarget(targetInput)
	if err != nil {
		return errors.Errorf("Error reading target: %v", err)
	}

	eval, err := evaluator.NewEvaluator(targetData)
	if err != nil {
		return errors.Errorf("Error creating evaluator: %v", err)
	}
	results := eval.Evaluate(validations)
	printer := printer.NewPrinter(eval)
	printer.PrintResults(results, supressObjects)
	return getErrors(results)
}

func unmarshalData(input string, output interface{}) (string, error) {
	//convert input to a byte slice
	configData := []byte(input)

	format, err := helpers.UnmarshalData(configData, output)
	if err != nil {
		configData, err := os.ReadFile(input)
		if err != nil {
			return "", errors.Errorf("Error reading data: %v", err)
		}
		format, err = helpers.UnmarshalData(configData, output)
		if err != nil {
			return "", errors.Errorf("Error parsing validations YAML: %v", err)
		}
	}

	return format, nil
}

func readTarget(input string) (*models.TargetData, error) {
	var targetObject map[string]interface{}
	var format string
	format, err := unmarshalData(input, &targetObject)
	if err != nil {
		return nil, errors.Errorf("Error parsing target data: %v", err)
	}
	return &models.TargetData{
		Data:   map[string]interface{}{"object": targetObject},
		Format: format,
	}, nil
}

func getErrors(results []models.EvaluationResult) error {
	multiErr := &multierror.Error{Errors: []error{}}
	for _, result := range results {
		if result.ValidationError != nil {
			multiErr.Errors = append(multiErr.Errors, fmt.Errorf("expression: %s\n\t  error: %v", result.Expression, result.ValidationError))
		}
	}
	if len(multiErr.Errors) == 0 {
		return nil
	}
	return printer.FmtError(multiErr.ErrorOrNil())
}
