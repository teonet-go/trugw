// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Golang net.Conn server sample application.
// This sample application listen unix socket connections and use splitter
// package to split and(or) combine incpming messages.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/teonet-go/trugw/splitter"
)

func main() {
	listen()
}

// listen connections
func listen() {
	addr := "/tmp/trugw"
	socket, err := net.Listen("unix", addr)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}
	log.Printf("start listen at: %s\n", addr)
	defer socket.Close()

	for {
		conn, _ := socket.Accept()
		c := splitter.New(conn, true)
		go process(c)
	}
}

func process(c net.Conn) {
	log.Printf("incoming connection accepted: %v\n", c)

	// Create slice to read messages. The size of the slice should be anaf to
	// read whole message. Message will be truncated if it size more than buffer
	buf := make([]byte, 128)
	for {
		l, err := c.Read(buf)
		if err != nil {
			fmt.Printf("LISTEN: Error: %v\n", err)
			if err == io.EOF {
				break
			}
		}

		data := buf[:l]
		fmt.Printf("LISTEN: received %d bytes from %v\n", len(data), c.LocalAddr())
		fmt.Printf("LISTEN: %v\n", string(data))

	}
	c.Close()

	log.Printf("incoming connection closed: %v\n", c)
}
