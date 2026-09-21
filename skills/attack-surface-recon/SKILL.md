---
name: attack-surface-recon
description: >-
  Reconnaissance/attack-surface mapping: passive whois/amass/crt.sh/FOFA/Shodan, active subfinder/httpx/naabu/katana/nuclei, DNS geography/CDN/Nginx catch-all/BT Panel/UniApp fingerprinting. First action at the start; write discoveries to the project blackboard. Use when starting recon, asset mapping, fingerprinting, or CDN/DNS bypass discovery.
metadata:
  tags: [penetration-testing, red-team]
---

## Reconnaissance / Attack-Surface Mapping

```
=== Reconnaissance/attack-surface mapping(first action, 60% of operational time; do not jump into exploitation; immediately upsert_project_fact for ports/fingerprints) ===
Passive first(do not touch the target; execute all of these):
  terminal: whois {domain} | grep -iE 'org|name|email|registrant'
  terminal: amass intel -asn {ASN} -d {domain}
  browser_navigate: https://crt.sh/?q=%.{domain} (enumerate subdomains through certificate transparency)
  terminal: curl -s "https://crt.sh/?q=%.{domain}&output=json" | python3 -c "import sys,json;[print(x['name_value']) for x in json.load(sys.stdin)]" | sort -u
  terminal: dig +short {domain} @114.114.114.114; dig +short {domain} @8.8.8.8 (differences indicate a CDN)
  terminal: curl -s "https://web.archive.org/cdx/search/cdx?url=*.{domain}/*&output=text&fl=original&collapse=urlkey" | head -100
  browser_navigate: https://fofa.info/result?qbase64=$(echo -n 'domain="{domain}"' | base64)
  browser_navigate: https://www.shodan.io/search?query=hostname:{domain}
Active pipeline(execute step by step):
  terminal: subfinder -d {domain} -silent | tee subs.txt; amass enum -d {domain} -passive >> subs.txt; sort -u subs.txt -o subs.txt
  terminal: cat subs.txt | dnsx -silent | tee alive_subs.txt
  terminal: cat alive_subs.txt | httpx -silent -title -status-code -tech-detect -cdn | tee httpx_out.txt
  terminal: naabu -l alive_subs.txt -top-ports 1000 -silent | tee ports.txt; nmap -sCV -iL <(head -20 alive_subs.txt) -oN nmap_out.txt
  terminal: katana -list alive_subs.txt -silent -d 3 | tee urls.txt; gau {domain} >> urls.txt
  terminal: cat urls.txt | grep -iE '\.(js|json)$' | httpx -silent -mc 200 | while read u; do curl -s "$u" | grep -oiE '(api|secret|key|token|password|aws|endpoint)[^"]*'; done
  terminal: ffuf -u https://{target}/FUZZ -w /usr/share/seclists/Discovery/Web-Content/raft-medium-directories.txt -mc 200,301,302,403 -fs {catch-all-size}
  terminal: nuclei -l alive_subs.txt -t /root/nuclei-templates/ -severity medium,high,critical -o nuclei_results.txt
  → when httpx identifies a framework/version, immediately trigger the `component-vuln-intel` search sequence
🚨DNS geographic restriction bypass: domestic CDNs(987dns/dnspod/Alibaba Cloud) often return 0.0.0.0 to overseas DNS; resolve through 114.114.114.114 to obtain the real IP. Compare multiple DNS servers: dig @1.1.1.1 vs @114.114.114.114 vs @8.8.8.8; differences indicate a CDN geographic policy. The real IP is hidden in domestic DNS results.
🚨CDN response-filter bypass(GCCDN/PWS and similar): for particular URL patterns(admin paths/actuator), the CDN returns 502 instead of forwarding the backend response:
  Fingerprint: Server: PWS/x.x.x.x | Via: 1.1 PS-XXX-XXXX:N (W) | X-Px: ms CS-XXX-XXXXnone(origin) | wildcard *.gccdn.net CNAME
  502≠nonexistent: 502=backend response filtered by CDN, 404=path truly nonexistent, 403=CDN route rejected. Enumerate live endpoints using 502/404/403 differences
  Path traversal confirmation: /api/v1/v/..;/..;/actuator/env → 502(exists but filtered) vs /api/v1/v/..;/..;/xxxx → 404(Tomcat natively nonexistent)
  Multiple CDN-node ports: scan ports 3000-9200 on the CDN IP(8000/8080/9090 often proxy the same backend but rules may differ; 9200 may be a standalone nginx)
  Bypass-attempt matrix(all failed in this session; record to prevent repetition): URL encoding/double encoding/case changes/semicolon suffix/TE-space smuggling/HTTP pipelining/Accept-Encoding/Range
  Extract useful information: even with 502, POST actuator/env may still execute in the backend(blind property write)→trigger with refresh(if refresh returns 404, the chain is broken)
🚨Nginx Catch-All trap: a PHP CMS rewrite for all routes→every path returns 200 with the same size=catch-all(not file existence!). Identify it by comparing three random paths with identical sizes. Test real live paths with a .php suffix(nginx `location ~ \.php$` passes directly to FastCGI; nonexistent paths return a true 404).
🚨Pure-FTPd puredb corruption≠authentication bypass: "Unable to read the indexed puredb file"=the virtual-user database is corrupted; every username(including anonymous) triggers a 421 disconnect. PAM system-user authentication also goes through puredb first→all attempts fail. Not an attack surface; skip it.
🚨Lateral movement across sites on the same IP: reverse-lookup the IP(SecurityTrails/VirusTotal) to find other domains on the server→break into a weak site(standalone nginx configuration returning a true 404 is enumerable)→a BT Panel managing all sites enables lateral movement across them.
Fingerprinting priority: once httpx/wappalyzer identifies a framework version → immediately switch to `component-vuln-intel` online intelligence. Map the attack surface first, then use Reason to choose the richest path; do not tunnel vision on the first bug.
🚨Wildcard DNS/CDN identification: all subdomains(including nonexistent random subdomains) resolve to the CDN=wildcard DNS(*.domain → CDN CNAME). Subdomain enumeration is meaningless here; all IPs are CDN nodes.
  Verify: dig random123456.target.com → if it also has an A record pointing to a known CDN IP = wildcard. Obtain the real origin IP through other channels(historical DNS/certificate transparency/other domains on the same server/mail headers/SSRF)
🚨BT Panel identification: a 32-character hex hash cookie(e.g. 789d6d2de16419a1e8bbfea926c8996e) + 302 redirect to /login + nginx reverse proxy + "entry validation failed" page = a BT Panel random security entry
  Windows version: IIS+nginx coexist, panel on port 8888, phpMyAdmin possibly on port 888, /index.html may leak site configuration(reverse-proxy target URL/site name)
  Attack: brute-force the entry path(8 random alphanumeric characters) | historical CVEs(API-key unauthenticated access/phar deserialization) | reuse default panel FTP/MySQL accounts
🚨Nginx 404 differential fingerprint: when an nginx catch-all rewrite masks real paths, compare response sizes—real nginx 404(146B)=file exists on disk but the location block denies it; catch-all(large page)=route not matched. Batch-probe CMS marker paths(/caches/configs/, /phpcms/modules/, /statics/, etc.); response-size differences identify the CMS type. See references/nginx-404-differential-fingerprinting.md
🚨UniApp/DCloud APK reverse-engineering fast path: assets/apps/__UNI__*/www/app-service.js = all business logic(JS obfuscated); assets/zlsioh.dat = encrypted configuration(API domain); lib/arm64-v8a/libirzlemnr.so = decryption library. Deobfuscate JS with RC4+base64(extract the a0G array+a0m decoder); search strings in the .so for a random 16-character string=possible key, and for the `xnumber.number.number.numberx` format=hardcoded IP(usually DCloud HTTPDNS infrastructure rather than the target API). See references/uniapp-dcloud-apk-reversing.md
```
