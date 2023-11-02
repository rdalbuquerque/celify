package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

func ExtractObject(input string) string {
	re := regexp.MustCompile(`object(?:\.[^.()\s+]+)*`)
	return re.FindString(input)
}

func GetErrorStr(input string) string {
	return fmt.Sprintf("%s|%s %s", colorRed, colorReset, input)
}

func getSeparatorStr() string {
	return fmt.Sprintf("%s|%s------------------------------------------%s", colorRed, colorYellow, colorReset)
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

func PrintEvaluatedObject(objStr string) {
	fmt.Print(GetErrorStr("Evaluated object:"))
	fmt.Printf("%s\n%s%s\n", colorBlue, objStr, colorReset)
}

func UnmarshalData(data []byte, target interface{}) error {
	if err := yaml.Unmarshal(data, target); err != nil {
		if err := json.Unmarshal(data, target); err != nil {
			return errors.Errorf("Error unmarshalling target: %v", err)
		}
	}
	return nil
}

func MarshalData(target interface{}) ([]byte, error) {
	yaml, err := yaml.Marshal(target)
	if err != nil {
		json, err := json.Marshal(target)
		if err != nil {
			return nil, errors.Errorf("Error marshalling target: %v", err)
		}
		return json, nil
	}
	return yaml, nil
}
