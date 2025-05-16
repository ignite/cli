package app

// Options ...
type Options struct {
	AppName          string
	AppPath          string
	ProtoDir         string
	GitHubPath       string
	BinaryNamePrefix string
	ModulePath       string
	AddressPrefix    string
	CoinType         uint32
	DefaultDenom     string
	// IncludePrefixes is used to filter the files to include from the generator
	IncludePrefixes []string
	IsChainMinimal  bool
}
