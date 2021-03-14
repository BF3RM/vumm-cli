package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

type Packager interface {
	SetFileFilter(filter FileFilter)
	Compress(src string, writer io.Writer) error
	Decompress(reader io.Reader, dest string) error
}

func NewPackager() Packager {
	return &tarGzPackager{}
}

type tarGzPackager struct {
	filter FileFilter
}

func (p *tarGzPackager) SetFileFilter(filter FileFilter) {
	p.filter = filter
}

func (p tarGzPackager) Compress(src string, writer io.Writer) error {
	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return tarballToWriter(src, tarWriter, p.filter)
}

func (p tarGzPackager) Decompress(reader io.Reader, dest string) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to decompress archive: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	return untarFromReader(tarReader, dest)
}
