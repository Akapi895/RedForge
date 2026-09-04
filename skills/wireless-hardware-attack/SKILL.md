---
name: wireless-hardware-attack
description: >-
  Wireless/hardware: WiFi PMKID/WPS/Evil Twin, BLE, Zigbee, NFC, SDR, UART/JTAG/SPI, side channels, and fault injection. Use when attacking WiFi, BLE, RFID, SDR, or hardware interfaces.
metadata:
  tags: [penetration-testing, red-team]
---

## Wireless / Hardware Attacks

```
=== Wireless/hardware ===
WiFi: aircrack-ng/PMKID (hcxdumptool+hashcat)/WPS reaver/Evil Twin | BLE: gatttool GATT enumeration/unauthenticated read-write/Just Works
Zigbee killerbee | NFC/RFID proxmark3 cloning/MIFARE mfoc | LoRa/SDR rtl-sdr+gnuradio replay
Hardware: UART baud-rate scanning for a shell | JTAG/SWD firmware read/write | SPI flash dump | side channels DPA/timing | fault injection (voltage/clock glitches to bypass authentication) | binwalk -Me
```
