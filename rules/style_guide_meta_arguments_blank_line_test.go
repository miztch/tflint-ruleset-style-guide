package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideMetaArgumentsBlankLineRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing blank line after count",
			Content: `resource "aws_instance" "example" {
  count = 1
  ami   = "ami-123"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 12},
					},
				},
			},
		},
		{
			Name: "missing blank line after for_each",
			Content: `resource "aws_instance" "example" {
  for_each = toset(["a", "b"])
  ami      = "ami-123"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 22},
					},
				},
			},
		},
		{
			Name: "missing blank line after source",
			Content: `module "example" {
  source  = "./modules/example"
  version = "1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 30},
					},
				},
			},
		},
		{
			Name: "missing blank line after provider",
			Content: `resource "aws_instance" "example" {
  provider = aws.west
  ami      = "ami-123"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 11},
					},
				},
			},
		},
		{
			Name: "missing blank line after provider in data source",
			Content: `data "aws_ami" "ubuntu" {
  provider    = aws.west
  most_recent = true
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 11},
					},
				},
			},
		},
		{
			Name: "missing blank line after providers",
			Content: `module "example" {
  source    = "./modules/example"
  providers = {
    aws = aws.west
  }
  version = "1.0.0"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be followed by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 12},
					},
				},
			},
		},
		{
			Name: "no issue when count is the only attribute",
			Content: `resource "aws_instance" "example" {
  count = 1
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count is the last item",
			Content: `resource "aws_instance" "example" {
  ami   = "ami-123"

  count = 1
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Leading meta argument should not be the last item in the block",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 12},
					},
				},
			},
		},
		{
			Name: "for_each is the last item",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"

  for_each = toset(["a", "b"])
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Leading meta argument should not be the last item in the block",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 22},
					},
				},
			},
		},
		{
			Name: "no issue when provider is the last item",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"

  provider = aws.west
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when provider is the last item in data source",
			Content: `data "aws_ami" "ubuntu" {
  most_recent = true

  provider = aws.west
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when providers is the last item",
			Content: `module "example" {
  source = "./modules/example"

  providers = {
    aws = aws.west
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when source is the last item",
			Content: `module "example" {
  version = "1.0.0"

  source = "./modules/example"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when source is a regular attribute in resource (not a meta-argument)",
			Content: `resource "aws_s3_object" "example" {
  bucket = "my-bucket"
  key    = "file.txt"
  source = "path/to/file.txt"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when provider is a regular attribute in action block",
			Content: `resource "aws_codepipeline" "example" {
  name     = "example"
  role_arn = "arn:aws:iam::123456789012:role/example"

  stage {
    name = "Source"

    action {
      name     = "Source"
      category = "Source"
      owner    = "AWS"
      provider = "CodeCommit"
    }
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when two leading meta arguments are consecutive",
			Content: `resource "aws_instance" "example" {
  for_each = toset(["a"])
  count    = 1

  ami = "ami-123"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "missing blank line before lifecycle",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"
  lifecycle {
    prevent_destroy = true
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 12},
					},
				},
			},
		},
		{
			Name: "missing blank line before connection",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"
  connection {
    type = "ssh"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 12},
					},
				},
			},
		},
		{
			Name: "missing blank line before provisioner",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"
  provisioner "local-exec" {
    command = "echo hello"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 14},
					},
				},
			},
		},
		{
			Name: "missing blank line before depends_on",
			Content: `resource "aws_instance" "example" {
  ami        = "ami-123"
  depends_on = [aws_vpc.example]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 13},
					},
				},
			},
		},
		{
			Name: "missing blank line before consecutive trailing meta arguments",
			Content: `resource "aws_instance" "example" {
  ami        = "ami-123"
  depends_on = [aws_vpc.example]
  lifecycle {
    prevent_destroy = true
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 13},
					},
				},
				{
					Rule:    NewStyleGuideMetaArgumentsBlankLineRule(),
					Message: "Meta argument should be preceded by a blank line",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 12},
					},
				},
			},
		},
		{
			Name: "no issue when trailing meta argument is the only content",
			Content: `resource "aws_instance" "example" {
  lifecycle {
    prevent_destroy = true
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when all blank lines are present",
			Content: `resource "aws_instance" "example" {
  count = 1

  ami = "ami-123"

  lifecycle {
    prevent_destroy = true
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue when comment line precedes trailing meta argument",
			Content: `resource "aws_instance" "example" {
  ami = "ami-123"
  # This is required to prevent accidental deletion
  lifecycle {
    prevent_destroy = true
  }
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideMetaArgumentsBlankLineRule()

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
