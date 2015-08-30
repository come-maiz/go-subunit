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

package subunit_test

import (
	"bytes"
	"testing"

	"github.com/subunit"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&SubunitSuite{})

type SubunitSuite struct {
	stream *subunit.StreamResultToBytes
	output bytes.Buffer
}

func (s *SubunitSuite) SetUpSuite(c *check.C) {
	s.stream = &subunit.StreamResultToBytes{Output: &s.output}
}

func (s *SubunitSuite) SetUpTest(c *check.C) {
	s.output.Reset()
}

func (s *SubunitSuite) TestPacketMustContainSignature(c *check.C) {
	s.stream.Status("dummytest", "dummystatus")
	signature := s.output.Next(1)[0]
	c.Check(int(signature), check.Equals, 0xb3,
		check.Commentf("Wrong signature"))
}

func (s *SubunitSuite) TestPackegMustContainVersion2Falg(c *check.C) {
	s.stream.Status("dummytest", "dummystatus")
	s.output.Next(1)
	flags := s.output.Next(2)
	version := flags[0] >> 4
	c.Check(version, check.Equals, uint8(0x2), check.Commentf("Wrong version"))
}
