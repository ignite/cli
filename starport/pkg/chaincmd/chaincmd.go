package chaincmd

const (
	commandStart = "start"
	commandInit = "init"
	commandKeys = "keys"
	commandAddGenesisAccount = "add-genesis-account"
	commandGentx = "gentx"
	commandCollectGentxs = "collect-gentxs"
	commandValidateGenesis = "validate-genesis"
	commandShowNodeID = "show-node-id"

	optionHome = "--home"
	optionKeyringBackend = "--keyring-backend"
)

type ChainCmd struct {
	appCmd string
	homeDir   string
	keyringBackend string
}
func New(appName string, homeDir string, keyringBackend string) ChainCmd {
	return ChainCmd{
		appCmd: appName,
		homeDir: homeDir,
		keyringBackend: keyringBackend,
	}
}

func (c ChainCmd) StartCommand() []string {
	return c.withFlags([]string{c.appCmd, commandStart})
}

func (c ChainCmd) InitCommand() []string {
	return c.withFlags([]string{c.appCmd, commandInit})
}

func (c ChainCmd) AddKeyCommand() []string {
	return c.withFlags([]string{c.appCmd, commandKeys})
}

func (c ChainCmd) ShowKeyCommand() []string {
	return c.withFlags([]string{c.appCmd, commandKeys})
}

func (c ChainCmd) AddGenesisAccountCommand() []string {
	return c.withFlags([]string{c.appCmd, commandAddGenesisAccount})
}

func (c ChainCmd) GentxCommand() []string {
	return c.withFlags([]string{c.appCmd, commandGentx})
}

func (c ChainCmd) CollectGentxsCommand() []string {
	return c.withFlags([]string{c.appCmd, commandCollectGentxs})
}

func (c ChainCmd) ValidateGenesisCommand() []string {
	return c.withFlags([]string{c.appCmd, commandValidateGenesis})
}

func (c ChainCmd) ShowNodeIDCommand() []string {
	return c.withFlags([]string{c.appCmd, "tendermint", commandShowNodeID})
}


func (c ChainCmd) withFlags (command []string) []string {
	// Attach home
	if c.homeDir != "" {
		command = append(command, []string{optionHome, c.homeDir}...)
	}
	// Attach keyring backend
	if c.homeDir != "" {
		command = append(command, []string{optionKeyringBackend, c.keyringBackend}...)
	}


	return command
}

