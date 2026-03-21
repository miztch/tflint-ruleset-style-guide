package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestStyleGuideTypeVariablesExceptAnyRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "variable with type any",
			Content: `variable "test" {
  type = any
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 10},
						End:      hcl.Pos{Line: 2, Column: 13},
					},
				},
			},
		},
		{
			Name: "variable with type list of any",
			Content: `variable "test" {
  type = list(any)
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 10},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: "variable with type set of any",
			Content: `variable "test" {
  type = set(any)
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 10},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
			},
		},
		{
			Name: "variable with type map with any",
			Content: `variable "test" {
  type = map(any)
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 14},
						End:      hcl.Pos{Line: 2, Column: 17},
					},
				},
			},
		},
		{
			Name: "variable with type object with any",
			Content: `variable "test" {
  type = object({
	nested = any
  })
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 14},
					},
				},
			},
		},
		{
			Name: "variable with list of object with any",
			Content: `variable "test" {
  type = list(object({
	nested = any
  }))
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 11},
						End:      hcl.Pos{Line: 3, Column: 14},
					},
				},
			},
		},
		{
			Name: "variable with type tuple with any",
			Content: `variable "test" {
  type = tuple([string, any, number])
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 25},
						End:      hcl.Pos{Line: 2, Column: 28},
					},
				},
			},
		},
		{
			Name: "variable with type optional any",
			Content: `variable "test" {
  type = optional(any)
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 19},
						End:      hcl.Pos{Line: 2, Column: 22},
					},
				},
			},
		},
		{
			Name: "variable with type optional any with default",
			Content: `variable "test" {
  type = optional(any, "default")
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 19},
						End:      hcl.Pos{Line: 2, Column: 22},
					},
				},
			},
		},
		{
			Name: "variable with object containing optional any",
			Content: `variable "test" {
  type = object({
	nested = optional(any)
  })
}`,
			Expected: helper.Issues{
				{
					Rule:    NewStyleGuideTypeVariablesExceptAnyRule(),
					Message: "Using 'any' as variable type should be avoided",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 20},
						End:      hcl.Pos{Line: 3, Column: 23},
					},
				},
			},
		},
		{
			Name: "variable with tuple without any",
			Content: `variable "test" {
  type = tuple([string, number, bool])
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable with optional without any",
			Content: `variable "test" {
  type = object({
	nested = optional(string, "default")
  })
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable with specific type",
			Content: `variable "test" {
  type = string
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "variable without type",
			Content: `variable "test" {
  description = "A variable without a type"
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewStyleGuideTypeVariablesExceptAnyRule()

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
