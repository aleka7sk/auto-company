# Curated integrations

The machine-readable registry lives at `internal/assets/registry/integrations.json` and is exposed by:

```bash
autoco integrations --profile <profile> --agent <claude|codex>
```

## Roles

| Integration | Role | Default policy |
|---|---|---|
| GitHub Spec Kit | artifact flow, analysis, convergence | optional for SaaS; do not duplicate Auto Company source of truth |
| gstack | founder/engineering/design/QA critique | invoke selected roles only |
| Superpowers | TDD, debugging, review, worktrees, verification | selective by risk; avoid full ceremony for tiny tasks |
| UI UX Pro Max | visual exploration, UI patterns, charts | after user flow and information architecture |
| Interface Design | durable dashboard/app design decisions | save one accepted direction and reuse it |
| Expo Skills | current Expo/EAS practices | preferred specialist for Expo profile |
| Callstack Agent Skills | React Native performance/testing/device QA | complements Expo |
| Wondel Skills | selected product/UX/architecture frameworks | optional; install only needed skills |
| Product Manager Skills | deep PM reference | reference-only by default due non-commercial/share-alike license |

## Update process

1. Review upstream changes and release notes.
2. Inspect install scripts, hooks, MCP, telemetry, and permissions.
3. Test in a disposable repository.
4. Pin a tag or commit for shared use.
5. Update the registry notes and compatibility evidence.
6. Keep an uninstall path.
