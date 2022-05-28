package common

// Validator is the interface defining all the common methods for the Validator struct across all supported versions
type Validator interface {
	GetName() string
	GetBonded() string
}

// Config is the interface defining all the common methods for the ConfigYaml struct across all supported versions
type Config interface {
	Clone() Config
	GetVersion() Version
	GetFaucet() Faucet
	ListAccounts() []Account
	ListValidators() []Validator
	GetClient() Client
	GetBuild() Build

	GetGenesis() map[string]interface{}

	// Keep this deprecated method to be backward-compatible.
	GetHost() Host

	// Keep this deprecated method to be backward-compatible.
	GetInit() Init

	// ConvertNext converts the instance of Config from the current version to the next version.
	ConvertNext() (Config, error)
}
