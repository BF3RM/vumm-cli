package workspace

import (
	"encoding/json"
	"github.com/Masterminds/semver"
	"os"
	"path/filepath"
)

type WorkspaceConfig struct {
	file *os.File

	Mods map[string]*semver.Constraints `json:"mods"`
}

var loadedConfig *WorkspaceConfig

// GetConfig will either load or create a new config if it does not exist yet
func GetConfig() (*WorkspaceConfig, error) {
	if loadedConfig != nil {
		return loadedConfig, nil
	}

	configFilePath := filepath.Join(workspaceRoot, "vumm.json")
	configFile, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the config
			configFile, err := os.Create(configFilePath)
			if err != nil {
				return nil, err
			}

			loadedConfig = &WorkspaceConfig{
				file: configFile,
				Mods: map[string]*semver.Constraints{},
			}
			loadedConfig.Save()

			return loadedConfig, nil
		}
		return nil, err
	}

	loadedConfig = &WorkspaceConfig{
		file: configFile,
	}
	err = json.NewDecoder(configFile).Decode(loadedConfig)
	if err != nil {
		return nil, err
	}

	return loadedConfig, nil
}

// AddMod adds a new mod to the config
func (c *WorkspaceConfig) AddMod(mod string, version *semver.Constraints) {
	c.Mods[mod] = version
}

// RemoveMod removes a mod from the config
func (c *WorkspaceConfig) RemoveMod(mod string) {
	delete(c.Mods, mod)
}

// Save saves the config
func (c *WorkspaceConfig) Save() error {
	encoder := json.NewEncoder(c.file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(c)
}

// Close saves and closes the config file
func (c *WorkspaceConfig) Close() error {
	// Always unload and close the file no matter what
	defer func() { loadedConfig = nil }()
	defer c.Close()

	return c.Save()
}
