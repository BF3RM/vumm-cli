package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type InstalledMod struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GetInstalledMod(name string) (*InstalledMod, error) {
	manifestPath := filepath.Join(GetModsPath(), name, "mod.json")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()

	var mod = new(InstalledMod)
	err = json.NewDecoder(manifestFile).Decode(mod)
	if err != nil {
		return nil, err
	}

	return mod, nil
}
