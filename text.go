package jstream

import (
	"encoding/json"
	"io"
)

func Text(commit func(string) error) *text {
	return &text{commit: commit}
}

type text struct {
	commit func(string) error
}

func (t text) String() string {
	return "Text"
}

func (t *text) Next(dec *json.Decoder) (err error) {
	var tok json.Token
	if tok, err = dec.Token(); err != nil {
		return err
	}
	if err = t.commit(tok.(string)); err != nil {
		return err
	}
	return io.EOF
}
