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
	return findArgumentOrderViolationsByRank(block, func(item bodyItem) (int, bool) {
		rank, ok := ranks[item.name]
		return rank, ok
	})
}

// findArgumentOrderViolationsByRank reports arguments in the block that appear
// after an argument they should precede, according to the given rank function.
// Items for which the rank function returns false are ignored.
func findArgumentOrderViolationsByRank(block *hclsyntax.Block, rank func(bodyItem) (int, bool)) []argumentOrderViolation {
	items := collectItems(block.Body)

	var seen []bodyItem
	var violations []argumentOrderViolation

	for _, item := range items {
		itemRank, ok := rank(item)
		if !ok {
			continue
		}

		for _, prev := range seen {
			if prevRank, _ := rank(prev); prevRank > itemRank {
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
