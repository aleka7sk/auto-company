# Runbook: existing project

## Goal

Add Auto Company without replacing existing product truth or allowing old concept documents to override working behavior.

## 1. Create a clean branch

```bash
git status
git switch -c chore/auto-company-bootstrap
```

## 2. Initialize in brownfield mode

```powershell
autoco init `
  --target . `
  --profile expo-mobile `
  --name "Belcanto Product" `
  --idea "Improve the current working student experience without reviving superseded concepts"
```

The CLI detects existing repository content as brownfield and preserves generated artifacts unless `--force` is supplied.

## 3. Edit source-of-truth rules first

In `.auto-company/product/product-brief.md` and `feature-contract.md`, record:

- current working behavior;
- approved active baseline;
- superseded documents and concepts;
- compatibility constraints;
- real users and data already present;
- migration/non-regression expectations.

## 4. Map the existing system

Before proposing architecture or UI, inspect:

- repository structure and build commands;
- current domain/API/data boundaries;
- navigation and user flows;
- authentication/authorization;
- design tokens/components;
- tests, CI, deployment, observability;
- known defects and technical debt relevant to the requested change.

Do not perform a full-repository rewrite or “modernization” unless separately approved.

## 5. Validate incremental scope

Each new slice must state:

- behavior preserved;
- behavior changed;
- data migration;
- feature flag/rollout;
- regression tests;
- rollback.

## 6. Review bootstrap diff

```bash
git diff -- .auto-company CLAUDE.md AGENTS.md .gitignore
autoco validate --target .
```

Commit the control artifacts before feature implementation.
