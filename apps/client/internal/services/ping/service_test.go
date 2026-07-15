package ping

import "testing"

func TestBuildCommandTargetValidation(t *testing.T) {
	oldBin := pingBin
	oldOutstanding := pingSupportsOutstanding
	pingBin = "/usr/bin/ping"
	pingSupportsOutstanding = false
	t.Cleanup(func() {
		pingBin = oldBin
		pingSupportsOutstanding = oldOutstanding
	})

	if _, _, err := BuildCommand("internet.example", 4); err != nil {
		t.Fatalf("valid domain rejected: %v", err)
	}
	for _, target := range []string{"host\nname", "-c", "host name"} {
		if _, _, err := BuildCommand(target, 4); err == nil {
			t.Errorf("invalid target %q was accepted", target)
		}
	}
}
