package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestValue(t *testing.T) {
	var got store
	for i, tc := range []testCaseT{
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body:     `123`,
			expected: []any{123.0},
		},
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body:     `123.5`,
			expected: []any{123.5},
		},
		{
			str: current.Value(func(v string, isNull bool) {
				got.push(v)
			}),
			body:     `"foo bar"`,
			expected: []any{"foo bar"},
		},
		{
			str: current.Value(func(v any, isNull bool) {
				got.push(v)
			}),
			body:     `null`,
			expected: []any{nil},
		},
		{
			str: current.Value(func(v int, isNull bool) {
				if isNull {
					got.push(nil)
				} else {
					got.push(v)
				}
			}),
			body:     `null`,
			expected: []any{nil},
		},
		{
			str: current.Value(func(v, isNull bool) {
				got.push(v)
			}),
			body:     `true`,
			expected: []any{true},
		},
		{
			str: current.Value(func(v, isNull bool) {
				got.push(v)
			}),
			body:     `false`,
			expected: []any{false},
		},
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body: `"123"`,
			err:  `invalid token at offset 5 decoded by Value[float64], "123" is not a float64`,
		},
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body: "{}",
			err:  `invalid token at offset 1 decoded by Value[float64], "{" is not a float64`,
		},
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body:     `1[`,
			expected: []any{1.0},
		},
		{
			str: current.Value(func(v float64, isNull bool) {
				got.push(v)
			}),
			body: `"foo`,
			err:  "unexpected EOF",
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
