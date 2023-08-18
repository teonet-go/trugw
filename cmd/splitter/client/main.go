// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Golang net.Conn client sample application.
// This sample application connect to unix socket and use splitter package to
// split outgoing messages.
package main

import (
	"fmt"
	"net"
	"os"

	"github.com/teonet-go/trugw/splitter"
)

func main() {
	dial()
}

func dial() {
	addr := "/tmp/trugw"
	conn, err := net.Dial("unix", addr)
	if err != nil {
		fmt.Printf("Failed to dial: %v\n", err)
		os.Exit(1)
	}

	c := splitter.New(conn, false)
	for i := 1; i <= 10; i++ {
		if n, err := c.Write([]byte(fmt.Sprintf("Test message %d", i))); err != nil {
			fmt.Printf("DIAL: Write error: %v\n", err)
		} else {
			fmt.Printf("DIAL: Success, sent %v bytes\n", n)
		}
	}
	conn.Close()
}
