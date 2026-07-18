package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideOrderedResourceArgumentsRule warns when arguments in a resource
// or data block are not in the recommended order.
type StyleGuideOrderedResourceArgumentsRule struct {
	tflint.DefaultRule
}

// NewStyleGuideOrderedResourceArgumentsRule creates a new rule.
func NewStyleGuideOrderedResourceArgumentsRule() *StyleGuideOrderedResourceArgumentsRule {
	return &StyleGuideOrderedResourceArgumentsRule{}
}

// Name returns the rule name.
func (r *StyleGuideOrderedResourceArgumentsRule) Name() string {
	return "style_guide_ordered_resource_arguments"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideOrderedResourceArgumentsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideOrderedResourceArgumentsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideOrderedResourceArgumentsRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// resourceArgumentRank returns the recommended order rank for an item in a
// resource or data block. The style guide's ordering does not cover the
// provider meta-argument or provisioner/connection blocks, so they are ignored.
func resourceArgumentRank(item bodyItem) (int, bool) {
	switch item.name {
	case "count", "for_each":
		return 1, true
	case "provider", "provisioner", "connection":
		return 0, false
	case "lifecycle":
		return 4, true
	case "depends_on":
		return 5, true
	}
	if item.isBlock {
		return 3, true
	}
	return 2, true
}

// Message returns the rule message for a misplaced argument.
func (r *StyleGuideOrderedResourceArgumentsRule) Message(name, before string) string {
	return fmt.Sprintf("'%s' should be defined before '%s' (recommended order: count/for_each, non-block arguments, block arguments, lifecycle, depends_on)", name, before)
}

// Check checks whether resource and data block arguments are in the
// recommended order.
func (r *StyleGuideOrderedResourceArgumentsRule) Check(runner tflint.Runner) error {
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
			if block.Type != BlockTypeResource && block.Type != BlockTypeData {
				continue
			}

			for _, v := range findArgumentOrderViolationsByRank(block, resourceArgumentRank) {
				if err := runner.EmitIssue(r, r.Message(v.name, v.before), v.rng); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
