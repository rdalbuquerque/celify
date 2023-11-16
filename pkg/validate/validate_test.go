package validate

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"celify/pkg/printer"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
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
	validations   string
	target        string
	expectedError error
}{
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
`,
		target: `
foo: bar
`,
		expectedError: nil,
	},
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
  MessageExpression: "'object.foo should be bar but was ' + object.foo"
`,
		target: `
foo: baz
`,
		expectedError: &multierror.Error{Errors: []error{printer.FmtError(errors.Errorf("Expression: object.foo == 'bar'\n\t  Error: 'object.foo should be bar but was baz'"))}},
	},
	{
		validations: `validations:
- expression: "object.foo == 'bar'"
  MessageExpression: "'foo should be bar but was ' + object.foo"
`,
		target: `{
	"foo": "baz"
}
`,
		expectedError: &multierror.Error{Errors: []error{printer.FmtError(errors.Errorf("Expression: object.foo == 'bar'\n\t  Error: 'foo should be bar but was baz'"))}},
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
		expectedError: nil,
	},
}

func TestValidateWithRawData(t *testing.T) {
	for _, tc := range validateTestCases {
		err := Validate(tc.validations, tc.target)
		if err != nil && tc.expectedError == nil {
			t.Errorf("Expected no error, got %v", err)
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
		err = Validate(validationsFile.Name(), targetFile.Name())
		if err != nil && tc.expectedError == nil {
			t.Errorf("Expected no error, got %v", err)
		}

	}
}

var validateSingleExpressionTestCases = []struct {
	expression    string
	target        string
	errorExpected bool
}{
	{
		expression: `object.foo == 'bar'`,
		target: `
foo: bar
`,
		errorExpected: false,
	},
	{
		expression: `object.foo == 'bar'`,
		target: `
foo: baz
`,
		errorExpected: true,
	},
	{
		expression: `has(object.foo.bar)`,
		target: `{
	"foo": {
		"bar": "baz"
	}
}
`,
		errorExpected: false,
	},
}

func TestValidateSingleExpressionWithRawData(t *testing.T) {
	for _, tc := range validateSingleExpressionTestCases {
		err := ValidateSingleExpression(tc.expression, tc.target)
		if err != nil && !tc.errorExpected {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.errorExpected {
			t.Errorf("Expected error, got none")
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
		err = ValidateSingleExpression(tc.expression, targetFile.Name())
		if err != nil && !tc.errorExpected {
			t.Errorf("Expected no error, got %v", err)
		}
		if err == nil && tc.errorExpected {
			t.Errorf("Expected error, got none")
		}
	}
}
