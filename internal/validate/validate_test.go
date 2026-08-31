package validate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aleka7sk/auto-company/internal/scaffold"
)

func TestRunAcceptsInitializedStructureWithWarnings(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	if _, err := scaffold.Initialize(scaffold.Options{
		Target: target, ProfileID: "fullstack-saas", ProjectName: "Test SaaS",
		Idea: "Solve a real problem", Agent: "both",
	}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	report, err := Run(target)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if report.HasErrors() {
		t.Fatalf("initialized project has structural errors: %+v", report.Findings)
	}
	if len(report.Findings) == 0 {
		t.Fatalf("expected TODO warnings in fresh templates")
	}
}

func TestRunReportsMissingRequiredArtifact(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	if _, err := scaffold.Initialize(scaffold.Options{
		Target: target, ProfileID: "expo-mobile", ProjectName: "Mobile",
		Idea: "Help users practice", Agent: "claude",
	}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	missing := filepath.Join(target, ".auto-company", "ux", "ux-spec.md")
	if err := os.Remove(missing); err != nil {
		t.Fatal(err)
	}

	report, err := Run(target)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !report.HasErrors() {
		t.Fatalf("missing required artifact did not produce an error: %+v", report.Findings)
	}
}

func TestRunReportsUninitializedProject(t *testing.T) {
	t.Parallel()

	report, err := Run(t.TempDir())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !report.HasErrors() {
		t.Fatalf("uninitialized project did not produce an error")
	}
}
