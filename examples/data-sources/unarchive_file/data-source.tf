data "unarchive_file" "zip" {
  type        = "zip"
  source_file = "a.zip"
  pattern     = "**/*.txt"
  output_dir  = "${path.root}/.terraform/tmp"
}
