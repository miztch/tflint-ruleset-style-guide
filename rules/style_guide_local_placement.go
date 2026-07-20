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

	defined, positionViolations := collectLocalDefinitions(files)
	for _, rng := range positionViolations {
		if err := runner.EmitIssue(r, r.MessagePosition(), rng); err != nil {
			return err
		}
	}

	referenceViolations, err := findMultiFileLocalReferences(runner, defined)
	if err != nil {
		return err
	}
	for _, v := range referenceViolations {
		msg := r.MessageMultiFile(v.name, filepath.Base(v.definedIn), filepath.Base(v.referencedIn))
		if err := runner.EmitIssue(r, msg, v.rng); err != nil {
			return err
		}
	}

	return nil
}

// collectLocalDefinitions scans every file for locals blocks outside
// locals.tf. It returns the file each local name is defined in (for the
// multi-file reference check) and the ranges of locals blocks that appear
// after another block type in their file, violating both of the guide's
// placements.
func collectLocalDefinitions(files map[string]*hcl.File) (defined map[string]string, positionViolations []hcl.Range) {
	defined = map[string]string{}

	for filename, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}
		if filepath.Base(filename) == "locals.tf" {
			continue
		}

		seenOtherBlock := false
		for _, block := range body.Blocks {
			if block.Type != BlockTypeLocals {
				seenOtherBlock = true
				continue
			}

			if seenOtherBlock {
				positionViolations = append(positionViolations, block.DefRange())
			}

			for name := range block.Body.Attributes {
				defined[name] = filename
			}
		}
	}

	return defined, positionViolations
}

// localReferenceViolation describes a local value defined outside locals.tf
// but referenced from a file other than the one that defines it.
type localReferenceViolation struct {
	name         string
	definedIn    string
	referencedIn string
	rng          hcl.Range
}

// findMultiFileLocalReferences walks every expression in the module and
// reports, once per (local, referencing file) pair, a reference to a local
// from a file other than its defining file.
func findMultiFileLocalReferences(runner tflint.Runner, defined map[string]string) ([]localReferenceViolation, error) {
	if len(defined) == 0 {
		return nil, nil
	}

	var violations []localReferenceViolation
	reported := map[string]bool{}

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

			violations = append(violations, localReferenceViolation{
				name:         attr.Name,
				definedIn:    definedIn,
				referencedIn: referencedIn,
				rng:          traversal.SourceRange(),
			})
		}
		return nil
	}))
	if diags.HasErrors() {
		return nil, diags
	}

	return violations, nil
}
