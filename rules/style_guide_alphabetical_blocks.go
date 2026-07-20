package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideAlphabeticalBlocksRule warns when variable or output blocks
// within a file are not in alphabetical order by name.
type StyleGuideAlphabeticalBlocksRule struct {
	tflint.DefaultRule
}

// NewStyleGuideAlphabeticalBlocksRule creates a new rule.
func NewStyleGuideAlphabeticalBlocksRule() *StyleGuideAlphabeticalBlocksRule {
	return &StyleGuideAlphabeticalBlocksRule{}
}

// Name returns the rule name.
func (r *StyleGuideAlphabeticalBlocksRule) Name() string {
	return "style_guide_alphabetical_blocks"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideAlphabeticalBlocksRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *StyleGuideAlphabeticalBlocksRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideAlphabeticalBlocksRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Message returns the rule message for a misplaced block.
func (r *StyleGuideAlphabeticalBlocksRule) Message(blockType, name, before string) string {
	return fmt.Sprintf("'%s' should be defined before '%s' (%s blocks should be in alphabetical order)", name, before, blockType)
}

// Check checks whether variable and output blocks are in alphabetical order,
// per file. Ordering is checked independently for each block type and is
// not restricted to files named variables.tf / outputs.tf.
func (r *StyleGuideAlphabeticalBlocksRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, blockType := range []string{BlockTypeVariable, BlockTypeOutput} {
			var blocks []*hclsyntax.Block
			for _, block := range body.Blocks {
				if block.Type == blockType {
					blocks = append(blocks, block)
				}
			}

			for _, v := range findBlockOrderViolations(blocks) {
				if err := runner.EmitIssue(r, r.Message(blockType, v.name, v.before), v.rng); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
