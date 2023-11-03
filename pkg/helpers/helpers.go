package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/tidwall/pretty"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"
)

func ExtractObject(input string) string {
	re := regexp.MustCompile(`object(?:\.[^.()\s+]+)*`)
	return re.FindString(input)
}

func GetErrorStr(input string) string {
	return fmt.Sprintf("%s|%s %s", colorRed, colorReset, input)
}

func GetMultilineErrorStr(input string) string {
	errLines := strings.Split(input, "\n")
	for i, line := range errLines {
		errLines[i] = GetErrorStr(line)
	}
	return strings.Join(errLines, "")
}

func PrintMultilineError(input string) {
	fmt.Print(GetMultilineErrorStr(input))
}

func PrintEvaluatedObject(objStr, format string) {
	fmt.Println(GetErrorStr("Evaluated object:"))
	if format == "yaml" {
		c := color.New(color.FgGreen)
		c.Println(objStr)
	} else if format == "json" {
		fmt.Println(string(pretty.Color(pretty.Pretty([]byte(objStr)), nil)))
	}
}

func UnmarshalData(data []byte, target interface{}) (string, error) {
	if err := json.Unmarshal(data, target); err != nil {
		if err := yaml.Unmarshal(data, target); err != nil {
			return "", errors.Errorf("Error unmarshalling target: %v", err)
		}
		return "yaml", nil
	}
	return "json", nil
}

// MarshalData marshals the target into YAML or JSON, returning the formatted bytes, and the format used ('yaml' or 'json')
func MarshalData(target interface{}, format string) ([]byte, error) {
	if format == "yaml" {
		return yaml.Marshal(target)
	} else if format == "json" {
		return json.Marshal(target)
	} else {
		return nil, errors.Errorf("Invalid format '%s' provided", format)
	}
}
