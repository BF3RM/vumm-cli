package publish

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Packager interface {
	SetIgnorer(ignorer *FileIgnorer)
	Make(src string, writer io.Writer) error
}

type tarGzPackager struct {
	ignorer *FileIgnorer
}

func NewTarGZPackager() Packager {
	return &tarGzPackager{}
}

func (f *tarGzPackager) SetIgnorer(ignorer *FileIgnorer) {
	f.ignorer = ignorer
}

func (f tarGzPackager) Make(src string, writer io.Writer) error {
	if cw, ok := writer.(io.WriteCloser); ok {
		defer cw.Close()
	}

	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return f.tarball(src, tarWriter)
}

func (f tarGzPackager) tarball(src string, tarWriter *tar.Writer) error {
	//srcInfo, err := os.Stat(src)
	//if err != nil {
	//	return fmt.Errorf("%s: stat: %v", src, err)
	//}

	//var baseDir string
	//if srcInfo.IsDir() {
	//	baseDir = filepath.Base(src)
	//}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking to %s: %v", path, err)
		}

		// Ignore root directory
		if path == src {
			return nil
		}

		if f.ignorer != nil && f.ignorer.Matches(path) {
			return nil
		}

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return fmt.Errorf("%s: making header: %v", path, err)
		}

		// Rewrite baseDir
		//if baseDir != "" {
		header.Name = strings.TrimPrefix(path, src+string(filepath.Separator))
		//}

		if info.IsDir() {
			header.Name += "/"
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("%s: writing header: %v", path, err)
		}

		// We done, nothing left to do for directories
		if info.IsDir() {
			return nil
		}

		if header.Typeflag == tar.TypeReg {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("%s: open: %v", path, err)
			}
			defer file.Close()

			_, err = io.CopyN(tarWriter, file, info.Size())
			if err != nil && err != io.EOF {
				return fmt.Errorf("%s: copying contents: %v", path, err)
			}
		}
		return nil
	})
}
