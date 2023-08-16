// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tru unix socket gateway client.
//
// If you can't link the tru package to your application than use this
// standalone unix socket client to communicate with any tru peers.
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/teonet-go/trugw/trugw"
)

var addr = flag.String("a", ":7070", "tru peer address")
var numMessages = flag.Int("n", 10, "number of messages to send")
var nomsg = flag.Bool("nomsg", false, "don't show send receive messages")

func main() {
	fmt.Printf("Tru unix socket gateway client\n")
	flag.Parse()
	dial()
}

func dial() {
	conn, err := trugw.Dial("tru", *addr)
	if err != nil {
		fmt.Printf("failed to dial: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("start send %d messages ...\n", *numMessages)

	const bufSize = 256

	var wg sync.WaitGroup
	wg.Add(*numMessages)

	// Reader
	go func() {
		for {
			data := make([]byte, bufSize)
			n, err := conn.Read(data)
			if err != nil {
				// fmt.Println("read error: ", err)
				// if err == io.EOF {
				// 	break
				// }
				break
			}
			if !*nomsg {
				fmt.Printf("read %v bytes: %s\n", n, data[:n])
			}
			wg.Done()
		}
		fmt.Printf("connection closed\n")
		os.Exit(2)
	}()

	// Sender
	for i := 1; i <= *numMessages; i++ {
		data := []byte(fmt.Sprintf("Test message %d", i))
		if n, err := conn.Write(data); err != nil {
			fmt.Printf("write error: %v\n", err)
		} else {
			if !*nomsg {
				fmt.Printf("sent %v bytes: %s\n", n, data)
			}
		}
	}
	fmt.Printf("send all messages done\n")

	wg.Wait()
	fmt.Printf("seceive all messages done\n")
	conn.Close()

	select {}
}
