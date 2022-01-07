package main

import (
	"errors"
	"os"
	"path/filepath"
)

var www_root string
var gDMDevID string

func exePath() (string, error) {
	exe_path, err := os.Executable()
	if err != nil {
		return "", err
	}
	path, err := filepath.EvalSymlinks(filepath.Dir(exe_path))
	if err != nil {
		return "", err
	}
	if len(path) == 0 {
		return path, errors.New("directory of current executable file is empty")
	}

	return path, nil
}
