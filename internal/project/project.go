package project

import (
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/pkg/ignorer"
	"os"
	"path/filepath"
)

var defaultIgnoresLines = []string{
	// Git
	".git",
	".gitignore",
	".github/",

	// Node
	"node_modules/",
	"package.json",
	"package-lock.json",
	"yarn.lock",

	// Editors
	".vscode/",
	".idea/",
	".editorconfig",
	"*.iml",

	// VU (ui folders, these should be compiled to a ui.vuic anyways)
	"^[Uu][Ii]/",
	"^[Ww]eb[Uu][Ii]/",
	".vummignore", // yea lets also remove ourselves
}

type Project struct {
	Metadata  common.ModMetadata
	Directory string
	Ignorer   ignorer.Ignorer
}

func Load(projectPath string) (*Project, error) {
	log.WithField("file", "mod.json").Info("loading metadata")
	metadata, err := common.LoadModMetadata(filepath.Join(projectPath, "mod.json"))
	if err != nil {
		return nil, err
	}
	log.WithField("file", "mod.json").Debug("loaded metadata")

	log.WithField("file", ".vummignore").Info("loading ignore file")
	fileIgnorer, err := ignorer.CompileIgnorerFromFile(filepath.Join(projectPath, ".vummignore"), defaultIgnoresLines...)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		fileIgnorer = ignorer.CompileIgnorerFromLines(defaultIgnoresLines...)
		log.Debug("no .vummignore found, using default")
	} else {
		log.WithField("file", ".vummignore").Debug("loaded ignore file")
	}

	return &Project{
		Metadata:  metadata,
		Directory: projectPath,
		Ignorer:   fileIgnorer,
	}, nil
}
