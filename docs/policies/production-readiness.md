# Production readiness policy

## Definition

A feature is production-ready only when evidence supports safe use under its intended scope. The label is not inferred from code generation, build success, or an agent's confidence.

## Required evidence categories

### Product and requirements

- Feature Contract approved.
- Delivered FR/BR/NFR/AC identifiers mapped to code and tests.
- Out-of-scope items remain out of scope.
- No superseded requirement was revived.

### Security and data

- authentication and authorization boundaries tested;
- tenant/organization/user data isolation tested when applicable;
- sensitive data classification and retention defined;
- secret handling, dependency risks, input validation, rate limits, and auditability reviewed;
- privacy/compliance questions recorded for qualified review where needed.

### Correctness and resilience

- unit, integration, contract, and critical-path E2E/device tests pass;
- retries, timeouts, idempotency, concurrency, partial failure, and degraded dependencies addressed where relevant;
- background jobs and webhooks are observable and replay-safe where relevant;
- migrations are backward-compatible or have a controlled maintenance strategy;
- backup and restore have evidence, not only documentation.

### Experience

- loading, empty, error, permission-denied, offline/degraded, long-content, localization, and responsive/native states handled as applicable;
- keyboard, screen-reader, contrast, focus, dynamic text, and reduced-motion requirements considered;
- analytics/dashboard visuals use truthful scales, labels, definitions, and empty/error states;
- real primary flow and a meaningful failure flow were exercised on target browser/device.

### Operations

- structured logs exclude secrets and unnecessary personal data;
- metrics, traces, dashboards, alerts, ownership, and runbook exist for critical behavior;
- preview/staging evidence exists;
- rollout is staged or feature-flagged when risk justifies it;
- rollback trigger and executable rollback procedure are documented;
- capacity/cost assumptions are recorded.

## Severity policy

| Severity | Meaning | Release decision |
|---|---|---|
| Blocker | data loss, security boundary failure, unusable core flow, irreversible unsafe operation | prohibited |
| High | likely major customer or operational harm | prohibited unless formally risk-accepted by owner and qualified reviewer |
| Medium | meaningful but bounded defect or maintainability risk | owner decision with explicit record |
| Low | polish or low-impact issue | may release with backlog item |

## Required human decisions

AI cannot approve:

- production release;
- destructive migration/data deletion;
- material privacy/compliance interpretation;
- payment or pricing policy;
- security-risk acceptance;
- protected branch merge when repository policy requires review.

## Evidence packet

`.auto-company/evidence/release-evidence.md` must contain:

- commit/base/head identifiers;
- changed scope;
- requirement traceability;
- exact commands and results;
- screenshots/device/browser evidence where relevant;
- migration/rollout/rollback evidence;
- findings and remaining risks;
- reviewer recommendation;
- owner acceptance record.
