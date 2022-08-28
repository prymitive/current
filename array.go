package current

import (
	"encoding/json"
	"fmt"
)

func Array[T any](dst *T, commit func()) *array[T] {
	return &array[T]{
		dst:    dst,
		commit: commit,
	}
}

type array[T any] struct {
	dst    *T
	commit func()
}

func (a array[T]) String() string {
	return fmt.Sprintf("Array[%T]", *new(T))
}

func (a *array[T]) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, arrayStart, a); err != nil {
		return err
	}

	for dec.More() {
		err = dec.Decode(a.dst)
		if err != nil {
			return err
		}
		a.commit()
	}

	if err = requireToken(dec, arrayEnd, a); err != nil {
		return err
	}

	return nil
}
