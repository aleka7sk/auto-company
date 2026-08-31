package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aleka7sk/auto-company/internal/doctor"
	"github.com/aleka7sk/auto-company/internal/model"
	"github.com/aleka7sk/auto-company/internal/registry"
	"github.com/aleka7sk/auto-company/internal/scaffold"
	"github.com/aleka7sk/auto-company/internal/validate"
)

const Version = "0.1.0"

var ErrValidation = errors.New("validation failed")

func Run(args []string, out, errOut io.Writer) error {
	if out == nil {
		out = io.Discard
	}
	if errOut == nil {
		errOut = io.Discard
	}
	if len(args) == 0 {
		printUsage(out)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(out)
		return nil
	case "version", "--version", "-v":
		fmt.Fprintf(out, "autoco %s\n", Version)
		return nil
	case "init":
		return runInit(args[1:], out, errOut)
	case "doctor":
		return runDoctor(args[1:], out, errOut)
	case "validate":
		return runValidate(args[1:], out, errOut)
	case "prompt":
		return runPrompt(args[1:], out, errOut)
	case "profiles":
		return runProfiles(out)
	case "integrations":
		return runIntegrations(args[1:], out, errOut)
	default:
		printUsage(errOut)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runInit(args []string, out, errOut io.Writer) error {
	set := flag.NewFlagSet("init", flag.ContinueOnError)
	set.SetOutput(errOut)
	target := set.String("target", ".", "project directory")
	profile := set.String("profile", "fullstack-saas", "delivery profile")
	name := set.String("name", "", "project name (defaults to target directory name)")
	idea := set.String("idea", "", "customer problem or product idea")
	agent := set.String("agent", "both", "claude, codex, or both")
	force := set.Bool("force", false, "replace generated artifacts; managed instruction blocks are always updated")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *agent != "claude" && *agent != "codex" && *agent != "both" {
		return fmt.Errorf("unsupported agent %q", *agent)
	}
	result, err := scaffold.Initialize(scaffold.Options{
		Target: *target, ProfileID: *profile, ProjectName: *name,
		Idea: *idea, Agent: *agent, Force: *force,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Auto Company initialized\n")
	fmt.Fprintf(out, "  target:  %s\n", result.Target)
	fmt.Fprintf(out, "  mode:    %s\n", result.Mode)
	fmt.Fprintf(out, "  profile: %s\n", result.Profile.ID)
	fmt.Fprintf(out, "  created/updated: %d\n", len(result.Created))
	if len(result.Skipped) > 0 {
		fmt.Fprintf(out, "  preserved existing artifacts: %d\n", len(result.Skipped))
	}
	fmt.Fprintln(out, "\nNext:")
	fmt.Fprintln(out, "  1. Review .auto-company/prompts/start.md")
	fmt.Fprintln(out, "  2. Start Claude Code or Codex in the project root")
	fmt.Fprintln(out, "  3. Invoke /auto-company:create-saas <your idea> in Claude Code")
	return nil
}

func runDoctor(args []string, out, errOut io.Writer) error {
	set := flag.NewFlagSet("doctor", flag.ContinueOnError)
	set.SetOutput(errOut)
	target := set.String("target", ".", "project directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	checks, err := doctor.Run(*target)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Auto Company doctor")
	for _, check := range checks {
		status := "OK"
		if !check.Found {
			status = "MISSING"
		}
		required := "optional"
		if check.Required {
			required = "required"
		}
		fmt.Fprintf(out, "  %-8s %-24s %-8s %s\n", status, check.Name, required, check.Detail)
	}
	if doctor.HasMissingRequired(checks) {
		return errors.New("required tools are missing")
	}
	return nil
}

func runValidate(args []string, out, errOut io.Writer) error {
	set := flag.NewFlagSet("validate", flag.ContinueOnError)
	set.SetOutput(errOut)
	target := set.String("target", ".", "project directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	report, err := validate.Run(*target)
	if err != nil {
		return err
	}
	if len(report.Findings) == 0 {
		fmt.Fprintln(out, "PASS: Auto Company project contract is structurally valid")
		return nil
	}
	for _, finding := range report.Findings {
		fmt.Fprintf(out, "%s: %s — %s\n", strings.ToUpper(finding.Level), finding.Path, finding.Message)
	}
	if report.HasErrors() {
		return ErrValidation
	}
	return nil
}

func runPrompt(args []string, out, errOut io.Writer) error {
	set := flag.NewFlagSet("prompt", flag.ContinueOnError)
	set.SetOutput(errOut)
	target := set.String("target", ".", "project directory")
	if err := set.Parse(args); err != nil {
		return err
	}
	abs, err := filepath.Abs(*target)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(filepath.Join(abs, ".auto-company", "prompts", "start.md"))
	if err != nil {
		return fmt.Errorf("read start prompt: %w; run `autoco init` first", err)
	}
	_, err = out.Write(content)
	return err
}

func runProfiles(out io.Writer) error {
	profiles, err := scaffold.ListProfiles()
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		fmt.Fprintf(out, "%-18s %s\n", profile.ID, profile.Description)
	}
	return nil
}

func runIntegrations(args []string, out, errOut io.Writer) error {
	set := flag.NewFlagSet("integrations", flag.ContinueOnError)
	set.SetOutput(errOut)
	profile := set.String("profile", "fullstack-saas", "delivery profile")
	agent := set.String("agent", "claude", "claude or codex")
	jsonOutput := set.Bool("json", false, "print machine-readable JSON")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *agent != "claude" && *agent != "codex" {
		return fmt.Errorf("unsupported agent %q", *agent)
	}
	if _, err := scaffold.LoadProfile(*profile); err != nil {
		return err
	}
	items, err := registry.Filter(*profile, *agent)
	if err != nil {
		return err
	}
	if *jsonOutput {
		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(data))
		return nil
	}
	fmt.Fprintf(out, "Curated integrations for profile=%s agent=%s\n", *profile, *agent)
	fmt.Fprintln(out, "Auto Company does not install third-party code automatically. Review and pin before team rollout.")
	fmt.Fprintln(out)
	for _, item := range items {
		label := "optional"
		if item.Recommended {
			label = "recommended"
		}
		fmt.Fprintf(out, "[%s] %s — %s\n", label, item.Name, item.Role)
		fmt.Fprintf(out, "  repo: %s | license: %s | maturity: %s\n", item.Repository, item.License, item.Maturity)
		for _, command := range item.Install[*agent] {
			fmt.Fprintf(out, "  %s\n", command)
		}
		if item.Notes != "" {
			fmt.Fprintf(out, "  note: %s\n", item.Notes)
		}
		fmt.Fprintln(out)
	}
	return nil
}

func printUsage(out io.Writer) {
	fmt.Fprintln(out, `Auto Company — controlled AI product delivery

Usage:
  autoco init [flags]          Initialize a new or existing product repository
  autoco doctor [flags]        Check local tools and project setup
  autoco validate [flags]      Validate required product-delivery artifacts
  autoco prompt [flags]        Print the generated starting prompt
  autoco profiles              List product profiles
  autoco integrations [flags]  Show curated external skills/plugins
  autoco version               Print version

Examples:
  autoco init --target ./my-saas --profile fullstack-saas --name "My SaaS" --idea "Reduce manual booking follow-up"
  autoco doctor --target ./my-saas
  autoco integrations --profile expo-mobile --agent claude`)
}

// Keep model imported in generated documentation tools that build only selected files.
var _ model.Profile
