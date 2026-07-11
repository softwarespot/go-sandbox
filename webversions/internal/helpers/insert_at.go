package helpers

import "sort"

type InsertAtOp struct {
	Value  string
	PosIdx int
}

// InsertAt inserts multiple values into a string at the provided original
// positions.
// It applies inserts from highest position to lowest so the original indexes
// remain valid.
// When positions are equal, later operations are inserted first so
// earlier operations remain earlier in the final string.
func InsertAt(s string, ops ...InsertAtOp) string {
	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].PosIdx == ops[j].PosIdx {
			return i > j
		}
		return ops[i].PosIdx > ops[j].PosIdx
	})
	for _, op := range ops {
		s = insertAt(s, op)
	}
	return s
}

func insertAt(s string, op InsertAtOp) string {
	posIdx := op.PosIdx
	switch {
	case posIdx < 0:
		posIdx = 0
	case posIdx > len(s):
		posIdx = len(s)
	}
	return s[:posIdx] + op.Value + s[posIdx:]
}
