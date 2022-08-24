package jstream_test

import (
	"testing"

	"github.com/prymitive/jstream"
)

func TestObject(t *testing.T) {
	var got store
	for _, tc := range []testCaseT{
		{
			name:     "no keys",
			iter:     jstream.Object(func() {}),
			body:     `{"name": "bob"}`,
			expected: []any{},
		},
		{
			name: "]",
			iter: jstream.Object(func() {}),
			body: `]`,
			err:  "invalid character ']' looking for beginning of value",
		},
		{
			name: "{",
			iter: jstream.Object(func() {}),
			body: `{`,
		},
		{
			name: "{}",
			iter: jstream.Object(func() {}),
			body: `{}`,
		},
		{
			name: "name / match",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob"}`,
			expected: []any{"bob"},
		},
		{
			name: "name, age - missing quote",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("age", jstream.Number(func(i float64) {
					got.push(i)
				})),
			),
			body: `{"name": "bob", age: 4}`,
			err:  "invalid character 'a' looking for beginning of object key string",
		},
		{
			name: "name / no match",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
			),
			body:     `{"fullname": "bob"}`,
			expected: []any{},
		},
		{
			name: "name, email / match",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("email", jstream.Text(func(s string) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob", "email": "bob@example.com"}`,
			expected: []any{"bob", "bob@example.com"},
		},
		{
			name: "name / match, email / no match",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("email", jstream.Text(func(s string) {
					got.push(s)
				})),
			),
			body:     `{"name": "bob", "emails": "bob@example.com"}`,
			expected: []any{"bob"},
		},
		{
			name: "name, age",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("age", jstream.Number(func(i float64) {
					got.push(i)
				})),
			),
			body:     `{"name": "bob", "age": 4}`,
			expected: []any{"bob", 4.0},
		},
		{
			name: "name, age - order",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("age", jstream.Number(func(i float64) {
					got.push(i)
				})),
			),
			body:     `{"age": 4, "name": "bob"}`,
			expected: []any{4.0, "bob"},
		},
		{
			name: "name, age / bad number",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("age", jstream.Number(func(i float64) {
					got.push(i)
				})),
			),
			body: `{"name": "bob", "age": "foo"}`,
			err:  `invalid token at offset 28 decoded by Number, "foo" is not a float64`,
		},
		{
			name: "name, age, emails",
			iter: jstream.Object(
				func() {},
				jstream.Key("name", jstream.Text(func(s string) {
					got.push(s)
				})),
				jstream.Key("age", jstream.Number(func(i float64) {
					got.push(i)
				})),
				jstream.Key("emails", jstream.Array(func(s *string) {
					got.push(*s)
				})),
			),
			body:     `{"name": "bob", "emails": ["one", "two"], "age": 0}`,
			expected: []any{"bob", "one", "two", 0.0},
		},
		{
			name: "user -> {}",
			iter: jstream.Object(
				func() {},
				jstream.Key("user", jstream.Object(
					func() {},
				)),
			),
			body: `{"user": []}`,
			err:  "invalid token at offset 10 decoded by Object{}, expected {, got [",
		},
		{
			name: "user -> {age,email}",
			iter: jstream.Object(
				func() {},
				jstream.Key("user", jstream.Object(
					func() {},
					jstream.Key("age", jstream.Number(func(f float64) {})),
					jstream.Key("email", jstream.Text(func(s string) {})),
				)),
			),
			body: `{"user": []}`,
			err:  "invalid token at offset 10 decoded by Object{age,email}, expected {, got [",
		},
	} {
		runTestCase(t, tc, &got)
	}
}
