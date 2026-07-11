package main

import (
	"encoding/json"
	"sort"
)

type InsertManyAtOp struct {
	Value  string
	PosIdx int
}

// InsertManyAt inserts multiple values into a string at the provided original
// positions.
// It applies inserts from highest position to lowest so the original indexes
// remain valid. When positions are equal, later ops are inserted first so
// earlier ops remain earlier in the final string.
func InsertManyAt(s string, ops []InsertManyAtOp) string {
	if len(ops) == 0 {
		return s
	}
	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].PosIdx == ops[j].PosIdx {
			// When positions are equal, insert later ops first so earlier ops remain earlier in the result.
			return i > j
		}
		return ops[i].PosIdx > ops[j].PosIdx
	})
	for _, op := range ops {
		s = InsertAt(s, op.Value, op.PosIdx)
	}
	return s
}

func InsertAt(s, insert string, posIdx int) string {
	if posIdx < 0 {
		posIdx = 0
	}
	if posIdx > len(s) {
		posIdx = len(s)
	}
	return s[:posIdx] + insert + s[posIdx:]
}

func ToJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
