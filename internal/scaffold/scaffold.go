package scaffold

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/aleka7sk/auto-company/internal/assets"
	"github.com/aleka7sk/auto-company/internal/model"
)

const (
	managedStart = "<!-- AUTO-COMPANY:START -->"
	managedEnd   = "<!-- AUTO-COMPANY:END -->"
	gitStart     = "# AUTO-COMPANY:START"
	gitEnd       = "# AUTO-COMPANY:END"
)

type Options struct {
	Target      string
	ProfileID   string
	ProjectName string
	Idea        string
	Agent       string
	Force       bool
	Now         time.Time
}

type Result struct {
	Target  string
	Profile model.Profile
	Mode    string
	Created []string
	Skipped []string
}

type templateData struct {
	ProjectName string
	Idea        string
	ProfileID   string
	ProfileName string
	Mode        string
	Agent       string
	Date        string
}

func Initialize(opts Options) (Result, error) {
	var result Result
	if opts.Target == "" {
		opts.Target = "."
	}
	if opts.ProfileID == "" {
		opts.ProfileID = "fullstack-saas"
	}
	if opts.Agent == "" {
		opts.Agent = "both"
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}

	absTarget, err := filepath.Abs(opts.Target)
	if err != nil {
		return result, fmt.Errorf("resolve target: %w", err)
	}
	if err := os.MkdirAll(absTarget, 0o755); err != nil {
		return result, fmt.Errorf("create target: %w", err)
	}

	profile, err := LoadProfile(opts.ProfileID)
	if err != nil {
		return result, err
	}
	if strings.TrimSpace(opts.ProjectName) == "" {
		opts.ProjectName = filepath.Base(absTarget)
	}
	if strings.TrimSpace(opts.Idea) == "" {
		opts.Idea = "[TODO: describe the customer problem and desired outcome]"
	}

	mode, err := detectMode(absTarget)
	if err != nil {
		return result, err
	}
	data := templateData{
		ProjectName: opts.ProjectName,
		Idea:        opts.Idea,
		ProfileID:   profile.ID,
		ProfileName: profile.DisplayName,
		Mode:        mode,
		Agent:       opts.Agent,
		Date:        opts.Now.Format("2006-01-02"),
	}

	manifest := model.Manifest{
		SchemaVersion: "1.0",
		Project: model.Project{
			Name:        opts.ProjectName,
			Idea:        opts.Idea,
			Profile:     profile.ID,
			Mode:        mode,
			TargetAgent: opts.Agent,
		},
		Governance: model.Governance{
			ApprovalMode:       "strategic-only",
			SourceOfTruth:      []string{".auto-company/manifest.json", ".auto-company/product/feature-contract.md", ".auto-company/architecture/ADR-0001.md", ".auto-company/ux/ux-spec.md", ".auto-company/delivery/implementation-plan.md"},
			MaxResearchAgents:  2,
			MaxReviewCycles:    2,
			AllowAutoMerge:     false,
			AllowProductionRun: false,
		},
		State:     model.State{Phase: "discovery", Status: "initialized"},
		CreatedAt: opts.Now,
		UpdatedAt: opts.Now,
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return result, fmt.Errorf("encode manifest: %w", err)
	}
	profileBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return result, fmt.Errorf("encode profile: %w", err)
	}

	staticFiles := map[string][]byte{
		".auto-company/manifest.json":                append(manifestBytes, '\n'),
		".auto-company/profile.json":                 append(profileBytes, '\n'),
		".auto-company/quality-gates.json":           mustRead("templates/quality-gates.json.tmpl"),
		".auto-company/claude-settings.example.json": mustRead("templates/claude-settings.example.json.tmpl"),
	}
	for path, content := range staticFiles {
		if err := writeManaged(absTarget, path, content, opts.Force, &result); err != nil {
			return result, err
		}
	}

	templates := map[string]string{
		".auto-company/README.md":                       "templates/project-readme.md.tmpl",
		".auto-company/product/product-brief.md":        "templates/product-brief.md.tmpl",
		".auto-company/product/feature-contract.md":     "templates/feature-contract.md.tmpl",
		".auto-company/architecture/ADR-0001.md":        "templates/architecture-decision.md.tmpl",
		".auto-company/ux/ux-spec.md":                   "templates/ux-spec.md.tmpl",
		".auto-company/delivery/implementation-plan.md": "templates/implementation-plan.md.tmpl",
		".auto-company/evidence/release-evidence.md":    "templates/release-evidence.md.tmpl",
		".auto-company/prompts/start.md":                "templates/start-prompt.md.tmpl",
	}
	paths := make([]string, 0, len(templates))
	for path := range templates {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		content, err := render(templates[path], data)
		if err != nil {
			return result, err
		}
		if err := writeManaged(absTarget, path, content, opts.Force, &result); err != nil {
			return result, err
		}
	}

	instructions, err := render("templates/agent-instructions.md.tmpl", data)
	if err != nil {
		return result, err
	}
	for _, path := range []string{"CLAUDE.md", "AGENTS.md"} {
		if err := upsertBlock(filepath.Join(absTarget, path), managedStart, managedEnd, string(instructions)); err != nil {
			return result, fmt.Errorf("update %s: %w", path, err)
		}
		result.Created = append(result.Created, path+" (managed block)")
	}

	gitIgnore := strings.Join([]string{
		"# Local Auto Company execution artifacts",
		".auto-company/runs/",
		".auto-company/local/",
		"",
		"# Secrets",
		".env",
		".env.*",
		"!.env.example",
	}, "\n")
	if err := upsertBlock(filepath.Join(absTarget, ".gitignore"), gitStart, gitEnd, gitIgnore); err != nil {
		return result, fmt.Errorf("update .gitignore: %w", err)
	}
	result.Created = append(result.Created, ".gitignore (managed block)")

	result.Target = absTarget
	result.Profile = profile
	result.Mode = mode
	sort.Strings(result.Created)
	sort.Strings(result.Skipped)
	return result, nil
}

