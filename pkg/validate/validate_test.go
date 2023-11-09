package validate

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"testing"
)

func TestReadTarget(t *testing.T) {
	testCases := []struct {
		input    string
		expected *models.TargetData
	}{
		{
			input: `foo: bar`,
			expected: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": "bar",
					},
				},
				Format: "yaml",
			},
		},
		{
			input: `{
	"contoso": {
		"foo": "bar"
	}
}`,
			expected: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"contoso": map[string]interface{}{
							"foo": "bar",
						},
					},
				},
				Format: "json",
			},
		},
	}
	for _, tc := range testCases {
		targetData, err := readTarget(tc.input)
		if err != nil {
			t.Errorf("Error reading target data: %v", err)
			t.FailNow()
		}
		if !helpers.CompareInterfaces(targetData, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, targetData)
		}
	}
}

var validateTestCases = []struct {
	validations              string
	target                   string
	validationError          bool
	expectedValidationResult bool
}{
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
`,
		target: `
foo: bar
`,
		validationError:          false,
		expectedValidationResult: true,
	},
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
  errorMessage: "'object.foo should be bar but was ' + object.foo"
`,
		target: `
foo: baz
`,
		validationError:          false,
		expectedValidationResult: false,
	},
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
  errorMessage: "'foo should be bar but was ' + object.foo"
`,
		target: `{
	"foo": "baz"
}
`,
		validationError:          false,
		expectedValidationResult: false,
	},
	{
		validations: `validations:
- expression: "has(object.foo.bar)"
`,
		target: `{
	"foo": {
		"bar": "baz"
	}
}
`,
		validationError:          false,
		expectedValidationResult: true,
	},
}

func TestValidateWithRawData(t *testing.T) {
	for _, tc := range validateTestCases {
		result, err := Validate(tc.validations, tc.target)
		if err != nil && !tc.validationError {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.validationError {
			t.Errorf("Expected error, got none")
		}
		if result != tc.expectedValidationResult {
			t.Errorf("Expected %v, got %v", tc.expectedValidationResult, result)
		}
	}
}

func TestValidateWithFiles(t *testing.T) {
	for _, tc := range validateTestCases {
		validationsFile, err := helpers.CreateTempFile(tc.validations)
		if err != nil {
			t.Errorf("Error creating validations file: %v", err)
			t.FailNow()
		}
		targetFile, err := helpers.CreateTempFile(tc.target)
		if err != nil {
			t.Errorf("Error creating target file: %v", err)
			t.FailNow()
		}
		result, err := Validate(validationsFile.Name(), targetFile.Name())
		if err != nil && !tc.validationError {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.validationError {
			t.Errorf("Expected error, got none")
		}
		if result != tc.expectedValidationResult {
			t.Errorf("Expected %v, got %v", tc.expectedValidationResult, result)
		}
	}
}

var validateSingleExpressionTestCases = []struct {
	expression              string
	target                  string
	validationError         bool
	expectedValidationValue bool
}{
	{
		expression: `object.foo == 'bar'`,
		target: `
foo: bar
`,
		validationError:         false,
		expectedValidationValue: true,
	},
	{
		expression: `object.foo == 'bar'`,
		target: `
foo: baz
`,
		validationError:         false,
		expectedValidationValue: false,
	},
	{
		expression: `has(object.foo.bar)`,
		target: `{
	"foo": {
		"bar": "baz"
	}
}
`,
		validationError:         false,
		expectedValidationValue: true,
	},
}

func TestValidateSingleExpressionWithRawData(t *testing.T) {
	for _, tc := range validateSingleExpressionTestCases {
		result, err := ValidateSingleExpression(tc.expression, tc.target)
		if err != nil && !tc.validationError {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.validationError {
			t.Errorf("Expected error, got none")
		}
		if result != tc.expectedValidationValue {
			t.Errorf("Expected %v, got %v", tc.expectedValidationValue, result)
		}
	}
}

func TestValidateSingleExpressionWithFiles(t *testing.T) {
	for _, tc := range validateSingleExpressionTestCases {
		targetFile, err := helpers.CreateTempFile(tc.target)
		if err != nil {
			t.Errorf("Error creating target file: %v", err)
			t.FailNow()
		}
		result, err := ValidateSingleExpression(tc.expression, targetFile.Name())
		if err != nil && !tc.validationError {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.validationError {
			t.Errorf("Expected error, got none")
		}
		if result != tc.expectedValidationValue {
			t.Errorf("Expected %v, got %v", tc.expectedValidationValue, result)
		}
	}
}
