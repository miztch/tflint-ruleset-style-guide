## style_guide_block_placement

Warns when `terraform` or `provider` blocks are placed outside their conventional files.

- `provider` blocks must be defined in `providers.tf`.
- `terraform` blocks must be defined in `terraform.tf`, unless they contain a nested `backend` block, in which case they must be defined in `backend.tf`.
  - A `terraform` block that mixes `required_version` / `required_providers` with a `backend` block is treated as a backend configuration block, since it can only satisfy one file name — split it into two `terraform` blocks (one per file) for full compliance.
- **Disabled by default.** The de facto community standard is `versions.tf` for the `terraform` block (used by terraform-aws-modules and most popular modules; the convention originates from `terraform 0.12upgrade` generating `versions.tf`), so most existing code would violate this rule. It is opt-in for those who want strict guide compliance.

### Example

```hcl
# main.tf
terraform {
  required_version = ">= 1.0"
}

provider "aws" {
  region = "us-east-1"
}
```

```
$ tflint
2 issue(s) found:

Warning: 'terraform' block should be defined in 'terraform.tf' (style_guide_block_placement)

  on main.tf line 2:
   2: terraform {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_block_placement.md

Warning: 'provider' block should be defined in 'providers.tf' (style_guide_block_placement)

  on main.tf line 6:
   6: provider "aws" {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_block_placement.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#file-names

> - `terraform.tf` contains `terraform` block... to configure Terraform itself, such as the required providers and Terraform version.
> - `providers.tf` contains `provider` blocks.
> - `backend.tf` contains `backend` configuration.

### How To Fix

Move the blocks into their conventional files:

```hcl
# terraform.tf
terraform {
  required_version = ">= 1.0"
}
```

```hcl
# providers.tf
provider "aws" {
  region = "us-east-1"
}
```
