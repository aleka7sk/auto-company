---
name: project-doctor
description: Diagnose Auto Company project setup, plugin discovery, required tools, artifact structure, and conflicting orchestrators. Invoke explicitly when setup is broken or before the first product run.
disable-model-invocation: true
---

# Project doctor

Run read-only checks first:

```text
autoco doctor --target .
autoco validate --target .
claude plugin details auto-company
```

Then inspect:

1. `.auto-company/manifest.json` exists and its profile is valid.
2. `CLAUDE.md` and `AGENTS.md` contain exactly one Auto Company managed block.
3. The plugin exposes the expected skills and four specialist agents.
4. Node, npm/npx, Git, and Go are in `PATH`; Claude Code or Codex is available for the selected agent.
5. No full competing orchestrator is automatically forcing a second product lifecycle.
6. External integrations are reviewed, licensed appropriately, and pinned before shared/team use.
7. Secret and destructive-action guards are active.

Report exact commands and outputs. Do not “fix” missing tools by downloading or executing third-party installers without approval.
