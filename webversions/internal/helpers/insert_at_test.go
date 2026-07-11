package helpers

import (
	"testing"

	testhelpers "webversions/test-helpers"
)

func TestInsertAt(t *testing.T) {
	t.Run("descending positions", func(t *testing.T) {
		s := "abcdef"
		ops := []InsertAtOp{
			{Value: "X", PosIdx: 2},
			{Value: "Y", PosIdx: 5},
		}
		got := InsertAt(s, ops...)
		want := "abXcdeYf"
		testhelpers.AssertEqual(t, got, want)
	})

	t.Run("same position preserves input order", func(t *testing.T) {
		s := "abc"
		ops := []InsertAtOp{
			{Value: "1", PosIdx: 1},
			{Value: "2", PosIdx: 1},
		}
		got := InsertAt(s, ops...)
		want := "a12bc"
		testhelpers.AssertEqual(t, got, want)
	})

	t.Run("skips negative positions using InsertAt semantics", func(t *testing.T) {
		s := "abc"
		ops := []InsertAtOp{
			{Value: "X", PosIdx: -1},
			{Value: "Y", PosIdx: 1},
		}
		got := InsertAt(s, ops...)
		want := "XaYbc"
		testhelpers.AssertEqual(t, got, want)
	})
}
