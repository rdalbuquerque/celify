package evaluator

import (
	"celify/pkg/models"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
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

func TestPrintEvaluatedObject(t *testing.T) {
	old := os.Stdout
	defer func() { os.Stdout = old }()

	r, w, _ := os.Pipe()
	os.Stdout = w

	eval, err := NewEvaluator(&models.TargetData{
		Data: map[string]interface{}{
			"object": map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
					"qux": "quux",
				},
			},
		},
		Format: "yaml",
	})
	if err != nil {
		t.Errorf("Error creating evaluator: %v", err)
	}

	eval.printEvaluatedObject("object.foo", "yaml")

	w.Close()

	out, _ := io.ReadAll(r)

	expected := `| Evaluated object:
| bar: baz
| qux: quux
| 
`

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected, string(out), false)

	os.Stdout = old
	fmt.Println(dmp.DiffPrettyText(diffs))
	if string(out) != expected {
		t.Errorf("Expected \n%v, got \n%v", expected, string(out))
	}
}

// stripANSI removes ANSI color/style codes from the string
func stripANSI(str string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(str, "")
}
