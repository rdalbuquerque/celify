package validate

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"testing"
)

func TestReadValidations(t *testing.T) {
	testCases := []struct {
		input    string
		expected []models.ValidationRule
	}{
		{
			input: `
validations:
- expression: "object.foo == 'bar'"
- expression: "object.foo == 'baz'"
  errorMessage: "'foo should be baz but was' + object.foo"
`,
			expected: []models.ValidationRule{
				{
					Expression: "object.foo == 'bar'",
				},
				{
					Expression:   "object.foo == 'baz'",
					ErrorMessage: "'foo should be baz but was' + object.foo",
				},
			},
		},
		{
			input: `
validations:
- expression: "size(object.foo) == 1"
- expression: "has(object.foo)"
  errorMessage: "foo should be present"
`,
			expected: []models.ValidationRule{
				{
					Expression: "size(object.foo) == 1",
				},
				{
					Expression:   "has(object.foo)",
					ErrorMessage: "foo should be present",
				},
			},
		},
	}
	for _, tc := range testCases {
		validations, err := readValidations(tc.input)
		if err != nil {
			t.Errorf("Error reading validations: %v", err)
		}
		if !helpers.CompareInterfaces(validations, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, validations)
		}
	}
}

func TestReadValidationsWithInvalidYAML(t *testing.T) {
	input := `	- expression: "object.foo == 'bar'"
	- expression: "object.foo == 'baz'"
	  errorMessage: "'foo should be baz but was' + object.foo"`
	_, err := readValidations(input)
	if err == nil {
		t.Errorf("Expected error reading validations")
	}
}

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

func TestValidate(t *testing.T) {
	testCases := []struct {
		validations              string
		target                   string
		validationError          bool
		expectedValidationResult bool
	}{
		{
			validations: `
validations:
- expression: "object.foo == 'bar'"
`,
			target: `
foo: bar
`,
			validationError:          false,
			expectedValidationResult: true,
		},
		{
			validations: `
validations:
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
			validations: `
validations:
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
			validations: `
validations:
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

	for _, tc := range testCases {
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

func TestValidateSingleExpression(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
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
