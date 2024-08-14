package utils

import (
	"fmt"
	"os"
)

// GetABI reads the ABI file and returns its content.
func GetABI() ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %v", err)
	}
	// Convert the relative path to an absolute path
	relativePath := cwd + "/ABI/bamla.json"

	abiFile, err := os.ReadFile(relativePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %v", err)
	}

	return abiFile, nil
}
