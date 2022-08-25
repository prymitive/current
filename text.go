package current

import (
	"encoding/json"
	"fmt"
	"io"
)

func Text(commit func(string)) *text {
	return &text{commit: commit}
}

type text struct {
	commit func(string)
}

func (t text) String() string {
	return "Text"
}

func (t *text) Stream(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}
	if v, ok := tok.(string); ok {
		t.commit(v)
	} else {
		return fmt.Errorf("%w at offset %d decoded by %s, %v is not a string", ErrInvalidToken, dec.InputOffset(), t, tok)
	}
	return io.EOF
}
