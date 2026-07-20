package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideAlphabeticalBlocksRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "variable blocks in alphabetical order",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type = number
				}

				variable "instance_name" {
				  type = string
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "variable blocks not in alphabetical order",
			Content: heredoc.Doc(`
				variable "instance_name" {
				  type = string
				}

				variable "instance_count" {
				  type = number
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideAlphabeticalBlocksRule(),
					Message: "'instance_count' should be defined before 'instance_name' (variable blocks should be in alphabetical order)",
				},
			},
		},
		{
			Name: "output blocks in alphabetical order",
			Content: heredoc.Doc(`
				output "instance_arn" {
				  value = aws_instance.web.arn
				}

				output "instance_id" {
				  value = aws_instance.web.id
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "output blocks not in alphabetical order",
			Content: heredoc.Doc(`
				output "instance_id" {
				  value = aws_instance.web.id
				}

				output "instance_arn" {
				  value = aws_instance.web.arn
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideAlphabeticalBlocksRule(),
					Message: "'instance_arn' should be defined before 'instance_id' (output blocks should be in alphabetical order)",
				},
			},
		},
		{
			Name: "variable and output blocks are ordered independently",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type = number
				}

				output "instance_arn" {
				  value = aws_instance.web.arn
				}

				variable "instance_name" {
				  type = string
				}

				output "instance_id" {
				  value = aws_instance.web.id
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "multiple misplaced blocks",
			Content: heredoc.Doc(`
				variable "z" {
				  type = string
				}

				variable "a" {
				  type = string
				}

				variable "b" {
				  type = string
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideAlphabeticalBlocksRule(),
					Message: "'a' should be defined before 'z' (variable blocks should be in alphabetical order)",
				},
				{
					Rule:    NewStyleGuideAlphabeticalBlocksRule(),
					Message: "'b' should be defined before 'z' (variable blocks should be in alphabetical order)",
				},
			},
		},
		{
			Name: "single variable block",
			Content: heredoc.Doc(`
				variable "instance_count" {
				  type = number
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "other block types are not checked",
			Content: heredoc.Doc(`
				resource "aws_instance" "web" {
				  instance_type = "t3.micro"
				}

				resource "aws_instance" "app" {
				  instance_type = "t3.micro"
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "comparison is byte-wise (case-sensitive)",
			Content: heredoc.Doc(`
				variable "Zeta" {
				  type = string
				}

				variable "alpha" {
				  type = string
				}
			`),
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideAlphabeticalBlocksRule()

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
