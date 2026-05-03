package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetSlug(cwd string) string {
	cmd := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err == nil {
		toplevel := strings.TrimSpace(string(out))
		if toplevel != "" {
			return filepath.Base(toplevel)
		}
	}
	home, _ := os.UserHomeDir()
	switch cwd {
	case home:
		return "[home]"
	case "/":
		return "[root]"
	default:
		return filepath.Base(cwd)
	}
}
