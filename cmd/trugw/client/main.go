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
	"log"
	"os"
	"sync"

	"github.com/teonet-go/trugw/trugw"
)

var sockAddr =  os.TempDir() + "/trugw.sock"
var addr = flag.String("a", ":7070", "tru peer address")
var numMessages = flag.Int("n", 10, "number of messages to send")
var nomsg = flag.Bool("nomsg", false, "don't show send receive messages")

func main() {
	fmt.Printf("Tru unix socket gateway client\n")
	flag.Parse()

	// log.SetOutput(io.Discard)

	dial(sockAddr)
}

func dial(sockAddr string) {
	conn, err := trugw.Dial(sockAddr, *addr)
	if err != nil {
		log.Printf("failed to dial: %v\n", err)
		os.Exit(1)
	}
	log.Printf("connection established\n")

	const bufSize = 256

	var wg sync.WaitGroup
	wg.Add(*numMessages)

	// Reader
	go func() {
		for {
			data := make([]byte, bufSize)
			n, err := conn.Read(data)
			if err != nil {
				// log.Println("read error: ", err)
				// if err == io.EOF {
				// 	break
				// }
				break
			}
			if !*nomsg {
				log.Printf("read %v bytes: %s\n", n, data[:n])
			}
			wg.Done()
		}
		log.Printf("connection closed\n")
		os.Exit(2)
	}()

	// Sender
	log.Printf("start send %d messages ...\n", *numMessages)
	for i := 1; i <= *numMessages; i++ {
		data := []byte(fmt.Sprintf("Test message %d", i))
		if n, err := conn.Write(data); err != nil {
			log.Printf("write error: %v\n", err)
		} else {
			if !*nomsg {
				log.Printf("sent %v bytes: %s\n", n, data)
			}
		}
	}
	log.Printf("send all messages done\n")

	wg.Wait()
	log.Printf("seceive all messages done\n")
	conn.Close()

	select {}
}
