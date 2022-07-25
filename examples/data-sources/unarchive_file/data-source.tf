data "unarchive_file" "zip" {
  type        = "zip"
  source_file = "archive.zip"
  pattern     = "**/*.txt"
  output_dir  = "${path.root}/.terraform/tmp"
}
