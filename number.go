package current

import (
	"encoding/json"
	"fmt"
	"io"
)

func Number(commit func(float64)) *number {
	return &number{commit: commit}
}

type number struct {
	commit func(float64)
}

func (n number) String() string {
	return "Number"
}

func (n *number) Next(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}
	if v, ok := tok.(float64); ok {
		n.commit(v)
	} else {
		return fmt.Errorf("%w at offset %d decoded by %s, %q is not a float64", ErrInvalidToken, dec.InputOffset(), n, tok)
	}
	return io.EOF
}
