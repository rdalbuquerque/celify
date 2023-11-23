package helpers

import (
	"testing"
)

func TestExtractObject(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "object",
			expected: []string{"object"},
		},
		{
			input:    "object.foo",
			expected: []string{"object.foo"},
		},
		{
			input:    "object.foo.bar",
			expected: []string{"object.foo.bar"},
		},
		{
			input:    "size(object.foo.bar.baz) > 0",
			expected: []string{"object.foo.bar.baz"},
		},
		{
			input:    "has(object.foo.bar.baz[0])",
			expected: []string{"object.foo.bar.baz[0]"},
		},
		{
			input:    "object.foos.all(foo => foo.bar == 'baz')",
			expected: []string{"object.foos"},
		},
		{
			input:    "object.foos[1].map(n, n*n)",
			expected: []string{"object.foos[1]"},
		},
		{
			input:    "object.foos[1].map(n, n*n) && object.bar == 'baz'",
			expected: []string{"object.foos[1]", "object.bar"},
		},
		{
			input:    "object.foos[1].map(n, n*n) || object.bar == 'baz'",
			expected: []string{"object.foos[1]", "object.bar"},
		},
		{
			input:    "object.foos[1].map(n, n*n) || object.bar == 'baz' && size(object.foos) > 0",
			expected: []string{"object.foos[1]", "object.bar", "object.foos"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := ExtractObjects(tc.input)
			if !CompareInterfaces(actual, tc.expected) {
				t.Errorf("Expected '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}

func TestMarshalData(t *testing.T) {
	testCases := []struct {
		input    interface{}
		format   string
		expected string
	}{
		{
			input:    map[string]string{"foo": "bar"},
			format:   "yaml",
			expected: "foo: bar\n",
		},
		{
			input:    []map[string]string{{"foo": "bar"}, {"baz": "qux"}},
			format:   "yaml",
			expected: "- foo: bar\n- baz: qux\n",
		},
		{
			input:  map[string]string{"foo": "bar"},
			format: "json",
			expected: `{
  "foo": "bar"
}
`,
		},
		{
			input:  []map[string]string{{"foo": "bar"}, {"baz": "qux"}},
			format: "json",
			expected: `[
  {
    "foo": "bar"
  },
  {
    "baz": "qux"
  }
]
`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.format, func(t *testing.T) {
			actual, err := MarshalData(tc.input, tc.format)
			if err != nil {
				t.Errorf("Error marshalling data: %v", err)
			}
			if string(actual) != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, actual)
			}
		})
	}
}

func TestUnmarshalData(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{
			input:    `foo: bar`,
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			input:    `{"foo": "bar"}`,
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			input: `---
foo:
  meta: bar
  bar:
  - baz
  - qux
---
`,
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"meta": "bar",
					"bar": []interface{}{
						"baz",
						"qux",
					},
				},
			},
		},
		{
			input: `{
  "foo": {
	"meta": "bar",
	"bar": [
	  "baz",
	  "qux"
	]
  }
}`,
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"meta": "bar",
					"bar": []interface{}{
						"baz",
						"qux",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var actual map[string]interface{}
			_, err := UnmarshalData([]byte(tc.input), &actual)
			if err != nil {
				t.Errorf("Error unmarshalling data: %v", err)
			}
			if !CompareInterfaces(actual, tc.expected) {
				t.Errorf("Expected '%v', got '%v'", tc.expected, actual)
			}
		})
	}
}
