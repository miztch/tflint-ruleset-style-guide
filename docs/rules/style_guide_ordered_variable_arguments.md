## style_guide_ordered_variable_arguments

Warns when arguments in a `variable` block are not in the recommended order:

1. `type`
2. `description`
3. `default`
4. `sensitive`
5. `validation` block(s)

Arguments not listed above (such as `nullable`) are not checked.

### Example

```hcl
variable "instance_count" {
  default     = 1
  description = "Number of instances"
  type        = number
}
```

```
$ tflint
2 issue(s) found:

Warning: 'description' should be defined before 'default' (recommended order: type, description, default, sensitive, validation) (style_guide_ordered_variable_arguments)

  on main.tf line 3:
   3:   description = "Number of instances"

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_ordered_variable_arguments.md

Warning: 'type' should be defined before 'default' (recommended order: type, description, default, sensitive, validation) (style_guide_ordered_variable_arguments)

  on main.tf line 4:
   4:   type        = number

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_ordered_variable_arguments.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#variables

> We recommend following a consistent order for variable parameters:
>
> 1. Type
> 1. Description
> 1. Default (optional)
> 1. Sensitive (optional)
> 1. Validation blocks

### How To Fix

Reorder the arguments:

```hcl
variable "instance_count" {
  type        = number
  description = "Number of instances"
  default     = 1
}
```
