package v1

import (
	xyaml "github.com/ignite/cli/ignite/pkg/yaml"
)

// Validator holds info related to validator settings.
type Validator struct {
	// Name is the name of the validator.
	Name string `yaml:"name"`

	// Bonded is how much the validator has staked.
	Bonded string `yaml:"bonded"`

	// App overwrites appd's config/app.toml configs.
	App xyaml.Map `yaml:"app,omitempty"`

	// Config overwrites appd's config/config.toml configs.
	Config xyaml.Map `yaml:"config,omitempty"`

	// Client overwrites appd's config/client.toml configs.
	Client xyaml.Map `yaml:"client,omitempty"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home,omitempty"`

	// Gentx overwrites appd's config/gentx.toml configs.
	Gentx *Gentx `yaml:"gentx,omitempty"`
}

// Gentx holds info related to Gentx settings.
type Gentx struct {
	// Amount is the amount for the current Gentx.
	Amount string `yaml:"amount"`

	// Moniker is the validator's (optional) moniker.
	Moniker string `yaml:"moniker"`

	// Home is directory for config and data.
	Home string `yaml:"home"`

	// KeyringBackend is keyring's backend.
	KeyringBackend string `yaml:"keyring-backend"`

	// ChainID is the network chain ID.
	ChainID string `yaml:"chain-id"`

	// CommissionMaxChangeRate is the maximum commission change rate percentage (per day).
	CommissionMaxChangeRate string `yaml:"commission-max-change-rate"`

	// CommissionMaxRate is the maximum commission rate percentage
	CommissionMaxRate string `yaml:"commission-max-rate"`

	// CommissionRate is the initial commission rate percentage.
	CommissionRate string `yaml:"commission-rate"`

	// Details is the validator's (optional) details.
	Details string `yaml:"details"`

	// SecurityContact is the validator's (optional) security contact email.
	SecurityContact string `yaml:"security-contact"`

	// Website is the validator's (optional) website.
	Website string `yaml:"website"`

	// AccountNumber is the account number of the signing account (offline mode only).
	AccountNumber int `yaml:"account-number"`

	// BroadcastMode is the transaction broadcasting mode (sync|async|block) (default "sync").
	BroadcastMode string `yaml:"broadcast-mode"`

	// DryRun is a boolean determining whether to ignore the --gas flag and perform a simulation of a transaction.
	DryRun bool `yaml:"dry-run"`

	// FeeAccount is the fee account pays fees for the transaction instead of deducting from the signer
	FeeAccount string `yaml:"fee-account"`

	// Fee is the fee to pay along with transaction; eg: 10uatom.
	Fee string `yaml:"fee"`

	// From is the name or address of private key with which to sign.
	From string `yaml:"from"`

	// From is the gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000).
	Gas string `yaml:"gas"`

	// GasAdjustment is the adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1).
	GasAdjustment string `yaml:"gas-adjustment"`

	// GasPrices is the gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom).
	GasPrices string `yaml:"gas-prices"`

	// GenerateOnly is a boolean determining whether to build an unsigned transaction and write it to STDOUT.
	GenerateOnly bool `yaml:"generate-only"`

	// Identity is the (optional) identity signature (ex. UPort or Keybase).
	Identity string `yaml:"identity"`

	// IP is the node's public IP (default "192.168.1.64").
	IP string `yaml:"ip"`

	// KeyringDir is the client Keyring directory; if omitted, the default 'home' directory will be used.
	KeyringDir string `yaml:"keyring-dir"`

	// Ledger is a boolean determining whether to use a connected Ledger device.
	Ledger bool `yaml:"ledger"`

	// KeyringDir is the minimum self delegation required on the validator.
	MinSelfDelegation string `yaml:"min-self-delegation"`

	// Node is <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657").
	Node string `yaml:"node"`

	// NodeID is the node's NodeID.
	NodeID string `yaml:"node-id"`

	// Note is the note to add a description to the transaction (previously --memo).
	Note string `yaml:"note"`

	// Offline is a boolean determining the offline mode (does not allow any online functionality).
	Offline bool `yaml:"offline"`

	// Output is the output format (text|json) (default "json").
	Output string `yaml:"output"`

	// OutputDocument writes the genesis transaction JSON document to the given file instead of the default location.
	OutputDocument string `yaml:"output-document"`

	// PubKey is the validator's Protobuf JSON encoded public key.
	PubKey string `yaml:"pubkey"`

	// Sequence is the sequence number of the signing account (offline mode only).
	Sequence uint `yaml:"sequence"`

	// SignMode is the choose sign mode (direct|amino-json), this is an advanced feature.
	SignMode string `yaml:"sign-mode"`

	// TimeoutHeight sets a block timeout height to prevent the tx from being committed past a certain height.
	TimeoutHeight uint `yaml:"timeout-height"`
}
