// Test vectors sourced from:
//   EIP-198: https://github.com/ethereum/go-ethereum/blob/master/core/vm/testdata/precompiles/modexp.json

package modexp

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
	"time"
)

func TestModExp_EIP198_Samples(t *testing.T) {
	cases := []struct {
		name     string
		inputHex string
		wantHex  string
	}{
		{
			name:     "eip_example1",
			inputHex: "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002003fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2efffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f",
			wantHex:  "0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name:     "nagydani-1-square",
			inputHex: "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000040e09ad9675465c53a109fac66a445c91b292d2bb2c5268addb30cd82f80fcb0033ff97c80a5fc6f39193ae969c6ede6710a6b7ac27078a06d90ef1c72e5c85fb502fc9e1f6beb81516545975218075ec2af118cd8798df6e08a147c60fd6095ac2bb02c2908cf4dd7c81f11c289e4bce98f3553768f392a80ce22bf5c4f4a248c6b",
			wantHex:  "60008f1614cc01dcfb6bfb09c625cf90b47d4468db81b5f8b7a39d42f332eab9b2da8f2d95311648a8f243f4bb13cfb3d8f7f2a3c014122ebb3ed41b02783adc",
		},
		{
			name:     "nagydani-1-pow0x10001",
			inputHex: "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000040e09ad9675465c53a109fac66a445c91b292d2bb2c5268addb30cd82f80fcb0033ff97c80a5fc6f39193ae969c6ede6710a6b7ac27078a06d90ef1c72e5c85fb5010001fc9e1f6beb81516545975218075ec2af118cd8798df6e08a147c60fd6095ac2bb02c2908cf4dd7c81f11c289e4bce98f3553768f392a80ce22bf5c4f4a248c6b",
			wantHex:  "c36d804180c35d4426b57b50c5bfcca5c01856d104564cd513b461d3c8b8409128a5573e416d0ebe38f5f736766d9dc27143e4da981dfa4d67f7dc474cbee6d2",
		},
		{
			name:     "nagydani-3-qube",
			inputHex: "000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000100c9130579f243e12451760976261416413742bd7c91d39ae087f46794062b8c239f2a74abf3918605a0e046a7890e049475ba7fbb78f5de6490bd22a710cc04d30088179a919d86c2da62cf37f59d8f258d2310d94c24891be2d7eeafaa32a8cb4b0cfe5f475ed778f45907dc8916a73f03635f233f7a77a00a3ec9ca6761a5bbd558a2318ecd0caa1c5016691523e7e1fa267dd35e70c66e84380bdcf7c0582f540174e572c41f81e93da0b757dff0b0fe23eb03aa19af0bdec3afb474216febaacb8d0381e631802683182b0fe72c28392539850650b70509f54980241dc175191a35d967288b532a7a8223ce2440d010615f70df269501944d4ec16fe4a3cb03d7a85909174757835187cb52e71934e6c07ef43b4c46fc30bbcd0bc72913068267c54a4aabebb493922492820babdeb7dc9b1558fcf7bd82c37c82d3147e455b623ab0efa752fe0b3a67ca6e4d126639e645a0bf417568adbb2a6a4eef62fa1fa29b2a5a43bebea1f82193a7dd98eb483d09bb595af1fa9c97c7f41f5649d976aee3e5e59e2329b43b13bea228d4a93f16ba139ccb511de521ffe747aa2eca664f7c9e33da59075cc335afcd2bf3ae09765f01ab5a7c3e3938ec168b74724b5074247d200d9970382f683d6059b94dbc336603d1dfee714e4b447ac2fa1d99ecb4961da2854e03795ed758220312d101e1e3d87d5313a6d052aebde75110363d",
			wantHex:  "1b280ecd6a6bf906b806d527c2a831e23b238f89da48449003a88ac3ac7150d6a5e9e6b3be4054c7da11dd1e470ec29a606f5115801b5bf53bc1900271d7c3ff3cd5ed790d1c219a9800437a689f2388ba1a11d68f6a8e5b74e9a3b1fac6ee85fc6afbac599f93c391f5dc82a759e3c6c0ab45ce3f5d25d9b0c1bf94cf701ea6466fc9a478dacc5754e593172b5111eeba88557048bceae401337cd4c1182ad9f700852bc8c99933a193f0b94cf1aedbefc48be3bc93ef5cb276d7c2d5462ac8bb0c8fe8923a1db2afe1c6b90d59c534994a6a633f0ead1d638fdc293486bb634ff2c8ec9e7297c04241a61c37e3ae95b11d53343d4ba2b4cc33d2cfa7eb705e",
		},
		{
			// 1024-byte base and modulus: exactly at maxOperandLen.
			name:     "nagydani-5-pow0x10001",
			inputHex: "000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000400c5a1611f8be90071a43db23cc2fe01871cc4c0e8ab5743f6378e4fef77f7f6db0095c0727e20225beb665645403453e325ad5f9aeb9ba99bf3c148f63f9c07cf4fe8847ad5242d6b7d4499f93bd47056ddab8f7dee878fc2314f344dbee2a7c41a5d3db91eff372c730c2fdd3a141a4b61999e36d549b9870cf2f4e632c4d5df5f024f81c028000073a0ed8847cfb0593d36a47142f578f05ccbe28c0c06aeb1b1da027794c48db880278f79ba78ae64eedfea3c07d10e0562668d839749dc95f40467d15cf65b9cfc52c7c4bcef1cda3596dd52631aac942f146c7cebd46065131699ce8385b0db1874336747ee020a5698a3d1a1082665721e769567f579830f9d259cec1a836845109c21cf6b25da572512bf3c42fd4b96e43895589042ab60dd41f497db96aec102087fe784165bb45f942859268fd2ff6c012d9d00c02ba83eace047cc5f7b2c392c2955c58a49f0338d6fc58749c9db2155522ac17914ec216ad87f12e0ee95574613942fa615898c4d9e8a3be68cd6afa4e7a003dedbdf8edfee31162b174f965b20ae752ad89c967b3068b6f722c16b354456ba8e280f987c08e0a52d40a2e8f3a59b94d590aeef01879eb7a90b3ee7d772c839c85519cbeaddc0c193ec4874a463b53fcaea3271d80ebfb39b33489365fc039ae549a17a9ff898eea2f4cb27b8dbee4c17b998438575b2b8d107e4a0d66ba7fca85b41a58a8d51f191a35c856dfbe8aef2b00048a694bbccff832d23c8ca7a7ff0b6c0b3011d00b97c86c0628444d267c951d9e4fb8f83e154b8f74fb51aa16535e498235c5597dac9606ed0be3173a3836baa4e7d756ffe1e2879b415d3846bccd538c05b847785699aefde3e305decb600cd8fb0e7d8de5efc26971a6ad4e6d7a2d91474f1023a0ac4b78dc937da0ce607a45974d2cac1c33a2631ff7fe6144a3b2e5cf98b531a9627dea92c1dc82204d09db0439b6a11dd64b484e1263aa45fd9539b6020b55e3baece3986a8bffc1003406348f5c61265099ed43a766ee4f93f5f9c5abbc32a0fd3ac2b35b87f9ec26037d88275bd7dd0a54474995ee34ed3727f3f97c48db544b1980193a4b76a8a3ddab3591ce527f16d91882e67f0103b5cda53f7da54d489fc4ac08b6ab358a5a04aa9daa16219d50bd672a7cb804ed769d218807544e5993f1c27427104b349906a0b654df0bf69328afd3013fbe430155339c39f236df5557bf92f1ded7ff609a8502f49064ec3d1dbfb6c15d3a4c11a4f8acd12278cbf68acd5709463d12e3338a6eddb8c112f199645e23154a8e60879d2a654e3ed9296aa28f134168619691cd2c6b9e2eba4438381676173fc63c2588a3c5910dc149cf3760f0aa9fa9c3f5faa9162b0bf1aac9dd32b706a60ef53cbdb394b6b40222b5bc80eea82ba8958386672564cae3794f977871ab62337cf010001e30049201ec12937e7ce79d0f55d9c810e20acf52212aca1d3888949e0e4830aad88d804161230eb89d4d329cc83570fe257217d2119134048dd2ed167646975fc7d77136919a049ea74cf08ddd2b896890bb24a0ba18094a22baa351bf29ad96c66bbb1a598f2ca391749620e62d61c3561a7d3653ccc8892c7b99baaf76bf836e2991cb06d6bc0514568ff0d1ec8bb4b3d6984f5eaefb17d3ea2893722375d3ddb8e389a8eef7d7d198f8e687d6a513983df906099f9a2d23f4f9dec6f8ef2f11fc0a21fac45353b94e00486f5e17d386af42502d09db33cf0cf28310e049c07e88682aeeb00cb833c5174266e62407a57583f1f88b304b7c6e0c84bbe1c0fd423072d37a5bd0aacf764229e5c7cd02473460ba3645cd8e8ae144065bf02d0dd238593d8e230354f67e0b2f23012c23274f80e3ee31e35e2606a4a3f31d94ab755e6d163cff52cbb36b6d0cc67ffc512aeed1dce4d7a0d70ce82f2baba12e8d514dc92a056f994adfb17b5b9712bd5186f27a2fda1f7039c5df2c8587fdc62f5627580c13234b55be4df3056050e2d1ef3218f0dd66cb05265fe1acfb0989d8213f2c19d1735a7cf3fa65d88dad5af52dc2bba22b7abf46c3bc77b5091baab9e8f0ddc4d5e581037de91a9f8dcbc69309be29cc815cf19a20a7585b8b3073edf51fc9baeb3e509b97fa4ecfd621e0fd57bd61cac1b895c03248ff12bdbc57509250df3517e8a3fe1d776836b34ab352b973d932ef708b14f7418f9eceb1d87667e61e3e758649cb083f01b133d37ab2f5afa96d6c84bcacf4efc3851ad308c1e7d9113624fce29fab460ab9d2a48d92cdb281103a5250ad44cb2ff6e67ac670c02fdafb3e0f1353953d6d7d5646ca1568dea55275a050ec501b7c6250444f7219f1ba7521ba3b93d089727ca5f3bbe0d6c1300b423377004954c5628fdb65770b18ced5c9b23a4a5a6d6ef25fe01b4ce278de0bcc4ed86e28a0a68818ffa40970128cf2c38740e80037984428c1bd5113f40ff47512ee6f4e4d8f9b8e8e1b3040d2928d003bd1c1329dc885302fbce9fa81c23b4dc49c7c82d29b52957847898676c89aa5d32b5b0e1c0d5a2b79a19d67562f407f19425687971a957375879d90c5f57c857136c17106c9ab1b99d80e69c8c954ed386493368884b55c939b8d64d26f643e800c56f90c01079d7c534e3b2b7ae352cefd3016da55f6a85eb803b85e2304915fd2001f77c74e28746293c46e4f5f0fd49cf988aafd0026b8e7a3bab2da5cdce1ea26c2e29ec03f4807fac432662b2d6c060be1c7be0e5489de69d0a6e03a4b9117f9244b34a0f1ecba89884f781c6320412413a00c4980287409a2a78c2cd7e65cecebbe4ec1c28cac4dd95f6998e78fc6f1392384331c9436aa10e10e2bf8ad2c4eafbcf276aa7bae64b74428911b3269c749338b0fc5075ad",
			wantHex:  "5a0eb2bdf0ac1cae8e586689fa16cd4b07dfdedaec8a110ea1fdb059dd5253231b6132987598dfc6e11f86780428982d50cf68f67ae452622c3b336b537ef3298ca645e8f89ee39a26758206a5a3f6409afc709582f95274b57b71fae5c6b74619ae6f089a5393c5b79235d9caf699d23d88fb873f78379690ad8405e34c19f5257d596580c7a6a7206a3712825afe630c76b31cdb4a23e7f0632e10f14f4e282c81a66451a26f8df2a352b5b9f607a7198449d1b926e27036810368e691a74b91c61afa73d9d3b99453e7c8b50fd4f09c039a2f2feb5c419206694c31b92df1d9586140cb3417b38d0c503c7b508cc2ed12e813a1c795e9829eb39ee78eeaf360a169b491a1d4e419574e712402de9d48d54c1ae5e03739b7156615e8267e1fb0a897f067afd11fb33f6e24182d7aaaaa18fe5bc1982f20d6b871e5a398f0f6f718181d31ec225cfa9a0a70124ed9a70031bdf0c1c7829f708b6e17d50419ef361cf77d99c85f44607186c8d683106b8bd38a49b5d0fb503b397a83388c5678dcfcc737499d84512690701ed621a6f0172aecf037184ddf0f2453e4053024018e5ab2e30d6d5363b56e8b41509317c99042f517247474ab3abc848e00a07f69c254f46f2a05cf6ed84e5cc906a518fdcfdf2c61ce731f24c5264f1a25fc04934dc28aec112134dd523f70115074ca34e3807aa4cb925147f3a0ce152d323bd8c675ace446d0fd1ae30c4b57f0eb2c23884bc18f0964c0114796c5b6d080c3d89175665fbf63a6381a6a9da39ad070b645c8bb1779506da14439a9f5b5d481954764ea114fac688930bc68534d403cff4210673b6a6ff7ae416b7cd41404c3d3f282fcd193b86d0f54d0006c2a503b40d5c3930da980565b8f9630e9493a79d1c03e74e5f93ac8e4dc1a901ec5e3b3e57049124c7b72ea345aa359e782285d9e6a5c144a378111dd02c40855ff9c2be9b48425cb0b2fd62dc8678fd151121cf26a65e917d65d8e0dacfae108eb5508b601fb8ffa370be1f9a8b749a2d12eeab81f41079de87e2d777994fa4d28188c579ad327f9957fb7bdecec5c680844dd43cb57cf87aeb763c003e65011f73f8c63442df39a92b946a6bd968a1c1e4d5fa7d88476a68bd8e20e5b70a99259c7d3f85fb1b65cd2e93972e6264e74ebf289b8b6979b9b68a85cd5b360c1987f87235c3c845d62489e33acf85d53fa3561fe3a3aee18924588d9c6eba4edb7a4d106b31173e42929f6f0c48c80ce6a72d54eca7c0fe870068b7a7c89c63cdda593f5b32d3cb4ea8a32c39f00ab449155757172d66763ed9527019d6de6c9f2416aa6203f4d11c9ebee1e1d3845099e55504446448027212616167eb36035726daa7698b075286f5379cd3e93cb3e0cf4f9cb8d017facbb5550ed32d5ec5400ae57e47e2bf78d1eaeff9480cc765ceff39db500",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.inputHex)%2 != 0 {
				t.Fatalf("odd-length hex string")
			}

			input, err := hex.DecodeString(tc.inputHex)
			if err != nil {
				t.Fatalf("invalid hex: %v", err)
			}
			want, _ := hex.DecodeString(tc.wantHex)

			if len(input) < 96 {
				t.Fatalf("input too short")
			}

			baseLen := new(big.Int).SetBytes(input[:32]).Int64()
			expLen := new(big.Int).SetBytes(input[32:64]).Int64()
			modLen := new(big.Int).SetBytes(input[64:96]).Int64()

			total := 96 + baseLen + expLen + modLen
			if int64(len(input)) != total {
				t.Fatalf("length mismatch: got=%d expected=%d (base=%d exp=%d mod=%d)",
					len(input), total, baseLen, expLen, modLen)
			}

			offset := int64(96)

			base := input[offset : offset+baseLen]
			offset += baseLen

			exp := input[offset : offset+expLen]
			offset += expLen

			mod := input[offset : offset+modLen]

			got := X_modExp(base, exp, mod)

			if !bytes.Equal(got, want) {
				t.Fatalf("\n got : %x\n want: %x", got, want)
			}
		})
	}
}

