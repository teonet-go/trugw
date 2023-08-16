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

	"github.com/teonet-go/trugw/trugw"
)

var nomsg = flag.Bool("nomsg", false, "don't show send receive messages")

func main() {
	fmt.Printf("Tru unix socket gateway server\n")
	flag.Parse()

	err := listen()
	if err != nil {
		log.Println("can't start listening, error:", err)
	}
}

// listen connections
func listen() error {
	listener, err := trugw.Listen("tru", "")
	if err != nil {
		return err
	}

	for {
		conn, _ := listener.Accept()
		go process(conn)
	}
}

func process(conn net.Conn) {
	log.Printf("incoming connection accepted: %v\n", conn)

	// Create slice to read messages. The size of the slice should be anaf to
	// read whole message. Message will be truncated if it size more than buffer
	buf := make([]byte, 256)
	for {
		l, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			if err == io.EOF {
				break
			}
			break
		}

		if !*nomsg {
			data := buf[:l]
			fmt.Printf("got %d bytes from linux socket: %s\n",
				len(data), string(data))
		}

	}
	conn.Close()

	log.Printf("incoming connection closed: %v\n", conn)
}
