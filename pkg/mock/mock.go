package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func JSONModelFromFilePath(file string, intoStruct interface{}) error {
	jsonFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read from file: %w", err)
	}

	if err = json.Unmarshal(byteValue, &intoStruct); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func JSONStringFromFilePath(file string) (string, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", fmt.Errorf("failed to read from file: %w", err)
	}

	return string(byteValue), nil
}
