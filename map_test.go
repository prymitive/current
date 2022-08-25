package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestMap(t *testing.T) {
	var got store
	for i, tc := range []testCaseT{
		{
			str:  current.Map(func(k string, v int) {}),
			body: `"foo"`,
			err:  "invalid token at offset 5 decoded by Map[int], expected {, got foo",
		},
		{
			str:  current.Map(func(k string, v int) {}),
			body: `[]`,
			err:  "invalid token at offset 1 decoded by Map[int], expected {, got [",
		},
		{
			str:  current.Map(func(k string, v int) {}),
			body: `{{}}`,
			err:  "invalid character '{'",
		},
		{
			str:  current.Map(func(k string, v int) {}),
			body: `{}`,
		},
		{
			str: current.Map(func(k string, v int) {
				got.push(map[string]int{k: v})
			}),
			body:     `{"foo": 1, "bar": 2}`,
			expected: []any{map[string]int{"foo": 1}, map[string]int{"bar": 2}},
		},
		{
			str: current.Map(func(k string, v int) {
				got.push(map[string]int{k: v})
			}),
			body: `{"foo": 1, "bar": "2"}`,
			err:  "json: cannot unmarshal string into Go value of type int",
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
