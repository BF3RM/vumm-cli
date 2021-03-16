package project

import (
	"github.com/apex/log"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/pkg/ignorer"
	"os"
	"path/filepath"
)

type Project struct {
	Metadata  common.ModMetadata
	Directory string
	Ignorer   ignorer.FileIgnorer
}

func Load(projectPath string) (*Project, error) {
	log.WithField("file", "mod.json").Info("loading metadata")
	metadata, err := common.LoadModMetadata(filepath.Join(projectPath, "mod.json"))
	if err != nil {
		return nil, err
	}
	log.WithField("file", "mod.json").Debug("loaded metadata")

	log.WithField("file", ".vummignore").Info("loading ignore file")
	fileIgnorer, err := ignorer.CompileFileIgnorer(filepath.Join(projectPath, ".vummignore"))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		fileIgnorer = ignorer.NOOP()
		log.Debug("no .vummignore found")
	} else {
		log.WithField("file", ".vummignore").Debug("loaded ignore file")
	}

	return &Project{
		Metadata:  metadata,
		Directory: projectPath,
		Ignorer:   fileIgnorer,
	}, nil
}
