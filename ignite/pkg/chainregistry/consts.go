package chainregistry

type NetworkType string

const (
	// NetworkTypeMainnet is the mainnet network type.
	NetworkTypeMainnet NetworkType = "mainnet"

	// NetworkTypeTestnet is the testnet network type.
	NetworkTypeTestnet NetworkType = "testnet"

	// NetworkTypeDevnet is the devnet network type.
	NetworkTypeDevnet NetworkType = "devnet"
)

type ChainType string

const (
	// ChainTypeCosmos is the cosmos chain type.
	ChainTypeCosmos ChainType = "cosmos"

	// ChainTypeEip155 is the eip155 chain type.
	ChainTypeEip155 ChainType = "eip155"
)

type ChainStatus string

const (
	// ChainStatusActive is the live chain status.
	ChainStatusActive ChainStatus = "live"

	// ChainStatusUpcoming is the upcoming chain status.
	ChainStatusUpcoming ChainStatus = "upcoming"

	// ChainStatusKilled is the inactive chain status.
	ChainStatusKilled ChainStatus = "killed"
)

type KeyAlgos string

const (
	// KeyAlgoSecp256k1 is the secp256k1 key algorithm.
	KeyAlgoSecp256k1 KeyAlgos = "secp256k1"

	// KeyAlgosEthSecp256k1 is the secp256k1 key algorithm with ethereum compatibility.
	KeyAlgosEthSecp256k1 KeyAlgos = "ethsecp256k1"

	// KeyAlgoEd25519 is the ed25519 key algorithm.
	KeyAlgoEd25519 KeyAlgos = "ed25519"
)
