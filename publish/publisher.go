package publish

import (
	"bytes"
	"fmt"
	"github.com/vumm/cli/common"
	"github.com/vumm/cli/registry"
	"github.com/vumm/cli/tar"
	"os"
	"path/filepath"
)

type Publisher struct {
	cwd      string
	metadata common.ModMetadata
	packager tar.Packager
}

func NewPublisher() (Publisher, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Publisher{}, err
	}

	publisher := Publisher{
		cwd: cwd,
	}

	if err = publisher.loadMetadata(); err != nil {
		return Publisher{}, fmt.Errorf("failed loading metadata: %v", err)
	}

	if err = publisher.loadPackager(); err != nil {
		return Publisher{}, fmt.Errorf("failed loading packager: %v", err)
	}

	return publisher, nil
}

func (p *Publisher) Publish() error {
	fmt.Println("Compressing mod to archive")

	var buf bytes.Buffer
	err := p.packager.Compress(p.cwd, &buf)
	if err != nil {
		return err
	}
	fmt.Printf("Compressed mod to archive of %d bytes\n", buf.Len())

	err = registry.PublishMod(p.metadata, "latest", &buf)
	if err != nil {
		return err
	}
	fmt.Printf("Published %s\n", p.metadata)
	return nil
}

func (p *Publisher) loadMetadata() (err error) {
	p.metadata, err = common.LoadModMetadata(filepath.Join(p.cwd, "mod.json"))
	return err
}

func (p *Publisher) loadPackager() error {
	p.packager = tar.NewPackager()

	ignorer, err := CompileFileIgnorer(p.cwd)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	p.packager.SetFileFilter(func(filePath string) bool {
		return !ignorer.Matches(filePath)
	})

	return nil
}