func TestModExpKnownAnswers(t *testing.T) {
	cases := []struct {
		name string
		base []byte
		exp  []byte
		mod  []byte
		want []byte
	}{
		{"2^10 mod 1000", []byte{2}, []byte{10}, []byte{0x03, 0xe8}, []byte{0x00, 0x18}}, // 24
		{"7^2 mod 4", []byte{7}, []byte{2}, []byte{4}, []byte{0x01}},                     // 49 mod 4 = 1
		{"0^0 mod 5", []byte{0}, []byte{0}, []byte{5}, []byte{0x01}},                     // 0^0 = 1 by convention
		{"3^0 mod 7", []byte{3}, nil, []byte{7}, []byte{0x01}},                           // empty exp == 0
		{"empty base", nil, []byte{5}, []byte{7}, []byte{0}},                             // 0^5 = 0
		{"modulus 1", []byte{2}, []byte{3}, []byte{1}, []byte{0}},                        // anything mod 1 = 0
		{"modulus 0", []byte{2}, []byte{3}, []byte{0, 0, 0}, []byte{0, 0, 0}},            // m=0 => len(m) zero bytes
		{"empty modulus", []byte{2}, []byte{3}, nil, nil},                                // len(m) == 0 => empty
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := X_modExp(tc.base, tc.exp, tc.mod)
			if !bytes.Equal(got, tc.want) {
				t.Fatalf("got %x, want %x", got, tc.want)
			}
		})
	}
}

