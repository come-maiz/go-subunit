package subunit_test

import (
	"bytes"
	"testing"

	"launchpad.net/subunit"

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
