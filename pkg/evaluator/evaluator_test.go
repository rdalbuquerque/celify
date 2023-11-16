package evaluator

import (
	"celify/pkg/helpers"
	"celify/pkg/models"
	"reflect"
	"testing"
)

var evalTests = []struct {
	targetData  *models.TargetData
	validations models.ValidationConfig
	expected    []models.EvaluationResult
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
		validations: models.ValidationConfig{
			Validations: []models.ValidationRule{
				{
					Expression: "object.foo == 'bar'",
				},
			},
		},
		expected: []models.EvaluationResult{
			{
				Expression:       "object.foo == 'bar'",
				ValidationResult: helpers.BoolPtr(true),
				EvaluatedObject:  "bar",
			},
		},
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
		validations: models.ValidationConfig{
			Validations: []models.ValidationRule{
				{
					Expression:        "object.foo == 'baz'",
					MessageExpression: "'foo should be baz but was ' + object.foo",
				},
			},
		},
		expected: []models.EvaluationResult{
			{
				Expression:        "object.foo == 'baz'",
				ValidationResult:  helpers.BoolPtr(false),
				EvaluatedObject:   "bar",
				FailedRule:        "object.foo == 'baz'",
				MessageExpression: "foo should be baz but was bar",
			},
		},
	},
}

func TestEvaluate(t *testing.T) {
	for _, tc := range evalTests {
		eval, err := NewEvaluator(tc.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		result := eval.Evaluate(tc.validations)
		if err != nil {
			t.Errorf("Error evaluating expression: %v", err)
		}
		// compare result to expected
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, result)
		}
	}
}

func TestEvaluateSingleExpression(t *testing.T) {
	for _, test := range evalTests {
		eval, err := NewEvaluator(test.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		result := eval.EvaluateRule(test.validations.Validations[0])
		if err != nil {
			t.Errorf("Error evaluating expression: %v", err)
		}
		if !reflect.DeepEqual(result, test.expected[0]) {
			t.Errorf("Expected %v, got %v", test.expected[0], result)
		}
	}
}

func TestExecuteEvaluation(t *testing.T) {
	var execEvalTests = []struct {
		targetData         *models.TargetData
		expression         string
		expectedReturnType reflect.Type
		errorExpected      bool
	}{
		{
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expression:         "object.foo == 'bar'",
			expectedReturnType: BoolType,
			errorExpected:      false,
		},
		{
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expression:         "string(object.foo)",
			expectedReturnType: StringType,
			errorExpected:      false,
		},
		{
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": []string{"bar", "baz"},
					},
				},
			},
			expression:         "object.foo",
			expectedReturnType: AnyType,
			errorExpected:      false,
		},
		{
			targetData: &models.TargetData{
				Data: map[string]interface{}{
					"object": map[string]interface{}{
						"foo": "bar",
					},
				},
			},
			expression:         "object",
			expectedReturnType: IntType,
			errorExpected:      true,
		},
	}
	for _, tc := range execEvalTests {
		eval, err := NewEvaluator(tc.targetData)
		if err != nil {
			t.Errorf("Error creating evaluator: %v", err)
		}
		_, err = eval.executeEvaluation(tc.expression, tc.expectedReturnType)
		if err != nil && !tc.errorExpected {
			t.Errorf("Error executing evaluation: %v", err)
		}
		if err == nil && tc.errorExpected {
			t.Errorf("Expected error but got none")
		}
	}
}
