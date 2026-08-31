package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aleka7sk/auto-company/internal/model"
)

func TestInitializeGreenfieldAndIdempotentManagedBlocks(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	now := time.Date(2026, 8, 31, 10, 0, 0, 0, time.UTC)
	opts := Options{
		Target:      target,
		ProfileID:   "fullstack-saas",
		ProjectName: "Booking OS",
		Idea:        "Reduce lost booking leads",
		Agent:       "both",
		Now:         now,
	}

	result, err := Initialize(opts)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if result.Mode != "greenfield" {
		t.Fatalf("Mode = %q, want greenfield", result.Mode)
	}

	manifestData, err := os.ReadFile(filepath.Join(target, ".auto-company", "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest model.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.Project.Name != "Booking OS" || manifest.Project.Profile != "fullstack-saas" {
		t.Fatalf("unexpected manifest project: %+v", manifest.Project)
	}
	if manifest.Governance.AllowAutoMerge || manifest.Governance.AllowProductionRun {
		t.Fatalf("unsafe defaults: %+v", manifest.Governance)
	}

	briefPath := filepath.Join(target, ".auto-company", "product", "product-brief.md")
	brief, err := os.ReadFile(briefPath)
	if err != nil {
		t.Fatalf("read product brief: %v", err)
	}
	if !strings.Contains(string(brief), "Reduce lost booking leads") {
		t.Fatalf("product brief does not contain idea")
	}

	custom := append(brief, []byte("\nCUSTOM OWNER NOTE\n")...)
	if err := os.WriteFile(briefPath, custom, 0o644); err != nil {
		t.Fatalf("write custom brief: %v", err)
	}

	second, err := Initialize(opts)
	if err != nil {
		t.Fatalf("second Initialize() error = %v", err)
	}
	if len(second.Skipped) == 0 {
		t.Fatalf("second Initialize() did not preserve generated artifacts")
	}
	preserved, err := os.ReadFile(briefPath)
	if err != nil {
		t.Fatalf("read preserved brief: %v", err)
	}
	if !strings.Contains(string(preserved), "CUSTOM OWNER NOTE") {
		t.Fatalf("existing product artifact was overwritten")
	}

	for _, file := range []string{"CLAUDE.md", "AGENTS.md"} {
		data, err := os.ReadFile(filepath.Join(target, file))
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		if count := strings.Count(string(data), managedStart); count != 1 {
			t.Fatalf("%s managed block count = %d, want 1", file, count)
		}
	}
}

func TestInitializeDetectsBrownfield(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "package.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Initialize(Options{
		Target: target, ProfileID: "saas-web", ProjectName: "Existing",
		Idea: "Improve an existing product", Agent: "claude",
	})
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if result.Mode != "brownfield" {
		t.Fatalf("Mode = %q, want brownfield", result.Mode)
	}
	if _, err := os.Stat(filepath.Join(target, "package.json")); err != nil {
		t.Fatalf("existing project file was affected: %v", err)
	}
}

func TestInitializeReusesExistingManifestForManagedBlocks(t *testing.T) {
	t.Parallel()

	target := t.TempDir()
	if _, err := Initialize(Options{
		Target:      target,
		ProfileID:   "saas-web",
		ProjectName: "Original SaaS",
		Idea:        "Original customer problem",
		Agent:       "claude",
	}); err != nil {
		t.Fatalf("first Initialize() error = %v", err)
	}

	result, err := Initialize(Options{
		Target:      target,
		ProfileID:   "expo-mobile",
		ProjectName: "Conflicting Project",
		Idea:        "Conflicting idea",
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("second Initialize() error = %v", err)
	}
	if result.Profile.ID != "saas-web" {
		t.Fatalf("Profile = %q, want existing saas-web", result.Profile.ID)
	}

	manifestData, err := os.ReadFile(filepath.Join(target, ".auto-company", "manifest.json"))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest model.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	if manifest.Project.Name != "Original SaaS" || manifest.Project.Profile != "saas-web" || manifest.Project.TargetAgent != "claude" {
		t.Fatalf("existing manifest changed unexpectedly: %+v", manifest.Project)
	}

	for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
		content, err := os.ReadFile(filepath.Join(target, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		text := string(content)
		if !strings.Contains(text, "Project: **Original SaaS**") || !strings.Contains(text, "Profile: `saas-web`") {
			t.Fatalf("%s managed block does not use existing manifest: %s", name, text)
		}
		if strings.Contains(text, "Conflicting Project") || strings.Contains(text, "expo-mobile") {
			t.Fatalf("%s managed block contains conflicting re-init arguments: %s", name, text)
		}
	}
}

func TestLoadProfileRejectsTraversalAndUnknownProfile(t *testing.T) {
	t.Parallel()

	for _, id := range []string{"../secret", `..\\secret`, "missing-profile"} {
		if _, err := LoadProfile(id); err == nil {
			t.Fatalf("LoadProfile(%q) returned nil error", id)
		}
	}
}
