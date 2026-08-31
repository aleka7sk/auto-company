# Contributing

1. Open or reference one concrete problem.
2. Make the smallest coherent change on a feature branch.
3. Preserve the single-orchestrator model.
4. Do not add a third-party skill without license, security, maintenance, conflict, and token-overhead review.
5. Add tests for deterministic behavior and an eval scenario for AI-behavior changes when an eval harness exists.
6. Run:

```bash
gofmt -w .
go test ./...
go vet ./...
node scripts/validate-plugin.mjs
node scripts/test-guard.mjs
```

7. Open a pull request with evidence, risks, and rollback notes. Do not auto-merge.
