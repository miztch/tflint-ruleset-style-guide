package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideTypeRepetitionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "whole type repeated",
			Content: `resource "aws_iam_role" "iam_role" {
  name = "test-iam-role"
}`,
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
			Content: `resource "aws_codepipeline" "codepipeline" {
  name = "test-pipeline"
}`,
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
			Content: `resource "aws_iam_role" "lambda_role" {
  name = "test-lambda-role"
}

data "aws_instance" "web_instance" {
  instance_id = "i-instanceid"
}`,
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
			Content: `resource "aws_iam_role" "lambda-role" {
  name = "test-lambda-role"
}`,
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
			Content: `resource "aws_iam_role" "lambda.role" {
  name = "test-lambda-role"
}`,
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
			Content: `resource "aws_iam_role" "lambda" {
  name = "test-lambda-role"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issue for camel case",
			Content: `resource "aws_iam_role" "lambdaRole" {
  name = "test-lambda-role"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideTypeRepetitionRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
