package updater

import (
	"github.com/creativeprojects/go-selfupdate"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/config"
	"os"
	"path/filepath"
	"time"
)

const updateInterval = time.Hour * 24 * 7 // 1 week

func PeriodicCheckForUpdates() (*selfupdate.Release, bool, error) {
	if !shouldCheck() {
		return nil, false, nil
	}

	return CheckForUpdates()
}

func shouldCheck() bool {
	// Don't check for updates on development versions
	if !common.IsRelease() {
		return false
	}

	checkFile := filepath.Join(config.GetPath(), "last-update-check")

	nextCheck := time.Now().Add(-updateInterval)
	var lastChecked time.Time

	stat, err := os.Stat(checkFile)
	if err != nil {
		lastChecked = nextCheck.Add(-time.Millisecond)
	} else {
		lastChecked = stat.ModTime()
	}

	if nextCheck.After(lastChecked) {
		f, _ := os.Create(checkFile)
		f.Close()
		return true
	}

	return false
}
