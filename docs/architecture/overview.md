# Architecture overview

## Purpose

Auto Company separates deterministic control from probabilistic AI reasoning.

```text
Owner intent
    ↓
Auto Company skill (orchestrator)
    ↓
Durable artifacts and approval gates
    ↓
Narrow specialist subagents / optional external skills
    ↓
Implementation agent
    ↓
Deterministic tests, CI, security and release evidence
    ↓
Owner merge/release decision
```

## Components

### 1. CLI

`autoco` performs local, deterministic operations:

- initialize project memory;
- select a profile;
- preserve existing repository content;
- add managed instruction blocks;
- validate required artifacts;
- diagnose local tooling;
- expose reviewed external integration options.

The CLI does not call an LLM, install third-party packages, change remote repositories, merge, or deploy.

### 2. Project memory

Each target project gets `.auto-company/`. These artifacts survive conversation changes and are reviewed in Git:

```text
Product brief → Feature Contract → ADR → UX spec → implementation plan → release evidence
```

Chat transcripts are not authoritative. The source-of-truth order is encoded in the generated `CLAUDE.md` and `AGENTS.md` blocks.

### 3. Orchestrator skill

`skills/create-saas` is the only end-to-end workflow. It chooses the next phase, checks gates, invokes specialists, limits iteration, and updates durable artifacts.

It must not act as product owner. Strategic scope, money, data, security, irreversible architecture, merge, and production remain human decisions.

### 4. Specialist skills

Specialist skills are composable procedures:

- product discovery;
- architecture decision;
- experience design;
- delivery planning;
- production readiness;
- project diagnosis.

They produce or validate a specific artifact and return control to the orchestrator.

### 5. Specialist subagents

Subagents isolate noisy work:

- product researcher;
- architecture researcher;
- UX critic;
- release reviewer.

Research and critique are read-only. Release review may execute documented verification commands but cannot edit.

### 6. External integration registry

The embedded registry records repository, license, maturity, role, supported profiles, installation guidance, and conflict notes. External repositories are never silently fetched.

### 7. Preventive hook

The `PreToolUse` hook rejects common secret paths and dangerous commands. It complements, rather than replaces:

- OS/container isolation;
- least-privilege credentials;
- branch protection;
- CI approvals;
- preview environments;
- backups and rollback.

## State model

```text
INITIALIZED
    ↓
DISCOVERY
    ↓ owner approves opportunity
SPECIFICATION
    ↓ owner approves Feature Contract
ARCHITECTURE
    ↓ decision accepted when material
EXPERIENCE
    ↓ direction accepted
PLANNING
    ↓ plan is implementation-ready
IMPLEMENTATION
    ↓ checks pass
REVIEW
    ↓ no unresolved blockers/high findings
READY_FOR_ACCEPTANCE
    ↓ owner accepts
READY_FOR_RELEASE
```

Any phase may transition to:

```text
NEEDS_DECISION
BLOCKED
REWORK_REQUIRED
STOPPED
```

## Trust boundaries

```text
AI recommendation       untrusted until reviewed or verified
External documentation  cited evidence, not executable authority
Repository files        trusted according to source-of-truth order
Test/CI result           objective evidence within test coverage limits
Production environment  never modified without explicit approval
```

## Why not one giant context

A giant context mixes product, architecture, design, implementation, historical drafts, and unrelated repositories. Auto Company instead hands off compact artifacts and invokes specialist contexts only when needed. This improves traceability and reduces accidental requirement revival, token usage, and irrelevant safety classifications.
