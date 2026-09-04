---
name: active-directory-attack
description: >-
  Internal-domain attacks: BloodHound, Kerberoast, ADCS ESC1/ESC8, NTLM Relay, coercion, DACL, DCSync, Zerologon/NoPac/PrintNightmare, mitm6, LLMNR, and Linux internal networks. Use when attacking Active Directory, ADCS, NTLM relay, or an internal domain.
metadata:
  tags: [penetration-testing, red-team]
---

## Internal-Domain Attacks

```
=== Internal domain (the real battleground in 2023+) ===
Reconnaissance: BloodHound (SharpHound collection → attack paths) | Kerberos: GetNPUsers (AS-REP) / GetUserSPNs (Kerberoast)
🚨 ADCS (Certipy full workflow): certipy find -vulnerable | ESC1: request a domain-controller certificate with a specified SAN | ESC8: relay to the CA to obtain a DC certificate
🚨 NTLM Relay (more important than PtH; PtH is often blocked by EDR): ntlmrelayx -t ldap --escalate-user / -t http CA --adcs (ESC8) / RBCD
🚨 Coercion: PetitPotam (MS-EFSRPC) / coercer across protocols / printerbug → feed into relay
🚨 DACL abuse: WriteDACL → grant yourself DCSync rights | Shadow credentials: certipy shadow (GenericWrite is sufficient; does not change the password or leave an obvious trace)
DCSync: secretsdump -just-dc → krbtgt hash → Golden Ticket
🚨 High-impact domain CVEs (test first; a hit may lead directly to domain admin): Zerologon (CVE-2020-1472, clear the DC machine-account password → DCSync) | NoPac (CVE-2021-42278/42287, rename a machine account to request a DC TGT) | PrintNightmare (CVE-2021-34527, spooler RCE / malicious-driver loading) | EternalBlue (MS17-010, direct RCE through legacy SMBv1)
🚨 IPv6/mitm6 (mandatory in dual-stack internal networks): mitm6 hijacks DHCPv6 + DNS → WPAD → ntlmrelayx to LDAP/ADCS (stealthier than LLMNR and preferred in modern internal networks)
LLMNR/NBT-NS poisoning: Responder captures NetNTLMv2 → crack with hashcat or relay
Linux internal network: unauthenticated Redis (CONFIG SET dir to write an SSH key) | NFS showmount | Docker 2375
```
