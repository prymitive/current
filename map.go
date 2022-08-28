package jstream

import (
	"encoding/json"
	"fmt"
	"io"
)

func Map[T any](commit func(k string, v T) error) *jmap[T] {
	return &jmap[T]{commit: commit}
}

type jmap[T any] struct {
	pos    position
	commit func(k string, v T) error
}

func (m jmap[T]) String() string {
	// nolint: gocritic
	return fmt.Sprintf("Map[%T]", *new(T))
}

func (m *jmap[T]) Next(dec *json.Decoder) (err error) {
	switch m.pos {
	case posFirst:
		if err = requireToken(dec, mapStart, m); err != nil {
			return err
		}
		m.pos = posDecoding
	case posDecoding:
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
			if err = m.commit(key, val); err != nil {
				return err
			}
		}
	case posLast:
		if err = requireToken(dec, mapEnd, m); err != nil {
			return err
		}
		return io.EOF
	}
	return nil
}
