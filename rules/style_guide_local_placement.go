package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideLocalPlacementRule warns when local values are not placed
// according to the style guide: in locals.tf if referenced across multiple
// files, or at the top of their file if specific to a single file.
type StyleGuideLocalPlacementRule struct {
	tflint.DefaultRule
}

// NewStyleGuideLocalPlacementRule creates a new rule.
func NewStyleGuideLocalPlacementRule() *StyleGuideLocalPlacementRule {
	return &StyleGuideLocalPlacementRule{}
}

// Name returns the rule name.
func (r *StyleGuideLocalPlacementRule) Name() string {
	return "style_guide_local_placement"
}

// Enabled returns whether the rule is enabled by default.
// This rule is disabled by default because co-locating locals right above
// the resources that use them is a reasonable, cohesion-friendly style that
// this rule would flag. Opt-in for strict guide compliance.
func (r *StyleGuideLocalPlacementRule) Enabled() bool {
	return false
}

// Severity returns the rule severity.
func (r *StyleGuideLocalPlacementRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideLocalPlacementRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// MessagePosition returns the rule message for a locals block that is
// neither in locals.tf nor at the top of its file.
func (r *StyleGuideLocalPlacementRule) MessagePosition() string {
	return "'locals' block should be defined in 'locals.tf', or moved to the top of the file if specific to this file"
}

// MessageMultiFile returns the rule message for a local value defined
// outside locals.tf but referenced from a file other than the one that
// defines it.
func (r *StyleGuideLocalPlacementRule) MessageMultiFile(name, definedIn, referencedIn string) string {
	return fmt.Sprintf("'local.%s' is defined in '%s' but referenced from '%s'; locals referenced from multiple files should be defined in 'locals.tf'", name, definedIn, referencedIn)
}

// Check checks whether local values are placed according to the style
// guide's two recommended placements: locals.tf for locals referenced
// across multiple files, or the top of the file for locals specific to a
// single file.
func (r *StyleGuideLocalPlacementRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	// local name -> file it's defined in, for locals defined outside locals.tf.
	defined := map[string]string{}

	for filename, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}
		base := filepath.Base(filename)
		if base == "locals.tf" {
			continue
		}

		seenOtherBlock := false
		for _, block := range body.Blocks {
			if block.Type != BlockTypeLocals {
				seenOtherBlock = true
				continue
			}

			if seenOtherBlock {
				if err := runner.EmitIssue(r, r.MessagePosition(), block.DefRange()); err != nil {
					return err
				}
			}

			for name := range block.Body.Attributes {
				defined[name] = filename
			}
		}
	}

	if len(defined) == 0 {
		return nil
	}

	reported := map[string]bool{}
	var emitErr error

	diags := runner.WalkExpressions(tflint.ExprWalkFunc(func(expr hcl.Expression) hcl.Diagnostics {
		for _, traversal := range expr.Variables() {
			if len(traversal) < 2 || traversal.RootName() != "local" {
				continue
			}
			attr, ok := traversal[1].(hcl.TraverseAttr)
			if !ok {
				continue
			}

			definedIn, ok := defined[attr.Name]
			if !ok {
				continue
			}

			referencedIn := traversal.SourceRange().Filename
			if referencedIn == definedIn {
				continue
			}

			key := attr.Name + "|" + referencedIn
			if reported[key] {
				continue
			}
			reported[key] = true

			if err := runner.EmitIssue(r, r.MessageMultiFile(attr.Name, filepath.Base(definedIn), filepath.Base(referencedIn)), traversal.SourceRange()); err != nil {
				emitErr = err
				return nil
			}
		}
		return nil
	}))
	if emitErr != nil {
		return emitErr
	}
	if diags.HasErrors() {
		return diags
	}

	return nil
}
