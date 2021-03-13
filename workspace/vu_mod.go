package workspace

import (
	"encoding/json"
	"github.com/vumm/cli/common"
	"os"
	"path/filepath"
)

func IsModInstalled(name string) bool {
	manifestPath := filepath.Join(GetModsPath(), name, "mod.json")
	_, err := os.Stat(manifestPath)
	if err != nil {
		return false
	}

	return true
}

func GetInstalledMod(name string) (*common.ModMetadata, error) {
	manifestPath := filepath.Join(GetModsPath(), name, "mod.json")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()

	var mod = new(common.ModMetadata)
	err = json.NewDecoder(manifestFile).Decode(mod)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
