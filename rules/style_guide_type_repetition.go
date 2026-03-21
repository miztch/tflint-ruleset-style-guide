package rules

import (
	"fmt"
	"strings"

	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideTypeRepetitionRule warns when a resource's or data source's name
// repeats part of its type
// Example: resource "aws_iam_role" "lambda_role" → repeats "role"
// Example: data "aws_iam_role" "lambda_role"     → repeats "role"
type StyleGuideTypeRepetitionRule struct {
	tflint.DefaultRule
}

// NewStyleGuideTypeRepetitionRule creates a new rule.
func NewStyleGuideTypeRepetitionRule() *StyleGuideTypeRepetitionRule {
	return &StyleGuideTypeRepetitionRule{}
}

// Name returns the rule name.
func (r *StyleGuideTypeRepetitionRule) Name() string {
	return "style_guide_type_repetition"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideTypeRepetitionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideTypeRepetitionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideTypeRepetitionRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Message returns the rule message based on its type
func (r *StyleGuideTypeRepetitionRule) Message(blockType string) string {
	switch blockType {
	case "resource":
		return "Resource name should not repeat its resource type"
	case "data":
		return "Data source name should not repeat its data source type"
	default:
		panic(fmt.Sprintf("terraform_type_repetition: unexpected block type %q", blockType))
	}
}

// Check checks whether the rule is satisfied.
func (r *StyleGuideTypeRepetitionRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}

	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, block := range body.Blocks {
		if len(block.Labels) < 2 {
			continue
		}
		typeName := block.Labels[0]
		name := block.Labels[1]
		if hasRepeatedTypeSegment(typeName, name) {
			rng := block.DefRange
			msg := r.Message(block.Type)

			if err := runner.EmitIssue(r, msg, rng); err != nil {
				return err
			}
		}
	}

	return nil
}

func hasRepeatedTypeSegment(typeName, name string) bool {
	typeSegs := strings.Split(strings.ToLower(typeName), "_")
	nameSegs := splitName(strings.ToLower(name))

	nameSet := make(map[string]struct{}, len(nameSegs))
	for _, ns := range nameSegs {
		nameSet[ns] = struct{}{}
	}

	for _, ts := range typeSegs {
		if _, ok := nameSet[ts]; ok {
			return true
		}
	}

	return false
}

func splitName(name string) []string {
	return strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})
}
