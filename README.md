# MinecraftProxyGo

This is a Minecraft proxy written in Go. I wrote it to start up a real Computer using Wake-on-LAN if the server is not running. If the server is running, it will connect the joining players to the server.

## Requirements

-   wakeonlan
    -   `sudo apt install wakeonlan`

## Usage

Copy the `config.example.json` to `config.json` and edit it to your needs.

Compile the program with `go build` and run it with `./minecraftproxy`.

## Huge thanks to

A huge thanks goes to [wiki.vg](https://wiki.vg) for their great documentation of the Minecraft protocol.
It helped me a lot to understand the protocol and to implement it.
