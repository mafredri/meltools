# Notes

## Known commands

| Command              | Description                                                                  | Usage                                  |
| -------------------- | ---------------------------------------------------------------------------- | -------------------------------------- |
| `p`                  | Process list and memory consumption                                          | `p`                                    |
| `ip`                 | Show IP configuration (IP, GW, MASK, MAC)                                    | `ip`                                   |
| `resolve`            | Resolve hostname, dumps a `hostent` struct                                   | `resolve=domain.com`                   |
| `debug enable`       | Enable debug mode                                                            | `debug enable`                         |
| `debug disable`      | Disable debug mode                                                           | `debug disable`                        |
| `debug info`         | Show debug info (e.g. model, firmware, boot, etc.)                           | `debug info`                           |
| `debug clear`        | Clear debug mode                                                             | `debug clear`                          |
| `debug reset`        | Reset debug mode                                                             | `debug reset`                          |
| `log start`          | Start log capture                                                            | `log start`                            |
| `log stop`           | Stop log capture                                                             | `log stop`                             |
| `log get`            | Print captured logs                                                          | `log get`                              |
| `log level`          | Set log level filter (default: `0x6c` = `0b1101100`, i.e. 3 levels excluded) | `log level=0xff`                       |
| `log type`           | Set log type filter (default: `0x00` = no logs)                              | `log type=0xff`                        |
| `log status`         | Show log status (`0/100,0x6C,0x00` = `[captured/limit],[level],[type]`)      | `log status`                           |
| `antenna_info`       | Show antenna info (shows `0:main antenna`)                                   | `antenna_info`                         |
| `atpp`               | Ping utility, enter `atpp` for help                                          | `atpp=192.168.1.1,-n,1`                |
| `atws`               | WiFi scan, outputs all SSIDs seen by the adapter                             | `atws`                                 |
| `flash erase`        | Erase flash memory (untested, likely DANGEROUS!)                             |                                        |
| `flash write`        | Write to flash (untested, likely DANGEROUS!)                                 |                                        |
| `flash sector read`  | Read flash sectors (128 bytes per count)                                     | `flash sector read=[hex addr],[count]` |
| `flash sector write` | Write flash sector (untested, likely DANGEROUS!)                             |                                        |
| `db`                 | (unknown/undocumented)                                                       |                                        |
| `mb`                 | (unknown/undocumented)                                                       |                                        |
| `iwpriv`             | (unknown/undocumented)                                                       |                                        |

(There are likely some additional commands, but I have not explored those as many cause the telnet connection to break.)

## Enable debug logging

```
./meltelnet 192.168.1.100
Connected to 192.168.1.100:23
log info
0/100,0x6C,0x00
log level=0xff
log type=0xff
2001/01/01_00:11:59 [Ea]HTTPC serverSignalSendMain_fd httpc_getResponsePoll ERROR(2)
2001/01/01_00:11:59 [Ea]HTTPC ERROR
2001/01/01_00:11:59 [Ip]task_netctrl_udprecv_fv (594):exit recvfrom=32
2001/01/01_00:11:59 [Ip]task_netctrl_udprecv_fv (591):enter recvfrom
2001/01/01_00:11:59 [Ii]FC 42 1 30 10 4 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
2001/01/01_00:11:59 [Ip]task_netctrl_udprecv_fv (594):exit recvfrom=14
2001/01/01_00:11:59 [Ip]task_netctrl_udprecv_fv (591):enter recvfrom
2001/01/01_00:11:59 [Ip]task_netctrl_udprecv_fv (591):enter recvfrom
2001/01/01_00:12:00 [Ii]FC 42 1 30 10 6 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
2001/01/01_00:12:00 [Ii]FC 62 1 30 10 6 0 0 0 [CENSORED]
2001/01/01_00:12:00 [Da]led_prictrl_log_fv (256):LED_STATUS_SV_SYNCHRO( 3):OFF
2001/01/01_00:12:00 [Ii]FC 42 1 30 10 9 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
2001/01/01_00:12:00 [Da]led_prictrl_log_fv (256):LED_STATUS_SV_SYNCHRO( 3):ON
```

Restore original logging settings:

```
log level=0x6c
log type=0x00
```

The `log start/stop/get` can be used to capture and buffer up to 100 log entries, but it seems pointless since the logs are printed in real-time.

## Enable ECHONET

Enable ECHONET by sending `ECHONET: "ON"` in the CSV data.

```
<CSV>
    <CONNECT>ON</CONNECT>
    <ECHONET>ON</ECHONET>
</CSV>
```

(Connect might not be needed, but it doesn't hurt.)

```console
$ curl -H 'Content-Type: text/xml' -d '<?xml version="1.0" encoding="UTF-8"?><ESV>7WVvmfhMYzGVi70nyFhmKEy9Jo3Hg3994vi9y1kEgDFWd/1ch9RWDUgY4HgsvMHFvP93fQ30AvEJCNcd0GTwPID0F8V5eyMVj/qAQCXFqYrRtJh8MIpm2/h7jZ2SsPj0</ESV>' http://192.168.2.189/smart
<?xml version="1.0" encoding="UTF-8"?><ESV>[encrypted content]</ESV>
```

Decrypted response example (zeroed out for privacy):

```xml
<LSV>
    <MAC>00:00:00:00:00:00</MAC>
    <SERIAL>0000000000</SERIAL>
    <CONNECT>ON</CONNECT>
    <STATUS>NORMAL</STATUS>
    <PROFILECODE>
        <VALUE>fc7b013010c900000000000000000000000000000000</VALUE>
        <VALUE>fc7b013010cd00000000000000000000000000000000</VALUE>
        <VALUE>fc7b013010ce00000000000000000000000000000000</VALUE>
        <VALUE>fc7b013010cf00000000000000000000000000000000</VALUE>
        <VALUE>fc7b013010d100000000000000000000000000000000</VALUE>
    </PROFILECODE>
    <DATDATE>2025/07/27 08:18:26</DATDATE>
    <CODE>
        <VALUE>fc620130100200000000000000000000000000000000</VALUE>
        <VALUE>fc620130100300000000000000000000000000000000</VALUE>
        <VALUE>fc620130100400000000000000000000000000000000</VALUE>
        <VALUE>fc620130100500000000000000000000000000000000</VALUE>
        <VALUE>fc620130100600000000000000000000000000000000</VALUE>
        <VALUE>fc620130100900000000000000000000000000000000</VALUE>
    </CODE>
    <APP_VER>37.00</APP_VER>
    <SSL_LIMIT>20371231</SSL_LIMIT>
    <RSSI>-32</RSSI>
    <LED>
        <LED1>0:1,0:1</LED1>
        <LED2>1:5,0:45</LED2>
        <LED3>0:1,0:1</LED3>
        <LED4>1:2,0:2,1:2,0:44</LED4>
    </LED>
    <ECHONET>ON</ECHONET>
</LSV>
```
