package rules

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// StyleGuideMetaArgumentsBlankLineRule warns when meta arguments are not
// separated from other arguments by a blank line
type StyleGuideMetaArgumentsBlankLineRule struct {
	tflint.DefaultRule
}

// NewStyleGuideMetaArgumentsBlankLineRule creates a new rule.
func NewStyleGuideMetaArgumentsBlankLineRule() *StyleGuideMetaArgumentsBlankLineRule {
	return &StyleGuideMetaArgumentsBlankLineRule{}
}

// Name returns the rule name.
func (r *StyleGuideMetaArgumentsBlankLineRule) Name() string {
	return "style_guide_meta_arguments_blank_line"
}

// Enabled returns whether the rule is enabled by default.
func (r *StyleGuideMetaArgumentsBlankLineRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *StyleGuideMetaArgumentsBlankLineRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *StyleGuideMetaArgumentsBlankLineRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// leadingMessage returns the message for leading meta argument violations.
func (r *StyleGuideMetaArgumentsBlankLineRule) leadingMessage() string {
	return "Meta argument should be followed by a blank line"
}

// trailingMessage returns the message for trailing meta argument violations.
func (r *StyleGuideMetaArgumentsBlankLineRule) trailingMessage() string {
	return "Meta argument should be preceded by a blank line"
}

// leadingMetaArgs are meta arguments that must be followed by a blank line
var leadingMetaArgs = map[string]struct{}{
	"count":     {},
	"for_each":  {},
	"source":    {},
	"provider":  {}, // resource blocks: specify which provider to use
	"providers": {}, // module blocks: pass provider configurations to child modules
}

// trailingMetaArgs are meta arguments that must be preceded by a blank line
var trailingMetaArgs = map[string]struct{}{
	"lifecycle":   {},
	"connection":  {},
	"provisioner": {},
	"depends_on":  {},
}

// bodyItem represents a single attribute or block within an HCL body,
// unified so they can be sorted by line number.
type bodyItem struct {
	name      string
	startLine int
	endLine   int
	isBlock   bool
}

// Check checks whether meta arguments are properly separated by blank lines.
func (r *StyleGuideMetaArgumentsBlankLineRule) Check(runner tflint.Runner) error {
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
			if err := r.checkBlock(runner, block); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkBlock inspects a single block for blank line violations around meta arguments.
func (r *StyleGuideMetaArgumentsBlankLineRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	items := collectItems(block.Body)
	if len(items) < 2 {
		return nil
	}

	for i, item := range items {
		// Check leading meta args: must be followed by a blank line
		if _, ok := leadingMetaArgs[item.name]; ok {
			next := items[i+1] // safe: leading meta arg is never the last item if len >= 2
			if i+1 < len(items) && !hasBlankLineBetween(item.endLine, next.startLine) {
				if _, nextIsAlsoLeading := leadingMetaArgs[next.name]; !nextIsAlsoLeading {
					rng := attrOrBlockRange(block, item)
					msg := r.leadingMessage()
					if err := runner.EmitIssue(r, msg, rng); err != nil {
						return err
					}
				}
			}
		}

		// Check trailing meta args: must be preceded by a blank line
		if _, ok := trailingMetaArgs[item.name]; ok {
			if i > 0 {
				prev := items[i-1]
				if !hasBlankLineBetween(prev.endLine, item.startLine) {
					rng := attrOrBlockRange(block, item)
					msg := r.trailingMessage()
					if err := runner.EmitIssue(r, msg, rng); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// collectItems gathers all attributes and blocks from an HCL body into a
// unified slice sorted by start line.
func collectItems(body *hclsyntax.Body) []bodyItem {
	var items []bodyItem

	for name, attr := range body.Attributes {
		items = append(items, bodyItem{
			name:      name,
			startLine: attr.Range().Start.Line,
			endLine:   attr.Range().End.Line,
			isBlock:   false,
		})
	}

	for _, block := range body.Blocks {
		items = append(items, bodyItem{
			name:      block.Type,
			startLine: block.Range().Start.Line,
			endLine:   block.Range().End.Line,
			isBlock:   true,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].startLine < items[j].startLine
	})

	return items
}

// hasBlankLineBetween returns true if there is at least one blank line between
// the end of one item and the start of the next.
func hasBlankLineBetween(endLine, nextStartLine int) bool {
	return nextStartLine > endLine+1
}

// attrOrBlockRange returns the HCL range for a bodyItem within a given block.
func attrOrBlockRange(block *hclsyntax.Block, item bodyItem) hcl.Range {
	if item.isBlock {
		for _, b := range block.Body.Blocks {
			if b.Type == item.name && b.Range().Start.Line == item.startLine {
				return b.OpenBraceRange
			}
		}
	}
	for _, attr := range block.Body.Attributes {
		if attr.Name == item.name {
			return attr.NameRange
		}
	}
	return block.DefRange()
}
