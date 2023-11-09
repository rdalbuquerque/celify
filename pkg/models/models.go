package models

type ValidationRule struct {
	Expression   string `yaml:"expression"`
	ErrorMessage string `yaml:"errorMessage"`
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
