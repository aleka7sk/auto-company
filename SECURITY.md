# Security policy

## Supported version

The current `main` branch is supported while Auto Company is pre-1.0.

## Security model

Auto Company combines instructions, local tooling, hooks, and external AI coding agents. Skills execute with the permissions of their host agent. The included hook reduces accidental risk but is not an operating-system sandbox.

## Protected operations

The default guard blocks typical access to secret files and commands associated with:

- force pushes and destructive Git cleanup;
- direct pushes to `main` or `master`;
- pull request merge;
- infrastructure destruction;
- namespace or release deletion;
- database/schema drops;
- package publishing;
- app-store submission and common production deployment commands;
- broad filesystem deletion.

Production release, destructive migration, data deletion, permission-model changes, credentials, billing, and protected branch operations require explicit human approval and should use constrained credentials.

## Reporting

Do not open a public issue containing credentials, exploit details against a live system, or private customer data. Contact the repository owner privately through the GitHub account associated with this repository.

## Third-party integrations

Auto Company does not vendor or auto-install the external repositories listed in its registry. Review the exact release/commit, scripts, hooks, MCP configuration, network calls, telemetry, and license before installation. Pin shared/team installations.
