data "unarchive_file" "zip" {
  type        = "zip"
  source_file = "archive.zip"
  patterns    = ["**/*.{js,ts}"]
  excludes    = ["**/*.test.*"]
  output_dir  = "${path.root}/.terraform/tmp"
}
