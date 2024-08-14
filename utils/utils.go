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
	fmt.Println("Current working directory:", cwd)

	abiFile, err := os.ReadFile("../ABI/bamla.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %v", err)
	}

	return abiFile, nil
}
