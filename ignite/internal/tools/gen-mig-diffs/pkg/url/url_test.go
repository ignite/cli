package url

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchesScpLike(t *testing.T) {
	examples := []string{
		// Most-extended case
		"git@github.com:james/bond",
		// Most-extended case with port
		"git@github.com:22:james/bond",
		// Most-extended case with numeric path
		"git@github.com:007/bond",
		// Most-extended case with port and numeric "username"
		"git@github.com:22:007/bond",
		// Single repo path
		"git@github.com:bond",
		// Single repo path with port
		"git@github.com:22:bond",
		// Single repo path with port and numeric repo
		"git@github.com:22:007",
		// Repo path ending with .git and starting with _
		"git@github.com:22:_007.git",
		"git@github.com:_007.git",
		"git@github.com:_james.git",
		"git@github.com:_james/bond.git",
	}

	for _, url := range examples {
		fmt.Println(url)
	}
}

func TestFindScpLikeComponents(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		want       URL
		wantString string
		err        error
	}{
		{
			name:       "https protocol",
			url:        "https://github.com/james/bond",
			wantString: "https://github.com/james/bond.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "https protocol with .git",
			url:        "https://github.com/james/bond.git",
			wantString: "https://github.com/james/bond.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "http protocol",
			url:        "http://github.com/james/bond",
			wantString: "http://github.com/james/bond.git",
			want: URL{
				Protocol: "http",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "http protocol with port",
			url:        "http://github.com/james/bond:8080",
			wantString: "http://github.com/james/bond.git",
			want: URL{
				Protocol: "http",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "https  with numeric path",
			url:        "https://github.com/007/bond",
			wantString: "https://github.com/007/bond.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "007/bond",
			},
		},
		{
			name:       "https with single repo path",
			url:        "https://github.com/bond",
			wantString: "https://github.com/bond.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "bond",
			},
		},
		{
			name:       "https repo path ending with .git and starting with _",
			url:        "https://github.com/_007.git",
			wantString: "https://github.com/_007.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "_007",
			},
		},
		{
			name:       "https repo path ending with .git and starting with _",
			url:        "https://github.com/_james.git",
			wantString: "https://github.com/_james.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "_james",
			},
		},
		{
			name:       "https repo path ending with .git and starting with _",
			url:        "https://github.com/_james/bond.git",
			wantString: "https://github.com/_james/bond.git",
			want: URL{
				Protocol: "https",
				Host:     "github.com",
				Path:     "_james/bond",
			},
		},
		{
			name:       "most-extended case",
			url:        "git@github.com:james/bond",
			wantString: "git@github.com:james/bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "most-extended case with port",
			url:        "git@github.com:22:james/bond",
			wantString: "git@github.com:james/bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "james/bond",
			},
		},
		{
			name:       "most-extended case with numeric path",
			url:        "git@github.com:007/bond",
			wantString: "git@github.com:007/bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "007/bond",
			},
		},
		{
			name:       "most-extended case with port and numeric path",
			url:        "git@github.com:22:007/bond",
			wantString: "git@github.com:007/bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "007/bond",
			},
		},
		{
			name:       "single repo path",
			url:        "git@github.com:bond",
			wantString: "git@github.com:bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "bond",
			},
		},
		{
			name:       "single repo path with port",
			url:        "git@github.com:22:bond",
			wantString: "git@github.com:bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "bond",
			},
		},
		{
			name:       "single repo path with port and numeric path",
			url:        "git@github.com:22:007",
			wantString: "git@github.com:007.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "007",
			},
		},
		{
			name:       "repo path ending with .git and starting with _",
			url:        "git@github.com:22:_007.git",
			wantString: "git@github.com:_007.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "_007",
			},
		},
		{
			name:       "repo path ending with .git, number and starting with _",
			url:        "git@github.com:_007.git",
			wantString: "git@github.com:_007.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "_007",
			},
		},
		{
			name:       "repo path ending with .git and starting with _",
			url:        "git@github.com:_james.git",
			wantString: "git@github.com:_james.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "_james",
			},
		},
		{
			name:       "repo path with .git and starting with _",
			url:        "git@github.com:_james/bond.git",
			wantString: "git@github.com:_james/bond.git",
			want: URL{
				Protocol: "ssh",
				Host:     "github.com",
				Path:     "_james/bond",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.url)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err.Error(), err.Error())
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
			require.EqualValues(t, tt.wantString, got.String())
		})
	}
}
