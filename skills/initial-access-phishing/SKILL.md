---
name: initial-access-phishing
description: >-
  Initial access/phishing/social engineering: password spraying, AiTM, device codes, OAuth-consent phishing, payloads, and vishing. Use when needing initial access, phishing, AiTM, device codes, or social engineering.
metadata:
  tags: [penetration-testing, red-team]
---

## Initial Access / Phishing / Social Engineering

```
=== Initial access/phishing/social engineering (the first step from external to internal) ===
Password spraying: spray weak passwords (Season+Year!/CompanyName123) across all accounts (1–2 attempts per account to avoid lockout) | Sources: breach corpora / OSINT email formats (f.last@) / default credentials
AiTM phishing (bypassing MFA): evilginx3/Modlishka reverse-proxies the real site → man-in-the-middle steals session cookies + tokens (bypassing MFA) | configure the phishlet domain + certificate
Device-code phishing (Azure/M365): device-code flow → induce the victim to enter the code → obtain access/refresh tokens (without the password/MFA) | TokenTactics/AADInternals
Phishing payloads: gophish mailer + landing page | payloads: lnk/iso/macros/HTA/OneNote | bypass gateways with password-protected ZIPs/cloud-drive links/HTML smuggling
OAuth-consent phishing (illicit consent): malicious app requests excessive scopes (Mail.Read/offline_access) → victim consents → long-lived access token
Social-engineering preparation: OSINT on organization structure/vendors/password habits (LinkedIn/recruiting/GitHub) | pretexts (IT support/vendor/HR) | vishing (deepfake voice)
→ After initial access lands, switch immediately to `redteam-opsec` (do not get caught by EDR at the first step), then feed acquired capabilities into `capability-primitive-search` to assemble the chain.
```
