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

func (s *SubunitReferenceSuite) SetUpSuite(c *check.C) {
	if !isSubunitInstalled() {
		c.Skip("subunit is not installed")
	}
}

var referencetests = []struct {
	id     string
	status string
}{
	// Status tests.
	{"existing-test", "exists"},
	{"progressing-test", "inprogress"},
	{"successful-test", "success"},
	{"unexpected-successful-test", "uxsuccess"},
	{"skipped-test", "skip"},
	{"failed-test", "fail"},
	{"expected-failed-test", "xfail"},

	// Different test id lengths.
	{"test-id (1 byte)", "exists"},
	{"test-id-with-63-chars (1 byte____)" + strings.Repeat("_", 63-34), "exists"},
	{"test-id-with-64-chars (2 bytes___)" + strings.Repeat("_", 64-34), "exists"},
	{"test-id-with-16383-chars (2 bytes)" + strings.Repeat("_", 16383-34), "exists"},
	{"test-id-with-16384-chars (3 bytes)" + strings.Repeat("_", 16384-34), "exists"},
	// We can't test IDs with more length bytes through the command line.
}

func (s *SubunitReferenceSuite) TestReference(c *check.C) {
	for _, t := range referencetests {
		var goOutput bytes.Buffer
		stream := &subunit.StreamResultToBytes{Output: &goOutput}
		err := stream.Status(subunit.Event{TestID: t.id, Status: t.status})
		c.Check(err, check.IsNil, check.Commentf("Error running the go version of subunit"))

		cmd := exec.Command("python", "-c", fmt.Sprintf(
			// FIXME the runnable flag must be a parameter. --elopio - 2015-08-31
			"import subunit; import sys; subunit.StreamResultToBytes(sys.stdout).status(test_id=%q, test_status=%q, runnable=False)",
			t.id, t.status))
		pythonOutput, err := cmd.Output()
		c.Check(err, check.IsNil, check.Commentf("Error runninng the python version of subunit"))

		c.Check(goOutput.Bytes(), check.DeepEquals, pythonOutput)
	}
}
