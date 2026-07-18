package rules

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// bodyItem represents a single attribute or block within an HCL body,
// unified so they can be sorted by line number.
type bodyItem struct {
	name    string
	rng     hcl.Range
	isBlock bool
}

// collectItems gathers all attributes and blocks from an HCL body into a
// unified slice sorted by start line.
func collectItems(body *hclsyntax.Body) []bodyItem {
	var items []bodyItem

	for name, attr := range body.Attributes {
		items = append(items, bodyItem{
			name:    name,
			rng:     attr.Range(),
			isBlock: false,
		})
	}

	for _, block := range body.Blocks {
		items = append(items, bodyItem{
			name:    block.Type,
			rng:     block.Range(),
			isBlock: true,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].rng.Start.Line < items[j].rng.Start.Line
	})

	return items
}

// attrOrBlockRange returns the HCL range for a bodyItem within a given block.
func attrOrBlockRange(block *hclsyntax.Block, item bodyItem) hcl.Range {
	if item.isBlock {
		for _, b := range block.Body.Blocks {
			if b.Type == item.name && b.Range().Start.Line == item.rng.Start.Line {
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
