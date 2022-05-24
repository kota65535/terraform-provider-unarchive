package provider

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const TestArchive = "commons-text-1.9.jar"

func TestAccDataSourceUnarchiveFile(t *testing.T) {
	td := createTempDir(t)
	defer os.RemoveAll(td)
	copyFile(TestArchive, filepath.Join(td, TestArchive))
	cwd, _ := os.Getwd()
	os.Chdir(td)
	defer os.Chdir(cwd)

	filenames := []string{
		"META-INF/LICENSE.txt",
		"META-INF/NOTICE.txt",
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUnarchiveFileMinimalConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccExtractedFilesExists(filenames, "."),
					resource.TestCheckResourceAttr("data.unarchive_file.basic", "output_files.#", "140"),
				),
			},
			{
				Config: testAccDataSourceUnarchiveFileAllConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccExtractedFilesExists(filenames, "out"),
					resource.TestCheckResourceAttr("data.unarchive_file.basic", "output_files.#", "2"),
					resource.TestCheckResourceAttr("data.unarchive_file.basic", "output_files.0", filenames[0]),
					resource.TestCheckResourceAttr("data.unarchive_file.basic", "output_files.1", filenames[1]),
				),
			},
		},
	})
}

func testAccDataSourceUnarchiveFileMinimalConfig() string {
	return fmt.Sprintf(`
	data "unarchive_file" "basic" {
		type        = "zip"
        source_file = "commons-text-1.9.jar"
	}
	`)
}

func testAccDataSourceUnarchiveFileAllConfig() string {
	return fmt.Sprintf(`
	data "unarchive_file" "basic" {
		type        = "zip"
        source_file = "commons-text-1.9.jar"
		pattern     = "**/*.txt"
        output_dir  = "out"
	}
	`)
}

func testAccExtractedFilesExists(filenames []string, dir string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, f := range filenames {
			fi, err := os.Stat(filepath.Join(dir, f))
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return errors.New("not file")
			}
		}
		return nil
	}
}

func createTempDir(t *testing.T) string {
	tmp, err := ioutil.TempDir("", "tf")
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
