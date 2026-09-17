// Package params provides a set of functions for setting and reading arbitrary
// realm-local parameters that can be called from any realm.
// It may or may not affect system behavior.
package params

func SetString(key string, val string)
func SetBool(key string, val bool)
func SetInt64(key string, val int64)
func SetUint64(key string, val uint64)
func SetBytes(key string, val []byte)
func SetStrings(key string, val []string)
func UpdateParamStrings(key string, val []string, add bool)

// GetXxx reads back a value previously stored under key with the matching
// SetXxx. It returns (value, true) when the realm-local key is set, and
// (zero, false) when nothing is stored there.
//
// Read a key with the getter that matches the setter which wrote it. Stored
// values carry no type tag, so a wrong getter is never reported as false: it
// either panics, ending the transaction, or returns a value the bool cannot be
// told apart from a correct read.
//
// Which of the two happens turns on the stored bytes, not on the pair of
// types, so there is no safe subset to learn. Ways a mismatch reads back
// cleanly, all returning true:
//
//   - GetBytes returns what was stored, unparsed. Every other setter's value
//     arrives as its JSON: a string keeps its quotes, a []string reads
//     ["a","b"].
//   - SetStrings(key, nil) stores JSON null, and null reads back as the zero
//     value of any type at all. SetStrings(key, []string{}) stores [] and
//     panics instead, so the two empty slices behave oppositely.
//   - int64 and uint64 are stored as quoted digits. They read as each other
//     while the value fits, and as a string either way; a string of digits
//     likewise reads back as an int64.
//
// gno test does not reproduce any of this. Its param store holds a Go value
// per key and type-asserts on read, so a mismatch returns false there. Code
// written against that behaves differently once deployed.
func GetString(key string) (string, bool)
func GetBool(key string) (bool, bool)
func GetInt64(key string) (int64, bool)
func GetUint64(key string) (uint64, bool)
func GetBytes(key string) ([]byte, bool)
func GetStrings(key string) ([]string, bool)
