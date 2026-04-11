package rules

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideTypeRepetitionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "whole type repeated",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "iam_role" {
				  name = "test-iam-role"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 35},
					},
				},
			},
		},
		{
			Name: "single word type repeated",
			Content: heredoc.Doc(`
				resource "aws_codepipeline" "codepipeline" {
				  name = "test-pipeline"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 43},
					},
				},
			},
		},
		{
			Name: "partially type repeated",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "lambda_role" {
				  name = "test-lambda-role"
				}

				data "aws_instance" "web_instance" {
				  instance_id = "i-instanceid"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 38},
					},
				},
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Data source name should not repeat its data source type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 1},
						End:      hcl.Pos{Line: 5, Column: 35},
					},
				},
			},
		},
		{
			Name: "kebab case",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "lambda-role" {
				  name = "test-lambda-role"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 38},
					},
				},
			},
		},
		{
			Name: "dot separated",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "lambda.role" {
				  name = "test-lambda-role"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 38},
					},
				},
			},
		},
		{
			Name: "no issues",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "lambda" {
				  name = "test-lambda-role"
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "no issue for camel case",
			Content: heredoc.Doc(`
				resource "aws_iam_role" "lambdaRole" {
				  name = "test-lambda-role"
				}
			`),
			Expected: helper.Issues{},
		},
		{
			Name: "provider prefix repeated in name warns by default",
			Content: heredoc.Doc(`
				resource "aws_s3_bucket" "aws_backup" {
				  bucket = "my-backup-bucket"
				}
			`),
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeRepetitionRule(),
					Message: "Resource name should not repeat its resource type",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 1, Column: 1},
						End:      hcl.Pos{Line: 1, Column: 38},
					},
				},
			},
		},
		{
			Name: "provider prefix repetition ignored with config",
			Content: heredoc.Doc(`
				resource "aws_s3_bucket" "aws_backup" {
				  bucket = "my-backup-bucket"
				}
			`),
			Config: heredoc.Doc(`
				rule "style_guide_type_repetition" {
				  enabled                   = true
				  ignored_provider_prefixes = ["aws"]
				}
			`),
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideTypeRepetitionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			files := map[string]string{"main.tf": tc.Content}
			if tc.Config != "" {
				files[".tflint.hcl"] = tc.Config
			}
			runner := helper.TestRunner(t, files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