// TestModExpLargeModReduction exercises the cometblszk use case: reducing a
// 32-byte digest modulo (BN254_R - 1) by raising to the first power.
func TestModExpLargeModReduction(t *testing.T) {
	bn254RMinusOne, _ := hex.DecodeString(
		"30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000000",
	)

	// A value below the modulus must pass through unchanged.
	small, _ := hex.DecodeString(
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	)
	got := X_modExp(small, []byte{1}, bn254RMinusOne)
	if !bytes.Equal(got, small) {
		t.Fatalf("x < m should pass through: got %x", got)
	}

	// A value above the modulus must actually be reduced.
	big1, _ := hex.DecodeString(
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	)
	got = X_modExp(big1, []byte{1}, bn254RMinusOne)
	expected := new(big.Int).Mod(
		new(big.Int).SetBytes(big1),
		new(big.Int).SetBytes(bn254RMinusOne),
	)
	expectedBytes := make([]byte, len(bn254RMinusOne))
	expected.FillBytes(expectedBytes)
	if !bytes.Equal(got, expectedBytes) {
		t.Fatalf("reduction mismatch:\n got  %x\n want %x", got, expectedBytes)
	}
}

// ones returns n bytes of 0xFF.
func ones(n int) []byte { return bytes.Repeat([]byte{0xFF}, n) }

