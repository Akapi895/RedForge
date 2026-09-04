---
name: redteam-opsec
description: >-
  OPSEC covert-operation discipline: IP blacklist bypass, rate/timing control, traffic obfuscation, minimal footprint, anti-forensics, and gradual exposure. Use when maintaining stealth, bypassing IP bans, or planning covert red-team ops.
metadata:
  tags: [penetration-testing, red-team]
---

## OPSEC / Covert-Operation Discipline (Evasion / Stability / Stealth)

```
Core principle: gaining access ≠ operating safely. Discovery resets the operation to zero. Before every action, ask "what does this look like to the defender?"
🚨IP blacklist bypass(try immediately after a ban): X-Forwarded-For: random IP (when the CDN trusts this header, bypass the application-layer blacklist directly; when receiving code:10010/blacklist, try XFF/X-Real-IP/CF-Connecting-IP/True-Client-IP one by one)
  Verify: a normal request returns "blacklisted" + your IP → add the XFF header and receive a normal 200 → include XFF in all subsequent requests
  Advanced: change the XFF value every N requests(to avoid the new IP being banned too) | some CDNs trust only the first XFF value, others trust the last
Rate and timing: rate-limit scans(nuclei -rl / nmap -T2 --max-rate / ffuf -p delay) to avoid WAF bans + IDS threshold alerts | keep high-risk actions low-frequency with random jitter | avoid both business peaks and late night to match the target's routine
Traffic obfuscation: blend into normal business traffic(common UA/Referer/legitimate paths) | probe defenses first(process list/known EDR/SIEM agent indicators) before heavy tooling → prefer quiet techniques when monitoring exists; use automated batches only without monitoring
Minimal footprint: prefer in-memory execution without writing to disk(DDexec/memfd/reflective loading) | webshell strong password+unusual path+disguised functionality | route tunnels over 443/DNS like common outbound traffic | delete tools immediately after use(/dev/shm ramdisk, leave no remnants)
Anti-forensics: command history unset HISTFILE / set +o history | selectively delete only your log entries(truncating everything triggers alerts) | touch -r using a neighboring file as reference to preserve mtime | do not touch monitoring/audit services(stopping them is itself an alert)
Gradual exposure: passive reconnaissance(certificate transparency/passive DNS/search engines/asset engines)→confirm no strong monitoring→active scanning→only then deploy an exploit. Never actively touch the target when public intelligence can provide the answer. At every escalation ask "is exposure worth it?"
> Stealth is not perfectionism; it is a red team's survival skill. One reckless full-port, full-speed scan can zero out the entire operation.
```

