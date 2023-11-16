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
	Expression        string
	ValidationResult  *bool
	EvaluatedObject   interface{}
	FailedRule        string
	MessageExpression string
	ValidationError   error
}
