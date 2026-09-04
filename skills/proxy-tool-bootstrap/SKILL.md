---
name: proxy-tool-bootstrap
description: >-
  Find proxies and bootstrap tools: SOCKS5/HTTP/Tor route-switching, Python tool bootstrapping, generated wordlists, and OOB infrastructure. Use when blocked by 403/429/WAF/timeouts, missing tools, or when OOB confirmation is needed.
metadata:
  tags: [penetration-testing, red-team]
---

## Find Proxies + Bootstrap Tools (switch routes when blocked; write the missing tool yourself)

```
🔴 Proxy (after refusal/rate limiting/timeout → first response is not retrying, but switching routes. Failing to switch routes means giving up and triggers blackboard rule 2; see `pentest-blackboard`):
  Execute in order; proceed only after the previous step fails:
  ① Identify the target region: terminal: curl -s "http://ip-api.com/json/{TARGET_IP}" | python3 -c "import sys,json;d=json.load(sys.stdin);print(d['country'],d['countryCode'])"
  ② Obtain SOCKS5 proxies in that region: terminal: curl -s "https://api.proxyscrape.com/v2/?request=displayproxies&protocol=socks5&country={countryCode}&timeout=5000" | head -20
  ③ Verify that the proxy can reach the target: terminal: curl --socks5 {PROXY_IP:PORT} --connect-timeout 8 -s -o /dev/null -w "%{http_code}" {TARGET_URL}
  ④ 200=available; immediately retry the blocked request through this proxy. Non-200=try the next proxy and repeat ③.
  ⑤ If all SOCKS5 proxies fail, obtain HTTP proxies: terminal: curl -s "https://api.proxyscrape.com/v2/?request=displayproxies&protocol=http&country={countryCode}&timeout=5000" | head -20
  ⑥ Verify the HTTP proxy: terminal: curl --proxy http://{PROXY_IP:PORT} --connect-timeout 8 -s -o /dev/null -w "%{http_code}" {TARGET_URL}
  ⑦ If every proxy fails, use Tor: terminal: curl --socks5 127.0.0.1:9050 --connect-timeout 15 {TARGET_URL}
  Add proxy arguments consistently: curl --socks5 / sqlmap --proxy=socks5://{P} / nmap --proxies socks5://{P} / nuclei -proxy socks5://{P} / ffuf -x socks5://{P}
  Rotation: on 429/403, immediately switch to the next proxy; switch proactively every 20 requests (a new IP may also be banned) | Cloudflare → proxy pool + random 2–5 second request interval
  🚨 HTTP proxy vs SOCKS5: an HTTP proxy may insert its own error page (502/cannot display this page), making it impossible to distinguish the target's actual response. Use SOCKS5 (--socks5) for probing.
    ProxyScrape SOCKS5: https://api.proxyscrape.com/v2/?request=displayproxies&protocol=socks5&country=CN,JP&timeout=5000
    Verify SOCKS5: curl --socks5 IP:PORT --connect-timeout 5 TARGET | HTTP proxies are suitable only for confirmed-reachable targets and anonymity/rotation
Tool bootstrapping (which X || implement it in Python):
  no nmap → socket port scan | no ffuf → requests directory brute force | no sqlmap → manual payload checks | no hydra → requests brute force | no nuclei → requests against known payloads
  Simulate complex tools with Python: requests+bs4 crawler / base64/hex encoding / hashlib hashes / pycryptodome encryption / scapy sniffing
Generated wordlists: derive variants from target domain/company name | extract keywords from web pages | username + year + special-character combinations | service default credentials
OOB infrastructure (all blind vulnerabilities depend on it): interactsh-client for an oast.fun domain | or a VPS python3 -m http.server/nc to observe callbacks | ngrok/cloudflared tunnel
  → Only an observed OOB callback (DNS query/HTTP request) counts as confirmation → write a Fact
```
