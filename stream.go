package current

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidToken = errors.New("invalid token")

type ErrUnexpectedToken struct {
	offset   int64
	str      Streamer
	expected json.Token
	got      json.Token
}

func (ut ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("%s at offset %d decoded by %s, expected %s, got %v", ErrInvalidToken, ut.offset, ut.str, ut.expected, ut.got)
}

var (
	arrayStart = json.Delim('[')
	arrayEnd   = json.Delim(']')
	mapStart   = json.Delim('{')
	mapEnd     = json.Delim('}')
)

func requireToken(dec *json.Decoder, expected json.Token, str Streamer) error {
	got, err := dec.Token()
	if err != nil {
		return err
	}
	if got != expected {
		return ErrUnexpectedToken{
			offset:   dec.InputOffset(),
			str:      str,
			expected: expected,
			got:      got,
		}
	}
	return nil
}

type Streamer interface {
	Stream(*json.Decoder) error
}

type NamedStreamer interface {
	Name() string
	Streamer
}

func Stream(dec *json.Decoder, str Streamer) (err error) {
	err = str.Stream(dec)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
