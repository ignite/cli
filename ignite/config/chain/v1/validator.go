package v1

import (
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
)

// Validator holds info related to validator settings.
type Validator struct {
	// Name is the name of the validator.
	Name string `yaml:"name" doc:"name of the validator"`

	// Bonded is how much the validator has staked.
	Bonded string `yaml:"bonded" doc:"how much the validator has staked"`

	// App overwrites appd's config/app.toml configs.
	App xyaml.Map `yaml:"app,omitempty" doc:"overwrites appd's config/app.toml configs"`

	// Config overwrites appd's config/config.toml configs.
	Config xyaml.Map `yaml:"config,omitempty" doc:"overwrites appd's config/config.toml configs"`

	// Client overwrites appd's config/client.toml configs.
	Client xyaml.Map `yaml:"client,omitempty" doc:"overwrites appd's config/client.toml configs"`

	// Home overwrites default home directory used for the app.
	Home string `yaml:"home,omitempty" doc:"overwrites default home directory used for the app"`

	// Gentx overwrites appd's config/gentx.toml configs.
	Gentx *Gentx `yaml:"gentx,omitempty" doc:"overwrites appd's config/gentx.toml configs"`
}

// Gentx holds info related to Gentx settings.
type Gentx struct {
	// Amount is the amount for the current Gentx.
	Amount string `yaml:"amount" doc:"the amount for the current Gentx"`

	// Moniker is the validator's (optional) moniker.
	Moniker string `yaml:"moniker" doc:"the validator's (optional) moniker"`

	// Home is directory for config and data.
	Home string `yaml:"home" doc:"directory for config and data"`

	// KeyringBackend is keyring's backend.
	KeyringBackend string `yaml:"keyring-backend" doc:"keyring's backend"`

	// ChainID is the network chain ID.
	ChainID string `yaml:"chain-id" doc:"network chain ID"`

	// CommissionMaxChangeRate is the maximum commission change rate percentage (per day).
	CommissionMaxChangeRate string `yaml:"commission-max-change-rate" doc:"maximum commission change rate percentage (per day)"`

	// CommissionMaxRate is the maximum commission rate percentage.
	CommissionMaxRate string `yaml:"commission-max-rate" doc:"maximum commission rate percentage"`

	// CommissionRate is the initial commission rate percentage.
	CommissionRate string `yaml:"commission-rate" doc:"initial commission rate percentage"`

	// Details is the validator's (optional) details.
	Details string `yaml:"details" doc:"validator's (optional) details"`

	// SecurityContact is the validator's (optional) security contact email.
	SecurityContact string `yaml:"security-contact" doc:"validator's (optional) security contact email"`

	// Website is the validator's (optional) website.
	Website string `yaml:"website" doc:"validator's (optional) website"`

	// AccountNumber is the account number of the signing account (offline mode only).
	AccountNumber int `yaml:"account-number" doc:"account number of the signing account (offline mode only)"`

	// BroadcastMode is the transaction broadcasting mode (sync|async|block) (default "sync").
	BroadcastMode string `yaml:"broadcast-mode" doc:"transaction broadcasting mode (sync|async|block) (default "sync")"`

	// DryRun is a boolean determining whether to ignore the --gas flag and perform a simulation of a transaction.
	DryRun bool `yaml:"dry-run" doc:"boolean determining whether to ignore the --gas flag and perform a simulation of a transaction"`

	// FeeAccount is the fee account pays fees for the transaction instead of deducting from the signer.
	FeeAccount string `yaml:"fee-account" doc:"fee account pays fees for the transaction instead of deducting from the signer"`

	// Fee is the fee to pay along with transaction; eg: 10uatom.
	Fee string `yaml:"fee" doc:"fee to pay along with transaction; eg: 10uatom"`

	// From is the name or address of private key with which to sign.
	From string `yaml:"from" doc:"name or address of private key with which to sign"`

	// From is the gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000).
	Gas string `yaml:"gas" doc:"gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default 200000)"`

	// GasAdjustment is the adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1).
	GasAdjustment string `yaml:"gas-adjustment" doc:"adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored  (default 1)"`

	// GasPrices is the gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom).
	GasPrices string `yaml:"gas-prices" doc:"gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)"`

	// GenerateOnly is a boolean determining whether to build an unsigned transaction and write it to STDOUT.
	GenerateOnly bool `yaml:"generate-only" doc:"boolean determining whether to build an unsigned transaction and write it to STDOUT"`

	// Identity is the (optional) identity signature (ex. UPort or Keybase).
	Identity string `yaml:"identity" doc:"identity signature (ex. UPort or Keybase)"`

	// IP is the node's public IP (default "192.168.1.64").
	IP string `yaml:"ip" doc:"node's public IP (default "192.168.1.64")"`

	// KeyringDir is the client Keyring directory; if omitted, the default 'home' directory will be used.
	KeyringDir string `yaml:"keyring-dir" doc:"client Keyring directory; if omitted, the default 'home' directory will be used"`

	// Ledger is a boolean determining whether to use a connected Ledger device.
	Ledger bool `yaml:"ledger" doc:"boolean determining whether to use a connected Ledger device"`

	// KeyringDir is the minimum self delegation required on the validator.
	MinSelfDelegation string `yaml:"min-self-delegation" doc:"minimum self delegation required on the validator"`

	// Node is <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657").
	Node string `yaml:"node" doc:"<host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")"`

	// NodeID is the node's NodeID.
	NodeID string `yaml:"node-id" doc:"node's NodeID"`

	// Note is the note to add a description to the transaction (previously --memo).
	Note string `yaml:"note" doc:"note to add a description to the transaction (previously --memo)"`

	// Offline is a boolean determining the offline mode (does not allow any online functionality).
	Offline bool `yaml:"offline" doc:"boolean determining the offline mode (does not allow any online functionality)"`

	// Output is the output format (text|json) (default "json").
	Output string `yaml:"output" doc:"output format (text|json) (default "json")"`

	// OutputDocument writes the genesis transaction JSON document to the given file instead of the default location.
	OutputDocument string `yaml:"output-document" doc:"writes the genesis transaction JSON document to the given file instead of the default location"`

	// PubKey is the validator's Protobuf JSON encoded public key.
	PubKey string `yaml:"pubkey" doc:"validator's Protobuf JSON encoded public key"`

	// Sequence is the sequence number of the signing account (offline mode only).
	Sequence uint `yaml:"sequence" doc:"sequence number of the signing account (offline mode only)"`

	// SignMode is the choose sign mode (direct|amino-json), this is an advanced feature.
	SignMode string `yaml:"sign-mode" doc:"choose sign mode (direct|amino-json), this is an advanced feature"`

	// TimeoutHeight sets a block timeout height to prevent the tx from being committed past a certain height.
	TimeoutHeight uint `yaml:"timeout-height" doc:"sets a block timeout height to prevent the tx from being committed past a certain height"`
}
