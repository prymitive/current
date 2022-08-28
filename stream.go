package current

import (
	"encoding/json"
	"fmt"
)

type ErrUnexpectedToken struct {
	offset int64
	str    Streamer
	msg    string
}

func (ut ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("invalid token at offset %d decoded by %s, %s", ut.offset, ut.str, ut.msg)
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
			offset: dec.InputOffset(),
			str:    str,
			msg:    fmt.Sprintf("expected %s, got %v", expected, got),
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
