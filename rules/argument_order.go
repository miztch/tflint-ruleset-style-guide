package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// argumentOrderViolation describes an argument defined after another argument
// it should precede.
type argumentOrderViolation struct {
	name   string
	before string
	rng    hcl.Range
}

// findArgumentOrderViolations reports arguments in the block that appear after
// an argument they should precede, according to the given rank map.
// Arguments not present in the rank map are ignored.
func findArgumentOrderViolations(block *hclsyntax.Block, ranks map[string]int) []argumentOrderViolation {
	items := collectItems(block.Body)

	var seen []bodyItem
	var violations []argumentOrderViolation

	for _, item := range items {
		rank, ok := ranks[item.name]
		if !ok {
			continue
		}

		for _, prev := range seen {
			if ranks[prev.name] > rank {
				violations = append(violations, argumentOrderViolation{
					name:   item.name,
					before: prev.name,
					rng:    attrOrBlockRange(block, item),
				})
				break
			}
		}

		seen = append(seen, item)
	}

	return violations
}
