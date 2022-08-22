package jstream

import (
	"encoding/json"
	"io"
)

func Map[T any](commit func(map[string]T) error) jmap[T] {
	return jmap[T]{}
}

type jmap[T any] struct {
	pos    position
	commit func(map[string]T) error
}

func (m *jmap[T]) Next(dec *json.Decoder) (err error) {
	data := map[string]T{}
	switch m.pos {
	case posFirst:
		if err = requireToken(dec, mapStart); err != nil {
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
				m.pos = posEOF
				break
			}
			key = tok.(string)

			if err = dec.Decode(&val); err != nil {
				return err
			}
			data[key] = val
		}
		m.pos = posLast
	case posLast:
		if err = requireToken(dec, mapEnd); err != nil {
			return err
		}
		m.pos = posEOF
	case posEOF:
		if err = m.commit(data); err != nil {
			return err
		}
		return io.EOF
	}
	return nil
}
