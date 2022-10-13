package owl

import (
	"os"

	"gopkg.in/yaml.v2"
)

func saveToYaml(path string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func loadFromYaml(path string, data interface{}) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, data)
}
