# Runbook: first greenfield project

## 1. Define the narrow problem

Use one sentence with a real actor, situation, current pain, and desired outcome. Avoid listing screens or technologies.

Example:

```text
Property owners receive non-guaranteed booking leads from several channels,
but sales managers respond inconsistently, so qualified demand is lost and
the owner cannot see where conversion breaks.
```

## 2. Initialize

```powershell
autoco init `
  --target C:\Users\aleka\projects\booking-os `
  --profile fullstack-saas `
  --name "Booking OS" `
  --idea "Property owners lose booking leads because follow-up is fragmented"
```

## 3. Diagnose

```powershell
cd C:\Users\aleka\projects\booking-os
autoco doctor --target .
autoco validate --target .
```

Warnings about TODO markers are expected before discovery is complete. Errors are not.

## 4. Start the plugin

```powershell
claude --plugin-dir C:\Users\aleka\projects\auto-company
```

Invoke:

```text
/auto-company:create-saas
```

Paste or reference `.auto-company/prompts/start.md`.

## 5. Approve only meaningful gates

You should approve:

- whether the problem is worth pursuing;
- the smallest valuable product scope;
- material data/security/architecture decisions;
- UX direction;
- merge/release.

Do not spend time choosing internal function names or routine component structure.

## 6. Install specialists only when needed

```powershell
autoco integrations --profile fullstack-saas --agent claude
```

Review and pin the selected repository. Do not enable a second complete lifecycle.

## 7. Build one vertical slice

A first slice should include the minimum frontend, backend, data, permissions, tests, observability, and demo path needed to validate one core value event.

## 8. Require evidence

Do not accept “done.” Require exact commands, results, diff, screenshots or device evidence, known limitations, rollout, and rollback.
