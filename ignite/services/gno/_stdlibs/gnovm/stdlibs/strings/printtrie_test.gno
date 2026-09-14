package strings_test

import (
	"strings"
	"testing"
)

func TestGenericTrieBuilding(t *testing.T) {
	testCases := []struct{ in, out string }{
		{"abc;abdef;abdefgh;xx;xy;z", `-
			a-
			.b-
			..c+
			..d-
			...ef+
			.....gh+
			x-
			.x+
			.y+
			z+
			`},
		{"abracadabra;abracadabrakazam;abraham;abrasion", `-
			a-
			.bra-
			....c-
			.....adabra+
			...........kazam+
			....h-
			.....am+
			....s-
			.....ion+
			`},
		{"aaa;aa;a;i;longerst;longer;long;xx;x;X;Y", `-
			X+
			Y+
			a+
			.a+
			..a+
			i+
			l-
			.ong+
			....er+
			......st+
			x+
			.x+
			`},
		{"foo;;foo;foo1", `+
			f-
			.oo+
			...1+
			`},
	}

	for _, tc := range testCases {
		keys := strings.Split(tc.in, ";")
		args := make([]string, len(keys)*2)
		for i, key := range keys {
			args[i*2] = key
		}

		got := strings.NewReplacer(args...).PrintTrie()
		// Remove tabs from tc.out
		wantbuf := make([]byte, 0, len(tc.out))
		for i := 0; i < len(tc.out); i++ {
			if tc.out[i] != '\t' {
				wantbuf = append(wantbuf, tc.out[i])
			}
		}
		want := string(wantbuf)

		if got != want {
			t.Errorf("PrintTrie(%q)\ngot\n%swant\n%s", tc.in, got, want)
		}
	}
}
