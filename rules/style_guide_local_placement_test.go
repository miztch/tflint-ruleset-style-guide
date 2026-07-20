package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideLocalPlacementRule(t *testing.T) {
	cases := []struct {
		Name     string
		Files    map[string]string
		Expected helper.Issues
	}{
		{
			Name: "locals block in locals.tf",
			Files: map[string]string{
				"locals.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}
				`),
				"main.tf": heredoc.Doc(`
					resource "aws_instance" "web" {
					  instance_type = local.name
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "locals block at the top of a single-purpose file",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}

					resource "aws_instance" "web" {
					  instance_type = local.name
					}
				`),
			},
			Expected: helper.Issues{},
		},
		{
			Name: "locals block after another block type",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					resource "aws_instance" "web" {
					  instance_type = local.name
					}

					locals {
					  name = "example"
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideLocalPlacementRule(),
					Message: "'locals' block should be defined in 'locals.tf', or moved to the top of the file if specific to this file",
				},
			},
		},
		{
			Name: "local defined outside locals.tf and referenced from another file",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}
				`),
				"outputs.tf": heredoc.Doc(`
					output "name" {
					  value = local.name
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideLocalPlacementRule(),
					Message: "'local.name' is defined in 'main.tf' but referenced from 'outputs.tf'; locals referenced from multiple files should be defined in 'locals.tf'",
				},
			},
		},
		{
			Name: "local referenced multiple times from the same other file is reported once",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}
				`),
				"outputs.tf": heredoc.Doc(`
					output "name" {
					  value = local.name
					}

					output "name_again" {
					  value = local.name
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideLocalPlacementRule(),
					Message: "'local.name' is defined in 'main.tf' but referenced from 'outputs.tf'; locals referenced from multiple files should be defined in 'locals.tf'",
				},
			},
		},
		{
			Name: "local referenced from multiple distinct other files is reported once per file",
			Files: map[string]string{
				"main.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}
				`),
				"outputs.tf": heredoc.Doc(`
					output "name" {
					  value = local.name
					}
				`),
				"extra.tf": heredoc.Doc(`
					resource "aws_instance" "extra" {
					  instance_type = local.name
					}
				`),
			},
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideLocalPlacementRule(),
					Message: "'local.name' is defined in 'main.tf' but referenced from 'outputs.tf'; locals referenced from multiple files should be defined in 'locals.tf'",
				},
				{
					Rule:    NewStyleGuideLocalPlacementRule(),
					Message: "'local.name' is defined in 'main.tf' but referenced from 'extra.tf'; locals referenced from multiple files should be defined in 'locals.tf'",
				},
			},
		},
		{
			Name: "local defined in locals.tf and referenced from multiple files is not checked",
			Files: map[string]string{
				"locals.tf": heredoc.Doc(`
					locals {
					  name = "example"
					}
				`),
				"main.tf": heredoc.Doc(`
					resource "aws_instance" "web" {
					  instance_type = local.name
					}
				`),
				"outputs.tf": heredoc.Doc(`
					output "name" {
					  value = local.name
					}
				`),
			},
			Expected: helper.Issues{},
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
	}

	rule := NewStyleGuideLocalPlacementRule()

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
