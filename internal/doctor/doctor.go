package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Check struct {
	Name     string
	Required bool
	Found    bool
	Detail   string
}

func Run(target string) ([]Check, error) {
	if target == "" {
		target = "."
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("resolve target: %w", err)
	}
	checks := []struct {
		name     string
		required bool
		args     []string
	}{
		{"git", true, []string{"--version"}},
		{"go", true, []string{"version"}},
		{"node", true, []string{"--version"}},
		{"npm", true, []string{"--version"}},
		{"npx", true, []string{"--version"}},
		{"claude", false, []string{"--version"}},
		{"codex", false, []string{"--version"}},
		{"gh", false, []string{"--version"}},
		{"uv", false, []string{"--version"}},
		{"python", false, []string{"--version"}},
	}
	result := make([]Check, 0, len(checks)+3)
	for _, item := range checks {
		result = append(result, commandCheck(item.name, item.required, item.args...))
	}
	result = append(result,
		pathCheck("git repository", false, filepath.Join(abs, ".git")),
		pathCheck("Auto Company manifest", false, filepath.Join(abs, ".auto-company", "manifest.json")),
		pathCheck("project instructions", false, filepath.Join(abs, "AGENTS.md")),
	)
	return result, nil
}

func HasMissingRequired(checks []Check) bool {
	for _, check := range checks {
		if check.Required && !check.Found {
			return true
		}
	}
	return false
}

func commandCheck(name string, required bool, args ...string) Check {
	path, err := exec.LookPath(name)
	if err != nil {
		return Check{Name: name, Required: required, Found: false, Detail: "not found in PATH"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, path, args...).CombinedOutput()
	detail := strings.TrimSpace(string(output))
	if len(detail) > 120 {
		detail = detail[:120] + "…"
	}
	if err != nil {
		if detail == "" {
			detail = err.Error()
		}
		return Check{Name: name, Required: required, Found: false, Detail: detail}
	}
	return Check{Name: name, Required: required, Found: true, Detail: detail}
}

func pathCheck(name string, required bool, path string) Check {
	_, err := os.Stat(path)
	if err == nil {
		return Check{Name: name, Required: required, Found: true, Detail: path}
	}
	return Check{Name: name, Required: required, Found: false, Detail: path}
}
