---
name: capability-primitive-search
description: >-
  Capability primitives + state-space search: combine primitives such as read/write/exec/ssrf into RCE equations A-F, map low-severity findings, search forward and backward, and realize cross-domain chains. Use when there is no single RCE, when chaining low-severity vulnerabilities, or when deriving novel attack chains.
metadata:
  tags: [penetration-testing, red-team]
---

## Capability Primitives + State-Space Search

```
Principle: RCE is not a single "vulnerability"; it emerges when a set of "capabilities" is assembled. When there is no single high-impact flaw, several informational/low-severity findings can still combine into code execution.
  Mindset to break: "we cannot find RCE/deserialization/upload, so this site is unexploitable". Your job is not to "find an RCE vulnerability"; it is to "assemble the primitives required for execution".
  Abstract each Fact into capability primitives(do not record "found vulnerability X"; record "what capability exists now + its limitations"):
    read(path) write(path) exec(cmd) ssrf(url) sqli redirect(url) eval_expr idor(id) cred(svc,priv) coerce_auth write_acl

RCE requires satisfying any one equation; break it into obtainable primitives and assemble them:
  A. Can write a file + the file is executed as code = RCE
  B. Can control configuration/env + configuration points to your code = RCE
  C. Can access an administrative interface + the panel has built-in execution = RCE(not a vulnerability; a feature!)
  D. Have credentials + the service has a legitimate execution entry point = RCE(abuse of legitimate functionality)
  E. Can read arbitrarily + obtain credentials + credentials log into an execution point = RCE
  F. Can control data + data flows into a dangerous sink(eval/template/SQL) = RCE

Low severity → primitive mapping(translate "low-value vulnerabilities" into puzzle pieces):
  Information disclosure(.git/backup/stack trace)→source/paths/keys→feed B/E/F | LFI/arbitrary read→read configuration keys→E, or poison logs→A
  SSRF(even GET-only)→reach internal Redis/Consul/K8s/cloud metadata→C; obtain temporary credentials from cloud metadata→D
  Weak/default/reused credentials→enter a backend with task/plugin/webhook/CI features→C | CORS/CSRF/XSS→use an administrator's browser to invoke execution features→C
  Controllable upload(even with extension restrictions)→combine with path traversal/parser differences/.htaccess→A | configuration write→modify templates/logs/connection strings→B
  SQLi(even read-only)→read hashes/keys→E, or OUTFILE→A | controllable template→SSTI→F | prototype pollution→pollute downstream properties→F

State-space search(search on your own when no existing chain is available): state=current capability set, action=use a capability to unlock a new capability, goal=Goal.
  Forward: for each capability ask "what can it unlock?"; for each pair of capabilities ask "what can they produce together?"
    (read+write=configuration change; ssrf+internal redis=RCE; sqli+FILE=webshell; coerce_auth+relay=domain SYSTEM; idor+massassign=modify another user's admin)
  Backward(use primarily when stuck): lock Goal=execute a command→choose the equation closest to the current state as a template→make each missing primitive a sub-goal→
    ask which low-severity finding/feature/information disclosure can provide it→if unavailable, recursively decompose or switch equations → forward and backward searches meet in the middle=the complete chain emerges → verify every segment

Breakthrough points(do not overlook them):
  ·"A feature is a primitive": backend task schedulers/plugins/template editors/SQL consoles/file managers/import-export/webhooks — legitimate features, but access provides ready-made execution/read/write primitives. To a pentester there is no "feature/vulnerability" distinction, only "capability".
  ·Cross-protocol jumps: SSRF's gopher/dict/file turns "can only send HTTP" into "reach Redis/send SMTP/read files".
  ·Credential reuse is universal glue: any password/key obtained anywhere should be sprayed across all services by default. Lateral movement is often faster than vertical escalation.
  ·Parser differences: upload validation/routing/reverse proxy interpret input differently→bypass gaps(double extensions/encoding/Host confusion).
  ·Time dimension: TOCTOU/predictable tokens/cache poisoning turn "occasional" into "stable" and non-exploitable into exploitable.
  ·Cross-domain realization: whenever you obtain a capability ask "what is it worth in another domain?"—capabilities are universal currency. Web SSRF→cloud metadata account takeover; hardcoded APK values→direct internal API access bypassing frontend authorization; supply chain→CI keys into production.
  ·Creative mode(when known combinations are exhausted): reassess capability boundaries(can you only read /var/log? What about /proc/self/environ?) | find equivalent RCE sinks(writing crontab/.bashrc/CI config/LD_PRELOAD/authorized_keys/systemd unit all=RCE) | use information(errors/timing/response length) as side channels | invert assumptions: list "what I thought was impossible" and ask "why is it impossible?" for each
> A derived chain is a hypothesis(record it as a tentative note/chain); only after actual execution + evidence write a confirmed Fact / record_vulnerability. The complete chain is valid only when every step has been verified.
```