// oddMod returns an n-byte modulus that is large and odd, so accepted calls take
// the Montgomery path rather than a shortcut.
func oddMod(n int) []byte {
	m := ones(n)
	m[n-1] = 0xFD
	return m
}

// huge returns an n-byte big-endian integer with only its top byte set. make()
// zero-fills for almost no allocation gas, so this is the cheap way for a realm
// to build an enormous operand — which is why operand length has to be bounded
// independently of what it costs to construct.
func huge(n int) []byte {
	b := make([]byte, n)
	b[0] = 0xFF
	return b
}

// TestModExpOperandLenBoundary checks each operand independently at
// maxOperandLen and at maxOperandLen+1. Only one operand is held at the limit
// per case: a symmetric all-at-1024 call is a real ~140ms modular
// exponentiation, and it proves nothing about a length boundary that a 1-byte
// exponent does not.
func TestModExpOperandLenBoundary(t *testing.T) {
	const lim = maxOperandLen
	cases := []struct {
		name           string
		base, exp, mod []byte
		wantLen        int
	}{
		{"base at limit", ones(lim), ones(1), oddMod(32), 32},
		{"exp at limit", ones(32), ones(lim), oddMod(32), 32},
		{"modulus at limit", ones(32), ones(1), oddMod(lim), lim},
		{"base over limit", ones(lim + 1), ones(1), oddMod(32), 0},
		{"exp over limit", ones(32), ones(lim + 1), oddMod(32), 0},
		{"modulus over limit", ones(32), ones(1), oddMod(lim + 1), 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// len(nil) == 0, so this covers the rejection cases too.
			if got := X_modExp(tc.base, tc.exp, tc.mod); len(got) != tc.wantLen {
				t.Fatalf("got %d bytes, want %d", len(got), tc.wantLen)
			}
		})
	}
}

