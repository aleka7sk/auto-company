---
name: create-saas
description: Orchestrate a new or existing SaaS or mobile product from an idea through product discovery, architecture, UX, implementation slices, and production evidence. Use only when the user explicitly asks to create a product end to end.
disable-model-invocation: true
argument-hint: "<product idea or change request>"
---

# Create a product with Auto Company

Treat `$ARGUMENTS` as the initial hypothesis, not as a complete specification.

## Operating mode

Work as a small, controlled product company. Do not simulate dozens of agents. Use the main context as orchestrator, delegate noisy read-only research to specialized subagents, and persist conclusions in `.auto-company/` after every phase.

Read first:

1. `CLAUDE.md` and/or `AGENTS.md`.
2. `.auto-company/manifest.json` and existing artifacts.
3. `references/state-machine.md` and `references/gates.md` in this skill.
4. The current repository structure and tests.

If the repository is not initialized, use the templates packaged with Auto Company or ask the user to run:

```text
autoco init --profile fullstack-saas --idea "$ARGUMENTS"
```

Do not delete or replace existing product documents. In a brownfield product, current working behavior and accepted decisions outrank generic framework advice.

## Workflow

### 1. Understand and challenge the product

Invoke the `auto-company:product-discovery` skill. Establish the real user problem, evidence, current workaround, smallest valuable wedge, business model hypothesis, risks, and measurable outcomes.

Do not ask the owner to design the whole product. Make a recommendation and surface only decisions whose consequences are material or hard to reverse.

### 2. Decide the architecture

Invoke `auto-company:architecture-decision`. Inspect the existing codebase first. Use current primary sources for external technical facts. Compare at least three realistic approaches and write an ADR with failure modes, operability, migration, and rollback.

### 3. Design the experience

Invoke `auto-company:experience-design`. Define journey, information architecture, screen inventory, state matrix, accessibility, and a coherent visual direction before UI implementation.

### 4. Prepare executable slices

Invoke `auto-company:delivery-plan`. Map every task and test to requirement IDs. Prefer vertical slices that produce demonstrable user value and leave the repository working.

### 5. Owner decision packet

Before major implementation, present one compact packet containing:

- recommended product wedge;
- architecture decision and alternatives;
- selected UX direction;
- material open decisions;
- first delivery slice;
- cost, risk, and dependency implications.

Proceed automatically on reversible implementation details. Pause only for business model, money, permissions, sensitive data, irreversible migration, material external dependency, or production release decisions.

### 6. Implement one production slice at a time

Use project-specific technical skills where installed. Do not allow BMAD, Spec Kit, Superpowers, gstack, or another framework to become a competing orchestrator; they are specialist libraries under this workflow.

For each slice:

1. Confirm requirement IDs.
2. Inspect existing patterns.
3. Add failing tests or an equivalent objective reproduction.
4. Implement the smallest complete behavior.
5. Run targeted and broad checks.
6. Update evidence and risks.
7. Stop scope expansion unless upstream artifacts are updated.

### 7. Verify independently

Invoke `auto-company:production-readiness`. A claim of completion without executable evidence is invalid.

Never merge, force-push, publish, submit to stores, delete production data, or deploy production without explicit owner approval.

## Context and token discipline

- Maximum two research subagents by default.
- Maximum two implementation agents, and only for independent work.
- Maximum two review/fix cycles before marking the task blocked.
- Do not preload every external skill or the whole repository.
- Use fresh subagent context for large searches and return structured summaries.
- After each phase, persist the decision and discard conversational detail.
