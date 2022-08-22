package jstream

import (
	"encoding/json"
	"io"
)

func Text(commit func(string)) text {
	return text{commit: commit}
}

type text struct {
	commit func(string)
}

func (t *text) Next(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}
	t.commit(tok.(string))
	return io.EOF
}
