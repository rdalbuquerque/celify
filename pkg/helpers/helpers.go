package helpers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

func ExtractObjects(input string) []string {
	re := regexp.MustCompile(`object(?:\.[^\s+\)]+)*`)
	matches := re.FindAllString(input, -1)
	objects := []string{}
	for _, match := range matches {
		if match != "" {
			parts := strings.Split(match, ".")
			lastPart := parts[len(parts)-1]
			if strings.ContainsAny(lastPart, "(") {
				match = strings.Join(parts[:len(parts)-1], ".")
			}

		}
		objects = append(objects, match)
	}
	return objects
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
		jsonStr, err := json.MarshalIndent(target, "", "  ")
		if err != nil {
			return nil, err
		}
		jsonStrln := fmt.Sprintf("%s\n", jsonStr)
		return []byte(jsonStrln), nil
	} else {
		return nil, errors.Errorf("Invalid format '%s' provided", format)
	}
}

func BoolPtr(b bool) *bool {
	return &b
}
