package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestNumber(t *testing.T) {
	var got store
	for i, tc := range []testCaseT{
		{
			body:     `123`,
			expected: []any{123.0},
		},
		{
			body:     `123.5`,
			expected: []any{123.5},
		},
		{
			body: `"123"`,
			err:  `invalid token at offset 5 decoded by Number, "123" is not a float64`,
		},
		{
			body: "{}",
			err:  `invalid token at offset 1 decoded by Number, "{" is not a float64`,
		},
		{
			body:     `1[`,
			expected: []any{1.0},
		},
		{
			body: `"foo`,
			err:  "unexpected EOF",
		},
	} {
		tc := tc
		tc.str = current.Number(func(v float64) {
			got.push(v)
		})
		runTestCase(t, i, tc, &got)
	}
}
