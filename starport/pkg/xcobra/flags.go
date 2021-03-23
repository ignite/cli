package xcobra

import (
	"github.com/spf13/pflag"
)

// FlagValueEmptiness represent empty/filled flag value statuses.
type FlagValueEmptiness int

const (
	// FlagsAllEmpty indicates that all compared flag values are empty.
	FlagsAllEmpty FlagValueEmptiness = iota

	// FlagsSomeEmpty indicates that some compared flag values are empty.
	FlagsSomeEmpty

	// FlagsAllFilled indicates that all compared flag values are filled.
	FlagsAllFilled
)

// CheckFlagValues checks the emptiness of given flags with names.
func CheckFlagValues(flags *pflag.FlagSet, names []string) (emptiness FlagValueEmptiness, err error) {
	var existsCount int

	for _, name := range names {
		flag := flags.Lookup(name)
		if flag == nil || flag.Value.String() == "" {
			continue
		}

		existsCount++
	}

	switch {
	case existsCount == 0:
		return FlagsAllEmpty, nil

	case existsCount != len(names):
		return FlagsSomeEmpty, nil

	default:
		return FlagsAllFilled, nil
	}
}
