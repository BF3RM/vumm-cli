package updater

import (
	"errors"
	"fmt"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/vumm/cli/internal/common"
	"log"
	"os"
	"runtime"
)

var latestVersion *selfupdate.Release
var updater *selfupdate.Updater

func init() {
	var err error
	updater, err = selfupdate.NewUpdater(selfupdate.Config{Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"}})
	if err != nil {
		panic(err)
	}
}

func CheckForUpdates() (*selfupdate.Release, bool, error) {
	latest, found, err := updater.DetectLatest("BF3RM/vumm-cli")
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch latest version: %v", err)
	}

	if !found {
		return nil, false, fmt.Errorf("failed to find latest version for %s", runtime.GOOS)
	}

	latestVersion = latest
	return latest, IsUpdateAvailable(), nil
}

func IsUpdateAvailable() bool {
	if latestVersion == nil || !common.IsRelease() {
		return false
	}

	return latestVersion.GreaterThan(common.GetVersion())
}

func SelfUpdate() (bool, error) {
	if latestVersion == nil {
		_, _, err := CheckForUpdates()
		if err != nil {
			return false, err
		}
	}

	if !IsUpdateAvailable() {
		return false, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return false, errors.New("could not locate executable path")
	}
	if err := updater.UpdateTo(latestVersion, exe); err != nil {
		return false, fmt.Errorf("error occurred while updating binary: %v", err)
	}
	log.Printf("Successfully updated to version %s", latestVersion.Version())
	return false, nil
}
