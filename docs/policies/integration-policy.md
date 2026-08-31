# External integration policy

## Principle

Auto Company is the top-level orchestrator. An external framework is either:

1. a specialist capability;
2. a source of patterns to adapt;
3. an alternative orchestrator that must not be enabled at the same time.

## Admission checklist

An integration cannot be recommended for team use until the following are recorded:

- exact repository and maintainer;
- license and commercial-use implications;
- reviewed release tag or commit SHA;
- purpose and the phase where it is allowed;
- files, scripts, hooks, MCP servers, and network calls it introduces;
- telemetry behavior;
- required credentials and permissions;
- context/token overhead;
- overlap with Auto Company and other integrations;
- uninstall and rollback procedure;
- one successful sandbox evaluation.

## Single-orchestrator rule

Do not run the complete lifecycle of BMAD, Spec Kit, Superpowers, gstack, and Auto Company simultaneously.

Allowed example:

```text
Auto Company create-saas
  → gstack-style founder critique
  → UI UX Pro Max visual exploration
  → Expo official implementation guidance
  → Superpowers systematic debugging
  → Auto Company production readiness
```

Disallowed example:

```text
Auto Company full lifecycle
  + BMAD full lifecycle
  + Spec Kit full lifecycle
  + Superpowers mandatory lifecycle
```

The disallowed configuration creates duplicate PRDs, competing plans, conflicting approval gates, and unnecessary context.

## Pinning

Personal exploration may use a current release after review. Shared/team/automation use must pin a tag or commit and record it in project documentation. Updates require re-running security, behavior, and compatibility checks.

## Licensing

MIT/Apache-style skills may still have trademark, bundled-asset, or dependency restrictions. Non-commercial or share-alike content must not be copied into a commercial toolkit without permission. Linking to a repository does not make its contents part of Auto Company.

## No silent installation

The CLI prints suggested commands but never executes them. The owner must choose the integration, inspect it, and install it explicitly.
