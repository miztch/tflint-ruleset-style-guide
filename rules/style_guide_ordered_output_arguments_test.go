package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideOrderedOutputArgumentsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "arguments in recommended order with type",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  type        = string
				  description = "The private IP address of the instance"
				  value       = aws_instance.web.private_ip
				  sensitive   = false
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "arguments in recommended order without type",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  description = "The private IP address of the instance"
				  value       = aws_instance.web.private_ip
				  sensitive   = false
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "description before type",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  description = "The private IP address of the instance"
				  type        = string
				  value       = aws_instance.web.private_ip
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedOutputArgumentsRule(),
					Message: "'type' should be defined before 'description' (recommended order: type, description, value, sensitive)",
				},
			},
		},
		{
			Name: "value before description",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  value       = aws_instance.web.private_ip
				  description = "The private IP address of the instance"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedOutputArgumentsRule(),
					Message: "'description' should be defined before 'value' (recommended order: type, description, value, sensitive)",
				},
			},
		},
		{
			Name: "sensitive before value",
			Content: heredoc.Doc(`
				output "db_password" {
				  description = "The database password"
				  sensitive   = true
				  value       = aws_db_instance.db.password
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedOutputArgumentsRule(),
					Message: "'value' should be defined before 'sensitive' (recommended order: type, description, value, sensitive)",
				},
			},
		},
		{
			Name: "unknown arguments are ignored",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  description = "The private IP address of the instance"
				  value       = aws_instance.web.private_ip
				  depends_on  = [aws_instance.web]
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "precondition block is ignored",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  description = "The private IP address of the instance"

				  precondition {
				    condition     = aws_instance.web.private_ip != ""
				    error_message = "IP address must not be empty."
				  }

				  value = aws_instance.web.private_ip
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "single argument",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  value = aws_instance.web.private_ip
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "multiple outputs with one misordered",
			Content: heredoc.Doc(`
				output "instance_ip_addr" {
				  description = "The private IP address of the instance"
				  value       = aws_instance.web.private_ip
				}

				output "instance_id" {
				  value       = aws_instance.web.id
				  description = "The ID of the instance"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideOrderedOutputArgumentsRule(),
					Message: "'description' should be defined before 'value' (recommended order: type, description, value, sensitive)",
				},
			},
		},
	}

	rule := NewStyleGuideOrderedOutputArgumentsRule()

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
