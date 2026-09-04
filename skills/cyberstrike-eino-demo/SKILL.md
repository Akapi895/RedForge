---
name: cyberstrike-eino-demo
description: >-
  Full-featured example skill package: SKILL.md + optional scripts/, references/, assets/, and other directories; validates Eino skills and HTTP package-internal paths(for authorized security testing and education only).
---

# CyberStrike × Eino Full-Featured Skill Demo

This package follows [Agent Skills](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview): **`SKILL.md` is the manifest + primary instructions**(there is no separate `SKILL.yaml`). The directory may contain arbitrary subdirectories such as **`scripts/`**, **`references/`**, and **`assets/`**(provided paths are safe and package depth/file-count limits are not exceeded); they are read by **`ListPackageFiles` / `resource_path`** and Eino local tools. See `FORMS.md` and `REFERENCE.md` for supplementary details.

## Overview

Use this to verify in one pass:

- HTTP `GET /api/skills` listing(`script_count`, `file_count`, `progressive`, etc. are derived/scanned results)
- `GET /api/skills/cyberstrike-eino-demo?depth=summary|full`
- `section=` maps to a short id for an **`##` heading** or ASCII heading in `SKILL.md`(for example, `## Payload Examples` commonly maps to `section=payload`)
- The multi-agent ADK **`skill`** tool(and optional local file tools) reads resources by relative path inside the package
- Eino `FilesystemSkillsRetriever` retrieves the package summary, `##` sections, and script entries

**Hard requirement**: every test must have written authorization and remain within the agreed scope and time window.

## Authorized Testing Workflow

1. **Confirm scope**: domains/IPs, endpoint list, prohibited actions(DoS, bulk data extraction, etc.).
2. **Record baseline**: perform read-only probes against agreed assets; save timestamps and raw request/response summaries.
3. **Test by category**: split tasks by vulnerability type; reconfirm authorization boundaries before high-risk actions.
4. **Evidence and reporting**: include reproduction steps, impact, and remediation for every finding; redact sensitive data.
5. **Closeout**: delete temporary accounts, clean test data, and hand over the report.

## Payload Examples

The following are **educational placeholders**; replace them with target context during authorized testing and never use them against unauthorized systems:

- SQLi probe(error-based): `"'`(observe whether database error information is disclosed)
- Reflected XSS(sanitized): `<script>alert(1)</script>` → in a lab it should be encoded or blocked by CSP
- Path traversal(read-only verification): `....//....//etc/passwd`(only in an authorized file-read scenario)

See `scripts/payloads.txt` for the detailed list.

## references/ and assets/

Used to verify that non-`scripts/` subdirectories are treated equivalently:

| Path | Purpose |
|------|------|
| `references/citations.md` | Citation and HTTP `resource_path` testing notes |
| `assets/README.txt` | Placeholder resource(can be replaced with a real binary for file-read limit testing) |

## Recommended Toolchain

| Stage | Example tools |
|------|-----------|
| Proxying and replay | Burp Suite, mitmproxy |
| Scanning and directory enumeration | ffuf, nuclei(reduce concurrency to respect authorization) |
| Vulnerability verification | Custom PoCs, official CLIs(sqlmap, etc.) only within authorized scope |
| Recording | Markdown + JSON snippet template(see `scripts/report-snippet.json`) |

## Checklist and Verification

- [ ] Written authorization and testing window saved
- [ ] Files under `scripts/` match references in the main text
- [ ] Index can be checked through Web or `GET /api/skills?...`; in multi-agent sessions, load by package with **`skill`** to save tokens
- [ ] When details are needed, fetch the full text through **`skill`**, or use HTTP `depth=full`, `section=<heading or short id>`
- [ ] When raw script content is needed, use local file tools or HTTP `resource_path=scripts/check-env.sh`
- [ ] `resource_path=references/citations.md` and `resource_path=assets/README.txt` are readable
