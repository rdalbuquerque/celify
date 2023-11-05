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

func ExtractObject(input string) string {
	re := regexp.MustCompile(`object(?:\.[^.()\s+]+)*`)
	return re.FindString(input)
}

func GetErrorStr() string {
	return color.New(color.FgRed).Sprint("|")
}

func PrintMultilineError(input string, color *color.Color) {
	errLines := strings.Split(input, "\n")
	for _, line := range errLines {
		fmt.Printf("%s %s\n", GetErrorStr(), color.Sprint(line))
	}
}

func PrintEvaluatedObject(objStr, format string) {
	fmt.Printf("%s %s", GetErrorStr(), color.New(color.Underline).Sprintln("Evaluated object:"))
	if format == "yaml" {
		c := color.New(color.FgBlue)
		PrintMultilineError(objStr, c)
	} else if format == "json" {
		PrintMultilineError(string(pretty.Color(pretty.Pretty([]byte(objStr)), nil)), color.New(color.Reset))
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
