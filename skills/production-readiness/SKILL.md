---
name: production-readiness
description: Independently verify an implementation against approved requirements and production gates using tests, runtime evidence, security checks, observability, rollout, and rollback. Use after a slice is implemented and before merge or release.
---

# Production readiness review

Act as an independent reviewer. Do not trust the implementer's summary. Read the feature contract, ADR, UX specification, implementation plan, git diff, and CI/test output.

Delegate objective repository inspection to `auto-company:release-reviewer` when the result would otherwise flood the main context.

Verify:

- every delivered requirement is traced to code and evidence;
- no approved rule was silently changed or omitted;
- authorization and tenant/data boundaries are tested;
- migrations are backward-compatible, recoverable, and safe under partial deployment;
- retries, idempotency, concurrency, timeouts, and degraded dependencies are handled where relevant;
- loading, empty, error, permission, offline/degraded, accessibility, responsive/native, and long-content UI states exist;
- logs avoid secrets and personal data; metrics, traces, alerts, backup, restore, rollout, rollback, and ownership are defined;
- the primary flow and at least one meaningful failure flow were run on the target environment/device.

Classify findings as blocker, high, medium, or low. Permit at most two fix/review cycles. After that, mark the slice `BLOCKED` with exact evidence.

Update `.auto-company/evidence/release-evidence.md`. Never merge or deploy production; present a release recommendation to the owner.
