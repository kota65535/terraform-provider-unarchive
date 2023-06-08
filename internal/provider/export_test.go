package provider

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const TestArchive = "test-archive.zip"

func setup(t *testing.T) (string, string) {
	td := createTempDir(t)
	copyFile(TestArchive, filepath.Join(td, TestArchive))
	cwd, _ := os.Getwd()
	fmt.Printf("temp dir: %s\n", td)
	os.Chdir(td)
	return td, cwd
}

func tearDown(t *testing.T, td, cwd string) {
	os.RemoveAll(td)
	os.Chdir(cwd)
}

func createTempDir(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	return tmp
}

func copyFile(src, dst string) {
	fin, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		log.Fatal(err)
	}
}
