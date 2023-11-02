package models

type ValidationRule struct {
	Expression   string `yaml:"expression"`
	ErrorMessage string `yaml:"errorMessage"`
}

type ValidationConfig struct {
	Validations []ValidationRule `yaml:"validations"`
}
