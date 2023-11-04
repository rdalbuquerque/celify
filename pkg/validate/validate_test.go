package validate

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"testing"
)

func TestReadValidations(t *testing.T) {
	//convert input to a byte slice
	input := `
validations:
- expression: "object.foo == 'bar'"
- expression: "object.foo == 'baz'"
  errorMessage: "'foo should be baz but was' + object.foo"
`
	expected := []models.ValidationRule{
		{
			Expression: "object.foo == 'bar'",
		},
		{
			Expression:   "object.foo == 'baz'",
			ErrorMessage: "'foo should be baz but was' + object.foo",
		},
	}
	validations, err := readValidations(input)
	if err != nil {
		t.Errorf("Error reading validations: %v", err)
	}
	if !helpers.CompareInterfaces(validations, expected) {
		t.Errorf("Expected %v, got %v", expected, validations)
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
	input := `
data:
  foo: bar
`
	expected := &models.TargetData{
		Data: map[string]interface{}{
			"object": map[string]interface{}{
				"data": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		Format: "yaml",
	}
	targetData, err := readTarget(input)
	if err != nil {
		t.Errorf("Error reading target data: %v", err)
	}
	if !helpers.CompareInterfaces(targetData, expected) {
		t.Errorf("Expected %v, got %v", expected, targetData)
	}
}
