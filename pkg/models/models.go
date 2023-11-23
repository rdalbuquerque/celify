package models

type ValidationRule struct {
	Expression        string `yaml:"expression"`
	MessageExpression string `yaml:"messageExpression"`
}

type ValidationConfig struct {
	Validations []ValidationRule `yaml:"validations"`
}

type TargetData struct {
	Data   map[string]interface{}
	Format string
}

type EvaluationResult struct {
	Expression       string
	EvaluatedObjects []EvaluatedObject
	ValidationError  error
}

type EvaluatedObject struct {
	Expression string
	Object     interface{}
}
