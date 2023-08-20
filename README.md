# Trugw

Golang trugw creates proxy connection to the Tru peers using unix socket.

If you can't link the [tru](https://github.com/teonet-go/tru) package to your application than use this standalone unix socket server to communicate with any tru peers.

The trugw contains two packages: trugw and splitter, and four examples in cmd folder.

[![GoDoc](https://godoc.org/github.com/teonet-go/trugw?status.svg)](https://godoc.org/github.com/teonet-go/trugw/)
[![Go Report Card](https://goreportcard.com/badge/github.com/teonet-go/trugw)](https://goreportcard.com/report/github.com/teonet-go/trugw)

## How to use

To execute whall example install and run next applications:

Install Tru and this Trugw repo

```shell
mkdir -p ~/go/src/github.com/teonet-go
cd ~/go/src/github.com/teonet-go
git clone https://github.com/teonet-go/tru
git clone https://github.com/teonet-go/trugw
```

Run on different consoles:

```shell
# Tru echo peer
cd ~/go/src/github.com/teonet-go/tru
go run -tags=truStat ./examples/trunet/
```

```shell
# Trugw server
cd ~/go/src/github.com/teonet-go/trugw
go run ./cmd/trugw/server/
```

### Golang client

```shell
# Trugw golang client
cd ~/go/src/github.com/teonet-go/trugw
go run ./cmd/trugw/client -n 5
```

### C, C+ clients

```shell
# Trugw C client
cd ~/go/src/github.com/teonet-go/trugw/cmd/trugw/client-c
gcc main.c ../../../trugw-c/trugw.c -o client-c
./client-c
```

```shell
# Trugw C++ client
cd ~/go/src/github.com/teonet-go/trugw/cmd/trugw/client-cpp
gcc main.cpp ../../../trugw-c/trugw.c -lstdc++  -o client-cpp
./client-cpp
```

## How it works

The `Tru echo peer` start Tru peer running at udp port 7070. It listening tru connection got messages and send answers it to sender.

The `Trugw server` start listening Unix Socket connection, got messages from `Trugw client` and resend it to `Tru echo peer` by Tru connection. When `Trugw server` got answer from `Tru echo peer` it resend it back to `Trugw client` by unix socket.

## License

[BSD](LICENSE)
