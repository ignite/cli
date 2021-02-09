package msgs

import (
	"fmt"
)

func ExampleDiscover() {
	fmt.Println(Discover("/home/ilker/Documents/code/src/github.com/tendermint/starport/local_test/test"))
	// outputs:
	//   map[github.com/test/test/x/test/types:[MsgUpdateUser MsgUpdateHello MsgCreateUser MsgCreateHello MsgDeleteHello MsgDeleteUser]] <nil>
}
