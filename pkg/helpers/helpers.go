package helpers

import (
	"encoding/json"
	"regexp"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

func ExtractObject(input string) string {
	re := regexp.MustCompile(`object(?:\.[^.()\s+]+)*`)
	return re.FindString(input)
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
