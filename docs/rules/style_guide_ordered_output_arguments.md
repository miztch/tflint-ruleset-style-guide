## style_guide_ordered_output_arguments

Warns when arguments in an `output` block are not in the recommended order:

1. `type`
2. `description`
3. `value`
4. `sensitive`

`output` blocks can declare an explicit `type` constraint since Terraform 1.15.

Arguments not listed above (such as `ephemeral`, `depends_on` and `precondition` blocks) are not checked.

### Example

```hcl
output "instance_ip_addr" {
  value       = aws_instance.web.private_ip
  description = "The private IP address of the instance"
}
```

```
$ tflint
1 issue(s) found:

Warning: 'description' should be defined before 'value' (recommended order: description, value, sensitive) (style_guide_ordered_output_arguments)

  on main.tf line 3:
   3:   description = "The private IP address of the instance"

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_ordered_output_arguments.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#outputs

> We recommend that you use the following order for your output parameters:
>
> 1. Type
> 1. Description
> 1. Value
> 1. Sensitive (optional)

### How To Fix

Reorder the arguments:

```hcl
output "instance_ip_addr" {
  description = "The private IP address of the instance"
  value       = aws_instance.web.private_ip
}
```
