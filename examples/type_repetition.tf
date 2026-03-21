resource "aws_iam_role" "lambda_role" {
  # "role" repeats the type "aws_iam_role"
  assume_role_policy = "{}"
}

data "aws_s3_bucket" "app_bucket" {
  # "bucket" repeats the type "aws_s3_bucket"
  bucket = "my-app-bucket"
}
