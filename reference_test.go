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
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/elopio/subunit"
	"gopkg.in/check.v1"
)

var _ = check.Suite(&SubunitReferenceSuite{})

type SubunitReferenceSuite struct {
	stream *subunit.StreamResultToBytes
}

func isSubunitInstalled() bool {
	cmd := exec.Command("python", "-c", "import subunit")
	err := cmd.Run()
	return err == nil
}

func makePythonArgs(e subunit.Event) string {
	args := fmt.Sprintf("test_id=%q, test_status=%q, runnable=False", e.TestID, e.Status)
	if !e.Timestamp.IsZero() {
		args += fmt.Sprintf(", timestamp=datetime.datetime(%d, %d, %d, %d, %d, %d, %d, iso8601.Utc())",
			e.Timestamp.Year(), e.Timestamp.Month(), e.Timestamp.Day(), e.Timestamp.Hour(),
			e.Timestamp.Minute(), e.Timestamp.Second(), e.Timestamp.Nanosecond()/1000)
	}
	return args
}

func (s *SubunitReferenceSuite) SetUpSuite(c *check.C) {
	if !isSubunitInstalled() {
		c.Skip("subunit is not installed")
	}
}

var referencetests = []subunit.Event{
	// Status tests.
	{TestID: "existing-test", Status: "exists"},
	{TestID: "progressing-test", Status: "inprogress"},
	{TestID: "successful-test", Status: "success"},
	{TestID: "unexpected-successful-test", Status: "uxsuccess"},
	{TestID: "skipped-test", Status: "skip"},
	{TestID: "failed-test", Status: "fail"},
	{TestID: "expected-failed-test", Status: "xfail"},

	// Different test id lengths.
	{TestID: "test-id (1 byte)", Status: "exists"},
	{TestID: "test-id-with-63-chars (1 byte____)" + strings.Repeat("_", 63-34), Status: "exists"},
	{TestID: "test-id-with-64-chars (2 bytes___)" + strings.Repeat("_", 64-34), Status: "exists"},
	{TestID: "test-id-with-16383-chars (2 bytes)" + strings.Repeat("_", 16383-34), Status: "exists"},
	{TestID: "test-id-with-16384-chars (3 bytes)" + strings.Repeat("_", 16384-34), Status: "exists"},
	// We can't test IDs with more length bytes through the command line.

	// Test with timestamp.
	// Round to microseconds because python's datetime does not accept nanoseconds.
	{TestID: "test-with-timestamp", Status: "success",
		Timestamp: time.Now().UTC().Round(time.Microsecond)},
}

func (s *SubunitReferenceSuite) TestReference(c *check.C) {
	for _, e := range referencetests {
		var goOutput bytes.Buffer
		stream := &subunit.StreamResultToBytes{Output: &goOutput}
		err := stream.Status(e)
		c.Check(err, check.IsNil, check.Commentf("Error running the go version of subunit: %s", err))

		cmd := exec.Command("python", "-c", fmt.Sprintf(
			// FIXME the runnable flag must be a parameter. --elopio - 2015-08-31
			"import datetime; import subunit; import sys; from subunit import iso8601; "+
				"subunit.StreamResultToBytes(sys.stdout).status(%s)",
			makePythonArgs(e)))
		pythonOutput, err := cmd.Output()
		c.Check(err, check.IsNil,
			check.Commentf("Error runninng the python version of subunit: %s", err))

		c.Check(goOutput.Bytes(), check.DeepEquals, pythonOutput,
			check.Commentf("Wrong stream for event %v", e))
	}
}
