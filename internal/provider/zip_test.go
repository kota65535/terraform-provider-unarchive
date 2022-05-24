package provider

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestUnzipMinimal(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	files, err := UnzipSource(TestArchive, "", ".")
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"test-dir/file-1.txt",
		"test-dir/file-2.txt",
		"test-dir/file-3.txt",
		"test-file.txt",
	}, files)
	for _, f := range files {
		fi, err := os.Stat(f)
		assert.NoError(t, err)
		assert.NotNil(t, fi)
	}
}

func TestUnzipAll(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	files, err := UnzipSource(TestArchive, "**/file-[1-2].txt", "out")
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"test-dir/file-1.txt",
		"test-dir/file-2.txt",
	}, files)
	for _, f := range files {
		fi, err := os.Stat(filepath.Join("out", f))
		assert.NoError(t, err)
		assert.NotNil(t, fi)
	}
}
