package plugin

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var (
	// ErrAppChainNotFound indicates that the plugin command is not running inside a blockchain app.
	ErrAppChainNotFound = errors.New("blockchain app not found")
	// ErrCmdNotFound indicates that the plugin command is not running exist in the blockchain app.
	ErrCmdNotFound = errors.New("blockchain CLI commands not found")
)

//go:generate mockery --srcpkg . --name Chainer --structname ChainerInterface --filename chainer.go --with-expecter
type Chainer interface {
	// AppPath returns the configured App's path.
	AppPath() string

	// ID returns the configured App's chain id.
	ID() (string, error)

	// ConfigPath returns the path to the App's config file.
	ConfigPath() string

	// RPCPublicAddress returns the configured App's rpc endpoint.
	RPCPublicAddress() (string, error)

	// Home returns the App's home dir.
	Home() (string, error)
}

// APIOption defines options for the client API.
type APIOption func(*apiOptions)

type apiOptions struct {
	cmd   *cobra.Command
	chain Chainer
}

// WithChain configures the chain to use for the client API.
func WithChain(c Chainer) APIOption {
	return func(o *apiOptions) {
		o.chain = c
	}
}

// WithCmd configures the chain CLI commands to use for the client API.
func WithCmd(cmd *cobra.Command) APIOption {
	return func(o *apiOptions) {
		o.cmd = cmd
	}
}

// NewClientAPI creates a new app ClientAPI.
func NewClientAPI(options ...APIOption) ClientAPI {
	o := apiOptions{}
	for _, apply := range options {
		apply(&o)
	}
	return clientAPI{o}
}

type clientAPI struct {
	o apiOptions
}

func (api clientAPI) GetChainInfo(context.Context) (*ChainInfo, error) {
	chain, err := api.getChain()
	if err != nil {
		return nil, err
	}

	chainID, err := chain.ID()
	if err != nil {
		return nil, err
	}

	rpc, err := chain.RPCPublicAddress()
	if err != nil {
		return nil, err
	}

	home, err := chain.Home()
	if err != nil {
		return nil, err
	}

	return &ChainInfo{
		ChainId:    chainID,
		AppPath:    chain.AppPath(),
		ConfigPath: chain.ConfigPath(),
		RpcAddress: rpc,
		Home:       home,
	}, nil
}

func (api clientAPI) RunCommand(ctx context.Context, command ...string) error {
	cmd, err := api.getCmd()
	if err != nil {
		return err
	}
	cmd.SetArgs(command)
	return cmd.ExecuteContext(ctx)
}

func (api clientAPI) getChain() (Chainer, error) {
	if api.o.chain == nil {
		return nil, ErrAppChainNotFound
	}
	return api.o.chain, nil
}

func (api clientAPI) getCmd() (*cobra.Command, error) {
	if api.o.cmd == nil {
		return nil, ErrCmdNotFound
	}
	return api.o.cmd, nil
}
