package validate

import (
	"os"

	"celify/pkg/evaluator"
	"celify/pkg/helpers"

	"celify/pkg/models"

	"github.com/pkg/errors"
)

func ValidateSingleExpression(expression, targetInput string) (bool, error) {
	targetData, err := readTarget(targetInput)
	if err != nil {
		return false, errors.Errorf("Error reading target: %v", err)
	}

	eval, err := evaluator.NewEvaluator(targetData)
	if err != nil {
		return false, errors.Errorf("Error creating evaluator: %v", err)
	}
	result, err := eval.EvaluateSingleExpression(expression)
	if err != nil {
		return false, errors.Errorf("Error evaluating expression: %v", err)
	}

	boolResult, ok := result.(bool)
	if !ok {
		return false, errors.New("result is not a boolean value")
	}
	return boolResult, nil
}

func Validate(validationInput, targetInput string) error {
	// Load validation rules
	var validations models.ValidationConfig
	_, err := unmarshalData(validationInput, &validations)
	if err != nil {
		return false, errors.Errorf("Error reading validations: %v", err)
	}

	// Load target YAML data
	targetData, err := readTarget(targetInput)
	if err != nil {
		return false, errors.Errorf("Error reading target: %v", err)
	}

	eval, err := evaluator.NewEvaluator(targetData)
	if err != nil {
		return false, errors.Errorf("Error creating evaluator: %v", err)
	}
	results := eval.Evaluate(validations)
	for _, result := range results {
		if result.ValidationError != nil {
			return errors.Errorf("Error validating data")
		}
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
