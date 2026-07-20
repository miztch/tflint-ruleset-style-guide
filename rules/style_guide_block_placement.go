package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideBlockPlacementRule warns when terraform, provider, or
// backend blocks are placed outside their conventional files.
type StyleGuideBlockPlacementRule struct {
	tflint.DefaultRule
}

// NewStyleGuideBlockPlacementRule creates a new rule.
func NewStyleGuideBlockPlacementRule() *StyleGuideBlockPlacementRule {
	return &StyleGuideBlockPlacementRule{}
}

// Name returns the rule name.
func (r *StyleGuideBlockPlacementRule) Name() string {
	return "style_guide_block_placement"
}

// Enabled returns whether the rule is enabled by default.
// This rule is disabled by default because the de facto community standard
// for the terraform block is versions.tf, not terraform.tf as the style
// guide literally recommends, so most existing code would violate it.
func (r *StyleGuideBlockPlacementRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *StyleGuideBlockPlacementRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideBlockPlacementRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Message returns the rule message for a misplaced block.
func (r *StyleGuideBlockPlacementRule) Message(blockType, want string, hasBackend bool) string {
	if hasBackend {
		return fmt.Sprintf("'%s' block should be defined in '%s' (contains a backend configuration)", blockType, want)
	}
	return fmt.Sprintf("'%s' block should be defined in '%s'", blockType, want)
}

// Check checks whether terraform, provider, and backend blocks are defined
// in their conventional files (terraform.tf, providers.tf, backend.tf).
func (r *StyleGuideBlockPlacementRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for filename, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}
		base := filepath.Base(filename)

		for _, block := range body.Blocks {
			switch block.Type {
			case BlockTypeProvider:
				if base != "providers.tf" {
					if err := runner.EmitIssue(r, r.Message(BlockTypeProvider, "providers.tf", false), block.DefRange()); err != nil {
						return err
					}
				}
			case BlockTypeTerraform:
				hasBackend := blockContainsBackend(block)
				want := "terraform.tf"
				if hasBackend {
					want = "backend.tf"
				}
				if base != want {
					if err := runner.EmitIssue(r, r.Message(BlockTypeTerraform, want, hasBackend), block.DefRange()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// blockContainsBackend reports whether the block has a nested backend block.
func blockContainsBackend(block *hclsyntax.Block) bool {
	for _, b := range block.Body.Blocks {
		if b.Type == BlockTypeBackend {
			return true
		}
	}
	return false
}
