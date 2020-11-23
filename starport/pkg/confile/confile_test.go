package confile

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	cases := []struct {
		name         string
		ec           EncodingCreator
		hellocontent string
	}{
		{"json", DefaultJSONEncodingCreator, `{"hello":"world"}`},
		{"yaml", DefaultYAMLEncodingCreator, `hello: world`},
		{"toml", DefaultTOMLEncodingCreator, `hello = "world"`},
	}

	type data struct {
		Hello string `json:"hello" yaml:"hello" toml:"hello"`
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			file, err := ioutil.TempFile("", "")
			require.NoError(t, err)
			defer func() {
				file.Close()
				os.Remove(file.Name())
			}()

			_, err = io.Copy(file, strings.NewReader(tt.hellocontent))
			require.NoError(t, err)

			cf := New(tt.ec, file.Name())
			var d data
			require.NoError(t, cf.Load(&d))
			require.Equal(t, "world", d.Hello)

			d.Hello = "cosmos"
			require.NoError(t, cf.Save(d))

			cf2 := New(tt.ec, file.Name())
			var d2 data
			require.NoError(t, cf2.Load(&d2))
			require.Equal(t, "cosmos", d2.Hello)
		})
	}

}
