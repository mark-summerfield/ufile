package ufile

import (
    "testing"
)

func Test001(t *testing.T) {
    expected := "Hello ufile v1.0.0\n"
    actual := Hello()
    if actual != expected {
        t.Errorf("expected %q, got %q", expected, actual)
    }
}
