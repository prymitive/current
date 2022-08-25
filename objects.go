package current

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func Objects(commit func(), str Streamer) *objects {
	return &objects{commit: commit, str: str}
}

type objects struct {
	commit func()
	str    Streamer
}

func (o objects) String() string {
	// nolint: gocritic
	return fmt.Sprintf("[]%s", o.str)
}

func (o *objects) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, arrayStart, o); err != nil {
		return err
	}

	var unErr ErrUnexpectedToken
	for {
		if err = Stream(dec, o.str); err != nil {
			if errors.As(err, &unErr) {
				if unErr.expected == mapStart && unErr.got == arrayEnd {
					return io.EOF
				}
			}
			return err
		}
	}
}
