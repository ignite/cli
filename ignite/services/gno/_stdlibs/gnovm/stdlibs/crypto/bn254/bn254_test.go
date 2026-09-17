// G1Add tests use verbatim vectors from:
//   go-ethereum bn256Add.json: https://github.com/ethereum/go-ethereum/blob/master/core/vm/testdata/precompiles/bn256Add.json
//   (test case names "chfast1", "chfast2", "cdetrio*")
//
// PairingCheck tests use verbatim vectors from:
//   go-ethereum bn256Pairing.json: https://github.com/ethereum/go-ethereum/blob/master/core/vm/testdata/precompiles/bn256Pairing.json
//   (test case names "jeff1"–"jeff6", "empty_data", "one_point", "two_point_match_*")
//
// G1Mul: go-ethereum has no separate bn256Mul.json; scalars 0/1/2 are verified via
// gnark-crypto and the G1Add identity (2*G == G+G). (EIP-196: https://eips.ethereum.org/EIPS/eip-196)

package bn254

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
)

// g1Generator returns the EIP-196 canonical G1 generator as a 64-byte (x|y)
// buffer: x=1, y=2.
func g1Generator() []byte {
	out := make([]byte, 64)
	out[31] = 1
	out[63] = 2
	return out
}

func fpModulusBytes() []byte {
	out := make([]byte, 32)
	fpModulus.FillBytes(out)
	return out
}

