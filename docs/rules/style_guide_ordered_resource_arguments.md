# style_guide_ordered_resource_arguments

Warns when arguments in a `resource` or `data` block are not in the recommended order:

1. `count` / `for_each`
2. Non-block arguments
3. Block arguments
4. `lifecycle` block
5. `depends_on`

The `provider` meta-argument and `provisioner` / `connection` blocks are not
covered by the style guide's ordering and are not checked.

### Example

```hcl
resource "aws_instance" "web" {
  ebs_block_device {
    device_name = "/dev/sdh"
  }

  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"
}
```

```
$ tflint
2 issue(s) found:

Warning: 'ami' should be defined before 'ebs_block_device' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on) (style_guide_ordered_resource_arguments)

  on main.tf line 6:
   6:   ami           = "ami-0c55b159cbfafe1f0"

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_ordered_resource_arguments.md

Warning: 'instance_type' should be defined before 'ebs_block_device' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on) (style_guide_ordered_resource_arguments)

  on main.tf line 7:
   7:   instance_type = "t3.micro"

Reference: https://github.com/miztch/tflint-ruleset-style-guide/blob/v0.2.0/docs/rules/style_guide_ordered_resource_arguments.md
```

### Why

https://developer.hashicorp.com/terraform/language/style#resource-order

> We recommend the following order for resource parameters:
>
> 1. If present, The `count` or `for_each` meta-argument.
> 1. Resource-specific non-block parameters.
> 1. Resource-specific block parameters.
> 1. If required, a `lifecycle` block.
> 1. If required, the `depends_on` parameter.

### How To Fix

Reorder the arguments:

```hcl
resource "aws_instance" "web" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.micro"

  ebs_block_device {
    device_name = "/dev/sdh"
  }
}
```
