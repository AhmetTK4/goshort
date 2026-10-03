package service

import "testing"

func TestGenerateShortCode(t *testing.T) {
	code := GenerateShortCode(6)
	if len(code) != 6 {
		t.Fatalf("expected length 6, got %d", len(code))
	}
	for _, c := range code {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Fatalf("invalid character in code: %c", c)
		}
	}
}

func TestValidateCustomCode(t *testing.T) {
	tests := []struct {
		name  string
		code  string
		valid bool
	}{
		{"valid lowercase", "abc", true},
		{"valid uppercase", "ABC", true},
		{"valid mixed", "Abc123", true},
		{"valid 20 chars", "abcdefghij1234567890", true},
		{"too short", "ab", false},
		{"too long", "abcdefghij1234567890x", false},
		{"contains dash", "abc-def", false},
		{"contains underscore", "abc_def", false},
		{"contains space", "abc def", false},
		{"contains special char", "abc@def", false},
		{"empty", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateCustomCode(tc.code)
			if result != tc.valid {
				t.Fatalf("ValidateCustomCode(%q) = %v, want %v", tc.code, result, tc.valid)
			}
		})
	}
}
