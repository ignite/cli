package protoutil

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/emicklei/proto"
	"github.com/stretchr/testify/require"
)

// Helpers:
// Only checks containment, not positioning.
func containsElement(f proto.Visitee, v proto.Visitee) bool {
	contains := false
	Apply(f, nil, func(c *Cursor) bool {
		if reflect.DeepEqual(c.Node(), v) {
			contains = true
			return false
		}
		return true
	})
	return contains
}

// parseStringProto takes a string, parses it into a proto.File, and returns a ProtoFile.
// Nodes can be created easily (newnode) by wrapping them correctly. (e.g field in a message)
func parseStringProto(s string) (*proto.Proto, error) {
	p, err := proto.NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Test that the changes from adding a list with starport scaffold list <Type>
// Relatively old files but still exercise some paths of the code.

var (
	genesisProto = `syntax = "proto3";
package cosmonaut.chainname.chainname;

import "gogoproto/gogo.proto";
import "chainname/params.proto";

option go_package = "github.com/cosmonaut/chainname/x/chainname/types";

// GenesisState defines the houhah module's genesis state.
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
}
`
	queryProto = `syntax = "proto3";
package cosmonaut.chainname.chainname;

import "gogoproto/gogo.proto";
import "google/api/annotations.proto";
import "cosmos/base/query/v1beta1/pagination.proto";
import "chainname/params.proto";

option go_package = "github.com/cosmonaut/chainname/x/houhah/types";

// Query defines the gRPC querier service.
service Query {
  // Parameters queries the parameters of the module.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/cosmonaut/chainname/chainname/params";
  }
}

// QueryParamsRequest is request type for the Query/Params RPC method.
message QueryParamsRequest {}

// QueryParamsResponse is response type for the Query/Params RPC method.
message QueryParamsResponse {
  // params holds all the parameters of this module.
  Params params = 1 [(gogoproto.nullable) = false];
}`
	txProto = `syntax = "proto3";
package cosmonautchainname.chainname;

option go_package = "github.com/cosmonaut/chainname/x/chainname/types";

// Msg defines the Msg service.
service Msg {}
`
)

// Test that the changes from adding a list with starport scaffold list <Type>
// are applied correctly to tx.proto
func TestAddEmptyList_tx(t *testing.T) {
	typename, modname := "Kirby", "chainname"
	f, err := parseStringProto(txProto)
	require.NoError(t, err)

	// 1) Add import for the new type (module/lowercase_typ)
	imp := NewImport(fmt.Sprintf("%s/%s.proto", modname, strings.ToLower(typename)))
	err = AddImports(f, true, imp)
	require.NoError(t, err)
	require.True(t, containsElement(f, imp))

	// 2) Add rpcs
	var rpcs []*proto.RPC
	for _, op := range []string{"Create", "Update", "Delete"} {
		rpc := NewRPC(op+typename, "Msg"+op+typename, "Msg"+op+typename+"Response")
		rpcs = append(rpcs, rpc)
	}
	Apply(f, nil, func(c *Cursor) bool {
		// Find the specific service and append.
		if m, ok := c.Node().(*proto.Service); ok {
			if m.Name == "Msg" {
				for _, rpc := range rpcs {
					Append(m, rpc)
				}
				return false // stop
			}
		}
		// Msg will be traversed first.
		// If it was empty, we just stop traversing.
		return true
	})
	for _, rpc := range rpcs {
		require.True(t, containsElement(f, rpc))
	}
	// Add the messages after service Msgs at the end of f.
	createtyp := NewMessage("MsgCreateKirby",
		WithFields(NewField("creator", "string", 1)))
	resp := NewMessage("MsgCreateKirbyResponse",
		WithFields(
			NewField("id", "uint64", 1),
		),
	)
	Append(f, createtyp, resp)
	require.True(t, containsElement(f, createtyp))
	require.True(t, containsElement(f, resp))

	updatetyp := NewMessage("MsgUpdateKirby",
		WithFields(
			NewField("creator", "string", 1),
			NewField("id", "uint64", 2),
		),
	)
	updateResp := NewMessage("MsgUpdateKirbyResponse")
	Append(f, updatetyp, updateResp)
	require.True(t, containsElement(f, updatetyp))
	require.True(t, containsElement(f, updateResp))

	deltyp := NewMessage("MsgDeleteKirby",
		WithFields(
			NewField("creator", "string", 1),
			NewField("id", "uint64", 2),
		),
	)
	delResp := NewMessage("MsgDeleteResponse")
	Append(f, deltyp, delResp)
	require.True(t, containsElement(f, deltyp))
	require.True(t, containsElement(f, delResp))
}

// Test that the changes from adding a list with starport scaffold list <Type>
// are applied correctly to genesis.proto
func TestAddEmptyList_genesis(t *testing.T) {
	typename, modname := "Kirby", "mod"
	f, err := parseStringProto(genesisProto)
	require.NoError(t, err)

	// 1) Add import for the new type (module/lowercase_typ)
	imp := NewImport(fmt.Sprintf("%s/%s.proto", modname, strings.ToLower(typename)))
	err = AddImports(f, true, imp)
	require.NoError(t, err)
	require.True(t, containsElement(f, imp))

	// 2) Add fields to GenesisState. Append.
	Apply(f, nil, func(c *Cursor) bool {
		if m, ok := c.Node().(*proto.Message); ok {
			if m.Name == "GenesisState" {
				lst := NewField(typename+"List", typename, 2,
					WithFieldOptions(NewOption("gogoproto.nullable", "false", Custom())),
					Repeated(),
				)
				field := NewField(typename+"Count", typename, 3)
				Append(m, lst, field)
				require.True(t, containsElement(f, lst))
				require.True(t, containsElement(f, field))
				return false
			}
		}
		return true
	})
}

func TestAddEmptyList_query(t *testing.T) {
	typename, modname := "Kirby", "mod"
	f, err := parseStringProto(queryProto)
	require.NoError(t, err)

	// 1) Add import for the new type (module/lowercase_typ)
	imp := NewImport(fmt.Sprintf("%s/%s.proto", modname, strings.ToLower(typename)))
	err = AddImports(f, true, imp)
	require.NoError(t, err)
	require.True(t, containsElement(f, imp))

	q, err := GetServiceByName(f, "Query")
	require.NoError(t, err)
	// Add the rpcs
	single := NewRPC(typename, "QueryGet"+typename+"Request", "QueryGet"+typename+"Response",
		WithRPCOptions(
			NewOption(
				"google.api.http",
				"/cosmonaut/chainname/chainname/"+typename+"/{id}",
				Custom(),
				SetField("get"),
			),
		),
	)
	all := NewRPC(typename+"All", "QueryAll"+typename+"Request", "QueryAll"+typename+"Response",
		WithRPCOptions(
			NewOption(
				"google.api.http",
				"/cosmonaut/chainname/chainname/"+typename,
				Custom(),
				SetField("get"),
			),
		),
	)
	Append(q, single, all)
	require.True(t, containsElement(f, single))
	require.True(t, containsElement(f, all))
}