// TestModExpRejectsOversizedShapes pins the three shapes that used to decouple
// the real cost from the charge, back when gas scaled on len(modulus) alone. Gas
// now scales on bits(exp) × words(modulus)², so these are priced rather than
// mispriced — but they still exceed maxOperandLen and must be refused before any
// big.Int work. This doubles as a timing guard: it would take ~40s of CPU if any
// shape were still accepted.
func TestModExpRejectsOversizedShapes(t *testing.T) {
	cases := []struct {
		name           string
		base, exp, mod []byte
	}{
		// Before the fix: ~5s of big.Int.Exp, charged ~6.2M gas (≈6.2ms).
		{"huge exponent, small modulus", huge(32), huge(500_000), oddMod(32)},
		// Before the fix: ~34s, charged ~789M gas (≈789ms).
		{"small exponent, huge modulus", huge(32), huge(256), huge(32 * 1024)},
		// Before the fix: ~61ms for the base reduction, charged a flat 828K gas.
		{"huge base, small modulus", huge(500_000), huge(4), oddMod(32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			if got := X_modExp(tc.base, tc.exp, tc.mod); got != nil {
				t.Fatalf("want nil, got %d bytes", len(got))
			}
			// Rejection happens before any big.Int work; anything near the
			// pre-cap timings above means the guard was bypassed.
			if d := time.Since(start); d > time.Second {
				t.Fatalf("rejection took %s — guard did not short-circuit", d)
			}
		})
	}
}

// TestModExpOutputLength: EIP-198 mandates output length = len(modulus),
// left-padded with zero bytes.
func TestModExpOutputLength(t *testing.T) {
	// 2 mod [0, 100]: result is 2, but padded to 2 bytes.
	got := X_modExp([]byte{2}, []byte{1}, []byte{0, 100})
	if len(got) != 2 {
		t.Fatalf("expected 2-byte output, got %d bytes", len(got))
	}
	if got[0] != 0 || got[1] != 2 {
		t.Fatalf("expected 0x0002, got %x", got)
	}
}
