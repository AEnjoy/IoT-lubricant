package file

import (
	"os"

	json "github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

func ReadFile(path string) (string, error) {
	content, err := ReadContentFile(path)
	return string(content), err
}

func ReadContentFile(filepath string) ([]byte, error) {
	return os.ReadFile(filepath)
}
func ReadYamlFile(filepath string, v any) error {
	content, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, v)
}
func ReadJsonFile(filepath string, v any) error {
	payload, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, v)
}
func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
