# Super Mario 3D World

Experimental Nextendo server for Super Mario 3D World + Bowser's Fury.

Two local Ryujinx-Nextendo instances reached the same level with game update 1.2.2 using the earlier test build. This cleaned-up version still needs another game test. Internet play and an unchanged official client have not been tested.

The main fix keeps the host's registered public port separate from its LAN UDP port during same-PC testing. The server also delays participation notifications and returns the existing player count on join.

Build with Go 1.23 or newer:

```sh
go build -o server .
```

Set the variables in `example.env` in your shell, then run `server`. Use the same account URL, signing secret and internal key as your Nextendo account service, plus your own TLS certificate and secure password. Route the game's auth hostname (`g2b306a00-lp1.s.n.srv.nintendo.net`) to the auth port.

Signed Nextendo tokens are required normally. For two instances on one PC, set `SM3DW_LOCAL_LAN=1` and list their account PIDs in `SM3DW_LOCAL_PIDS` if the client sends bare PIDs. This mode only starts on loopback and still checks the accounts. Run `go run ./cmd/nncs-local` for the matching local NAT responder; do not run it alongside another NNCS responder.

The NEX dependency points to a small branch of [my fork](https://github.com/alyeri/nextendo-nex/tree/feature/super-mario-3d-world), based on the current official core. Its new options are off by default. No emulator patches, game files, account data or packet captures are included.

The old test relay was removed. This version advertises no relay, which needs checking in the next game test. The shared upstream library still has diagnostic logs, so keep runtime logs private.

Based on [Nextendo Network](https://github.com/NextendoNetwork)'s server layout and account protocol. Original code remains under its license. Debugging and cleanup used AI assistance.
