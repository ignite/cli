package v1

import (
	"fmt"
	"strings"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	"github.com/ignite-hq/cli/ignite/pkg/xnet"
)

var (
	// GRPCPort is the default port number of GRPC.
	GRPCPort = 9090

	// GRPCWebPort is the default port number of GRPC-Web.
	GRPCWebPort = 9091

	// APIPort is the default port number of API.
	APIPort = 1317

	// RPCPort is the default port number of RPC.
	RPCPort = 26657

	// P2PPort is the default port number of P2P.
	P2PPort = 26656

	// PProfPort is the default port number of Prof.
	PProfPort = 6060
)

// DefaultConfig returns a config with default values.
func DefaultConfig() *Config {
	c := Config{BaseConfig: config.DefaultBaseConfig()}
	c.Version = 1
	return &c
}

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	config.BaseConfig `yaml:",inline"`

	Validators []Validator `yaml:"validators"`
}

func (c *Config) SetDefaults() error {
	if err := c.BaseConfig.SetDefaults(); err != nil {
		return err
	}

	if err := c.updateValidatorAddresses(); err != nil {
		return err
	}

	return nil
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() config.Converter {
	copy := *c
	return &copy
}

func (c *Config) updateValidatorAddresses() (err error) {
	// Margin to increase address ports
	margin := uint64(10)

	// Update empty address configuration fields for each validator
	for i := range c.Validators {
		// Define the validator here to be able to reference it during merge
		validator := &c.Validators[i]

		// Make sure the default Cosmos SDK and Tendermint config addresses are initialized
		validator.setDefaultAddresses()

		// Use default addresses for the first validator
		if i == 0 {
			continue
		}

		// Increase the ports for each address when the current validator is not the first.
		// The ports are increased using the previous validator addresses.
		prev := c.Validators[i-1]

		// Increase the Cosmos app config ports for the current validator
		for field, v := range validator.App {
			path := fmt.Sprintf("%s.address", field)
			prevAddr := getConfigValue(prev.App, path)
			m := v.(map[string]interface{})

			m["address"], err = xnet.IncreasePortBy(prevAddr, margin)
			if err != nil {
				return err
			}
		}

		// Increase the Tendermint config ports for the current validator
		for field, v := range validator.Config {
			// Skip the fields that are not a map
			m, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			path := fmt.Sprintf("%s.laddr", field)
			prevAddr := getConfigValue(prev.Config, path)

			m["laddr"], err = xnet.IncreasePortBy(prevAddr, margin)
			if err != nil {
				return err
			}
		}

		addr := prev.Config["pprof_laddr"].(string)

		validator.Config["pprof_laddr"], err = xnet.IncreasePortBy(addr, margin)
		if err != nil {
			return err
		}
	}

	return nil
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

// Validator holds info related to validator settings.
type Validator struct {
	// Name is the name of the validator.
	Name string `yaml:"name"`

	// Bonded is how much the validator has staked.
	Bonded string `yaml:"bonded"`

	// App overwrites appd's config/app.toml configs.
	App map[string]interface{} `yaml:"app"`

	// Config overwrites appd's config/config.toml configs.
	Config map[string]interface{} `yaml:"config"`

	// Client overwrites appd's config/client.toml configs.
	Client map[string]interface{} `yaml:"client"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home"`

	// KeyringBackend is the default keyring backend to use for blockchain initialization
	KeyringBackend string `yaml:"keyring-backend"`

	// Gentx overwrites appd's config/gentx.toml configs.
	Gentx *Gentx `yaml:"gentx"`
}

// GetGRPC returns the GRPC address.
func (v *Validator) GetGRPC() string {
	return getConfigValue(v.App, "grpc.address")
}

// GetGRPCWeb returns the GRPC web address.
func (v *Validator) GetGRPCWeb() string {
	return getConfigValue(v.App, "grpc-web.address")
}

// GetAPI returns the API address.
func (v *Validator) GetAPI() string {
	return getConfigValue(v.App, "api.address")
}

// GetProf returns the Prof address.
func (v *Validator) GetProf() string {
	return getConfigValue(v.Config, "pprof_laddr")
}

// GetP2P returns the P2P address.
func (v *Validator) GetP2P() string {
	return getConfigValue(v.Config, "p2p.laddr")
}

// GetRPC returns the RPC address.
func (v *Validator) GetRPC() string {
	return getConfigValue(v.Config, "rpc.laddr")
}

func (v *Validator) setDefaultAddresses() {
	v.App = map[string]interface{}{
		"grpc":     map[string]interface{}{"address": xnet.AnyIPv4Address(GRPCPort)},
		"grpc-web": map[string]interface{}{"address": xnet.AnyIPv4Address(GRPCWebPort)},
		"api":      map[string]interface{}{"address": xnet.AnyIPv4Address(APIPort)},
	}
	v.Config = map[string]interface{}{
		"rpc":         map[string]interface{}{"laddr": xnet.AnyIPv4Address(RPCPort)},
		"p2p":         map[string]interface{}{"laddr": xnet.AnyIPv4Address(P2PPort)},
		"pprof_laddr": xnet.AnyIPv4Address(PProfPort),
	}
}

func getConfigValue(cfg map[string]interface{}, path string) string {
	if cfg == nil {
		return ""
	}

	// Get the first path element and also the extra path elements in suffix
	key, suffix, _ := strings.Cut(path, ".")
	if key == "" {
		return ""
	}

	// Continue traversing the path when it contains a suffix
	if suffix != "" {
		var ok bool

		cfg, ok = cfg[key].(map[string]interface{})
		if !ok {
			// Handle the case where YAML decodes the nested maps.
			// Nested maps key types are decoded as interface{}.
			nv, ok := cfg[key].(map[interface{}]interface{})
			if !ok {
				// The path doesn't exist in the config
				return ""
			}

			// Convert the nested value to a map with string keys
			for k, v := range nv {
				// Skip key types that are not strings
				n, ok := k.(string)
				if !ok {
					continue
				}

				cfg[n] = v
			}
		}

		return getConfigValue(cfg, suffix)
	}

	// Get config value as string
	s, _ := cfg[key].(string)

	return s
}
