package jstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidToken = errors.New("invalid token")

type position int

const (
	posFirst position = iota
	posDecoding
	posLast
	posEOF
)

var (
	arrayStart = json.Delim('[')
	arrayEnd   = json.Delim(']')
	mapStart   = json.Delim('{')
	mapEnd     = json.Delim('}')
)

func requireToken(dec *json.Decoder, expected json.Token, iter Iterator) error {
	got, err := dec.Token()
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("%w at offset %d decoded by %s, expected %s, got %v", ErrInvalidToken, dec.InputOffset(), iter, expected, got)
	}
	return nil
}

type Iterator interface {
	Next(*json.Decoder) error
}

type NamedIterator interface {
	Name() string
	Iterator
}

func Stream(dec *json.Decoder, iter Iterator) (err error) {
	for {
		err = iter.Next(dec)
		switch {
		case errors.Is(err, io.EOF):
			return nil
		case err == nil:
			continue
		default:
			return err
		}
	}
}
