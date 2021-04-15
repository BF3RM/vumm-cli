package workspace

import (
	"bufio"
	"fmt"
	"github.com/apex/log"
	"os"
	"path/filepath"
	"strings"
)

type ModList struct {
	dir   string
	lines []string
}

func TryLoadModList(dir string) (*ModList, error) {
	modList := &ModList{
		dir:   dir,
		lines: []string{},
	}

	filePath := filepath.Join(dir, "ModList.txt")
	log.WithField("file", filePath).Debugf("loading mod list")
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warn("mod list not found, creating")
			return modList, nil
		}

		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		modList.lines = append(modList.lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return modList, nil
}

func (m *ModList) EnableMod(mod string) bool {
	mod = strings.ToLower(mod)

	enabled, idx := m.isModEnabled(mod)

	// not found, add new line
	if idx == -1 {
		m.lines = append(m.lines, mod)
		return true
	}

	m.lines[idx] = mod

	// disabled
	return !enabled
}

func (m *ModList) DisableMod(mod string) bool {
	mod = strings.ToLower(mod)

	enabled, idx := m.isModEnabled(mod)

	// not found, nothing to do
	if idx == -1 {
		return false
	}

	// enabled, lets disable
	if enabled {
		m.lines[idx] = "#" + mod
		return true
	}

	// Replace anyway to remove leftover spaces
	m.lines[idx] = mod
	return false
}

func (m ModList) isModEnabled(mod string) (bool, int) {
	for idx, line := range m.lines {
		parsedLine := strings.TrimLeft(strings.ToLower(line), " ")
		// already enabled
		if parsedLine == mod {
			return true, idx
		}

		// disabled, lets enable
		if strings.HasPrefix(parsedLine, "#") && strings.TrimLeft(parsedLine, "# ") == mod {
			return false, idx
		}
	}

	return false, -1
}

func (m *ModList) Save() error {
	log.WithField("file", "ModList.txt").Debugf("saving ModList.txt")

	file, err := os.OpenFile(filepath.Join(m.dir, "ModList.txt"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0x666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range m.lines {
		if _, err = writer.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
	}
	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}
