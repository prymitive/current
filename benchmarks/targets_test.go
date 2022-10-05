package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	iter "github.com/json-iterator/go"
	"github.com/prymitive/current"
	"github.com/stretchr/testify/require"
)

type Response struct {
	Status string `json:"status"`
	Data   struct {
		ActiveTargets []ActiveTarget `json:"activeTargets"`
	} `json:"data"`
}

type ActiveTarget struct {
	Labels     map[string]string `json:"labels"`
	ScrapePool string            `json:"scrapePool"`
	ScrapeURL  string            `json:"scrapeUrl"`
}

type parserFN func(io.Reader) ([]ActiveTarget, error)

func parseTargetsVanilla(r io.Reader) (targets []ActiveTarget, err error) {
	var dto Response
	err = json.NewDecoder(r).Decode(&dto)
	return dto.Data.ActiveTargets, err
}

func parseTargetsStream(r io.Reader) ([]ActiveTarget, error) {
	mapStart := json.Delim('{')

	dec := json.NewDecoder(r)

	// Expect { as the first token.
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if t != mapStart {
		return nil, fmt.Errorf("expected {, got %v", t)
	}

	results := []ActiveTarget{}
	var inData, inActiveTargets bool
	var aTarget ActiveTarget
	var key string
	// now we need to iterate over input until we reach data->activeTargets
	for dec.More() {
		// we reached the list of active targets so let's decode next target
		// and append it to results
		if inData && inActiveTargets {
			if err = dec.Decode(&aTarget); err != nil {
				return nil, err
			}
			results = append(results, aTarget)
			aTarget.Labels = map[string]string{}
			continue
		}

		// we didn't hit a target yet, decode current token to see what it is
		t, err = dec.Token()
		if err != nil {
			return nil, err
		}

		// if this is { then we're almost there ("data": {)
		if t == mapStart {
			continue
		}

		key = t.(string)
		switch key {
		case "data":
			inData = true
		case "activeTargets":
			inActiveTargets = true
			if _, err = dec.Token(); err != nil {
				return nil, err
			}
		}
	}

	return results, nil
}

func parseTargetsCurrent(r io.Reader) (targets []ActiveTarget, err error) {
	targets = []ActiveTarget{}
	var target ActiveTarget
	decoder := current.Object(
		current.Key("activeTargets", current.Array(&target, func() {
			targets = append(targets, target)
			target.Labels = map[string]string{}
		})),
	)

	dec := json.NewDecoder(r)
	if err = decoder.Stream(dec); err != nil {
		return nil, err
	}

	return targets, nil
}

func parseTargetsGoIter(r io.Reader) (targets []ActiveTarget, err error) {
	var dto Response
	err = iter.NewDecoder(r).Decode(&dto)
	return dto.Data.ActiveTargets, err
}

func BenchmarkTargets(b *testing.B) {
	b.ReportAllocs()

	for _, tc := range []struct {
		name string
		fn   parserFN
	}{
		{name: "vanilla", fn: parseTargetsVanilla},
		{name: "goiter", fn: parseTargetsGoIter},
		{name: "stream", fn: parseTargetsStream},
		{name: "current", fn: parseTargetsCurrent},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				r, err := os.Open("targets.json")
				require.NoError(b, err)
				require.NoError(b, err)
				b.StartTimer()
				targets, err := tc.fn(r)
				b.StopTimer()
				require.NoError(b, err)
				b.ReportMetric(float64(len(targets)), "targets")
			}
		})
	}
}
