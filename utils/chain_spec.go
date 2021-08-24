package utils

import (
	"os"
	"path/filepath"
)

func StoreChainSpec(filePath, chainSpec string) error {
	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(chainSpec), 0600)
}
