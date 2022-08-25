package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestText(t *testing.T) {
	var got store
	for i, tc := range []testCaseT{
		{
			body:     `"foo"`,
			expected: []any{"foo"},
		},
		{
			body: `foo`,
			err:  "invalid character 'o' in literal false (expecting 'a')",
		},
		{
			body: `123`,
			err:  "invalid token at offset 3 decoded by Text, 123 is not a string",
		},
		{
			body: `{}`,
			err:  "invalid token at offset 1 decoded by Text, { is not a string",
		},
		{
			str: current.Text(func(s string) {
				got.push(s)
			}),
			body: `"foo`,
			err:  "unexpected EOF",
		},
	} {
		tc := tc
		tc.str = current.Text(func(s string) {
			got.push(s)
		})
		runTestCase(t, i, tc, &got)
	}
}
