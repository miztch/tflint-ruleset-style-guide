package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideOrderedVariableArgumentsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "arguments in recommended order",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type        = number
				  description = "Number of instances"
				  default     = 1
				  sensitive   = false

				  validation {
				    condition     = var.instance_count > 0
				    error_message = "Must be positive."
				  }
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "description before type",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  description = "Number of instances"
				  type        = number
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'type' should be defined before 'description' (recommended order: type, description, default, sensitive, validation)",
				},
			},
		},
		{
			Name: "default before description",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type        = number
				  default     = 1
				  description = "Number of instances"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'description' should be defined before 'default' (recommended order: type, description, default, sensitive, validation)",
				},
			},
		},
		{
			Name: "validation block before default",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type        = number
				  description = "Number of instances"

				  validation {
				    condition     = var.instance_count > 0
				    error_message = "Must be positive."
				  }

				  default = 1
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'default' should be defined before 'validation' (recommended order: type, description, default, sensitive, validation)",
				},
			},
		},
		{
			Name: "multiple misplaced arguments",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  sensitive   = false
				  default     = 1
				  description = "Number of instances"
				  type        = number
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'default' should be defined before 'sensitive' (recommended order: type, description, default, sensitive, validation)",
				},
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'description' should be defined before 'sensitive' (recommended order: type, description, default, sensitive, validation)",
				},
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'type' should be defined before 'sensitive' (recommended order: type, description, default, sensitive, validation)",
				},
			},
		},
		{
			Name: "unknown arguments are ignored",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type        = number
				  nullable    = false
				  description = "Number of instances"
				  default     = 1
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "single argument",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type = number
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "multiple variables with one misordered",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type        = number
				  description = "Number of instances"
				}

				variable "instance_name" {
				  description = "Name of the instance"
				  type        = string
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedVariableArgumentsRule(),
					Message: "'type' should be defined before 'description' (recommended order: type, description, default, sensitive, validation)",
				},
			},
		},
		{
			Name: "other block types are not checked",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  instance_type = "t3.micro"
				  ami           = "ami-0c55b159cbfafe1f0"
				}
			`),
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideOrderedVariableArgumentsRule()

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
