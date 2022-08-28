package current

import (
	"encoding/json"
	"fmt"
)

func Value[T any](commit func(T)) *value[T] {
	return &value[T]{commit: commit}
}

type value[T any] struct {
	commit func(T)
}

func (n value[T]) String() string {
	return fmt.Sprintf("Value[%T]", *new(T))
}

func (n *value[T]) Stream(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}

	if tok == nil {
		var zero T
		n.commit(zero)
		return nil
	}

	if v, ok := tok.(T); ok {
		n.commit(v)
	} else {
		return fmt.Errorf("%w at offset %d decoded by %s, %q is not a %T", ErrInvalidToken, dec.InputOffset(), n, tok, *new(T))
	}
	return nil
}
