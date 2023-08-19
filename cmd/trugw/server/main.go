// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tru unix socket gateway server.
//
// If you can't link the tru package to your application than use this
// standalone unix socket server to communicate with any tru servers.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/teonet-go/trugw/trugw"
)

var sockAddr =  os.TempDir() + "/trugw.sock"
var nomsg = flag.Bool("nomsg", false, "don't show send receive messages")

func main() {
	fmt.Printf("Tru unix socket gateway server, sock path: %s\n", sockAddr)
	flag.Parse()

	// log.SetOutput(io.Discard)

	err := cleanup(sockAddr)
	if err != nil {
		log.Fatal("can't cleanup unix socket, error:", err)
	}

	err = listen(sockAddr)
	if err != nil {
		log.Fatal("can't start listening, error:", err)
	}
}

// cleanup unix socket file
func cleanup(sockAddr string) error {
	if _, err := os.Stat(sockAddr); err == nil {
		if err := os.RemoveAll(sockAddr); err != nil {
			return err
		}
	}
	return nil
}

// listen connections
func listen(sockAddr string) error {
	listener, err := trugw.Listen("tru", sockAddr)
	if err != nil {
		return err
	}

	for {
		conn, _ := listener.Accept()
		go process(conn)
	}
}

// process reads incoming messages and automatically resend it tru peer
func process(conn net.Conn) {
	log.Printf("incoming connection accepted: %v\n", conn)

	// Create slice to read messages. The size of the slice should be anaf to
	// read whole message. Message will be truncated if it size more than buffer
	buf := make([]byte, 256)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error: %v\n", err)
			if err == io.EOF {
				break
			}
			break
		}

		if !*nomsg {
			data := buf[:l]
			log.Printf("got %d bytes from unix socket: %s\n",
				len(data), string(data))
		}

	}
	conn.Close()

	log.Printf("incoming connection closed: %v\n", conn)
}
