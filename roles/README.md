# RedForge Roles

This directory is the active English role catalog for RedForge. It combines
the original conservative benchmark roles with domain-specific assessment
roles. Roles are prompts and policy context; server-side scope enforcement,
tool allowlists, approvals, and audit logging remain authoritative.

The catalog is designed for authorized black-box assessment across web,
network, cloud, container, and infrastructure targets. Discovery should be
read-only where possible. Intrusive, destructive, multi-host, persistence, and
data-access actions require explicit policy and human approval.

`CTF.yaml` is retained for isolated training use and disabled by default. The
legacy archive has been merged into this directory; duplicate `Default` and
`Information Gathering` roles were not retained because they are covered by
`Default` and `Reconnaissance`.
