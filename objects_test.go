package current_test

import (
	"testing"

	"github.com/prymitive/current"
	"github.com/stretchr/testify/require"
)

func TestObjects(t *testing.T) {
	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var got store
	var name string
	var age int
	var emails []string
	for i, tc := range []testCaseT{
		{
			str: current.Objects(
				func() {
				},
				current.Object(
					func() {
						u := user{
							Name:   name,
							Age:    age,
							Emails: emails,
						}
						got.push(u)
						name = ""
						age = 0
						emails = []string{}
					},
					current.Key("name", current.Text(func(s string) {
						name = s
					})),
					current.Key("age", current.Number(func(v float64) {
						age = int(v)
					})),
					current.Key("emails", current.Objects(
						func() {},
						current.Text(func(s string) {
							emails = append(emails, s)
						}),
					)),
				),
			),
			body: `[
				{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
				{"age": 23, "name": "alice", "extra": "ignore", "emails": ["alice@example.com", "alias@example.com"]},
				{"name": "deleted"},
				{"name": "bob", "age": 0, "emails": ["bob@example.com"]}
			]`,
			expected: []any{
				user{Name: "bob", Age: 0, Emails: []string{"bob@example.com"}},
				user{Name: "alice", Age: 23, Emails: []string{"alice@example.com", "alias@example.com"}},
				user{Name: "deleted", Age: 0, Emails: []string{}},
				user{Name: "bob", Age: 0, Emails: []string{"bob@example.com"}},
			},
		},
		{
			str: current.Objects(
				func() {
				},
				current.Object(
					func() {
						u := user{
							Name:   name,
							Age:    age,
							Emails: emails,
						}
						got.push(u)
						name = ""
						age = 0
						emails = []string{}
					},
					current.Key("name", current.Text(func(s string) {
						name = s
					})),
					current.Key("age", current.Number(func(v float64) {
						age = int(v)
					})),
					current.Key("emails", current.Array(func(s *string) {
						emails = append(emails, *s)
					})),
				),
			),
			body: `[x]`,
			err:  "invalid character 'x' looking for beginning of value",
		},
		{
			str: current.Objects(
				func() {
				},
				current.Object(
					func() {
						u := user{
							Name:   name,
							Age:    age,
							Emails: emails,
						}
						got.push(u)
						name = ""
						age = 0
						emails = []string{}
					},
					current.Key("name", current.Text(func(s string) {
						name = s
					})),
					current.Key("age", current.Number(func(v float64) {
						age = int(v)
					})),
					current.Key("emails", current.Array(func(s *string) {
						emails = append(emails, *s)
					})),
				),
			),
			body: `{}`,
			err:  "invalid token at offset 1 decoded by []Object{name,age,emails}, expected [, got {",
		},
		{
			str: current.Object(
				func() {},
				current.Key("status", current.Text(func(s string) {
					require.Equal(t, "success", s)
				})),
				current.Key("users", current.Objects(
					func() {
					},
					current.Object(
						func() {
							u := user{
								Name:   name,
								Age:    age,
								Emails: emails,
							}
							got.push(u)
							name = ""
							age = 0
							emails = []string{}
						},
						current.Key("name", current.Text(func(s string) {
							name = s
						})),
						current.Key("age", current.Number(func(v float64) {
							age = int(v)
						})),
						current.Key("emails", current.Array(func(s *string) {
							emails = append(emails, *s)
						})),
					),
				)),
			),
			body: `{
				"status": "success",
				"users": [
					{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
					{"age": 23, "name": "alice", "extra": "ignore", "emails": ["alice@example.com", "alias@example.com"]},
					{"name": "deleted"}
				]
			}`,
			expected: []any{
				user{Name: "bob", Age: 0, Emails: []string{"bob@example.com"}},
				user{Name: "alice", Age: 23, Emails: []string{"alice@example.com", "alias@example.com"}},
				user{Name: "deleted", Age: 0, Emails: []string{}},
			},
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
