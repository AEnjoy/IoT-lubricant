package file

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile(path string) (string, error) {
	fs, err := os.Open(path)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(fs)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ReadContentFile(filepath string) ([]byte, error) {
	fd, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	payload, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
func ReadYamlFile(filepath string, v any) error {
	content, err := ReadContentFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, v)
}
func ReadJsonFile(filepath string, v any) error {
	fd, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	payload, err := io.ReadAll(fd)
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
