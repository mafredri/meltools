# meltools

meltools is a set of command-line utilities for working with Mitsubishi WiFi adapters (like the MAC-577IF-E).

## Included tools

- **meltelnet** — A Telnet client tailored for Mitsubishi adapters. Handles the quirks of their protocol, supports command history (with up/down arrows), persistent history (`~/.meltelnet_history`), and optional session logging.
- **melsmart** — Utility for sending ESV/CSV-based local control (`/smart`) commands to the adapter.

## Usage

```console
Usage: meltelnet [options] <host>
  -log string
        path to log file

Usage: melsmart [options] <host>
  -enable-echonet
        Enable ECHONET (default false)
  -key string
        AES key
```

Example:

```console
$ ./meltelnet 192.168.1.100
Connected to 192.168.1.100:23

$ ./melsmart 192.168.1.100
```

Type `exit` or press Ctrl+D/Ctrl+C to quit meltelnet.

## Building

Just run:

```
make
```

This will build all available commands in the `cmd/` directory (like `meltelnet` and `melsmart`).

## License

MIT. See LICENSE file.

## Credits

This was made possible thanks to research done out in the open in [info about port 80 in device · Issue #2 · ncaunt/meldec](https://github.com/ncaunt/meldec/issues/2).

Special thanks to [@bacon-5665](https://github.com/bacon-5665) and [@dragonbane0](https://github.com/dragonbane0) for their contributions.

---

If you find bugs or have suggestions, open an issue or PR. This is a hobby project, so expect rough edges.
