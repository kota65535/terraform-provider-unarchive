package provider

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceUnarchiveFile(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	filenames := []string{
		"test-dir/file-1.txt",
		"test-dir/file-2.txt",
		"test-dir/file-3.txt",
		"test-file.txt",
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUnarchiveFileMinimalConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccExtractedFilesExists(filenames, "."),
					resource.TestCheckResourceAttr("data.unarchive_file.minimal", "output_files.#", "4"),
				),
			},
			{
				Config: testAccDataSourceUnarchiveFileWithPattern(),
				Check: resource.ComposeTestCheckFunc(
					testAccExtractedFilesExists(filenames[:2], "out"),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.#", "2"),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.0.name", filenames[0]),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.1.name", filenames[1]),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.0.path", path.Join("out", filenames[0])),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.1.path", path.Join("out", filenames[1])),
				),
			},
			{
				Config: testAccDataSourceUnarchiveFileWithAllConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccExtractedFilesExists(filenames[:2], "out"),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.#", "2"),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.0.name", filenames[1]),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.1.name", filenames[2]),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.0.path", path.Join("out", filenames[1])),
					resource.TestCheckResourceAttr("data.unarchive_file.all", "output_files.1.path", path.Join("out", filenames[2])),
				),
			},
		},
	})
}

func testAccDataSourceUnarchiveFileMinimalConfig() string {
	return fmt.Sprintf(`
	data "unarchive_file" "minimal" {
		type        = "zip"
        source_file = "test-archive.zip"
	}
	`)
}

func testAccDataSourceUnarchiveFileWithPattern() string {
	return fmt.Sprintf(`
	data "unarchive_file" "all" {
		type        = "zip"
        source_file = "test-archive.zip"
		pattern     = "**/file-[1-2].txt"
        output_dir  = "out"
	}
	`)
}

func testAccDataSourceUnarchiveFileWithAllConfig() string {
	return fmt.Sprintf(`
	data "unarchive_file" "all" {
		type        = "zip"
        source_file = "test-archive.zip"
		patterns     = ["**/file-[1-3].txt"]
		excludes     = ["**/file-1.txt"]
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
