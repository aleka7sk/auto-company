---
name: delivery-plan
description: Convert approved product, architecture, and UX artifacts into traceable vertical implementation slices with tests, migrations, rollout, and rollback. Use immediately before coding a significant feature.
---

# Delivery planning

1. Read all approved upstream artifacts and current code/tests.
2. Create a traceability row for every `BR`, `FR`, `NFR`, and `AC` that belongs to the selected slice.
3. Split work by user value and dependency, not by arbitrary frontend/backend departments. A slice should be demonstrable and leave the repository working.
4. Name exact modules/files only after inspecting the repository. Do not invent paths.
5. For each slice include data/API changes, permissions, UI states, tests, observability, migration, backward compatibility, preview/demo evidence, and rollback.
6. Identify tasks that can safely run in parallel. Keep dependent work sequential.
7. Define concrete commands for format, lint, type/build, unit, integration/contract, E2E/device, and security checks.
8. Update `.auto-company/delivery/implementation-plan.md`.

A plan is not implementation-ready when it contains unresolved strategic decisions, unverified file paths, generic “add tests” tasks, or no failure/rollback path.
