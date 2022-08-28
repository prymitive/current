package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestMap(t *testing.T) {
	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var got store
	for i, tc := range []testCaseT{
		{
			str:  current.Map(func(k string, v float64) {}),
			body: `"foo"`,
			err:  "invalid token at offset 5 decoded by Map[float64], expected {, got foo",
		},
		{
			str:  current.Map(func(k string, v float64) {}),
			body: `[]`,
			err:  "invalid token at offset 1 decoded by Map[float64], expected {, got [",
		},
		{
			str:  current.Map(func(k string, v float64) {}),
			body: `{{}}`,
			err:  "invalid character '{'",
		},
		{
			str:  current.Map(func(k string, v float64) {}),
			body: `{}`,
		},
		{
			str:  current.Map(func(k string, v float64) {}),
			body: `{]`,
			err:  "invalid character ']'",
		},
		{
			str: current.Map(func(k string, v float64) {
				got.push(map[string]float64{k: v})
			}),
			body:     `{"foo": 1, "bar": 2}`,
			expected: []any{map[string]float64{"foo": 1}, map[string]float64{"bar": 2}},
		},
		{
			str: current.Map(func(k string, v float64) {
				got.push(map[string]float64{k: v})
			}),
			body: `{"foo": 1, "bar": "2"}`,
			err:  `json: cannot unmarshal string into Go value of type float64`,
		},
		{
			str: current.Map(func(k string, v float64) {
				got.push(map[string]float64{k: v})
			}),
			body: `{"foo": 1, "bar": }}`,
			err:  "invalid character '}' looking for beginning of value",
		},
		{
			str: current.Map(func(k string, u user) {
				got.push(k)
				got.push(u)
			}),
			body: `{
				"bob": {"name": "bob", "age": 1, "emails": ["bob@example.com"]},
				"alice": {"name": "alice", "age": 2, "emails": []}
			}`,
			expected: []any{
				"bob",
				user{Name: "bob", Age: 1, Emails: []string{"bob@example.com"}},
				"alice",
				user{Name: "alice", Age: 2, Emails: []string{}},
			},
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
