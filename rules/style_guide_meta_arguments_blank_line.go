package rules

import (
	"slices"

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

// Message returns the appropriate error message for the given violation type.
func (r *StyleGuideMetaArgumentsBlankLineRule) Message(msgType string, args ...any) string {
	switch msgType {
	case "leading":
		return "Meta argument should be followed by a blank line"
	case "trailing":
		return "Meta argument should be preceded by a blank line"
	case "lastItem":
		return "Leading meta argument should not be the last item in the block"
	default:
		return ""
	}
}

// metaArgConfig holds configuration for a meta-argument
type metaArgConfig struct {
	shouldNotBeLast bool
	validBlocks     []string
}

// leadingMetaArgs are meta arguments that must be followed by a blank line
var leadingMetaArgs = map[string]metaArgConfig{
	"count": {
		shouldNotBeLast: true,
		validBlocks:     []string{BlockTypeResource, BlockTypeData, BlockTypeModule},
	},
	"for_each": {
		shouldNotBeLast: true,
		validBlocks:     []string{BlockTypeResource, BlockTypeData, BlockTypeModule},
	},
	"source": {
		shouldNotBeLast: false,
		validBlocks:     []string{BlockTypeModule},
	},
	"provider": {
		shouldNotBeLast: false,
		validBlocks:     []string{BlockTypeResource, BlockTypeData},
	},
	"providers": {
		shouldNotBeLast: false,
		validBlocks:     []string{BlockTypeModule},
	},
}

// trailingMetaArgs are meta arguments that must be preceded by a blank line
var trailingMetaArgs = map[string]metaArgConfig{
	"lifecycle": {
		validBlocks: []string{BlockTypeResource, BlockTypeData},
	},
	"connection": {
		validBlocks: []string{BlockTypeResource},
	},
	"provisioner": {
		validBlocks: []string{BlockTypeResource},
	},
	"depends_on": {
		validBlocks: []string{BlockTypeResource, BlockTypeData, BlockTypeModule},
	},
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
		if config, ok := leadingMetaArgs[item.name]; ok {
			// Skip if this meta-arg is not valid for this block type
			if !isValidForBlock(config, block.Type) {
				continue
			}

			if i+1 < len(items) {
				next := items[i+1]
				if !hasBlankLineBetween(item.rng.End.Line, next.rng.Start.Line) {
					if _, nextIsAlsoLeading := leadingMetaArgs[next.name]; !nextIsAlsoLeading {
						rng := attrOrBlockRange(block, item)
						msg := r.Message("leading")
						if err := runner.EmitIssueWithFix(r, msg, rng, func(f tflint.Fixer) error {
							return f.InsertTextAfter(item.rng, "\n")
						}); err != nil {
							return err
						}
					}
				}
			} else if config.shouldNotBeLast {
				// This leading meta arg should not be the last item
				rng := attrOrBlockRange(block, item)
				msg := r.Message("lastItem")
				if err := runner.EmitIssueWithFix(r, msg, rng, func(f tflint.Fixer) error {
					return tflint.ErrFixNotSupported
				}); err != nil {
					return err
				}
			}
		}

		// Check trailing meta args: must be preceded by a blank line
		if config, ok := trailingMetaArgs[item.name]; ok {
			// Skip if this meta-arg is not valid for this block type
			if !isValidForBlock(config, block.Type) {
				continue
			}

			if i > 0 {
				prev := items[i-1]
				if !hasBlankLineBetween(prev.rng.End.Line, item.rng.Start.Line) {
					rng := attrOrBlockRange(block, item)
					msg := r.Message("trailing")
					if err := runner.EmitIssueWithFix(r, msg, rng, func(f tflint.Fixer) error {
						return f.InsertTextBefore(item.rng, "\n")
					}); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// hasBlankLineBetween returns true if there is at least one blank line between
// the end of one item and the start of the next.
func hasBlankLineBetween(endLine, nextStartLine int) bool {
	return nextStartLine > endLine+1
}

// isValidForBlock checks if a meta-argument is valid for the given block type.
// Returns true if there are no restrictions or if the block type is in the allowed list.
func isValidForBlock(config metaArgConfig, blockType string) bool {
	if config.validBlocks == nil {
		return true
	}
	return slices.Contains(config.validBlocks, blockType)
}
