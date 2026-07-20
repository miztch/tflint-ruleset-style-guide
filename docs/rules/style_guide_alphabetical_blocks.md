## style_guide_alphabetical_blocks

Warns when `variable` or `output` blocks within a file are not in alphabetical order by name.

- Checked per file, independently for `variable` and `output` blocks. A file containing both is checked for each block type separately.
- Not restricted to files named `variables.tf` / `outputs.tf` — any file containing `variable` or `output` blocks is checked.
- Comparison is byte-wise (case-sensitive), so uppercase letters sort before lowercase ones (e.g. `"Zeta"` before `"alpha"`).
  - Block names are expected to be `snake_case` by convention; if your codebase mixes letter case, enable [`terraform_naming_convention`](https://github.com/terraform-linters/tflint-ruleset-terraform/blob/master/docs/rules/terraform_naming_convention.md) (from `tflint-ruleset-terraform`) to enforce `snake_case` first, so this rule's ordering matches what you'd expect.
- **Disabled by default.** This rule is noisy on existing codebases, so it is opt-in for those who want strict guide compliance.

### Example

```hcl
variable "instance_name" {
  type = string
}

variable "instance_count" {
  type = number
}
```

```
$ tflint
1 issue(s) found:

Warning: 'instance_count' should be defined before 'instance_name' (variable blocks should be in alphabetical order) (style_guide_alphabetical_blocks)

  on main.tf line 5:
   5: variable "instance_count" {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_alphabetical_blocks.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#file-names

> `variables.tf` contains variable declarations, sorted alphabetically
> `outputs.tf` contains outputs, sorted alphabetically

### How To Fix

Reorder the blocks alphabetically:

```hcl
variable "instance_count" {
  type = number
}

variable "instance_name" {
  type = string
}
```
