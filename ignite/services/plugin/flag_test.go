package plugin

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

const (
	flagString1      = "string_flag_1"
	flagString2      = "string_flag_2"
	flagString3      = "string_flag_3"
	flagStringSlice1 = "string_slice_flag_1"
	flagStringSlice2 = "string_slice_flag_2"
	flagStringSlice3 = "string_slice_flag_3"
	flagBool1        = "bool_flag_1"
	flagBool2        = "bool_flag_2"
	flagBool3        = "bool_flag_3"
	flagInt1         = "int_flag_1"
	flagInt2         = "int_flag_2"
	flagInt3         = "int_flag_3"
	flagUint1        = "uint_flag_1"
	flagUint2        = "uint_flag_2"
	flagUint3        = "uint_flag_3"
	flagInt641       = "int64_flag_1"
	flagInt642       = "int64_flag_2"
	flagInt643       = "int64_flag_3"
	flagUint641      = "uint64_flag_1"
	flagUint642      = "uint64_flag_2"
	flagUint643      = "uint64_flag_3"
	flagWrongType1   = "wrong_type_1"
	flagWrongType2   = "wrong_type_2"
	flagWrongType3   = "wrong_type_3"
)

var testFlags = Flags{
	{Name: flagString1, Value: "text_1", DefaultValue: "def_text_1", Type: FlagTypeString},
	{Name: flagString2, DefaultValue: "def_text_2", Type: FlagTypeString},
	{Name: flagString3, Type: FlagTypeString},

	{Name: flagStringSlice1, Value: "slice_1,slice_2", DefaultValue: "slice_1,slice_2,slice_3", Type: FlagTypeStringSlice},
	{Name: flagStringSlice2, DefaultValue: "slice_1,slice_2,slice_3", Type: FlagTypeStringSlice},
	{Name: flagStringSlice3, Type: FlagTypeStringSlice},

	{Name: flagInt1, Value: "-100", DefaultValue: "300", Type: FlagTypeInt},
	{Name: flagInt2, DefaultValue: "200", Type: FlagTypeInt},
	{Name: flagInt3, Type: FlagTypeInt},

	{Name: flagUint1, Value: "22", DefaultValue: "34", Type: FlagTypeUint},
	{Name: flagUint2, DefaultValue: "40", Type: FlagTypeUint},
	{Name: flagUint3, Type: FlagTypeUint},

	{Name: flagInt641, Value: "123", DefaultValue: "641", Type: FlagTypeInt64},
	{Name: flagInt642, DefaultValue: "344", Type: FlagTypeInt64},
	{Name: flagInt643, Type: FlagTypeInt64},

	{Name: flagUint641, Value: "123", DefaultValue: "433333", Type: FlagTypeUint64},
	{Name: flagUint642, DefaultValue: "100000", Type: FlagTypeUint64},
	{Name: flagUint643, Type: FlagTypeUint64},

	{Name: flagBool1, Value: "true", DefaultValue: "false", Type: FlagTypeBool},
	{Name: flagBool2, DefaultValue: "true", Type: FlagTypeBool},
	{Name: flagBool3, Type: FlagTypeBool},

	{Name: flagWrongType1, Value: "text_wrong", DefaultValue: "def_text", Type: FlagTypeUint64},
	{Name: flagWrongType2, DefaultValue: "text_wrong", Type: FlagTypeBool},
	{Name: flagWrongType3, Type: FlagTypeInt},
}

