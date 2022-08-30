package current

import (
	"encoding/json"
	"fmt"
)

// Array iterates over a list of objects and will unmarshal each of them onto
// dst. After each unmarshal it will call function commit().
// Use commit to append dst to the slice storing all results
// and to reset any field back to zero state.
// Resetting is especially important if you dst is a struct
// with slices or maps, as not resetting those would cause
// next unmarshall to append to already existing values.
//
// Example:
//
//	type User struct {
//		Name   string
//		Age    int
//		Emails []string
//	}
//
//	r := strings.NewReader(`[
//		{"name": "bob", "age": 40, "emails": ["bob1@example.com","bob2@example.com","bob3@example.com"]},
//	 	{"name": "alice", "age": 34, "emails": ["alice@example.com"]}
//	]`)
//
//	// Create a new json.Decoder for our reader
//	dec := json.NewDecoder(r)
//
//	// user will be used to unmarshal all elements
//	var user User
//	// users will store all decoded users
//	users := []User{}
//
//	streamer := current.Array(
//		&user,
//		func() {
//			users = append(got, user)
//			// We don't need to reset other fields since our JSON body
//			// always passes name & age, so those struct fields are overwritten
//			// but we do need to reset emails slice so it's always empty at the
//			// beginning of unmarshal.
//			user.Emails = []string{}
//		}
//	)
//
//	if err := str.Stream(dec); err != nil {
//	    panic(err)
//	}
func Array[T any](dst *T, commit func()) *array[T] {
	return &array[T]{
		dst:    dst,
		commit: commit,
	}
}

type array[T any] struct {
	dst    *T
	commit func()
	zero   T
}

func (a array[T]) String() string {
	return fmt.Sprintf("Array[%T]", a.zero)
}

func (a *array[T]) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, arrayStart, a); err != nil {
		return err
	}

	for dec.More() {
		err = dec.Decode(a.dst)
		if err != nil {
			return err
		}
		a.commit()
	}

	if err = requireToken(dec, arrayEnd, a); err != nil {
		return err
	}

	return nil
}
