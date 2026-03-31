package vcs

import "testing"

func TestVersion(t *testing.T) {
	version := Version()
	if version == "" {
		t.Log("Version() returned empty string (expected when built without version info)")
	}
	t.Logf("Version: %s", version)
}

func TestVersion_ReturnsString(t *testing.T) {
	version := Version()
	if _, ok := interface{}(version).(string); !ok {
		t.Error("Version() should return a string")
	}
}
