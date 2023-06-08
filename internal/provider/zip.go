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

func matches(filename string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matches, err := doublestar.Match(pattern, filename)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

func matchesWithExclude(filename string, patterns []string, excludes []string) (bool, error) {
	m, err := matches(filename, patterns)
	if err != nil {
		return false, err
	}
	if !m {
		return false, nil
	}
	m, err = matches(filename, excludes)
	if err != nil {
		return false, err
	}
	return !m, nil
}

func UnzipSource(source string, patterns []string, excludes []string, outputDir string) ([]string, error) {
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
		m, err := matchesWithExclude(f.Name, patterns, excludes)
		if err != nil {
			return nil, err
		}
		if !m {
			continue
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
	// Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(dst, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Create a destination file for unzipped content
	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Unzip the content of a file and copy it to the destination file
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
