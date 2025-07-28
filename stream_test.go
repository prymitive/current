package current_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prymitive/current"
)

type store struct {
	data []any
}

func (s *store) push(v any) {
	s.data = append(s.data, v)
}

func (s *store) reset() {
	s.data = []any{}
}

type testCaseT struct {
	str      current.Streamer
	body     string
	err      string
	expected []any
}

func runTestCase(t *testing.T, index int, tc testCaseT, got *store) {
	t.Run(fmt.Sprintf("%d: %s", index, tc.body), func(t *testing.T) {
		got.reset()
		dec := json.NewDecoder(strings.NewReader(tc.body))
		err := tc.str.Stream(dec)
		if tc.err != "" {
			require.EqualError(t, err, tc.err)
		} else {
			require.NoError(t, err)
			require.ElementsMatch(t, tc.expected, got.data)
		}
	})
}
