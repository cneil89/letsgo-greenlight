package validator

import (
	"testing"
)

func TestNew(t *testing.T) {
	v := New()
	if v.Errors == nil {
		t.Error("expected Errors map to be initialized, got nil")
	}
	if len(v.Errors) != 0 {
		t.Error("expected empty Errors map")
	}
}

func TestValidator_Valid(t *testing.T) {
	tests := []struct {
		name     string
		errors   map[string]string
		expected bool
	}{
		{"valid with no errors", map[string]string{}, true},
		{"invalid with errors", map[string]string{"email": "invalid"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{Errors: tt.errors}
			if got := v.Valid(); got != tt.expected {
				t.Errorf("Valid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidator_AddError(t *testing.T) {
	t.Run("adds new error", func(t *testing.T) {
		v := New()
		v.AddError("email", "invalid email")
		if v.Errors["email"] != "invalid email" {
			t.Errorf("expected error message 'invalid email', got '%s'", v.Errors["email"])
		}
	})

	t.Run("ignores duplicate key", func(t *testing.T) {
		v := New()
		v.AddError("email", "first error")
		v.AddError("email", "second error")
		if v.Errors["email"] != "first error" {
			t.Errorf("expected 'first error', got '%s'", v.Errors["email"])
		}
		if len(v.Errors) != 1 {
			t.Errorf("expected 1 error, got %d", len(v.Errors))
		}
	})
}

func TestValidator_Check(t *testing.T) {
	t.Run("adds error when condition is false", func(t *testing.T) {
		v := New()
		v.Check(false, "email", "invalid email")
		if v.Valid() {
			t.Error("expected validator to be invalid")
		}
	})

	t.Run("does not add error when condition is true", func(t *testing.T) {
		v := New()
		v.Check(true, "email", "invalid email")
		if !v.Valid() {
			t.Error("expected validator to be valid")
		}
	})
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with dot", "user.name@example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"invalid email no at", "testexample.com", false},
		{"invalid email no domain", "test@", false},
		{"invalid email empty", "", false},
		{"invalid email spaces", "test @example.com", false},
		{"invalid email starts with dash", "test@-example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(tt.value, EmailRX); got != tt.expected {
				t.Errorf("Matches(%q) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestPermittedValue(t *testing.T) {
	t.Run("value in permitted list", func(t *testing.T) {
		if !PermittedValue("red", "red", "green", "blue") {
			t.Error("expected true for value in list")
		}
	})

	t.Run("value not in permitted list", func(t *testing.T) {
		if PermittedValue("yellow", "red", "green", "blue") {
			t.Error("expected false for value not in list")
		}
	})

	t.Run("empty permitted list", func(t *testing.T) {
		if PermittedValue("red") {
			t.Error("expected false for empty permitted list")
		}
	})

	t.Run("integer values", func(t *testing.T) {
		if !PermittedValue(2, 1, 2, 3) {
			t.Error("expected true for integer in list")
		}
		if PermittedValue(4, 1, 2, 3) {
			t.Error("expected false for integer not in list")
		}
	})
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{"all unique", []string{"a", "b", "c"}, true},
		{"has duplicates", []string{"a", "b", "a"}, false},
		{"all same", []string{"a", "a", "a"}, false},
		{"empty slice", []string{}, true},
		{"single element", []string{"a"}, true},
		{"two same elements", []string{"a", "a"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.values); got != tt.expected {
				t.Errorf("Unique(%v) = %v, want %v", tt.values, got, tt.expected)
			}
		})
	}
}

func TestUnique_Integers(t *testing.T) {
	if !Unique([]int{1, 2, 3}) {
		t.Error("expected true for unique integers")
	}
	if Unique([]int{1, 2, 1}) {
		t.Error("expected false for duplicate integers")
	}
}
