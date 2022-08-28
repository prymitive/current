package current

import (
	"encoding/json"
	"fmt"
)

func Value[T any](commit func(T, bool)) *value[T] {
	return &value[T]{commit: commit}
}

type value[T any] struct {
	commit func(T, bool)
	zero   T
}

func (n value[T]) String() string {
	return fmt.Sprintf("Value[%T]", n.zero)
}

func (n *value[T]) Stream(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}

	if tok == nil {
		n.commit(n.zero, true)
		return nil
	}

	var v T
	var ok bool
	if v, ok = tok.(T); ok {
		n.commit(v, false)
	} else {
		return ErrUnexpectedToken{
			offset: dec.InputOffset(),
			str:    n,
			msg:    fmt.Sprintf("%q is not a %T", tok, n.zero),
		}
	}
	return nil
}
