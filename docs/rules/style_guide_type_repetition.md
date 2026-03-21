## style_guide_type_repetition

Warns when a resource or data source name repeats a word in its type.

### Example

```hcl
resource "aws_iam_role" "lambda_role" {
  # "role" is already in the type "aws_iam_role"
}

data "aws_s3_bucket" "app_bucket" {
  # "bucket" is already in the type "aws_s3_bucket"
}
```

```
$ tflint
2 issue(s) found:

Warning: Resource name should not repeat its resource type (style_guide_type_repetition)

  on test.tf line 1:
   1: resource "aws_iam_role" "lambda_role" {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.0/docs/rules/style_guide_type_repetition.md

Warning: Data source name should not repeat its data source type (style_guide_type_repetition)

  on test.tf line 5:
   5: data "aws_s3_bucket" "app_bucket" {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.0/docs/rules/style_guide_type_repetition.md
```

### Why

The resource address already includes the type, so repeating it in the name is redundant.

https://developer.hashicorp.com/terraform/language/style#resource-naming

> Do not include the resource type in the resource identifier since the resource address already includes it.

### How to fix

Example:

```hcl
resource "aws_iam_role" "lambda_role" {
  # "role" repeats the type "aws_iam_role"
  assume_role_policy = "{}"
}

data "aws_s3_bucket" "app_bucket" {
  # "bucket" repeats the type "aws_s3_bucket"
  bucket = "my-app-bucket"
}
```

Change this to:

```hcl
resource "aws_iam_role" "lambda" {
  assume_role_policy = "{}"
}

data "aws_s3_bucket" "app" {
  bucket = "my-app-bucket"
}
```

See https://developer.hashicorp.com/terraform/language/style#resource-naming for better naming.