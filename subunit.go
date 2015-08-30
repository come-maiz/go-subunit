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

type StreamResultToBytes struct {
	Output io.Writer
}

type packet struct {
	testID string
	status string
}

func (p *packet) write(writer io.Writer) error {
	var b bytes.Buffer
	b.WriteByte(signature)
	b.Write(p.makeFlags())
	_, err := writer.Write(b.Bytes())
	return err
}

func (p *packet) makeFlags() []byte {
	flags := make([]byte, 2, 2)
	flags[0] = version << 4
	if p.testID != "" {
		flags[0] = flags[0] | testIDPresent
	}
	flags[1] = flags[0] | status[p.status]
	return flags
}

func (s *StreamResultToBytes) Status(testID, testStatus string) error {
	p := packet{testID: testID, status: testStatus}
	return p.write(s.Output)
}
