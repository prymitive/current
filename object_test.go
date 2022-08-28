package current_test

import (
	"testing"

	"github.com/prymitive/current"
)

func TestObject(t *testing.T) {
	var got store
	var email string
	for i, tc := range []testCaseT{
		{
			str:      current.Object(func() {}),
			body:     `{"name": "bob"}`,
			expected: []any{},
		},
		{
			str:  current.Object(func() {}),
			body: `]`,
			err:  "invalid character ']' looking for beginning of value",
		},
		{
			str:  current.Object(func() {}),
			body: `{`,
			err:  "EOF",
		},
		{
			str:  current.Object(func() {}),
			body: `{}`,
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob"}`,
			expected: []any{"bob"},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i float64, isNull bool) {
					got.push(i)
				})),
			),
			body: `{"name": "bob", age: 4}`,
			err:  "invalid character 'a' looking for beginning of object key string",
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
			),
			body:     `{"fullname": "bob"}`,
			expected: []any{},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("email", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob", "email": "bob@example.com"}`,
			expected: []any{"bob", "bob@example.com"},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("email", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob", "emails": "bob@example.com"}`,
			expected: []any{"bob"},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i float64, isNull bool) {
					got.push(i)
				})),
			),
			body:     `{"name": "bob", "age": 4}`,
			expected: []any{"bob", 4.0},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i float64, isNull bool) {
					got.push(i)
				})),
			),
			body:     `{"age": 4, "name": "bob"}`,
			expected: []any{4.0, "bob"},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i int, isNull bool) {
					got.push(i)
				})),
			),
			body: `{"name": "bob", "age": "foo"}`,
			err:  `invalid token at offset 28 decoded by Value[int], "foo" is not a int`,
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i float64, isNull bool) {
					got.push(i)
				})),
				current.Key("emails", current.Array(&email, func() {
					got.push(email)
				})),
			),
			body:     `{"name": "bob", "emails": ["one", "two"], "age": 0}`,
			expected: []any{"bob", "one", "two", 0.0},
		},
		{
			str: current.Object(
				func() {},
				current.Key("name", current.Value(func(s string, isNull bool) {
					got.push(s)
				})),
				current.Key("age", current.Value(func(i float64, isNull bool) {
					got.push(i)
				})),
				current.Key("emails", current.Array(&email, func() {
					got.push(email)
				})),
			),
			body: `
					{"name": "bob", "emails": ["one", "two"], "age": 5},
					{"name": "not", "emails": ["three"], "age": 0}`,
			expected: []any{"bob", "one", "two", 5.0},
		},
		{
			str: current.Object(
				func() {},
				current.Key("user", current.Object(
					func() {},
				)),
			),
			body: `{"user": []}`,
			err:  "invalid token at offset 10 decoded by Object{}, expected {, got [",
		},
		{
			str: current.Object(
				func() {},
				current.Key("user", current.Object(
					func() {},
					current.Key("age", current.Value(func(f float64, isNull bool) {})),
					current.Key("email", current.Value(func(s string, isNull bool) {})),
				)),
			),
			body: `{"user": []}`,
			err:  "invalid token at offset 10 decoded by Object{age,email}, expected {, got [",
		},
	} {
		runTestCase(t, i, tc, &got)
	}
}
