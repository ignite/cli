package cometblszk

import "encoding/hex"

// BN254 scalar field modulus minus one, as a 32-byte big-endian value. Used
// by hashToField for the reduction step.
var primeRMinusOne = mustDecodeHex(
	"30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000000",
)

// Verifying-key constants, mirroring evm/contracts/lib/CometblsZKVerifier.sol.
// "_0" / "_1" suffixes match the Solidity naming convention: _0 is the real
// part of the Fp2 coordinate, _1 is the imaginary part. The EIP-197 pairing
// precompile expects the imaginary part first in its input packing.
var (
	alphaX = mustDecodeHex("245229d9b076b3c0e8a4d70bde8c1cccffa08a9fae7557b165b3b0dbd653e2c7")
	alphaY = mustDecodeHex("253ec85988dbb84e46e94b5efa3373b47a000b4ac6c86b2d4b798d274a182302")

	betaNegX0 = mustDecodeHex("2424bcc1f60a5472685fd50705b2809626e170120acaf441e133a2bd5e61d244")
	betaNegX1 = mustDecodeHex("07090a82e8fabbd39299be24705b92cf208ee8b3487f6f2b39ff27978a29a1db")
	betaNegY0 = mustDecodeHex("04ddc8d30d5c438ca34091c5d2c6ded571382cba2b3c4fdc4222df2938b4e51e")
	betaNegY1 = mustDecodeHex("25833b15e156ae01f2741f4f4120ddb466c52eb83a959f79eb99b23caa7fbf1d")

	gammaNegX0 = mustDecodeHex("1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed")
	gammaNegX1 = mustDecodeHex("198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2")
	gammaNegY0 = mustDecodeHex("1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d")
	gammaNegY1 = mustDecodeHex("275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec")

	deltaNegX0 = mustDecodeHex("02aca5d2a73f8d34e4b26eee3932365e6526c8d5e2f3347d679c2cb1867104dc")
	deltaNegX1 = mustDecodeHex("07b8dbefa90bde075a26318e5066db729155514e3c06b888d4e03c56d82c97e6")
	deltaNegY0 = mustDecodeHex("1696ccafaefe49a5d8bad8e79630e19b25e5392a203aff0042d0216f254806f5")
	deltaNegY1 = mustDecodeHex("2edb19cbb2b6ad0c98fdd7d1845500c26e497dc35e4cdc1cb02cc65dc4ba1bf2")

	// Gamma-abc (public input MSM) constants. Solidity names them CONSTANT,
	// PUB_0, PUB_1.
	constantX = mustDecodeHex("2f5d8a3817f21d3e453573c90c3cc47b7ff235fad7bdfbd59bbd6ae5d153273e")
	constantY = mustDecodeHex("147fa22142b1fd86ce75fc87230a0feac8765d02938784dcfc828d17d7e7c432")
	pub0X     = mustDecodeHex("2a81b98e1c997bd01a20893a08a46c6804493e838c1a0ff6c8c069ef5ab66b9a")
	pub0Y     = mustDecodeHex("276938ada8075cec20d4d6a1f157ec94cc7ba6207c98576e98c1ad9d6378fb6f")
	pub1X     = mustDecodeHex("179496ce140df89ce35c5ee7fb496efdffda5e5d3b95ff9116e2e5df96b36ab7")
	pub1Y     = mustDecodeHex("0326e7d44688ce5903676b7d646e46a5938c8e5fd8cd54e4d5aa3300649f3cfc")

	pedersenGX0 = mustDecodeHex("1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed")
	pedersenGX1 = mustDecodeHex("198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2")
	pedersenGY0 = mustDecodeHex("12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa")
	pedersenGY1 = mustDecodeHex("090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b")

	pedersenGRootSigmaNegX0 = mustDecodeHex("02aca5d2a73f8d34e4b26eee3932365e6526c8d5e2f3347d679c2cb1867104dc")
	pedersenGRootSigmaNegX1 = mustDecodeHex("07b8dbefa90bde075a26318e5066db729155514e3c06b888d4e03c56d82c97e6")
	pedersenGRootSigmaNegY0 = mustDecodeHex("1696ccafaefe49a5d8bad8e79630e19b25e5392a203aff0042d0216f254806f5")
	pedersenGRootSigmaNegY1 = mustDecodeHex("2edb19cbb2b6ad0c98fdd7d1845500c26e497dc35e4cdc1cb02cc65dc4ba1bf2")
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("cometblszk: invalid hex constant: " + err.Error())
	}
	return b
}
