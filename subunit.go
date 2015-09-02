// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

// Package subunit provides a writer of the Subunit v2 protocol.
package subunit

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
)

const (
	signature     byte = 0xb3
	version       byte = 0x2
	testIDPresent byte = 0x8
)

var status = map[string]byte{
	"exists":     0x1,
	"inprogress": 0x2,
	"success":    0x3,
	"uxsuccess":  0x4,
	"skip":       0x5,
	"fail":       0x6,
	"xfail":      0x7,
}

// StreamResultToBytes is an implementation of the StreamResult API that converts calls to bytes.
type StreamResultToBytes struct {
	Output io.Writer
}

type packet struct {
	testID string
	status string
}

func (p *packet) write(writer io.Writer) error {
	// PACKET := SIGNATURE FLAGES PACKET_LENGTH TIMESTAMP? TESTID? TAGS? MIME? FILECONTENT?
	//           ROUTING_CODE? CRC32

	flagsChan := make(chan []byte)
	go p.makeFlags(flagsChan)

	idChan := make(chan []byte)
	go p.makeTestID(idChan)

	// We construct a temporary buffer because we won't know the lenght until it's finished.
	// Then we insert the lenght.
	var bTemp bytes.Buffer
	bTemp.WriteByte(signature)
	bTemp.Write(<-flagsChan)
	bTemp.Write(<-idChan)

	// FIXME Support lenghts of 2, 3 and 4 bytes. --elopio - 2015-08-30
	length := bTemp.Len() + 1 + 4 // Add the size for the length itself and the CRC32.
	// Insert the length.
	var b bytes.Buffer
	b.Write(bTemp.Next(3)) // signature (1 byte) and flags (2 bytes)
	binary.Write(&b, binary.BigEndian, uint8(length))
	b.Write(bTemp.Next(bTemp.Len()))

	// Add the CRC32
	crc := crc32.ChecksumIEEE(b.Bytes())
	binary.Write(&b, binary.BigEndian, crc)

	_, err := writer.Write(b.Bytes())
	return err
}

func (p *packet) makeFlags(c chan<- []byte) {
	flags := make([]byte, 2, 2)
	flags[0] = version << 4
	if p.testID != "" {
		flags[0] = flags[0] | testIDPresent
	}
	flags[1] = flags[1] | status[p.status]
	c <- flags
}

func (p *packet) makeTestID(c chan<- []byte) {
	var testID bytes.Buffer
	if p.testID != "" {
		binary.Write(&testID, binary.BigEndian, uint8(len(p.testID)))
		testID.WriteString(p.testID)
	}
	c <- testID.Bytes()
}

// Status informs the result about a test status.
func (s *StreamResultToBytes) Status(testID, testStatus string) error {
	p := packet{testID: testID, status: testStatus}
	return p.write(s.Output)
}
