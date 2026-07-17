package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideOrderedOutputArgumentsRule warns when arguments in an output
// block are not in the recommended order.
type StyleGuideOrderedOutputArgumentsRule struct {
	tflint.DefaultRule
}

// NewStyleGuideOrderedOutputArgumentsRule creates a new rule.
func NewStyleGuideOrderedOutputArgumentsRule() *StyleGuideOrderedOutputArgumentsRule {
	return &StyleGuideOrderedOutputArgumentsRule{}
}

// Name returns the rule name.
func (r *StyleGuideOrderedOutputArgumentsRule) Name() string {
	return "style_guide_ordered_output_arguments"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideOrderedOutputArgumentsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideOrderedOutputArgumentsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideOrderedOutputArgumentsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// outputArgumentRanks is the recommended order of arguments in an output block.
var outputArgumentRanks = map[string]int{
	"description": 1,
	"value":       2,
	"sensitive":   3,
}

// Message returns the rule message for a misplaced argument.
func (r *StyleGuideOrderedOutputArgumentsRule) Message(name, before string) string {
	return fmt.Sprintf("'%s' should be defined before '%s' (recommended order: description, value, sensitive)", name, before)
}

// Check checks whether output block arguments are in the recommended order.
func (r *StyleGuideOrderedOutputArgumentsRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if block.Type != BlockTypeOutput {
				continue
			}

			for _, v := range findArgumentOrderViolations(block, outputArgumentRanks) {
				if err := runner.EmitIssue(r, r.Message(v.name, v.before), v.rng); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
