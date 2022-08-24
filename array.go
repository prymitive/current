package current

import (
	"encoding/json"
	"fmt"
	"io"
)

func Array[T any](commit func(*T)) *array[T] {
	return &array[T]{commit: commit}
}

type array[T any] struct {
	pos    position
	commit func(*T)
}

func (a array[T]) String() string {
	// nolint: gocritic
	return fmt.Sprintf("Array[%T]", *new(T))
}

func (a *array[T]) Next(dec *json.Decoder) (err error) {
	switch a.pos {
	case posFirst:
		if err = requireToken(dec, arrayStart, a); err != nil {
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
			a.commit(&elem)
		}
		a.pos = posLast
	case posLast:
		if err = requireToken(dec, arrayEnd, a); err != nil {
			return err
		}
		return io.EOF
	}
	return nil
}
