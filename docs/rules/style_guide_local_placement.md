## style_guide_local_placement

Warns when local values are not placed according to the style guide's two recommended placements.

- **Placement within a file**: `locals` blocks outside `locals.tf` must be at the top of their file, before any other block type. A `locals` block that appears after another block satisfies neither of the guide's two placements.
- **Multi-file references**: a local defined outside `locals.tf` must not be referenced from a file other than the one that defines it. If it is, it should be moved to `locals.tf` instead.
- The reverse direction is not checked: a local defined in `locals.tf` is always accepted, even if it's only referenced from a single file.
- `.tf.json` files are not analyzed for the placement-within-a-file check (hclsyntax-based, same constraint as the existing rules in this ruleset).
- **Disabled by default.** Co-locating locals right above the resources that use them is a reasonable, cohesion-friendly style that this rule would flag. It is opt-in for those who want strict guide compliance.

### Example

```hcl
# main.tf
resource "aws_instance" "web" {
  instance_type = local.instance_type
}

locals {
  instance_type = "t3.micro"
}
```

```
$ tflint
1 issue(s) found:

Warning: 'locals' block should be defined in 'locals.tf', or moved to the top of the file if specific to this file (style_guide_local_placement)

  on main.tf line 5:
   5: locals {

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_local_placement.md
```

```hcl
# main.tf
locals {
  instance_type = "t3.micro"
}
```

```hcl
# outputs.tf
output "instance_type" {
  value = local.instance_type
}
```

```
$ tflint
1 issue(s) found:

Warning: 'local.instance_type' is defined in 'main.tf' but referenced from 'outputs.tf'; locals referenced from multiple files should be defined in 'locals.tf' (style_guide_local_placement)

  on outputs.tf line 2:
   2:   value = local.instance_type

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_local_placement.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#local-values

> Local values in Terraform are much like local variables in a programming language, and can be used similarly. Where relevant, they should be declared in a `locals.tf` file, or defined near the top of the file that uses them.
>
> If the local value is only relevant to one file, define them in that file. If they are used across multiple files (like `main.tf`, or a set of files beginning with a common prefix), put them in a `locals.tf` file.

### How To Fix

Move the local into `locals.tf` if it's referenced from more than one file:

```hcl
# locals.tf
locals {
  instance_type = "t3.micro"
}
```

Otherwise, move the `locals` block to the top of the file that uses it:

```hcl
# main.tf
locals {
  instance_type = "t3.micro"
}

resource "aws_instance" "web" {
  instance_type = local.instance_type
}
```
