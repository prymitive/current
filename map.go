package current

import (
	"encoding/json"
	"fmt"
	"io"
)

func Map[T any](commit func(k string, v T)) *jmap[T] {
	return &jmap[T]{commit: commit}
}

type jmap[T any] struct {
	commit func(k string, v T)
}

func (m jmap[T]) String() string {
	// nolint: gocritic
	return fmt.Sprintf("Map[%T]", *new(T))
}

func (m *jmap[T]) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, mapStart, m); err != nil {
		return err
	}

	var tok json.Token
	var key string
	var val T
	for {
		if tok, err = dec.Token(); err != nil {
			return err
		}
		if tok == mapEnd {
			return io.EOF
		}
		key = tok.(string)

		if err = dec.Decode(&val); err != nil {
			return err
		}
		m.commit(key, val)
	}
}
