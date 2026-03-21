## style_guide_typed_variables_except_any

Warns when `any` is used as variable type, including inside composite types such as `list(any)` or `object({ key = any })`. 

### Example

```hcl
variable "config" {
  type = any
}

variable "items" {
  type = list(any)
}

variable "record" {
  type = object({
    name  = string
    value = any
  })
}
```

```
$ tflint
3 issue(s) found:

Warning: Using 'any' as variable type should be avoided (style_guide_typed_variables_except_any)

  on test.tf line 2:
   2:   type        = any

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.0/docs/rules/style_guide_typed_variables_except_any.md

Warning: Using 'any' as variable type should be avoided (style_guide_typed_variables_except_any)

  on test.tf line 6:
   6:   type        = list(any)

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.0/docs/rules/style_guide_typed_variables_except_any.md

Warning: Using 'any' as variable type should be avoided (style_guide_typed_variables_except_any)

  on test.tf line 12:
  12:     value = any

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.1.0/docs/rules/style_guide_typed_variables_except_any.md
```

### Why

Using `any` defeats the purpose of type constraints.

https://developer.hashicorp.com/terraform/language/expressions/type-constraints#dynamic-types-the-any-constraint

> Warning: `any` is very rarely the correct type constraint to use. **Do not use `any` just to avoid specifying a type constraint.** Always write an exact type constraint unless you are truly handling dynamic data.

### How To Fix

Use a proper type to the variable. See https://developer.hashicorp.com/terraform/language/values/variables#type-constraints for more details about types