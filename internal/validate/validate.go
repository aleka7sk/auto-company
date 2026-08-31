package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aleka7sk/auto-company/internal/model"
	"github.com/aleka7sk/auto-company/internal/scaffold"
)

type Finding struct {
	Level   string
	Path    string
	Message string
}

type Report struct {
	Findings []Finding
}

func (r Report) HasErrors() bool {
	for _, finding := range r.Findings {
		if finding.Level == "error" {
			return true
		}
	}
	return false
}

func Run(target string) (Report, error) {
	var report Report
	if target == "" {
		target = "."
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return report, fmt.Errorf("resolve target: %w", err)
	}
	manifestPath := filepath.Join(abs, ".auto-company", "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			report.Findings = append(report.Findings, Finding{"error", ".auto-company/manifest.json", "project is not initialized; run `autoco init`"})
			return report, nil
		}
		return report, fmt.Errorf("read manifest: %w", err)
	}
	var manifest model.Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		report.Findings = append(report.Findings, Finding{"error", ".auto-company/manifest.json", "invalid JSON: " + err.Error()})
		return report, nil
	}
	profile, err := scaffold.LoadProfile(manifest.Project.Profile)
	if err != nil {
		report.Findings = append(report.Findings, Finding{"error", ".auto-company/manifest.json", err.Error()})
		return report, nil
	}

	required := append([]string{}, profile.RequiredArtifacts...)
	required = append(required, "CLAUDE.md", "AGENTS.md")
	for _, relative := range required {
		path := filepath.Join(abs, filepath.FromSlash(relative))
		content, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				report.Findings = append(report.Findings, Finding{"error", relative, "required artifact is missing"})
				continue
			}
			return report, fmt.Errorf("read %s: %w", relative, err)
		}
		text := string(content)
		if strings.Contains(text, "[TODO") {
			report.Findings = append(report.Findings, Finding{"warning", relative, "contains unresolved TODO markers"})
		}
		if strings.Contains(text, "{{") || strings.Contains(text, "}}") {
			report.Findings = append(report.Findings, Finding{"warning", relative, "contains unresolved template markers"})
		}
	}
	if manifest.Governance.AllowAutoMerge {
		report.Findings = append(report.Findings, Finding{"warning", ".auto-company/manifest.json", "allowAutoMerge is enabled; keep it false until CI and branch protection are configured"})
	}
	if manifest.Governance.AllowProductionRun {
		report.Findings = append(report.Findings, Finding{"warning", ".auto-company/manifest.json", "allowProductionRun is enabled; production actions should require explicit owner approval"})
	}
	return report, nil
}
