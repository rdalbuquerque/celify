package validate

import (
	"os"

	"github.com/go-yaml/yaml"

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

	return result.(bool), nil
}

func Validate(validationInput, targetInput string) (bool, error) {
	// Load validation rules
	validations, err := readValidations(validationInput)
	if err != nil {
		return false, errors.Errorf("Error reading validations: %v\n", err)
	}

	// Load target YAML data
	targetData, err := readTarget(targetInput)
	if err != nil {
		return false, errors.Errorf("Error reading target: %v\n", err)
	}

	eval, err := evaluator.NewEvaluator(targetData)
	if err != nil {
		return false, errors.Errorf("Error creating evaluator: %v\n", err)
	}
	result, err := eval.Evaluate(validations)
	if err != nil {
		return false, errors.Errorf("Error evaluating expression: %v\n", err)
	}

	return result.(bool), nil
}

func readValidations(input string) ([]models.ValidationRule, error) {
	//convert input to a byte slice
	configData := []byte(input)

	var vals models.ValidationConfig
	err := yaml.Unmarshal(configData, &vals)
	if err != nil {
		configData, err := os.ReadFile(input)
		if err != nil {
			return nil, errors.Errorf("Error reading validations: %v\n", err)
		}
		err = yaml.Unmarshal(configData, &vals)
		if err != nil {
			return nil, errors.Errorf("Error parsing validations YAML: %v\n", err)
		}
	}

	return vals.Validations, nil
}

func readTarget(input string) (*models.TargetData, error) {
	targetData := []byte(input)

	var targetObject map[string]interface{}
	var format string
	format, err := helpers.UnmarshalData(targetData, &targetObject)
	if err != nil {
		targetData, err := os.ReadFile(input)
		if err != nil {
			return nil, errors.Errorf("Error reading target: %v\n", err)
		}
		format, err = helpers.UnmarshalData(targetData, &targetObject)
		if err != nil {
			return nil, errors.Errorf("Error parsing target YAML: %v\n", err)
		}
	}

	return &models.TargetData{
		Data:   map[string]interface{}{"object": targetObject},
		Format: format,
	}, nil
}
