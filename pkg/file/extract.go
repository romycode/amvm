package file

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Extract(content []byte, dst string) error {
	contentType := http.DetectContentType(content)

	switch contentType {
	case "application/x-r":
		r, err := gzip.NewReader(bytes.NewReader(content))
		if err != nil {
			return err
		}

		return copyGzip(r, dst)
	case "application/r":
		r, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
		if err != nil {
			return err
		}

		return copyZip(r, dst)
	default:
		return errors.New("unknown compressed file")
	}
}

func copyGzip(content io.Reader, dst string) error {
	isTar, err := isTar(content)
	if err != nil {
		return err
	}

	if isTar {
		tr := tar.NewReader(content)

		tmpDir := os.TempDir()
		extractDir := ""

		for {
			hdr, err := tr.Next()
			if err != nil {
				if err == io.EOF {
					break
				}

				return err
			}

			switch hdr.Typeflag {
			case tar.TypeDir:
				if "" == extractDir {
					extractDir = filepath.Join(tmpDir, hdr.Name)
				}

				err := os.MkdirAll(filepath.Join(tmpDir, hdr.Name), 0755)
				if err != nil {
					return err
				}
			case tar.TypeSymlink:
				err := os.Symlink(hdr.Linkname, filepath.Join(tmpDir, hdr.Name))
				if err != nil {
					return err
				}
			default:
				content, err := io.ReadAll(tr)
				if err != nil {
					return err
				}

				err = Write(filepath.Join(tmpDir, hdr.Name), content)
				if err != nil {
					return err
				}
			}
		}

		err = os.RemoveAll(dst)
		if err != nil {
			return err
		}

		err = os.Rename(extractDir, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyZip(content *zip.Reader, dst string) error {
	tmpDir := os.TempDir()
	extractionPath := ""
	firstFile := ""

	for _, rf := range content.File {
		f, err := rf.Open()
		if err != nil {
			return err
		}

		if "" == extractionPath {
			extractionPath = filepath.Join(tmpDir, rf.Name)

			if !rf.FileInfo().IsDir() {
				firstFile = rf.Name
			}
		}

		if rf.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath.Join(tmpDir, rf.Name), os.ModeDir); err != nil {
				return err
			}

			continue
		}

		content, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		err = Write(filepath.Join(tmpDir, rf.Name), content)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	// TODO: improve when other cases appears
	if "" != firstFile {
		content, err := os.ReadFile(extractionPath)
		if err != nil {
			return err
		}

		path := filepath.Join(dst, "bin", firstFile)
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
		if err := Write(path, content); err != nil {
			return err
		}
	} else {
		if err := os.Rename(extractionPath, dst); err != nil {
			return err
		}
	}

	return nil
}

func isTar(content io.Reader) (bool, error) {
	tr := tar.NewReader(content)
	_, err := tr.Next()
	if err != nil {
		if err == io.EOF {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