// TestG1AddKnownVectors uses verbatim Input/Expected pairs from
// go-ethereum/core/vm/testdata/precompiles/bn256Add.json.
// Per EIP-196, short inputs are right-padded with zeros; excess bytes are ignored.
func TestG1AddKnownVectors(t *testing.T) {
	cases := []struct {
		name     string
		input    string // hex, from JSON "Input" field (may be shorter than 128 bytes)
		expected string // hex, 64-byte G1 point
	}{
		// chfast1/2: non-trivial point addition
		{
			name:     "chfast1",
			input:    "18b18acfb4c2c30276db5411368e7185b311dd124691610c5d3b74034e093dc9063c909c4720840cb5134cb9f59fa749755796819658d32efc0d288198f3726607c2b7f58a84bd6145f00c9c2bc0bb1a187f20ff2c92963a88019e7c6a014eed06614e20c147e940f2d70da3f74c9a17df361706a4485c742bd6788478fa17d7",
			expected: "2243525c5efd4b9c3d3c45ac0ca3fe4dd85e830a4ce6b65fa1eeaee202839703301d1d33be6da8e509df21cc35964723180eed7532537db9ae5e7d48f195c915",
		},
		{
			name:     "chfast2",
			input:    "2243525c5efd4b9c3d3c45ac0ca3fe4dd85e830a4ce6b65fa1eeaee202839703301d1d33be6da8e509df21cc35964723180eed7532537db9ae5e7d48f195c91518b18acfb4c2c30276db5411368e7185b311dd124691610c5d3b74034e093dc9063c909c4720840cb5134cb9f59fa749755796819658d32efc0d288198f37266",
			expected: "2bd3e6d0f3b142924f5ca7b49ce5b9d54c4703d7ae5648e61d02268b1a0a9fb721611ce0a6af85915e2f1d70300909ce2e49dfad4a4619c8390cae66cefdb204",
		},
		// cdetrio1: 0+0 = 0 (128 bytes, both identity)
		{
			name:     "cdetrio1",
			input:    "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		// cdetrio4: empty input → both points zero-padded → 0
		{
			name:     "cdetrio4",
			input:    "",
			expected: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		// cdetrio6: 0 + G = G (identity first)
		{
			name:     "cdetrio6",
			input:    "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
			expected: "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		},
		// cdetrio7: G + 0 = G (identity second)
		{
			name:     "cdetrio7",
			input:    "000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		},
		// cdetrio8: 64-byte input (G only), second point zero-padded → G + 0 = G
		{
			name:     "cdetrio8",
			input:    "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
			expected: "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
		},
		// cdetrio11: G + G = 2G
		{
			name:     "cdetrio11",
			input:    "0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
			expected: "030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd315ed738c0e0a7c92e7845f96b2ae9c0a68a6a449e3538fc7ff3ebf7a5a18a2c4",
		},
		// cdetrio13: P + Q = R (non-trivial)
		{
			name:     "cdetrio13",
			input:    "17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa901e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c039730ea8dff1254c0fee9c0ea777d29a9c710b7e616683f194f18c43b43b869073a5ffcc6fc7a28c30723d6e58ce577356982d65b833a5a5c15bf9024b43d98",
			expected: "15bf2bb17880144b5d1cd2b1f46eff9d617bffd1ca57c37fb5a49bd84e53cf66049c797f9ce0d17083deb32b5e36f2ea2a212ee036598dd7624c168993d1355f",
		},
		// cdetrio14: P + (-P) = 0
		{
			name:     "cdetrio14",
			input:    "17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa901e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa92e83f8d734803fc370eba25ed1f6b8768bd6d83887b87165fc2434fe11a830cb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var input []byte
			if tc.input != "" {
				var err error
				input, err = hex.DecodeString(tc.input)
				if err != nil {
					t.Fatalf("hex decode input: %v", err)
				}
			}
			want, err := hex.DecodeString(tc.expected)
			if err != nil {
				t.Fatalf("hex decode expected: %v", err)
			}

			got := X_g1Add(input)
			if got == nil {
				t.Fatalf("X_g1Add returned nil")
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("mismatch:\n got  %x\n want %x", got, want)
			}
		})
	}
}

func TestG1AddInvalidInputs(t *testing.T) {
	// Coordinate exactly equal to the field modulus must be rejected.
	p := fpModulusBytes()
	input := make([]byte, 128)
	copy(input[0:32], p)
	input[63] = 2
	copy(input[64:128], g1Generator())
	if X_g1Add(input) != nil {
		t.Fatalf("expected nil for unreduced coordinate")
	}

	// Point not on curve: (1, 3) doesn't satisfy y^2 = x^3 + 3.
	input = make([]byte, 128)
	input[31] = 1
	input[63] = 3
	copy(input[64:128], g1Generator())
	if X_g1Add(input) != nil {
		t.Fatalf("expected nil for off-curve point")
	}
}

func TestG1MulKnownScalars(t *testing.T) {
	// 0 * G = identity (0, 0).
	input := make([]byte, 96)
	copy(input[0:64], g1Generator())
	got := X_g1Mul(input)
	want := make([]byte, 64)
	if !bytes.Equal(got, want) {
		t.Fatalf("0*G != 0: got %x", got)
	}

	// 1 * G = G.
	input = make([]byte, 96)
	copy(input[0:64], g1Generator())
	input[95] = 1
	got = X_g1Mul(input)
	if !bytes.Equal(got, g1Generator()) {
		t.Fatalf("1*G != G: got %x", got)
	}

	// 2 * G must match the G1Add result.
	addInput := make([]byte, 128)
	copy(addInput[0:64], g1Generator())
	copy(addInput[64:128], g1Generator())
	wantDouble := X_g1Add(addInput)

	input = make([]byte, 96)
	copy(input[0:64], g1Generator())
	input[95] = 2
	gotDouble := X_g1Mul(input)
	if !bytes.Equal(gotDouble, wantDouble) {
		t.Fatalf("2*G via mul != G+G via add:\n mul: %x\n add: %x", gotDouble, wantDouble)
	}
}

func TestG1MulInvalidInputs(t *testing.T) {
	if X_g1Mul(make([]byte, 95)) != nil {
		t.Fatalf("expected nil for wrong-length input")
	}
	if X_g1Mul(nil) != nil {
		t.Fatalf("expected nil for nil input")
	}

	input := make([]byte, 96)
	input[31] = 1
	input[63] = 3
	input[95] = 5
	if X_g1Mul(input) != nil {
		t.Fatalf("expected nil for off-curve point")
	}
}

// TestPairingCheckKnownVectors uses verbatim Input/Expected pairs from
// go-ethereum/core/vm/testdata/precompiles/bn256Pairing.json.
// Each entry is the EVM precompile encoding: concatenated (G1 64B | G2 128B) pairs.
// Expected is a 32-byte big-endian 1 (check passes) or 0 (check fails).
func TestPairingCheckKnownVectors(t *testing.T) {
	cases := []struct {
		name     string
		input    string // hex, from JSON "Input" field
		wantByte byte   // last byte of Expected: 0x01 = true, 0x00 = false
	}{
		// empty_data: product of zero pairings = 1 in GT (EIP-197 §3)
		{
			name:     "empty_data",
			input:    "",
			wantByte: 1,
		},
		// jeff1: 2-pair valid equality check
		{
			name:     "jeff1",
			input:    "1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 1,
		},
		// jeff2
		{
			name:     "jeff2",
			input:    "2eca0c7238bf16e83e7a1e6c5d49540685ff51380f309842a98561558019fc0203d3260361bb8451de5ff5ecd17f010ff22f5c31cdf184e9020b06fa5997db841213d2149b006137fcfb23036606f848d638d576a120ca981b5b1a5f9300b3ee2276cf730cf493cd95d64677bbb75fc42db72513a4c1e387b476d056f80aa75f21ee6226d31426322afcda621464d0611d226783262e21bb3bc86b537e986237096df1f82dff337dd5972e32a8ad43e28a78a96a823ef1cd4debe12b6552ea5f06967a1237ebfeca9aaae0d6d0bab8e28c198c5a339ef8a2407e31cdac516db922160fa257a5fd5b280642ff47b65eca77e626cb685c84fa6d3b6882a283ddd1198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 1,
		},
		// jeff3
		{
			name:     "jeff3",
			input:    "0f25929bcb43d5a57391564615c9e70a992b10eafa4db109709649cf48c50dd216da2f5cb6be7a0aa72c440c53c9bbdfec6c36c7d515536431b3a865468acbba2e89718ad33c8bed92e210e81d1853435399a271913a6520736a4729cf0d51eb01a9e2ffa2e92599b68e44de5bcf354fa2642bd4f26b259daa6f7ce3ed57aeb314a9a87b789a58af499b314e13c3d65bede56c07ea2d418d6874857b70763713178fb49a2d6cd347dc58973ff49613a20757d0fcc22079f9abd10c3baee245901b9e027bd5cfc2cb5db82d4dc9677ac795ec500ecd47deee3b5da006d6d049b811d7511c78158de484232fc68daf8a45cf217d1c2fae693ff5871e8752d73b21198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 1,
		},
		// jeff4
		{
			name:     "jeff4",
			input:    "2f2ea0b3da1e8ef11914acf8b2e1b32d99df51f5f4f206fc6b947eae860eddb6068134ddb33dc888ef446b648d72338684d678d2eb2371c61a50734d78da4b7225f83c8b6ab9de74e7da488ef02645c5a16a6652c3c71a15dc37fe3a5dcb7cb122acdedd6308e3bb230d226d16a105295f523a8a02bfc5e8bd2da135ac4c245d065bbad92e7c4e31bf3757f1fe7362a63fbfee50e7dc68da116e67d600d9bf6806d302580dc0661002994e7cd3a7f224e7ddc27802777486bf80f40e4ca3cfdb186bac5188a98c45e6016873d107f5cd131f3a3e339d0375e58bd6219347b008122ae2b09e539e152ec5364e7e2204b03d11d3caa038bfc7cd499f8176aacbee1f39e4e4afc4bc74790a4a028aff2c3d2538731fb755edefd8cb48d6ea589b5e283f150794b6736f670d6a1033f9b46c6f5204f50813eb85c8dc4b59db1c5d39140d97ee4d2b36d99bc49974d18ecca3e7ad51011956051b464d9e27d46cc25e0764bb98575bd466d32db7b15f582b2d5c452b36aa394b789366e5e3ca5aabd415794ab061441e51d01e94640b7e3084a07e02c78cf3103c542bc5b298669f211b88da1679b0b64a63b7e0e7bfe52aae524f73a55be7fe70c7e9bfc94b4cf0da1213d2149b006137fcfb23036606f848d638d576a120ca981b5b1a5f9300b3ee2276cf730cf493cd95d64677bbb75fc42db72513a4c1e387b476d056f80aa75f21ee6226d31426322afcda621464d0611d226783262e21bb3bc86b537e986237096df1f82dff337dd5972e32a8ad43e28a78a96a823ef1cd4debe12b6552ea5f",
			wantByte: 1,
		},
		// jeff5
		{
			name:     "jeff5",
			input:    "20a754d2071d4d53903e3b31a7e98ad6882d58aec240ef981fdf0a9d22c5926a29c853fcea789887315916bbeb89ca37edb355b4f980c9a12a94f30deeed30211213d2149b006137fcfb23036606f848d638d576a120ca981b5b1a5f9300b3ee2276cf730cf493cd95d64677bbb75fc42db72513a4c1e387b476d056f80aa75f21ee6226d31426322afcda621464d0611d226783262e21bb3bc86b537e986237096df1f82dff337dd5972e32a8ad43e28a78a96a823ef1cd4debe12b6552ea5f1abb4a25eb9379ae96c84fff9f0540abcfc0a0d11aeda02d4f37e4baf74cb0c11073b3ff2cdbb38755f8691ea59e9606696b3ff278acfc098fa8226470d03869217cee0a9ad79a4493b5253e2e4e3a39fc2df38419f230d341f60cb064a0ac290a3d76f140db8418ba512272381446eb73958670f00cf46f1d9e64cba057b53c26f64a8ec70387a13e41430ed3ee4a7db2059cc5fc13c067194bcc0cb49a98552fd72bd9edb657346127da132e5b82ab908f5816c826acb499e22f2412d1a2d70f25929bcb43d5a57391564615c9e70a992b10eafa4db109709649cf48c50dd2198a1f162a73261f112401aa2db79c7dab1533c9935c77290a6ce3b191f2318d198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 1,
		},
		// jeff6: pairing inequality — expected 0
		{
			name:     "jeff6",
			input:    "1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c103188585e2364128fe25c70558f1560f4f9350baf3959e603cc91486e110936198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 0,
		},
		// one_point: single (G1, G2) pair — not equal to 1, returns 0
		{
			name:     "one_point",
			input:    "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 0,
		},
		// two_point_match_2
		{
			name:     "two_point_match_2",
			input:    "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d",
			wantByte: 1,
		},
		// two_point_match_3
		{
			name:     "two_point_match_3",
			input:    "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002203e205db4f19b37b60121b83a7333706db86431c6d835849957ed8c3928ad7927dc7234fd11d3e8c36c59277c3e6f149d5cd3cfa9a62aee49f8130962b4b3b9195e8aa5b7827463722b8c153931579d3505566b4edf48d498e185f0509de15204bb53b8977e5f92a0bc372742c4830944a59b4fe6b1c0466e2a6dad122b5d2e030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd31a76dae6d3272396d0cbe61fced2bc532edac647851e3ac53ce1cc9c7e645a83198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantByte: 1,
		},
		// two_point_match_4
		{
			name:     "two_point_match_4",
			input:    "105456a333e6d636854f987ea7bb713dfd0ae8371a72aea313ae0c32c0bf10160cf031d41b41557f3e7e3ba0c51bebe5da8e6ecd855ec50fc87efcdeac168bcc0476be093a6d2b4bbf907172049874af11e1b6267606e00804d3ff0037ec57fd3010c68cb50161b7d1d96bb71edfec9880171954e56871abf3d93cc94d745fa114c059d74e5b6c4ec14ae5864ebe23a71781d86c29fb8fb6cce94f70d3de7a2101b33461f39d9e887dbb100f170a2345dde3c07e256d1dfa2b657ba5cd030427000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000021a2c3013d2ea92e13c800cde68ef56a294b883f6ac35d25f587c09b1b3c635f7290158a80cd3d66530f74dc94c94adb88f5cdb481acca997b6e60071f08a115f2f997f3dbd66a7afe07fe7862ce239edba9e05c5afff7f8a1259c9733b2dfbb929d1691530ca701b4a106054688728c9972c8512e9789e9567aae23e302ccd75",
			wantByte: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var input []byte
			if tc.input != "" {
				var err error
				input, err = hex.DecodeString(tc.input)
				if err != nil {
					t.Fatalf("hex decode: %v", err)
				}
			}

			got := X_pairingCheck(input)
			if len(got) != 32 {
				t.Fatalf("expected 32-byte output, got %d", len(got))
			}
			if got[31] != tc.wantByte {
				t.Fatalf("last byte: got 0x%02x, want 0x%02x", got[31], tc.wantByte)
			}
			// All bytes except the last must be zero.
			for i := range 31 {
				if got[i] != 0 {
					t.Fatalf("byte[%d] = 0x%02x, want 0x00", i, got[i])
				}
			}
		})
	}
}

func TestPairingCheckInvalidLengths(t *testing.T) {
	if X_pairingCheck(make([]byte, 191)) != nil {
		t.Fatalf("expected nil for non-192-multiple input")
	}
	if X_pairingCheck(make([]byte, 193)) != nil {
		t.Fatalf("expected nil for non-192-multiple input")
	}
}

func TestPairingCheckRejectsUnreducedCoordinates(t *testing.T) {
	_, _, _, g2Gen := bn254.Generators() //nolint:dogsled // gnark-crypto returns (g1Gen, g1GenAff, g2Gen, g2GenAff); only g2Gen is needed here.
	g2Marshal := g2Gen.Marshal()

	p := fpModulusBytes()
	input := make([]byte, 0, 192)
	input = append(input, p...)
	input = append(input, []byte{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}...)
	input = append(input, g2Marshal...)

	if X_pairingCheck(input) != nil {
		t.Fatalf("expected nil for unreduced G1.x coordinate")
	}
}

func TestG1AddBoundaryCoordinate(t *testing.T) {
	// p-1 is a reduced coordinate, but (p-1, 0) is not on the curve.
	pMinusOne := new(big.Int).Sub(fpModulus, big.NewInt(1))
	buf := make([]byte, 32)
	pMinusOne.FillBytes(buf)

	input := make([]byte, 128)
	copy(input[0:32], buf)

	if X_g1Add(input) != nil {
		t.Fatalf("expected nil for off-curve point with reduced coordinate")
	}
}
