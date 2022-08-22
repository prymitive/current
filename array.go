package jstream

import (
	"encoding/json"
	"io"
)

func Array[T any](commit func(*T) error) array[T] {
	return array[T]{commit: commit}
}

type array[T any] struct {
	pos    position
	commit func(*T) error
}

func (a *array[T]) Next(dec *json.Decoder) (err error) {
	switch a.pos {
	case posFirst:
		if err = requireToken(dec, arrayStart); err != nil {
			return err
		}
		a.pos = posDecoding
	case posDecoding:
		var elem T
		for dec.More() {
			err = dec.Decode(&elem)
			if err != nil {
				return err
			}
			if err = a.commit(&elem); err != nil {
				return err
			}
		}
		a.pos = posLast
	case posLast:
		if err = requireToken(dec, arrayEnd); err != nil {
			return err
		}
		a.pos = posEOF
	case posEOF:
		return io.EOF
	}
	return nil
}
