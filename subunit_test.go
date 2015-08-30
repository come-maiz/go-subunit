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
	c.Assert(int(signature), check.Equals, 0xb3,
		check.Commentf("Wrong signature"))
}

func (s *SubunitSuite) TestPackageMustContainVersion2Flag(c *check.C) {
	s.stream.Status("dummytest", "dummystatus")
	s.output.Next(1)
	flags := s.output.Next(2)
	version := flags[0] >> 4 // 4 first bits of the first byte.
	c.Assert(version, check.Equals, uint8(0x2), check.Commentf("Wrong version"))
}

func (s *SubunitSuite) TestWithoutIDPackageMustNotSetPresentFlag(c *check.C) {
	s.stream.Status("", "dummystatus")
	s.output.Next(1)
	flags := s.output.Next(2)
	testIDPresent := flags[0] & 0x8 // bit 11 of the first byte.
	c.Assert(testIDPresent, check.Equals, uint8(0x0),
		check.Commentf("Test ID present flag is set"))
}

func (s *SubunitSuite) TestWithIDPackageMustSetPresentFlag(c *check.C) {
	s.stream.Status("test-id", "dummystatus")
	s.output.Next(1)
	flags := s.output.Next(2)
	testIDPresent := flags[0] & 0x8 // bit 11 of the first byte.
	c.Assert(testIDPresent, check.Equals, uint8(0x8),
		check.Commentf("Test ID present flag is not set"))
}

var statustests = []struct {
	status string
	flag   byte
}{
	{"", 0x0},
	{"undefined", 0x0},
	{"exists", 0x1},
	{"inprogress", 0x2},
	{"success", 0x3},
	{"uxsuccess", 0x4},
	{"skip", 0x5},
	{"fail", 0x6},
	{"xfail", 0x7},
}

func (s *SubunitSuite) TestPackageStatusFlag(c *check.C) {
	for _, t := range statustests {
		s.stream.Status("dummytest", t.status)
		s.output.Next(1)
		flags := s.output.Next(2)
		testStatus := flags[1] & 0x7 // Last three bits of the second byte.
		c.Check(testStatus, check.Equals, t.flag, check.Commentf("Wrong status"))
	}
}
