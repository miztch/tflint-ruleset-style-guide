package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideOrderedVariableArgumentsRule warns when arguments in a variable
// block are not in the recommended order.
type StyleGuideOrderedVariableArgumentsRule struct {
	tflint.DefaultRule
}

// NewStyleGuideOrderedVariableArgumentsRule creates a new rule.
func NewStyleGuideOrderedVariableArgumentsRule() *StyleGuideOrderedVariableArgumentsRule {
	return &StyleGuideOrderedVariableArgumentsRule{}
}

// Name returns the rule name.
func (r *StyleGuideOrderedVariableArgumentsRule) Name() string {
	return "style_guide_ordered_variable_arguments"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideOrderedVariableArgumentsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideOrderedVariableArgumentsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideOrderedVariableArgumentsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// variableArgumentRanks is the recommended order of arguments in a variable block.
var variableArgumentRanks = map[string]int{
	"type":        1,
	"description": 2,
	"default":     3,
	"sensitive":   4,
	"validation":  5,
}

// Message returns the rule message for a misplaced argument.
func (r *StyleGuideOrderedVariableArgumentsRule) Message(name, before string) string {
	return fmt.Sprintf("'%s' should be defined before '%s' (recommended order: type, description, default, sensitive, validation)", name, before)
}

// Check checks whether variable block arguments are in the recommended order.
func (r *StyleGuideOrderedVariableArgumentsRule) Check(runner tflint.Runner) error {
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
			if block.Type != BlockTypeVariable {
				continue
			}

			for _, v := range findArgumentOrderViolations(block, variableArgumentRanks) {
				if err := runner.EmitIssue(r, r.Message(v.name, v.before), v.rng); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
