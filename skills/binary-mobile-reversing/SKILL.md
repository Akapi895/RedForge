---
name: binary-mobile-reversing
description: >-
  APK/EXE/binary: UniApp/DCloud/Flutter reverse engineering, certificate-pinning bypass, exported components, memory-corruption exploit chains, and IoT firmware. Use when reversing APK/EXE, UniApp/Flutter, native .so, or memory-corruption exploits.
metadata:
  tags: [penetration-testing, red-team]
---

## APK / EXE / Binary Reverse Engineering

```
=== APK/EXE/Binary ===
APK: apktool d / jadx | inspect exported components/deeplink/debuggable in Manifest | rg hardcoded keys+API addresses | .so strings/Ghidra | dynamic frida/objection
🚨UniApp/DCloud APK reverse engineering(H5 hybrid app; all business logic is in JS, not DEX):
  Identify: __UNI__XXXXXXX + assets/apps/ + uni-jsframework.js (NOT Flutter; no libapp.so) | manifest.json contains the "uni-app" field
  Core file: assets/apps/__UNI__XXX/www/app-service.js (800KB+ obfuscated JS=all business logic) | jadx shows only the shell(native plugin registration); the real logic is in JS
  JS deobfuscation(RC4+rotating string array): extract the large a0G() array(10,000+ elements)+a0m(idx,key) decoder(base64→RC4)+rotation IIFE→assemble executable JS→batch-decode with node→JSON dictionary(offset→plaintext)
    🔴Trap: the rotation IIFE may end with a comma(not a semicolon; it belongs to a larger expression)→node reports "Unexpected token"→extract using the actual terminator
  Encrypted configuration zlsioh.dat/dcloud3.dat: header(96B offset table: off40=block1 decompressed/44=compressed/60=block2 offset/64=block2 size/80=block3 offset)+block1(zlib→DEX)+block2(encrypted, contains the critical API domain list; requires .so decryption)+block3(zlib→AndroidX class mapping)
    zlib magic 78DA=unencrypted, decompress directly | no magic=encrypted block(key is in .so .rodata/assembly immediates)
  .so string deobfuscation: extract even-position characters→reverse(e.g. "mAojcl.dubdFiHaebP.nwywfwb"→"www.baidu.com") | strings -n8 lib*.so to find long strings containing dots
  ⚠️Trap: x-separated IP strings in .so(e.g. "x111.230.69.120x118.126.105.164")=DCloud HTTPDNS nodes(not the API backend! 140.205.11.x=Alibaba Cloud DNS, 111.230/118.126/106.52/42.193=Tencent Cloud) | resources.arsc errors are intentional decompilation protection(do not affect jadx) | fake PNG assets(1-8px)=integrity checks, not data | fcapp.run is only an APK download proxy
  Key API variables: $apiHost/$opHost(set dynamically at runtime, not hardcoded)/urlList/checkAvailableDomainList(HEAD domain/favicon.ico for connectivity; 200-399=available)
  High-value endpoints: /app_init /user/gustRegister(guest registration without verification) /user/getOssSts(OSS temporary credentials!) /uploadFile2
  Configuration: manifest.json(appId/version/nativePlugins→wrs-httpserver=embedded HTTP service!) | supplierconfig.json(vivo/xiaomi/huawei/oppo appid) | dcloud_uniplugins.json(plugin manifest)
  ChengZi SDK(site-building distribution platform): init3 returns XOR single-byte(key=0x96, see chengzi_decrypt.py) encrypted JSON→decrypt to obtain the real APK download link(fu field)+channel code
  Domain fallback: when $apiHost cannot be obtained statically→Android emulator(Waydroid/AVD)+tcpdump/mitmproxy to capture runtime DNS
  See: references/uniapp-apk-reverse-engineering.md, references/uniapp-apk-reversing.md, scripts/js_rc4_deobfuscate.js, scripts/chengzi_decrypt.py
🚨Flutter APK quick reverse engineering(no unpacking required): hardening(MogoSec/Bangbang/360) protects the DEX layer, but Flutter's lib/arm64-v8a/libapp.so(Dart AOT) is usually not encrypted
  strings -n8 libapp.so|grep 'https\?://' → all hardcoded URLs/API domains/S3 addresses/CDNs
  strings -n8 libapp.so|grep -E '^/' → API path matrix(/Member/Login, /Web/VideoList, /BBS/GetSTSToken, etc.)
  strings -n8 libapp.so|grep -iE 'key|secret|token|aws|bucket' → credential/key disclosure
  assets/config_*.xml + assets/data_*.dat → encrypted configuration(possibly containing domains/API addresses) | .DS_Store → macOS developer information disclosure | assets/xinstall* → channel-tracking SDK
🚨Native .so string deobfuscation(general pattern): obfuscated strings are in the .rodata section
  Common patterns: extract even-position characters+reverse(e.g. "mAojcl.dubdFiHaebP..."→extract even positions→reverse=plaintext domain) | XOR constant | RC4+base64
  Identify: strings -n8 lib*.so for equal-length 16-byte strings(AES key/IV)/strings containing dots(domains)/x-separated IP strings(HTTPDNS) | nm --dynamic for obfuscated exported symbols(16 random upper/lowercase characters)
  Verify: compare a known string("classes.dex"/"io.dcloud.application") with its obfuscated version→reverse-engineer the algorithm→batch-decode
🚨General mobile(non-UniApp/Flutter conventional apps):
  Certificate-pinning bypass(capture HTTPS): objection android sslpinning disable | frida universal-unpinning | edit smali to remove pinning and repackage
  Exported-component privilege issues: Activity/Service/Provider/Receiver with Manifest exported=true → cross-app invocation via drozer/adb am start | ContentProvider SQLi/path traversal
  Deep-link hijacking: unchecked scheme://→WebView loads arbitrary URL(XSS/file read) | parameter enters intent→component hijacking
  Insecure storage: /data/data/pkg/(plaintext shared_prefs token/unencrypted db/logcat disclosure) | world-readable sdcard
  WebView: addJavascriptInterface(<4.2 RCE)/file:// local read/setAllowFileAccess | static: MobSF all-in-one | iOS: unpack(frida-ios-dump)+inspect Mach-O with Ghidra+dynamic objection
EXE/PE: file/strings | decompile with Ghidra/IDA to find hardcoded values/encryption/network logic | dynamic analysis with x64dbg | vulnerabilities: overflow/format string/UAF/DLL hijacking
🚨Memory-corruption exploit chain(crash→RCE): inspect protections with checksec → determine primitive(stack overflow/UAF/format string=arbitrary read/write)
  → information disclosure to bypass ASLR → ROP chain(ROPgadget/pwntools ret2libc) → heap exploitation(tcache poison/__free_hook hijack)
  → arbitrary write to modify GOT/hook/vtable/exit_funcs → control-flow hijacking | IoT firmware: extract with binwalk -Me+debug with qemu-user(weak protections often allow direct stack overflow)
```
