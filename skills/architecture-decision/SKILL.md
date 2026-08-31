---
name: architecture-decision
description: Research and select a production-oriented architecture for an approved feature or product. Use after requirements are clear and before implementation of material technical decisions.
---

# Architecture decision

Do not implement production code while this skill is active.

1. Read the approved feature contract and inspect the existing repository, runtime, deployment, data model, tests, and operational constraints.
2. Delegate broad technical research to `auto-company:architecture-researcher` so raw search results do not pollute the main context.
3. Use current primary sources for claims about frameworks, cloud services, standards, security controls, limits, and compatibility.
4. Compare at least three realistic options. Include the simplest viable option; for a greenfield SaaS, a modular monolith is the default hypothesis, not a mandatory conclusion.
5. Evaluate product fit, delivery speed, team fit, total operational complexity, data consistency, security, failure behavior, cost, observability, migration, rollback, and reversibility.
6. Explicitly design tenant boundaries, authentication, authorization, idempotency, concurrency, retries, auditability, backups, and degradation where relevant.
7. Write or update `.auto-company/architecture/ADR-0001.md`. Cite sources, record access dates, and list uncertainties.
8. Stop for owner approval only when the decision affects money, sensitive data, permissions, irreversible migration, long-term platform lock-in, or production operations.

Never select an architecture because it is fashionable or because a skill recommended it without project evidence.
