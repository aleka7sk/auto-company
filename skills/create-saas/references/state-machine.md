# Auto Company state machine

```text
UNINITIALIZED
  -> DISCOVERY
  -> PRODUCT_DECISION
  -> ARCHITECTURE_DECISION
  -> UX_READY
  -> IMPLEMENTATION_READY
  -> IMPLEMENTING
  -> VERIFYING
  -> READY_FOR_ACCEPTANCE
  -> RELEASE_APPROVED
  -> RELEASED
```

Exceptional states:

- `BLOCKED`: missing dependency, unresolved high-impact decision, repeated failed verification, or budget/context limit.
- `REWORK_REQUIRED`: objective evidence contradicts an approved artifact.
- `SUPERSEDED`: a newer approved artifact replaces this work.
- `REJECTED`: owner decides not to continue.

State transitions require artifacts, not conversational claims.
