package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var cachedPath string

func init() {
	path := GetPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(fmt.Errorf("failed creating config path: %v", err))
		}
	}
}

func GetPath() string {
	if cachedPath != "" {
		return cachedPath
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed resolving config path: %v", err))
	}

	return filepath.Join(homedir, ".vumm")
}
