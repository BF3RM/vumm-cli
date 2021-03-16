package common

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"os"
	"strings"
)

var InternalMods = []string{"veniceext"}

func IsInternalMod(mod string) bool {
	for _, internalMod := range InternalMods {
		if strings.ToLower(mod) == internalMod {
			return true
		}
	}

	return false
}

// ModMetadata holds info of the mod.json
type ModMetadata struct {
	Name         string                        `json:"Name"`
	Version      *semver.Version               `json:"Version"`
	Dependencies map[string]*SemverConstraints `json:"Dependencies"`
}

func (m ModMetadata) String() string {
	return fmt.Sprintf("%s@%s", m.Name, m.Version)
}

func LoadModMetadata(metadataFile string) (ModMetadata, error) {
	file, err := os.Open(metadataFile)
	if err != nil {
		return ModMetadata{}, err
	}
	defer file.Close()

	var metadata ModMetadata
	err = json.NewDecoder(file).Decode(&metadata)
	if err != nil {
		return ModMetadata{}, err
	}

	// Normalize names (to lowercase)
	metadata.Name = strings.ToLower(metadata.Name)

	for dep, constraints := range metadata.Dependencies {
		normalizedDep := strings.ToLower(dep)

		if normalizedDep == dep {
			continue
		}
		delete(metadata.Dependencies, dep)

		metadata.Dependencies[normalizedDep] = constraints
	}

	return metadata, nil
}
