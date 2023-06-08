package provider

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMatchesExclude(t *testing.T) {
	m, err := matchesWithExclude("foo/bar.txt", []string{"**/*"}, []string{"*.txt"})
	assert.NoError(t, err)
	assert.True(t, m)
	m, err = matchesWithExclude("foo/bar.txt", []string{"**/*"}, []string{"**/bar.txt"})
	assert.NoError(t, err)
	assert.False(t, m)
	m, err = matchesWithExclude("foo/bar.txt", []string{"*.txt"}, []string{"**/bar.txt"})
	assert.NoError(t, err)
	assert.False(t, m)
	m, err = matchesWithExclude("foo/bar.txt", []string{}, []string{"**/bar.txt"})
	assert.NoError(t, err)
	assert.False(t, m)
}

func TestUnzipMinimal(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	files, err := UnzipSource(TestArchive, []string{"**/*"}, []string{}, ".")
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

	files, err := UnzipSource(TestArchive, []string{"**/file-[1-2].txt"}, []string{}, "out")
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
