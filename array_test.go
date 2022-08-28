package current_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/prymitive/current"
	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var elemS string
	var elemF float64
	var elemI int
	var elemU user
	var got store
	for i, tc := range []testCaseT{
		{
			str: current.Array(&elemS, func() {
				got.push(elemS)
			}),
			body:     `["bob@example.com", "bob@second.com"]`,
			expected: []any{"bob@example.com", "bob@second.com"},
		},
		{
			str: current.Array(&elemI, func() {
				got.push(elemI)
			}),
			body:     `[2,3,4,1]`,
			expected: []any{2, 3, 4, 1},
		},
		{
			str: current.Array(&elemI, func() {
				got.push(elemI)
			}),
			body: `[2,3,4,1[`,
			err:  "expected comma after array element",
		},
		{
			str: current.Array(&elemF, func() {
				got.push(elemF)
			}),
			body:     `[2,3,4,1`,
			expected: []any{2.0, 3.0, 4.0, 1.0},
			err:      "EOF",
		},
		{
			str: current.Array(&elemF, func() {
				got.push(elemF)
			}),
			body: `2,3,4,1`,
			err:  "invalid token at offset 1 decoded by Array[float64], expected [, got 2",
		},
		{
			str: current.Array(&elemF, func() {
				got.push(elemF)
			}),
			body: `[2,]`,
			err:  "invalid character ']' looking for beginning of value",
		},
		{
			str: current.Array(&elemF, func() {
				got.push(elemF)
			}),
			body: `{"foo":"bar"}`,
			err:  "invalid token at offset 1 decoded by Array[float64], expected [, got {",
		},
		{
			str: current.Array(&elemU, func() {
				got.push(elemU)
				elemU.Age = 0
				elemU.Emails = []string{}
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
			str: current.Array(&elemU, func() {
				got.push(elemU)
				elemU.Age = 0
				elemU.Emails = []string{}
			}),
			body: `[
	{"name": "bob", "age": 0, "emails": ["bob@example.com"]},
	[],
]`,
			err: "json: cannot unmarshal array into Go value of type current_test.user",
		},
		{
			str: current.Array(&elemU, func() {
				got.push(elemU)
				elemU.Age = 0
				elemU.Emails = []string{}
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

func BenchmarkArray(b *testing.B) {
	b.ReportAllocs()

	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var u user
	got := []user{}

	bob := user{
		Name:   "bob",
		Age:    40,
		Emails: []string{"bob1@example.com", "bob2@example.com", "bob3@example.com"},
	}

	for _, tc := range []struct {
		title    string
		str      current.Streamer
		body     string
		expected []user
	}{
		{
			title: "[bob*7]",
			str: current.Array(
				&u,
				func() {
					got = append(got, u)
					u.Age = 0
					u.Emails = []string{}
				},
			),
			body: `[
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
				{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]}
			]`,
			expected: []user{bob, bob, bob, bob, bob, bob, bob},
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
				got = got[:0]
			}
		})
	}
}
