package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	var out, errOut bytes.Buffer
	if err := Run([]string{"version"}, &out, &errOut); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got, want := out.String(), "autoco "+Version+"\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}

func TestUnknownCommand(t *testing.T) {
	t.Parallel()

	var out, errOut bytes.Buffer
	err := Run([]string{"unknown"}, &out, &errOut)
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("unexpected error = %v", err)
	}
	if !strings.Contains(errOut.String(), "Usage:") {
		t.Fatalf("usage not written to stderr")
	}
}

func TestInitPromptAndValidationFlow(t *testing.T) {
	t.Parallel()

	target := filepath.Join(t.TempDir(), "new-saas")
	var out, errOut bytes.Buffer
	err := Run([]string{
		"init", "--target", target, "--profile", "saas-web",
		"--name", "New SaaS", "--idea", "Reduce manual operations", "--agent", "both",
	}, &out, &errOut)
	if err != nil {
		t.Fatalf("init error = %v, stderr=%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), "Auto Company initialized") {
		t.Fatalf("unexpected init output: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	if err := Run([]string{"prompt", "--target", target}, &out, &errOut); err != nil {
		t.Fatalf("prompt error = %v", err)
	}
	if !strings.Contains(out.String(), "New SaaS") || !strings.Contains(out.String(), "Reduce manual operations") {
		t.Fatalf("generated prompt missing project data: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	err = Run([]string{"validate", "--target", target}, &out, &errOut)
	if err != nil && !errors.Is(err, ErrValidation) {
		t.Fatalf("validate error = %v", err)
	}
	if errors.Is(err, ErrValidation) {
		t.Fatalf("fresh initialized structure should have warnings, not errors: %s", out.String())
	}

	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); err != nil {
		t.Fatalf("AGENTS.md missing: %v", err)
	}
}
