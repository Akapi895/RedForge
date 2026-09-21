---
name: source-code-hunting
description: >-
  Source-code hunting: exposed .git repositories, dangerous-function grep, JS RC4 deobfuscation, semgrep/CodeQL, trufflehog, patch diffs, supply chain/CI. Use when hunting source leaks, secrets, JS deobfuscation, or supply-chain issues.
metadata:
  tags: [penetration-testing, red-team]
---

## Source-Code Hunting

```
=== Source-code hunting ===
.git exposure: git-dumper → git log -p --all (deleted sensitive files) | .svn/.DS_Store/composer.lock
Dangerous-function grep: exec/system/eval/unserialize/pickle.loads/render/curl_exec + hard-coded keys (sk-/ghp_/BEGIN RSA)
JS deobfuscation (RC4 + base64 string-array pattern): 1) extract the string array (var a0G=[...]) 2) find the decoder (a0m(idx,key) → RC4 decrypt + base64) 3) find the rotation IIFE (target offset) 4) rebuild the decoder in Node.js and batch-decode all strings → obtain plaintext variable names/API paths/configuration
  UniApp indicators: app-service.js (business logic, often 800KB+ and obfuscated) + zlsioh.dat (encrypted configuration, decrypted by native .so) + dcloud_uniplugins.json (plugin manifest)
Static scanning: semgrep --config=auto for a quick scan / build a CodeQL database and write queries (taint solving + variant-analysis methodology in `zero-day-discovery`)
Deep secret hunting: trufflehog/gitleaks across full Git history + Docker image layers + npm/PyPI tarballs + frontend bundles (--only-verified distinguishes active from dead keys)
Framework patch diff: diff versions before/after a patch → what was fixed = where the vulnerability is | dependency chain: composer.json/npm audit/pip-audit
Supply chain/CI: dependency confusion (claim an internal package name on a public registry) | GHA command injection `${{github.event.issue.title}}` | takeover of self-hosted runners | .npmrc/.pypirc credentials
```
