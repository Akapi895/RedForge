---
name: zero-day-discovery
description: >-
  0day discovery engine: variant analysis/patch gaps/differential testing/fuzzing/taint reasoning/N-day weaponization/hunter mindset. Use when public vulns are not found and you need to discover a 0day or weaponize an N-day.
metadata:
  tags: [penetration-testing, red-team]
---

## 0day Discovery Engine (Dig for vulnerabilities when none can be found online)

```
Core shift: from "matching known vulnerability databases" → "understanding how code/protocols work and inferring where they fail". 0day is not luck; it is a method.
Five primary paths:
  1.Variant analysis(most productive): take a CVE patch→extract the vulnerability pattern→grep the entire codebase for the pattern elsewhere→locations not covered by the patch=0day
  2.Patch gaps: read the patch's filtering logic; blacklists can almost always be bypassed(missed encoding/equivalent function/alias)→new CVE
  3.Differential testing: two components interpret the same input differently(WAF vs backend/validator vs executor)→request smuggling/SSRF bypass
  4.Fuzzing: write a harness(wrapping the function that processes untrusted input)+generate corpus/dictionary+triage crashes(reproducibility/controllability/exploitability)
     AFL++/libFuzzer(coverage-guided memory corruption discovery) boofuzz(protocol) radamsa(black-box) restler(REST API)
  5.Taint reasoning(strongest with source): source(parameters/Headers/deserialization fields) without an effective sanitizer reaching a sink(exec/SQL/template)=0day
     write CodeQL queries to automatically find reachable data flows / Semgrep / Joern
N-day weaponization(advisory released but no PoC online): reverse-engineer an exploit from the patch diff(bindiff/diaphora binary comparison,
  vendor regression tests are often a PoC prototype)→make it work in a local lab→target the system. Maximize the vulnerability window; one of the most valuable red-team capabilities.
Hunter mindset(ask layered questions of any code/endpoint/protocol; every "yes" is a 0day candidate):
  Trust boundaries: can the input be trusted? When does that assumption fail? | state and timing: can state change between two steps(TOCTOU)? Can the order be disrupted?
  Parsing and normalization: how many times is it parsed? Does normalize happen before or after validation? | boundary extremes: negative/0/huge values/type confusion/encoding/null bytes?
  Implicit capabilities: what does this feature give me "as a side effect"? | uniqueness: are IDs/tokens predictable? Is the "secret" really secret?
Break a component's "assumption checklist" one item at a time: file-upload assumptions "only images are uploaded/extensions are trusted/filenames contain no path/content is only data"→breaking each assumption reveals a vulnerability; combining them creates a chain.
0day validation(higher bar than a CVE): reproducible(minimal PoC)+clear root cause(which line/assumption)+provable impact(actual read/write/execution)+false-positive exclusion.
```

