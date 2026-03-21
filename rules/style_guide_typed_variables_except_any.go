package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideTypeVariablesExceptAnyRule warns when a variable is declared with type any or contains any.
type StyleGuideTypeVariablesExceptAnyRule struct {
	tflint.DefaultRule
}

// NewStyleGuideTypeVariablesExceptAnyRule creates a new rule.
func NewStyleGuideTypeVariablesExceptAnyRule() *StyleGuideTypeVariablesExceptAnyRule {
	return &StyleGuideTypeVariablesExceptAnyRule{}
}

// Name returns the rule name.
func (r *StyleGuideTypeVariablesExceptAnyRule) Name() string {
	return "style_guide_typed_variables_except_any"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideTypeVariablesExceptAnyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideTypeVariablesExceptAnyRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideTypeVariablesExceptAnyRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Message returns the rule message.
func (r *StyleGuideTypeVariablesExceptAnyRule) Message() string {
	return "Using 'any' as variable type should be avoided"
}

func (r *StyleGuideTypeVariablesExceptAnyRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}

	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{{
			Type:       "variable",
			LabelNames: []string{"name"},
			Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{{Name: "type"}},
			},
		}},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		attr, exists := variable.Body.Attributes["type"]
		if !exists {
			continue
		}

		// attr.Expr is hcl.Expression, so need to cast hclsyntax.Expression to analyze the type expression.
		syntaxExpr, ok := attr.Expr.(hclsyntax.Expression)
		if !ok {
			continue
		}

		ranges := collectAnyRanges(syntaxExpr)
		for _, rng := range ranges {
			if err := runner.EmitIssue(r, r.Message(), rng); err != nil {
				return err
			}
		}
	}

	return nil
}

// collectAnyRanges walks a type expression and returns the ranges that should
// be flagged. The logic mirrors what the test cases expect:
//
//   - `any`           → flag the whole expression range
//   - `list(any)`     → flag the whole expression range (not just the inner any)
//   - `set(any)`      → same as list
//   - `map(any)`      → recurse: flag the inner `any` range
//   - `object({...})` → recurse: flag any `any` values found inside
//   - `list(object({...}))` → recurse into the object
//   - `tuple([type, any, type])` → flag the inner `any` range
//   - `optional(any)` → flag the inner `any` range (not the whole optional(...))
//   - `optional(object({nested=any}))` → recurse into the object and flag the inner `any`
func collectAnyRanges(expr hclsyntax.Expression) []hcl.Range {
	switch e := expr.(type) {

	// A bare scope-traversal like `any`, `string`, `bool`, etc.
	case *hclsyntax.ScopeTraversalExpr:
		if len(e.Traversal) == 1 && e.Traversal.RootName() == "any" {
			return []hcl.Range{e.Range()}
		}
		return nil

	case *hclsyntax.FunctionCallExpr:
		name := e.Name
		args := e.Args

		switch name {
		case "list", "set":
			// Flag the whole expression if the sole argument is (or contains) any.
			// The tests show list(any) and set(any) use the outer range,
			// but list(object({nested=any})) recurses into the object.
			if len(args) == 1 {
				if isBareAny(args[0]) {
					// Flag the whole list(any)/set(any) expression.
					return []hcl.Range{e.Range()}
				}
				// Otherwise recurse (e.g. list(object({...}))).
				return collectAnyRanges(args[0])
			}

		case "map":
			// map(any) → flag the inner any; map(object({...})) → recurse.
			if len(args) == 1 {
				return collectAnyRanges(args[0])
			}

		case "object":
			// object takes a single ObjectConsExpr argument.
			if len(args) == 1 {
				return collectAnyRanges(args[0])
			}

		case "tuple":
			// tuple([type, type, ...]) — the [...] is a single TupleConsExpr arg
			if len(args) == 1 {
				if tce, ok := args[0].(*hclsyntax.TupleConsExpr); ok {
					var ranges []hcl.Range
					for _, ex := range tce.Exprs {
						ranges = append(ranges, collectAnyRanges(ex)...)
					}
					return ranges
				}
			}

		case "optional":
			// optional(type) or optional(type, default) — only check arg[0]
			if len(args) >= 1 {
				return collectAnyRanges(args[0])
			}
		}

		return nil

	// object({key = type, ...}) is represented as an ObjectConsExpr.
	case *hclsyntax.ObjectConsExpr:
		var ranges []hcl.Range
		for _, item := range e.Items {
			ranges = append(ranges, collectAnyRanges(item.ValueExpr)...)
		}
		return ranges
	}

	return nil
}

// isBareAny returns true only if expr is a plain `any` traversal.
func isBareAny(expr hclsyntax.Expression) bool {
	t, ok := expr.(*hclsyntax.ScopeTraversalExpr)
	return ok && len(t.Traversal) == 1 && t.Traversal.RootName() == "any"
}
