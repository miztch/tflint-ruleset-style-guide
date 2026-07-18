package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideOrderedResourceArgumentsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "arguments in recommended order",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  count = 2

				  ami           = "ami-0c55b159cbfafe1f0"
				  instance_type = "t3.micro"

				  ebs_block_device {
				    device_name = "/dev/sdh"
				  }

				  lifecycle {
				    create_before_destroy = true
				  }

				  depends_on = [aws_iam_role.web]
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "count after non-block argument",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ami = "ami-0c55b159cbfafe1f0"

				  count = 2
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'count' should be defined before 'ami' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "non-block argument after block argument",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ebs_block_device {
				    device_name = "/dev/sdh"
				  }

				  ami = "ami-0c55b159cbfafe1f0"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'ami' should be defined before 'ebs_block_device' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "block argument after lifecycle",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ami = "ami-0c55b159cbfafe1f0"

				  lifecycle {
				    create_before_destroy = true
				  }

				  ebs_block_device {
				    device_name = "/dev/sdh"
				  }
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'ebs_block_device' should be defined before 'lifecycle' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "lifecycle after depends_on",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ami = "ami-0c55b159cbfafe1f0"

				  depends_on = [aws_iam_role.web]

				  lifecycle {
				    create_before_destroy = true
				  }
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'lifecycle' should be defined before 'depends_on' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "dynamic block counts as a block argument",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  dynamic "ebs_block_device" {
				    for_each = var.devices
				    content {
				      device_name = ebs_block_device.value
				    }
				  }

				  ami = "ami-0c55b159cbfafe1f0"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'ami' should be defined before 'dynamic' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "dynamic block with for_each inside is not flagged when in order",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  count = 2

				  ami = "ami-0c55b159cbfafe1f0"

				  dynamic "ebs_block_device" {
				    for_each = var.devices
				    content {
				      device_name = ebs_block_device.value
				    }
				  }
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "one issue per misplaced argument even with multiple higher-rank predecessors",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  lifecycle {
				    create_before_destroy = true
				  }

				  depends_on = [aws_iam_role.web]

				  ami = "ami-0c55b159cbfafe1f0"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'ami' should be defined before 'lifecycle' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "ignored provider between misordered arguments",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ami = "ami-0c55b159cbfafe1f0"

				  provider = aws.us_east_1

				  count = 2
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'count' should be defined before 'ami' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "provisioner and connection after depends_on are ignored",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  ami = "ami-0c55b159cbfafe1f0"

				  depends_on = [aws_iam_role.web]

				  connection {
				    type = "ssh"
				    host = self.public_ip
				  }

				  provisioner "local-exec" {
				    command = "echo done"
				  }
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "data blocks are checked",
			Content: heredoc.Doc(`
				data "aws_ami" "ubuntu" {
				  most_recent = true

				  for_each = var.regions
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedResourceArgumentsRule(),
					Message: "'for_each' should be defined before 'most_recent' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
				},
			},
		},
		{
			Name: "provider, provisioner and connection are ignored",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  provider = aws.us_east_1

				  count = 2

				  ami = "ami-0c55b159cbfafe1f0"

				  connection {
				    type = "ssh"
				    host = self.public_ip
				  }

				  provisioner "local-exec" {
				    command = "echo done"
				  }

				  depends_on = [aws_iam_role.web]
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "other block types are not checked",
			Content: heredoc.Doc(`
				module "vpc" {
				  cidr_block = "10.0.0.0/16"

				  count = 2
				}
			`),
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideOrderedResourceArgumentsRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}

// TestStyleGuideOrderedResourceArgumentsRule_Ranges asserts the reported
// ranges: the name range for attributes, and the open brace range of the
// correct occurrence when multiple blocks share the same type.
func TestStyleGuideOrderedResourceArgumentsRule_Ranges(t *testing.T) {
	content := heredoc.Doc(`
		resource "aws_instance" "web" {
		  ebs_block_device {
		    device_name = "/dev/sdh"
		  }

		  ami = "ami-0c55b159cbfafe1f0"

		  lifecycle {
		    create_before_destroy = true
		  }

		  ebs_block_device {
		    device_name = "/dev/sdi"
		  }
		}
	`)

	rule := NewStyleGuideOrderedResourceArgumentsRule()
	expected := helper.Issues{
		{
			Rule:    rule,
			Message: "'ami' should be defined before 'ebs_block_device' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
			Range: hcl.Range{
				Filename: "main.tf",
				Start:    hcl.Pos{Line: 6, Column: 3},
				End:      hcl.Pos{Line: 6, Column: 6},
			},
		},
		{
			Rule:    rule,
			Message: "'ebs_block_device' should be defined before 'lifecycle' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)",
			Range: hcl.Range{
				Filename: "main.tf",
				Start:    hcl.Pos{Line: 12, Column: 20},
				End:      hcl.Pos{Line: 12, Column: 21},
			},
		},
	}

	runner := helper.TestRunner(t, map[string]string{"main.tf": content})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("Unexpected error occurred: %s", err)
	}

	helper.AssertIssues(t, expected, runner.Issues)
}
