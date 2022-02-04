package utils

import (
	"path"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, file, _, _ := runtime.Caller(0) // nolint:dogsled
	fileDir := filepath.Dir(file)
	return path.Join(fileDir, "..")
}
