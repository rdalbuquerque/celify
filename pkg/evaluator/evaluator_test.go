package evaluator

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"testing"
)

var evalTests = []struct {
	targetData  *models.TargetData
	validations []models.ValidationRule
	expected    any
}{
	{
		targetData: &models.TargetData{
			Data: map[string]interface{}{
				"object": map[string]interface{}{
					"foo": "bar",
				},
			},
			Format: "yaml",
		},
		validations: []models.ValidationRule{
			{
				Expression: "object.foo == 'bar'",
			},
		},
		expected: true,
	},
	{
		targetData: &models.TargetData{
			Data: map[string]interface{}{
				"object": map[string]interface{}{
					"foo": "bar",
				},
			},
			Format: "yaml",
		},
		validations: []models.ValidationRule{
			{
				Expression:   "object.foo == 'baz'",
				ErrorMessage: "'foo should be baz but was' + object.foo",
			},
		},
		expected: false,
	},
}

func TestEvaluate(t *testing.T) {
	for _, test := range evalTests {
		eval, err := NewEvaluator(test.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		result, err := eval.Evaluate(test.validations)
		if err != nil {
			t.Errorf("Error evaluating expression: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestEvaluateSingleExpression(t *testing.T) {
	for _, test := range evalTests {
		eval, err := NewEvaluator(test.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		result, err := eval.EvaluateSingleExpression(test.validations[0].Expression)
		if err != nil {
			t.Errorf("Error evaluating expression: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result)
		}
	}
}

func TestGetEvaluatedObject(t *testing.T) {
	var testCases = []struct {
		targetData *models.TargetData
		expression string
		expected   interface{}
	}{
		{
			expression: "object.foo > 1",
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": map[string]interface{}{
							"bar": "baz",
							"qux": "quux",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"bar": "baz",
				"qux": "quux",
			},
		},
		{
			expression: "size(object.foo)> 1",
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": map[string]interface{}{
							"bar": "baz",
							"qux": "quux",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"bar": "baz",
				"qux": "quux",
			},
		},
		{
			expression: "object.foo[0] > 1",
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": []interface{}{
							"bar",
							"baz",
						},
					},
				},
			},
			expected: "bar",
		},
	}

	for _, tc := range testCases {
		eval, err := NewEvaluator(tc.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		evalObj, err := eval.getEvaluatedObject(tc.expression)
		if err != nil {
			t.Errorf("Error getting evaluated object: %v", err)
		}
		if !helpers.CompareInterfaces(tc.expected, evalObj) {
			t.Errorf("Expected \n%v, got \n%v", tc.expected, evalObj)
		}
	}
}
