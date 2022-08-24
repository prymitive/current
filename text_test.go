package jstream_test

import (
	"testing"

	"github.com/prymitive/jstream"
)

func TestText(t *testing.T) {
	var got store
	for _, tc := range []testCaseT{
		{
			name: "foo",
			iter: jstream.Text(func(s string) {
				got.push(s)
			}),
			body:     `"foo"`,
			expected: []any{"foo"},
		},
		{
			name: "foo",
			iter: jstream.Text(func(s string) {
				got.push(s)
			}),
			body: `foo`,
			err:  "invalid character 'o' in literal false (expecting 'a')",
		},
		{
			name: "123",
			iter: jstream.Text(func(s string) {
				got.push(s)
			}),
			body: `123`,
			err:  "invalid token at offset 3 decoded by Text, 123 is not a string",
		},
		{
			name: "{}",
			iter: jstream.Text(func(s string) {
				got.push(s)
			}),
			body: `{}`,
			err:  "invalid token at offset 1 decoded by Text, { is not a string",
		},
		{
			name: "foo",
			iter: jstream.Text(func(s string) {
				got.push(s)
			}),
			body: `"foo`,
			err:  "unexpected EOF",
		},
	} {
		runTestCase(t, tc, &got)
	}
}
