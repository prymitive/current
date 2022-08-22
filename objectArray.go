package jstream

import (
	"encoding/json"
	"io"
)

func ObjectArray(obj *object) objectArray {
	return objectArray{obj: obj}
}

type objectArray struct {
	obj *object
	pos position
}

func (oa *objectArray) Next(dec *json.Decoder) (err error) {
	switch oa.pos {
	case posFirst:
		if err = requireToken(dec, arrayStart); err != nil {
			return err
		}
		oa.pos = posDecoding
	case posDecoding:
		if err = Stream(dec, oa.obj); err != nil {
			return err
		}
		oa.pos = posLast
	case posLast:
		if err = requireToken(dec, arrayEnd); err != nil {
			return err
		}
		oa.pos = posEOF
	case posEOF:
		return io.EOF
	}
	return nil
}
