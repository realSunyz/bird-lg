package validation

import "testing"

func TestValidateToolTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		target    string
		wantValue string
		wantError string
	}{
		{name: "public IPv4", target: "8.8.8.8", wantValue: "8.8.8.8"},
		{name: "domain", target: "Internet.Example", wantValue: "internet.example"},
		{name: "IPv6 loopback", target: "::1", wantError: "target_bogon_blocked"},
		{name: "IPv4-mapped loopback", target: "::ffff:127.0.0.1", wantError: "target_bogon_blocked"},
		{name: "IPv4 private", target: "192.168.1.1", wantError: "target_bogon_blocked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, errKey := ValidateToolTarget(tt.target)
			if value != tt.wantValue || errKey != tt.wantError {
				t.Fatalf("ValidateToolTarget(%q) = (%q, %q), want (%q, %q)", tt.target, value, errKey, tt.wantValue, tt.wantError)
			}
		})
	}
}
