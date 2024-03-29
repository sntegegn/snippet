package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func StringContains(t *testing.T, actual, expectedString string) {
	t.Helper()

	if !strings.Contains(actual, expectedString) {
		t.Errorf("expected to contain %v but got %v", expectedString, actual)
	}
}

func NilError(t *testing.T, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("expecting to get nil error but got %v", actual)
	}
}
