// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Splitter golang package implement Read and Write methods to splits and
// combine messages for net.Conn connection.
package splitter

import (
	"bytes"
	"encoding/binary"
	"net"
	"unsafe"

	"github.com/go-errors/errors"
)

// Conn is Splitter data structure and methods receiver
type Conn struct {
	net.Conn
	readBuf          *bytes.Buffer
	sliceCapacityErr bool
}

// dataLenType is type of message header
type dataLenType uint32

var ErrInsufficientCapacity = errors.New("insufficient input slice capacity")

// New creates new splitter object. If the sliceCapacityErr is true than error
// will returned in Read function if input slice insufficient capacitance. If
// the sliceCapacityErr is false and Read function input slice insufficient
// capacitance than data will be trancated.
func New(conn net.Conn, sliceCapacityErr bool) net.Conn {
	c := new(Conn)
	c.Conn = conn
	c.readBuf = new(bytes.Buffer)
	c.sliceCapacityErr = sliceCapacityErr
	return c
}

// Read reads data from the connection.
//
// Read reads data from the connection and split data messages.
// The size of the input slice argument should be enough to read whole message.
// Message will be truncated if it size more than slice capacity when the 'New'
// function sliceCapacityErr argument is false or error will be returnet if
// sliceCapacityErr is true.
//
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (c *Conn) Read(data []byte) (l int, err error) {

	// read reads from conn to buffer
	read := func() (int, error) {
		n, err := c.Conn.Read(data)
		if err != nil {
			return 0, err
		}
		return c.readBuf.Write(data[:n])
	}

	var dataLen dataLenType
	var dataLenSize = int(unsafe.Sizeof(dataLen))

	// Read data from socket and write it to buffer
	for {
		if buflen := c.readBuf.Len(); buflen >= dataLenSize {
			err := binary.Read(c.readBuf, binary.LittleEndian, &dataLen)
			if err == nil && c.readBuf.Len() >= int(dataLen) {
				// There is valid data in buffer, go to read it
				break
			} else {
				// Restore buffer: write dataLen and data to buffer
				// Read data
				data := make([]byte, buflen-dataLenSize)
				c.readBuf.Read(data)

				// Reset buffer and write dataLen and data to buffer
				c.readBuf.Reset()
				binary.Write(c.readBuf, binary.LittleEndian, dataLen)
				c.readBuf.Write(data)
			}
		}

		// Read data from socket and write it to buffer
		_, err := read()
		if err != nil {
			return 0, err
		}
		// log.Printf("read: %v bytes\n", n)
	}

	// Read message from buffer with invalid size.
	// If the output data slice is less than the length of the message data,
	// then read message from buffer and truncate it
	if dataCap := cap(data); int(dataLen) > dataCap {
		d := make([]byte, dataLen)
		_, err = c.readBuf.Read(d)
		if err != nil {
			return 0, err
		}
		if c.sliceCapacityErr {
			return 0, ErrInsufficientCapacity
		}
		copy(data, d[:dataCap])
		return dataCap, nil
	}

	// Read message from buffer with valid size.
	return c.readBuf.Read(data[:dataLen])
}

// Write writes data to the connection.
//
// Write adds message header and write message to connection.
//
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (c *Conn) Write(data []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, dataLenType(len(data)))
	buf.Write(data)
	return c.Conn.Write(buf.Bytes())
}
