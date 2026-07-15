package bird

import (
	"strings"
	"unicode"
)

func IsAllowedCommand(cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" || len(cmd) > 4096 || strings.ContainsRune(cmd, ';') || strings.IndexFunc(cmd, unicode.IsControl) >= 0 {
		return false
	}

	allowed := []string{
		"show route",
		"show protocols",
		"show protocol",
		"show status",
		"show memory",
		"show interfaces",
		"show ospf",
		"show bfd",
		"show roa",
		"show static",
		"show symbols",
	}

	for _, prefix := range allowed {
		if cmd == prefix || (strings.HasPrefix(cmd, prefix) && cmd[len(prefix)] == ' ') {
			return true
		}
	}

	return false
}
