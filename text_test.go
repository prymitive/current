package current_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/prymitive/current"
	"github.com/stretchr/testify/require"
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

func BenchmarkText(b *testing.B) {
	b.ReportAllocs()

	var got string

	for _, tc := range []struct {
		title    string
		str      current.Streamer
		body     string
		expected string
	}{
		{
			str: current.Text(func(s string) {
				got = s
			}),
			body:     `"foo"`,
			expected: "foo",
		},
		{
			str:      current.Text(func(s string) {}),
			body:     `"foo"`,
			expected: "foo",
		},
		{
			str: current.Text(func(s string) {
				got = s
			}),
			body:     `"foo bar baz alice bob"`,
			expected: "foo bar baz alice bob",
		},
	} {
		b.Run(tc.title, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				r := strings.NewReader(tc.body)
				dec := json.NewDecoder(r)
				var err error
				b.StartTimer()
				err = tc.str.Stream(dec)
				b.StopTimer()
				require.NoError(b, err)
				require.Equal(b, tc.expected, got)
			}
		})
	}
}
