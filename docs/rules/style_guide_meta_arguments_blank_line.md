## style_guide_meta_arguments_blank_line

Warns when meta-arguments are not separated from other arguments by a blank line.

- Leading meta-arguments must be followed by a blank line
  - `count`, `for_each`, `source`, `provider`, `providers`
- Trailing meta-arguments must be preceded by a blank line
  - `lifecycle`, `connection`, `provisioner`, `depends_on`
- Leading meta-arguments should not be the last item in a block (with exceptions)
  - `count` and `for_each` indicate incomplete configuration if they are the last item
  - `source`, `provider` and `providers` can validly be the last item in module blocks

### Example

#### Missing blank lines

```hcl
resource "aws_instance" "web" {
  count         = 2
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  lifecycle {
    prevent_destroy = true
  }
}
```

```
$ tflint
2 issue(s) found:

Warning: Meta argument should be followed by a blank line (style_guide_meta_arguments_blank_line)

  on test.tf line 2:
   2:   count         = 2

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.1/docs/rules/style_guide_meta_arguments_blank_line.md

Warning: Meta argument should be preceded by a blank line (style_guide_meta_arguments_blank_line)

  on test.tf line 5:
   5:   lifecycle {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.1/docs/rules/style_guide_meta_arguments_blank_line.md
```

#### Leading meta argument as last item (incomplete configuration)

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  count = 2
}
```

```
$ tflint
1 issue(s) found:

Warning: Leading meta argument should not be the last item in the block (style_guide_meta_arguments_blank_line)

  on test.tf line 5:
   5:   count = 2

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.1/docs/rules/style_guide_meta_arguments_blank_line.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#code-formatting

> For blocks that contain both arguments and "meta-arguments" (as defined by the Terraform language semantics), list meta-arguments first and separate them from other arguments with one blank line. Place meta-argument blocks last and separate them from other blocks with one blank line.

### How To Fix

Properly add blank lines:

Example:

```hcl
resource "aws_instance" "web" {
  count         = 2
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
  lifecycle {
    prevent_destroy = true
  }
}
```

Change this to:

```hcl
resource "aws_instance" "web" {
  count = 2

  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  lifecycle {
    prevent_destroy = true
  }
}
```