package jstream

import (
	"encoding/json"
	"io"
)

type Text struct {
	commit func(string)
}

func (t *Text) Next(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}
	t.commit(tok.(string))
	return io.EOF
}
