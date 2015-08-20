package subunit

import "io"

const signature byte = 0xb3

type StreamResultToBytes struct {
	Output io.Writer
}

func (s StreamResultToBytes) Status(test_id, test_status string) error {
	_, err := s.Output.Write([]byte{signature})
	return err
}
