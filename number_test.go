package jstream_test

import (
	"testing"

	"github.com/prymitive/jstream"
)

func TestNumber(t *testing.T) {
	var got store
	for _, tc := range []testCaseT{
		{
			name: "123",
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body:     `123`,
			expected: []any{123.0},
		},
		{
			name: "123.5",
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body:     `123.5`,
			expected: []any{123.5},
		},
		{
			name: `"123"`,
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body: `"123"`,
			err:  `invalid token at offset 5 decoded by Number, "123" is not a float64`,
		},
		{
			name: "{}",
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body: "{}",
			err:  `invalid token at offset 1 decoded by Number, "{" is not a float64`,
		},
		{
			name: "1[",
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body:     `1[`,
			expected: []any{1.0},
		},
		{
			name: "foo",
			iter: jstream.Number(func(v float64) {
				got.push(v)
			}),
			body: `"foo`,
			err:  "unexpected EOF",
		},
	} {
		runTestCase(t, tc, &got)
	}
}
