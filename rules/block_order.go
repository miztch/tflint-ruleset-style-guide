package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// blockOrderViolation describes a labeled block defined after another block
// it should alphabetically precede.
type blockOrderViolation struct {
	name   string
	before string
	rng    hcl.Range
}

// findBlockOrderViolations reports blocks whose first label sorts before an
// earlier block's first label, in violation of alphabetical order. Blocks
// without a label are ignored. Comparison is byte-wise (case-sensitive)
func findBlockOrderViolations(blocks []*hclsyntax.Block) []blockOrderViolation {
	var seen []*hclsyntax.Block
	var violations []blockOrderViolation

	for _, block := range blocks {
		if len(block.Labels) == 0 {
			continue
		}
		name := block.Labels[0]

		for _, prev := range seen {
			if prev.Labels[0] > name {
				violations = append(violations, blockOrderViolation{
					name:   name,
					before: prev.Labels[0],
					rng:    block.DefRange(),
				})
				break
			}
		}

		seen = append(seen, block)
	}

	return violations
}
