package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestArray(t *testing.T) {
	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var got store
	for i, tc := range []testCaseT{
		{
			str: current.Array(func(s *string) {
				got.push(*s)
			}),
			body:     `["bob@example.com", "bob@second.com"]`,
			expected: []any{"bob@example.com", "bob@second.com"},
		},
		{
			str: current.Array(func(n *float64) {
				got.push(int(*n))
			}),
			body:     `[2,3,4,1]`,
			expected: []any{2, 3, 4, 1},
		},
		{
			str: current.Array(func(n *float64) {
				got.push(int(*n))
			}),
			body: `[2,3,4,1[`,
			err:  "expected comma after array element",
		},
		{
			str: current.Array(func(n *float64) {
				got.push(*n)
			}),
			body:     `[2,3,4,1`,
			expected: []any{2.0, 3.0, 4.0, 1.0},
		},
		{
			str: current.Array(func(n *float64) {
				got.push(*n)
			}),
			body: `2,3,4,1`,
			err:  "invalid token at offset 1 decoded by Array[float64], expected [, got 2",
		},
		{
			str: current.Array(func(n *float64) {
				got.push(*n)
			}),
			body: `[2,]`,
			err:  "invalid character ']' looking for beginning of value",
		},
		{
			str: current.Array(func(n *float64) {
				got.push(*n)
			}),
			body: `{"foo":"bar"}`,
			err:  "invalid token at offset 1 decoded by Array[float64], expected [, got {",
		},
		{
			str: current.Array(func(u *user) {
				got.push(*u)
				u.Age = 0
				u.Emails = []string{}
			}),
			body: `[
	{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
	{"age": 23, "name": "alice", "extra": "ignore", "emails": ["alice@example.com", "alias@example.com"]},
	{"name": "deleted"}
]`,
			expected: []any{
				user{Name: "bob", Age: 0, Emails: []string{"bob@example.com"}},
				user{Name: "alice", Age: 23, Emails: []string{"alice@example.com", "alias@example.com"}},
				user{Name: "deleted", Age: 0, Emails: []string{}},
			},
		},
		{
			str: current.Array(func(u *user) {
				got.push(*u)
				u.Age = 0
				u.Emails = []string{}
			}),
			body: `[
	{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
	[],
]`,
			err: "json: cannot unmarshal array into Go value of type current_test.user",
		},
		{
			str: current.Array(func(u *user) {
				got.push(*u)
				u.Age = 0
				u.Emails = []string{}
			}),
			body: `[
	{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
	{[]},
]`,
			err: "invalid character '[' looking for beginning of object key string",
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