func TestFlags_GetBool(t *testing.T) {
	tests := []struct {
		name string
		key  string
		f    Flags
		want bool
		err  error
	}{
		{
			name: "flag with value",
			key:  flagBool1,
			f:    testFlags,
			want: true,
		},
		{
			name: "flag with default value",
			key:  flagBool2,
			f:    testFlags,
			want: true,
		},
		{
			name: "flag without value and default value",
			key:  flagBool3,
			f:    testFlags,
			err:  errors.New("strconv.ParseBool: parsing \"\": invalid syntax"),
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType1,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_UINT64 for key wrong_type_1"),
		},
		{
			name: "wrong flag value",
			key:  flagWrongType2,
			f:    testFlags,
			err:  errors.New("strconv.ParseBool: parsing \"text_wrong\": invalid syntax"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetBool(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetInt(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want int
		err  error
	}{
		{
			name: "flag with value",
			key:  flagInt1,
			f:    testFlags,
			want: -100,
		},
		{
			name: "flag with default value",
			key:  flagInt2,
			f:    testFlags,
			want: 200,
		},
		{
			name: "flag without value and default value",
			key:  flagInt3,
			f:    testFlags,
			err:  errors.New("strconv.Atoi: parsing \"\": invalid syntax"),
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType2,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_BOOL for key wrong_type_2"),
		},
		{
			name: "wrong flag value without default or value",
			key:  flagWrongType3,
			f:    testFlags,
			err:  errors.New("strconv.Atoi: parsing \"\": invalid syntax"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetInt(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetInt64(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want int64
		err  error
	}{
		{
			name: "flag with value",
			key:  flagInt641,
			f:    testFlags,
			want: 123,
		},
		{
			name: "flag with default value",
			key:  flagInt642,
			f:    testFlags,
			want: 344,
		},
		{
			name: "flag without value and default value",
			key:  flagInt643,
			f:    testFlags,
			err:  errors.New("strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType3,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_INT for key wrong_type_3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetInt64(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetString(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want string
		err  error
	}{
		{
			name: "flag with value",
			key:  flagString1,
			f:    testFlags,
			want: "text_1",
		},
		{
			name: "flag with default value",
			key:  flagString2,
			f:    testFlags,
			want: "def_text_2",
		},
		{
			name: "flag without value and default value",
			key:  flagString3,
			f:    testFlags,
			want: "",
		},
		{
			name: "invalid flag type",
			key:  flagInt1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeInt, flagInt1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType2,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_BOOL for key wrong_type_2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetString(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetStringSlice(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want []string
		err  error
	}{
		{
			name: "flag with default value",
			key:  flagStringSlice1,
			f:    testFlags,
			want: []string{"slice_1", "slice_2"},
		},
		{
			name: "flag with default value",
			key:  flagStringSlice2,
			f:    testFlags,
			want: []string{"slice_1", "slice_2", "slice_3"},
		},
		{
			name: "flag without value and default value",
			key:  flagStringSlice3,
			f:    testFlags,
			want: []string{},
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType1,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_UINT64 for key wrong_type_1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetStringSlice(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetUint(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want uint
		err  error
	}{
		{
			name: "flag with value",
			key:  flagUint1,
			f:    testFlags,
			want: 22,
		},
		{
			name: "flag with default value",
			key:  flagUint2,
			f:    testFlags,
			want: 40,
		},
		{
			name: "flag without value and default value",
			key:  flagUint3,
			f:    testFlags,
			err:  errors.New("strconv.ParseUint: parsing \"\": invalid syntax"),
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType1,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_UINT64 for key wrong_type_1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetUint(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_GetUint64(t *testing.T) {
	tests := []struct {
		name string
		f    Flags
		key  string
		want uint64
		err  error
	}{
		{
			name: "flag with value",
			key:  flagUint641,
			f:    testFlags,
			want: 123,
		},
		{
			name: "flag with default value",
			key:  flagUint642,
			f:    testFlags,
			want: 100000,
		},
		{
			name: "flag without value and default value",
			key:  flagUint643,
			f:    testFlags,
			err:  errors.New("strconv.ParseUint: parsing \"\": invalid syntax"),
		},
		{
			name: "invalid flag type",
			key:  flagString1,
			f:    testFlags,
			err:  errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name: "invalid flag",
			key:  "invalid_key",
			f:    testFlags,
			err:  errors.Wrap(ErrFlagNotFound, "invalid_key"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType1,
			f:    testFlags,
			err:  errors.New("strconv.ParseUint: parsing \"text_wrong\": invalid syntax"),
		},
		{
			name: "wrong flag type",
			key:  flagWrongType3,
			f:    testFlags,
			err:  errors.Wrap(ErrInvalidFlagType, "invalid flag type TYPE_FLAG_INT for key wrong_type_3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.GetUint64(tt.key)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlags_getValue(t *testing.T) {
	tests := []struct {
		name     string
		f        Flags
		key      string
		flagType FlagType
		convFunc func(v string) (interface{}, error)
		want     interface{}
		err      error
	}{
		{
			name:     "valid string conversion",
			f:        testFlags,
			key:      flagString1,
			flagType: FlagTypeString,
			convFunc: func(v string) (interface{}, error) { return v, nil },
			want:     "text_1",
		},
		{
			name:     "valid int conversion",
			f:        testFlags,
			key:      flagInt1,
			flagType: FlagTypeInt,
			convFunc: func(v string) (interface{}, error) { return strconv.Atoi(v) },
			want:     -100,
		},
		{
			name:     "invalid flag type",
			f:        testFlags,
			key:      flagString1,
			flagType: FlagTypeInt,
			convFunc: func(v string) (interface{}, error) { return v, nil },
			err:      errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", FlagTypeString, flagString1),
		},
		{
			name:     "flag not found",
			f:        testFlags,
			key:      "non_existing_flag",
			flagType: FlagTypeString,
			convFunc: func(v string) (interface{}, error) { return v, nil },
			err:      errors.Wrap(ErrFlagNotFound, "non_existing_flag"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.getValue(tt.key, tt.flagType, tt.convFunc)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_flagValue(t *testing.T) {
	tests := []struct {
		name string
		flag *Flag
		want string
	}{
		{
			name: "with value",
			flag: &Flag{Name: flagString1, Value: "actual_value", DefaultValue: "default_value"},
			want: "actual_value",
		},
		{
			name: "with default value",
			flag: &Flag{Name: flagString1, DefaultValue: "default_value"},
			want: "default_value",
		},
		{
			name: "without value and default value",
			flag: &Flag{Name: flagString1},
			want: "",
		},
		{
			name: "number without value and default value",
			flag: &Flag{Name: flagUint642, Type: FlagTypeUint64},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flagValue(tt.flag)
			require.Equal(t, tt.want, got)
		})
	}
}
