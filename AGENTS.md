# Auto Company repository instructions

## Mission

Build a controlled, evidence-driven product delivery toolkit. This repository is not an autonomous production deployer and must not claim that one prompt guarantees a production-ready product.

## Architecture

- `cmd/autoco` and `internal/*`: deterministic Go CLI.
- `internal/assets`: embedded project profiles, registry, and templates.
- `skills/*`: on-demand product delivery procedures.
- `agents/*`: narrow, mostly read-only specialist contexts.
- `hooks/*` and `scripts/guard-tool-use.mjs`: preventive safety layer.
- `docs/*`: human-readable governance and runbooks.

## Change rules

1. Keep Auto Company as the only top-level orchestrator.
2. External frameworks are optional specialists. Do not vendor them or silently install them.
3. Keep persistent instructions small; put procedural depth in skills/references.
4. Preserve existing target-project files. Only edit Auto Company managed blocks unless `--force` is explicit.
5. Never add automatic merge, force push, production deployment, store submission, destructive migration, or secret access.
6. New profile requirements must have deterministic validation and tests.
7. User-visible completion claims require actual command output.
8. Third-party additions require repository, license, role, maturity, install path, conflict notes, and pinning guidance in the registry.

## Required checks

```bash
gofmt -w .
go test ./...
go vet ./...
node scripts/validate-plugin.mjs
node scripts/test-guard.mjs
```
