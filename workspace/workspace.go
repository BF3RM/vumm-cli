package workspace

import (
	"errors"
	"os"
	"path/filepath"
)

var workspaceRoot string

func SetWorkspacePath(workspacePath string) error {
	var err error
	workspaceRoot, err = filepath.Abs(workspacePath)
	if err != nil {
		return err
	}

	stat, err := os.Stat(workspaceRoot)
	if os.IsNotExist(err) {
		return err
	}

	if !stat.IsDir() {
		return errors.New("expected workspace to be a directory")
	}

	return nil
}

func GetModsPath() string {
	return filepath.Join(workspaceRoot, "Mods")
}

//func GetInstalledMods() ([]InstalledMod, error) {
//	config, err := GetConfig()
//	if err != nil {
//		return nil, err
//	}
//
//	mods := make([]*InstalledMod, 0, len(config.Mods))
//	for modName, version := range config.Mods {
//		installedMod, err := GetInstalledMod(modName)
//		if err != nil {
//			fmt.Printf("Failed reading installed mod %s: %v", modName, err)
//			continue
//		}
//		if version != installedMod.Version {
//			fmt.Printf("Installed mod %s is outdated, expected version %s, is %s", modName, version, installedMod.Version)
//		}
//		mods = append(mods, installedMod)
//	}
//	fmt.Println(mods)
//
//	return nil, nil
//}
