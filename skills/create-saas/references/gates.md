# Gate policy

| Gate | Required evidence | Who may pass it |
|---|---|---|
| Product | Product brief, feature contract, metrics, risky assumptions | Owner for material scope; agent for reversible details |
| Architecture | Current-source research, three options, ADR, failure and rollback plan | Owner/architect for material data or platform decisions |
| UX | Journey, information architecture, states, accessibility, selected direction | Owner for direction; agent for component details |
| Implementation | Traceable tasks, commands, migration and rollback | Agent when no material decisions remain |
| Verification | Clean checks plus demonstrated primary and failure flows | Independent reviewer + objective tools |
| Production | Observability, backup, rollout, rollback, runbook, owner approval | Owner only |

Do not weaken a gate to make the status green. Mark the work blocked and explain exactly what is missing.
