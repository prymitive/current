# current - JSON streaming wrapper for encoding/json

## Unmarshalling JSON in Go 

Go has many ways of dealing with JSON using standard library.
[go.dev/blog/json](https://go.dev/blog/json) covers it pretty well.

The most common is to unmarshal JSON blob into a struct:

```Go
type User struct {
    Name      string `json:"name"`
    Nicknames []string `json:"nicknames"`
    Age       int `json:"age"`
}

var user User
var b []byte
b = get(...)
json.Unmarshal(b, &user)
```

This handles all of the complexity like dealing with different types etc.

When dealing with HTTP responses there is a Decoder we can use to read our response
body, so we don't have to:

```Go
r, err := http.Get(url)
defer r.Body.Close()
var user User
json.NewDecoder(r.Body).Decode(user)
```

In most cases this is enough, but what if you're dealing with a JSON content that
is a list of users instead just decoding a single user onto a single struct,
especially if that list of users is very long?

Calling `Unmarshal` on the entire list would need the whole response to be read
into memory and decoded in a single call, so it would require a lot of memory.

This is what `json.NewDecoder(r.Body).Decode(user)` would do - it will read
the whole response body and then call `Unmarshal` on it.

If that's just one user then it's perfectly fine, but if we had a huge list
of users to unmarshal this way and we have limited memory available, then
ideally we would want to avoid having to read the entire response into memory.

Let's assume this is our JSON response to parse:

```JSON
[
    {"name": "bob", "nicknames": ["b", "bobby"], "age": 5},
    {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
    {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
    {...}
]
```

The more users there are the bigger the response will be and the more memory will
be needed to parse it.

[go.dev/blog/json](https://go.dev/blog/json) mentions streaming, but the example
provided is for reading JSON objects pushed into a long lived connection, rather
that a single structure of a HTTP response, what it expects is:

```JSON
{"name": "bob", "nicknames": ["b", "bobby"], "age": 5},
{"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
{"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
{...}
```

If our HTTP response looked like that we could simply call `Decode` in a loop
until there is no more objects to decode:

```Go
dec := json.NewDecoder(r.Body)
for dec.More() {
    var user User
    err := dec.Decode(&user)
}
```

Each `Decode` call will handle next JSON object:

```JSON
{"name": "bob", "nicknames": ["b", "bobby"], "age": 5},          <-- Decode()
{"name": "...", "nicknames": ["...", "...", "..."], "age": ...}, <-- Decode()
{"name": "...", "nicknames": ["...", "...", "..."], "age": ...}, <-- Decode()
{...}
```

`More()` will peek at what would be the next token and when it finds `]` it will return
`false`.

The problem is that we have our array tokens wrapping all objects (`[...]`).
Luckily we can move ourselves in JSON stream by calling `Token` instead of `Decode`.
`Decode` call tries to decode next JSON token in a stream onto provided struct.
`Token` on the other hand simply reads the next token and returns it, it's then our
job to do something with it.
We can use `Token` to navigate in the stream until we're in the right place and then
start decoding users.

So with our users list we would want to read first `[` using `Token` call so that
we're in front of our first user, then we hand over decoding to our loop:

```JSON
[
    We want to be here <---
    {"name": "bob", "nicknames": ["b", "bobby"], "age": 5},
    {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
    {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
    {...}
]
```

```Go
for dec.More() {
    var user User
    err := dec.Decode(&user)
}
```

So our final code would be:

```Go
dec := json.NewDecoder(r.Body)

// first token should be our array opening
t, err := dec.Token()
if t != json.Delim('[') {
    panic("Expected [, got %s", t)
}

var users []user
for dec.More() {
    var user User
    err := dec.Decode(&user)
    // append user to a slice so we can do something with it
    users = append(users, user)
}

// we're done with last user, so we should get array end token next
t, err := dec.Token()
if t != json.Delim(']') {
    panic("Expected ], got %s", t)
}
```

Great success! This code is fairly simple to write and, because we don't load
the whole body at once into memory, we limit how much memory is needed to parse
even a very big file.

The only problem is that navigating JSON streams can be very error prone, especially
for deeply nested JSON blobs. Just imagine that our response is a lot less flat:

```JSON
{
    "status": "ok",
    "response": {
        "data": {
            "users": [
                {"name": "bob", "nicknames": ["b", "bobby"], "age": 5},
                {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
                {"name": "...", "nicknames": ["...", "...", "..."], "age": ...},
                {...}
            ]
        }
    }
}
```

We still need to navigate just before the first user object before we can start
our `Decode` loop, but now we need to issue multiple `Token` calls and keep track
where we are in the stream.

## github.com/prymitive/current

`current` is helper package that allows you to easily navigate in a JSON stream, so
you don't have to manually write all those `Token` calls and keep track of your
position. It also uses generics for decoding simple fields, like `"status": "ok"`
in the example above.

Thanks to current you can parse large JSON responses using a streaming decoder
and keep your memory usage low.

## Benchmarks

My initial experience with streaming JSON responses was when I tried to parse
a huge (170MB) JSON response from [Prometheus targets endpoint](https://prometheus.io/docs/prometheus/latest/querying/api/#targets).
Doing so required ~800MB of heap and caused my service to run out of memory
and crash. Writing manual streaming code allowed me to make it work with
very little memory.

Here is a benchmark using a similar sized targets response.
It's broken down by used parser, you can find complete code in the [benchmarks] folder.

```go
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
```

### encoding/json

```go
func parseTargetsVanilla(r io.Reader) (targets []ActiveTarget, err error) {
	var resp Response
	err = json.NewDecoder(r).Decode(&resp)
	return resp.Data.ActiveTargets, err
}
```

### encoding/json with manual streaming code

```go
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
```

### prymitive/current (this package)

```go
func parseTargetsCurrent(r io.Reader) (targets []ActiveTarget, err error) {
	targets = []ActiveTarget{}
	var target ActiveTarget
	decoder := current.Object(
		func() {},
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
```

### github.com/json-iterator/go

```go
func parseTargetsGoIter(r io.Reader) (targets []ActiveTarget, err error) {
	var resp Response
	err = iter.NewDecoder(r).Decode(&resp)
	return resp.Data.ActiveTargets, err
}
```

### Benchmark code

```go
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
```

### Benchmark results

Running this benchmark against a real-world response from Prometheus with 157466 targets (165MB JSON file):

```
BenchmarkTargets/vanilla-8         	       1	2062807784 ns/op	    157466 targets	830797392 B/op	 7181439 allocs/op
BenchmarkTargets/goiter-8          	       1	1467535850 ns/op	    157466 targets	432962160 B/op	14387618 allocs/op
BenchmarkTargets/stream-8          	       1	2158223629 ns/op	    157466 targets	304209216 B/op	 7181324 allocs/op
BenchmarkTargets/current-8         	       1	2099241602 ns/op	    157466 targets	304203704 B/op	 7181534 allocs/op
```

Passed through [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) to make
it more human readable:

```
name               time/op
Targets/vanilla-8  2.13s ±10%
Targets/goiter-8   1.34s ± 1%
Targets/stream-8   1.96s ± 1%
Targets/current-8  2.21s ±19%

name               targets
Targets/vanilla-8   157k ± 0%
Targets/goiter-8    157k ± 0%
Targets/stream-8    157k ± 0%
Targets/current-8   157k ± 0%

name               alloc/op
Targets/vanilla-8  831MB ± 0%
Targets/goiter-8   433MB ± 0%
Targets/stream-8   304MB ± 0%
Targets/current-8  304MB ± 0%

name               allocs/op
Targets/vanilla-8  7.18M ± 0%
Targets/goiter-8   14.4M ± 0%
Targets/stream-8   7.18M ± 0%
Targets/current-8  7.18M ± 0%
```

Running same benchmark against mock JSON file generated using
[examples/benchmarks/mock.go](examples/benchmarks/mock.go) shows similar results:

```
BenchmarkTargets/vanilla-8         	       1	2202727471 ns/op	    200000 targets	753842168 B/op	 7200237 allocs/op
BenchmarkTargets/goiter-8          	       1	1381775190 ns/op	    200000 targets	351971712 B/op	14794136 allocs/op
BenchmarkTargets/stream-8          	       1	2385830895 ns/op	    200000 targets	236299544 B/op	 7200079 allocs/op
BenchmarkTargets/current-8         	       1	2342428154 ns/op	    200000 targets	236299848 B/op	 7200085 allocs/op
```

Steps to reproduce:

```SHELL
cd benchmarks
go run mock.go
go test -run=none -bench=. -benchmem .
```
