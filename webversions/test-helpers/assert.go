package testhelpers

import (
	"reflect"
	"testing"
)

// AssertEqual checks if two values are equal. If they are not, it logs using t.Fatalf()
func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("AssertEqual: expected values to be equal, got:\n%+v\ncorrect:\n%+v", got, want)
	}
}