func LoadProfile(id string) (model.Profile, error) {
	var profile model.Profile
	if strings.ContainsAny(id, `/\\`) || id == "" {
		return profile, fmt.Errorf("invalid profile id %q", id)
	}
	data, err := assets.FS.ReadFile("profiles/" + id + ".json")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return profile, fmt.Errorf("unknown profile %q", id)
		}
		return profile, fmt.Errorf("read profile %q: %w", id, err)
	}
	if err := json.Unmarshal(data, &profile); err != nil {
		return profile, fmt.Errorf("parse profile %q: %w", id, err)
	}
	return profile, nil
}

func ListProfiles() ([]model.Profile, error) {
	entries, err := fs.ReadDir(assets.FS, "profiles")
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	profiles := make([]model.Profile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		profile, err := LoadProfile(strings.TrimSuffix(entry.Name(), ".json"))
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].ID < profiles[j].ID })
	return profiles, nil
}

func render(path string, data templateData) ([]byte, error) {
	content, err := assets.FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", path, err)
	}
	tmpl, err := template.New(filepath.Base(path)).Option("missingkey=error").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", path, err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return nil, fmt.Errorf("render template %s: %w", path, err)
	}
	return out.Bytes(), nil
}

func mustRead(path string) []byte {
	content, err := assets.FS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func writeManaged(root, relative string, content []byte, force bool, result *Result) error {
	path := filepath.Join(root, filepath.FromSlash(relative))
	if _, err := os.Stat(path); err == nil && !force {
		result.Skipped = append(result.Skipped, relative)
		return nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect %s: %w", relative, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent for %s: %w", relative, err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relative, err)
	}
	result.Created = append(result.Created, relative)
	return nil
}

func upsertBlock(path, start, end, body string) error {
	content, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	current := string(content)
	block := start + "\n" + strings.TrimSpace(body) + "\n" + end
	startIndex := strings.Index(current, start)
	endIndex := strings.Index(current, end)
	if startIndex >= 0 && endIndex >= startIndex {
		endIndex += len(end)
		current = current[:startIndex] + block + current[endIndex:]
	} else {
		if strings.TrimSpace(current) != "" {
			current = strings.TrimRight(current, "\r\n") + "\n\n"
		}
		current += block + "\n"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(current), 0o644)
}

func detectMode(target string) (string, error) {
	entries, err := os.ReadDir(target)
	if err != nil {
		return "", fmt.Errorf("inspect target: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" || name == ".auto-company" || name == ".gitignore" || name == "CLAUDE.md" || name == "AGENTS.md" {
			continue
		}
		return "brownfield", nil
	}
	return "greenfield", nil
}
