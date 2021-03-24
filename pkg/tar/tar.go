package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type FileFilter func(filePath string) bool

func tarballToWriter(src string, writer *tar.Writer, filter FileFilter) error {
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

		if filter != nil && !filter(header.Name) {
			return nil
		}

		err = writer.WriteHeader(header)
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

			_, err = io.CopyN(writer, file, info.Size())
			if err != nil && err != io.EOF {
				return fmt.Errorf("%s: copying contents: %v", path, err)
			}
		}
		return nil
	})
}

func untarFromReader(reader *tar.Reader, dest string) error {
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err = mkdir(filepath.Join(dest, header.Name))
		case tar.TypeReg:
			err = writeNewFile(filepath.Join(dest, header.Name), reader, header.FileInfo().Mode())
		case tar.TypeSymlink:
			err = writeNewSymbolicLink(filepath.Join(dest, header.Name), header.Linkname)
		default:
			return fmt.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func writeNewFile(filePath string, in io.Reader, fm os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", filePath, err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%s: creating new file: %v", filePath, err)
	}
	defer out.Close()

	err = out.Chmod(fm)
	if err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("%s: changing file mode: %v", filePath, err)
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%s: writing file: %v", filePath, err)
	}
	return nil
}

func writeNewSymbolicLink(filePath string, target string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", filePath, err)
	}

	err = os.Symlink(target, filePath)
	if err != nil {
		return fmt.Errorf("%s: making symbolic link for: %v", filePath, err)
	}

	return nil
}

func mkdir(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}
	return nil
}
