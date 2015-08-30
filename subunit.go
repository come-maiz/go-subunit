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

const signature byte = 0xb3

type StreamResultToBytes struct {
	Output io.Writer
}

type packet struct {
}

func (p *packet) write(writer io.Writer) error {
	var b bytes.Buffer
	b.WriteByte(signature)
	_, err := writer.Write(b.Bytes())
	return err
}

func (s *StreamResultToBytes) Status(test_id, test_status string) error {
	p := packet{}
	return p.write(s.Output)
}
