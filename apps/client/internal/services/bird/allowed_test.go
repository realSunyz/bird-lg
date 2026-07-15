package bird

import "testing"

func TestIsAllowedCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		command string
		allowed bool
	}{
		{command: "show status", allowed: true},
		{command: "show route all", allowed: true},
		{command: "show protocol bgp1", allowed: true},
		{command: "show status\ndown", allowed: false},
		{command: "show route\treload", allowed: false},
		{command: "show route ; down", allowed: false},
		{command: "show status-extra", allowed: false},
		{command: "show", allowed: false},
	}

	for _, tt := range tests {
		if got := IsAllowedCommand(tt.command); got != tt.allowed {
			t.Errorf("IsAllowedCommand(%q) = %v, want %v", tt.command, got, tt.allowed)
		}
	}
}
