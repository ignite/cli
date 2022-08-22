package v1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/imdario/mergo"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
)

// DefaultValidator defines the default values for the validator.
var (
	// DefaultPortMargin is the default incremental margin for the the port number.
	DefaultPortMargin = 10

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

	// PPROFPort is the default port number of Prof.
	PPROFPort = 6060

	// DefaultValidator is the default configuration of the validator
	DefaultValidator = Validator{
		App: map[string]interface{}{
			"grpc":     map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", GRPCPort)},
			"grpc-web": map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", GRPCWebPort)},
			"api":      map[string]interface{}{"address": fmt.Sprintf("0.0.0.0:%d", APIPort)},
		},
		Config: map[string]interface{}{
			"rpc":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", RPCPort)},
			"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("0.0.0.0:%d", P2PPort)},
			"pprof_laddr": fmt.Sprintf("0.0.0.0:%d", PPROFPort),
		},
	}
)

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	config.BaseConfig `yaml:",inline"`

	Validators []Validator `yaml:"validators"`
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() config.Converter {
	copy := *c
	return &copy
}

// FillValidatorsDefaults fills in the defaults values for the validators if they are missing.
func (c *Config) FillValidatorsDefaults(defaultValidator Validator) error {
	for i := range c.Validators {
		var validator Validator
		if i > 0 {
			previousValidatorPorts := Validator{
				App: map[string]interface{}{
					"grpc":     map[string]interface{}{"address": c.Validators[i-1].GetGRPC()},
					"grpc-web": map[string]interface{}{"address": c.Validators[i-1].GetGRPCWeb()},
					"api":      map[string]interface{}{"address": c.Validators[i-1].GetAPI()},
				},
				Config: map[string]interface{}{
					"rpc":         map[string]interface{}{"laddr": c.Validators[i-1].GetRPC()},
					"p2p":         map[string]interface{}{"laddr": c.Validators[i-1].GetP2P()},
					"pprof_laddr": c.Validators[i-1].GetProf(),
				},
			}

			validator = previousValidatorPorts.IncreasePort(DefaultPortMargin)
		} else {
			validator = defaultValidator
		}

		if err := c.Validators[i].FillDefaults(validator); err != nil {
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

// FillDefaults fills in the default values in the parameter defaultValidator.
func (v *Validator) FillDefaults(defaultValidator Validator) error {
	if err := mergo.Merge(v, defaultValidator); err != nil {
		return err
	}
	return nil
}

// GetGRPC returns the GRPC
func (v *Validator) GetGRPC() string {
	if v.App != nil {
		if val, ok := v.App["grpc"]; ok {
			return getValue(val, "address")
		}
	}
	return ""
}

// GetGRPCAddress returns the GRPC IP
func (v *Validator) GetGRPCAddress() string {
	grpc := v.GetGRPC()
	return getAddress(grpc)
}

// GetGRPCPort returns the GRPC port
func (v *Validator) GetGRPCPort() int {
	grpc := v.GetGRPC()
	return getPort(grpc)
}

// GetGRPCWeb returns the GRPCWeb
func (v *Validator) GetGRPCWeb() string {
	if v.App != nil {
		if val, ok := v.App["grpc-web"]; ok {
			return getValue(val, "address")
		}
	}
	return ""
}

// GetGRPCWebAddress returns the GRPCWeb IP
func (v *Validator) GetGRPCWebAddress() string {
	grpcweb := v.GetGRPCWeb()
	return getAddress(grpcweb)
}

// GetGRPCWebPort returns the GRPCWeb port
func (v *Validator) GetGRPCWebPort() int {
	grpcweb := v.GetGRPCWeb()
	return getPort(grpcweb)
}

// GetAPI returns the API
func (v *Validator) GetAPI() string {
	if v.App != nil {
		if val, ok := v.App["api"]; ok {
			return getValue(val, "address")
		}
	}
	return ""
}

// GetAPIAddress returns the API IP
func (v *Validator) GetAPIAddress() string {
	return getAddress(v.GetAPI())
}

// GetAPIPort returns the API port
func (v *Validator) GetAPIPort() int {
	return getPort(v.GetAPI())
}

// GetProf returns the Prof
func (v *Validator) GetProf() string {
	if v.Config != nil {
		if val, ok := v.Config["pprof_laddr"]; ok {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

// GetProfAddress returns the Prof IP
func (v *Validator) GetProfAddress() string {
	return getAddress(v.GetProf())
}

// GetProfPort returns the Prof port
func (v *Validator) GetProfPort() int {
	return getPort(v.GetProf())
}

// GetP2P returns the P2P
func (v *Validator) GetP2P() string {
	if v.Config != nil {
		if val, ok := v.Config["p2p"]; ok {
			return getValue(val, "laddr")
		}
	}
	return ""
}

// GetP2PAddress returns the P2P IP
func (v *Validator) GetP2PAddress() string {
	return getAddress(v.GetP2P())
}

// GetP2PPort returns the P2P port
func (v *Validator) GetP2PPort() int {
	return getPort(v.GetP2P())
}

// GetRPC returns the RPC
func (v *Validator) GetRPC() string {
	if v.Config != nil {
		if val, ok := v.Config["rpc"]; ok {
			return getValue(val, "laddr")
		}
	}
	return ""
}

// GetRPCAddress returns the RPC IP
func (v *Validator) GetRPCAddress() string {
	return getAddress(v.GetRPC())
}

// GetRPCPort returns the RPC port
func (v *Validator) GetRPCPort() int {
	return getPort(v.GetRPC())
}

// IncreasePort generates an validator with all the ports incremented by the value portIncrement.
func (v *Validator) IncreasePort(portIncrement int) Validator {
	result := Validator{
		App: map[string]interface{}{
			"grpc":     map[string]interface{}{"address": fmt.Sprintf("%s:%d", v.GetGRPCAddress(), v.GetGRPCPort()+portIncrement)},
			"grpc-web": map[string]interface{}{"address": fmt.Sprintf("%s:%d", v.GetGRPCWebAddress(), v.GetGRPCWebPort()+portIncrement)},
			"api":      map[string]interface{}{"address": fmt.Sprintf("%s:%d", v.GetAPIAddress(), v.GetAPIPort()+portIncrement)},
		},
		Config: map[string]interface{}{
			"rpc":         map[string]interface{}{"laddr": fmt.Sprintf("%s:%d", v.GetRPCAddress(), v.GetRPCPort()+portIncrement)},
			"p2p":         map[string]interface{}{"laddr": fmt.Sprintf("%s:%d", v.GetP2PAddress(), v.GetP2PPort()+portIncrement)},
			"pprof_laddr": fmt.Sprintf("%s:%d", v.GetProfAddress(), v.GetProfPort()+portIncrement),
		},
	}
	return result
}

func getValue(val interface{}, keyMap string) string {
	switch v := val.(type) {
	case map[string]interface{}:
		for key, address := range v {
			if key == keyMap {
				return fmt.Sprintf("%v", address)
			}
		}
	case map[interface{}]interface{}:
		for key, address := range v {
			if fmt.Sprintf("%v", key) == keyMap {
				return fmt.Sprintf("%v", address)
			}
		}
	}
	return ""
}

func getAddress(fullAddress string) string {
	if fullAddress == "" {
		return ""
	}
	index := strings.LastIndex(fullAddress, ":")
	return fullAddress[:index]
}

func getPort(fullAddress string) int {
	if fullAddress == "" {
		return 0
	}
	index := strings.LastIndex(fullAddress, ":")
	port, err := strconv.Atoi(fullAddress[index+1:])
	if err != nil {
		return 0
	}
	return port
}
