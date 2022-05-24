package provider

import (
	"archive/zip"
	"fmt"
	"github.com/bmatcuk/doublestar/v4"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func UnzipSource(source, pattern, outputDir string) ([]string, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if outputDir == "" {
		outputDir = "."
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, f := range reader.File {
		if pattern != "" {
			matches, err := doublestar.Match(pattern, f.Name)
			if err != nil {
				return nil, err
			}
			if !matches {
				continue
			}
		}
		if !f.FileInfo().IsDir() {
			ret = append(ret, f.Name)
		}
		err = UnzipFile(f, outputDir)
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(ret)

	return ret, nil
}

func UnzipFile(f *zip.File, dst string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(dst, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(dstFile, zippedFile); err != nil {
		return err
	}
	return nil
}
