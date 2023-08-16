// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Trugw golang package creates proxy connection to the tru peers using unix
// socket
package trugw

import (
	"fmt"
	"net"

	"github.com/teonet-go/tru"
	"github.com/teonet-go/trugw/splitter"
)

// A Listener is a generic network listener for stream-oriented protocols.
type Listner struct {
	net.Listener // Unix socket listener
}

// Conn is trugw Servers generic stream-oriented network connection.
type Conn struct {
	net.Conn              // Unix socket connect
	tru      *tru.Tru     // Tru object
	ch       *tru.Channel // Tru channel
}

// Dial connects to the address on the named network.
func Dial(network, address string) (net.Conn, error) {

	// Create linux socket connection
	addr := "/tmp/trugw"
	conn, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}
	conn = splitter.New(conn, false)

	// Send tru address to server
	conn.Write([]byte(address))

	return conn, nil
}

// Listen announces on the local network address.
func Listen(network, address string) (net.Listener, error) {

	// Start listening unix socket
	addr := "/tmp/trugw"
	listener, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}

	return &Listner{Listener: listener}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l Listner) Accept() (net.Conn, error) {

	// Wait next net.Listner connection
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	// Create new splitter and add it and Tru object to net.Conn
	conn = splitter.New(conn, false)
	conn = &Conn{Conn: conn}

	return conn, nil
}

// Close closes unix socket and tru connections.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c Conn) Close() error {
	if c.ch != nil {
		c.ch.Close()
	}
	return c.Conn.Close()
}

// Read reads data from the connection.
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(b []byte) (n int, err error) {

	// Create tru connection at first read
	if c.ch == nil {
		// Read tru peer address
		n, err := c.Conn.Read(b)
		if err != nil {
			return 0, err
		}
		address := string(b[:n])
		fmt.Printf("got TRU address: %s\n", address)

		// Create tru object
		tru, err := tru.New(0)
		if err != nil {
			return 0, err
		}

		// Create tru connection
		ch, err := tru.Connect(address, c.reader)
		if err != nil {
			return 0, err
		}
		c.ch = ch
		fmt.Printf("connected to TRU peer: %s\n", ch.Addr().String())
	}

	// Read message from unix socket, from client
	n, err = c.Conn.Read(b)
	if err != nil {
		return 0, err
	}

	// Resend message to TRU
	c.ch.WriteTo(b[:n])

	return
}

// reader receive TRU messages and send it to linux socket
func (c Conn) reader(ch *tru.Channel, pac *tru.Packet, err error) (processed bool) {
	if err != nil {
		fmt.Printf("got TRU err: %v\n", err)
		return
	}
	// fmt.Printf("got %d bytes from tru: %s\n", pac.Len(), pac.Data())

	// Resend message to unix socket
	c.Write(pac.Data())

	return true
}
