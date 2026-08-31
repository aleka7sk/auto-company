# Token and context policy

## Goals

- spend context on the current decision;
- preserve durable knowledge without replaying entire conversations;
- prevent uncontrolled subagent and review loops;
- make expensive research explicit;
- keep behavior reproducible across sessions.

## Context layers

### Always loaded

Keep this layer small:

- project identity and profile;
- source-of-truth precedence;
- critical security constraints;
- active feature and current phase;
- canonical test commands.

### Loaded on demand

- skill procedures;
- detailed references;
- relevant Feature Contract sections;
- the accepted ADR;
- the UX section for the current slice;
- focused code search results and file excerpts.

### Isolated specialist context

Use a subagent for:

- broad market research;
- comparing architecture alternatives;
- visual/UX critique;
- cumulative release review.

Return a compact evidence artifact, not the full subagent transcript.

## Default budgets

```yaml
research_agents: 0..2
parallel_implementation_agents: 0..2 only for independent slices
review_fix_cycles: 2
architecture_options: 3 realistic options
full_repository_reads: forbidden by default
```

The limits are workflow constraints, not API billing controls. API/runtime implementations should also enforce per-run monetary and turn budgets.

## Task routing

### Tiny

Examples: copy change, narrow style fix, small deterministic bug.

```text
inspect → implement → targeted test → diff review
```

No product discovery or multi-agent review unless risk changes.

### Normal

```text
compact contract → plan → implement → checks → review
```

### Risky

Examples: auth, tenancy, payments, personal data, migrations, concurrency, audio/video uploads, production infrastructure.

```text
contract → architecture options → approval → isolated implementation → deterministic checks → independent review → owner acceptance
```

## Compaction by artifacts

At each phase, write the accepted result to disk and discard obsolete working discussion:

```text
research transcript → product-brief.md
architecture discussion → ADR.md
UX exploration → ux-spec.md + design system
implementation conversation → commits + evidence
```

## Review admission

The first meaningful review is required. A repeated review is admitted only when:

- code changed in response to a finding;
- a failing verifier became green;
- a new risk class was discovered;
- the cumulative change crosses a security/data/public-API boundary.

Otherwise defer to the next meaningful boundary instead of spending another full pass.

## Cost visibility roadmap

The future runtime should record per run:

- model/provider/version;
- requests and tool calls;
- input, cached input, reasoning, and output tokens;
- estimated cost;
- elapsed time;
- retries and review cycles;
- produced artifacts and final status.
