package traceroute

import "testing"

func TestBuildCommandTargetValidation(t *testing.T) {
	oldBin := tracerouteBin
	oldFlags := tracerouteFlags
	tracerouteBin = "/usr/bin/traceroute"
	tracerouteFlags = []string{"-I"}
	t.Cleanup(func() {
		tracerouteBin = oldBin
		tracerouteFlags = oldFlags
	})

	if _, _, err := BuildCommand("internet.example"); err != nil {
		t.Fatalf("valid domain rejected: %v", err)
	}
	for _, target := range []string{"host\rname", "--help", "host name"} {
		if _, _, err := BuildCommand(target); err == nil {
			t.Errorf("invalid target %q was accepted", target)
		}
	}
}
