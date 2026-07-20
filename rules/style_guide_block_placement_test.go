package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideBlockPlacementRule(t *testing.T) {
	cases := []struct {
		Name     string
		Files    map[string]string
		Expected helper.Issues
	}{
		{
			Name: "terraform block in terraform.tf",
			Files: map[string]string{
				"terraform.tf": heredoc.Doc(`
					terraform {
					  required_version = ">= 1.0"
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "terraform block in main.tf",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					terraform {
					  required_version = ">= 1.0"
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'terraform' block should be defined in 'terraform.tf'",
				},
			},
		},
		{
			Name: "provider block in providers.tf",
			Files: map[string]string{
				"providers.tf": heredoc.Doc(`
					provider "aws" {
					  region = "us-east-1"
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "provider block in main.tf",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "aws" {
					  region = "us-east-1"
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'provider' block should be defined in 'providers.tf'",
				},
			},
		},
		{
			Name: "terraform block with backend in backend.tf",
			Files: map[string]string{
				"backend.tf": heredoc.Doc(`
					terraform {
					  backend "s3" {
					    bucket = "example"
					  }
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "terraform block with backend in terraform.tf",
			Files: map[string]string{
				"terraform.tf": heredoc.Doc(`
					terraform {
					  backend "s3" {
					    bucket = "example"
					  }
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'terraform' block should be defined in 'backend.tf' (contains a backend configuration)",
				},
			},
		},
		{
			Name: "terraform block with both required_providers and backend must be in backend.tf",
			Files: map[string]string{
				"terraform.tf": heredoc.Doc(`
					terraform {
					  required_version = ">= 1.0"

					  backend "s3" {
					    bucket = "example"
					  }
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'terraform' block should be defined in 'backend.tf' (contains a backend configuration)",
				},
			},
		},
		{
			Name: "de facto versions.tf convention is still flagged",
			Files: map[string]string{
				"versions.tf": heredoc.Doc(`
					terraform {
					  required_version = ">= 1.0"
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'terraform' block should be defined in 'terraform.tf'",
				},
			},
		},
		{
			Name: "other block types are not checked",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					resource "aws_instance" "web" {
					  instance_type = "t3.micro"
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "multiple provider blocks in the wrong file",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					provider "aws" {
					  region = "us-east-1"
					}

					provider "aws" {
					  alias  = "west"
					  region = "us-west-2"
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'provider' block should be defined in 'providers.tf'",
				},
				{
					Rule:    NewStyleGuideBlockPlacementRule(),
					Message: "'provider' block should be defined in 'providers.tf'",
				},
			},
		},
	}

	rule := NewStyleGuideBlockPlacementRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.Files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}
