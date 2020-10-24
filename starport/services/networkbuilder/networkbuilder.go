package networkbuilder

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/services/chain"
)

// Builder is network builder.
type Builder struct {
	ev events.Bus
}

type Option func(*Builder)

// CollectEvents collects events from Builder.
func CollectEvents(ev events.Bus) Option {
	return func(b *Builder) {
		b.ev = ev
	}
}

// New creates a Builder.
// TODO receive SPN info here.
func New(options ...Option) *Builder {
	b := &Builder{}
	for _, opt := range options {
		opt(b)
	}
	return b
}

// Init initializes blockchain from a gitURL and returns its genesis content.
func (b *Builder) Init(ctx context.Context, gitURL string) (genesis []byte, err error) {
	defer b.ev.Shutdown()

	appPath, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(appPath)

	b.ev.Send(events.New(events.StatusOngoing, "Pulling the blockchain"))
	if _, err := git.PlainCloneContext(ctx, appPath, false, &git.CloneOptions{
		URL: gitURL,
	}); err != nil {
		return nil, err
	}
	b.ev.Send(events.New(events.StatusDone, "Pulled the blockchain"))

	path, err := gomodulepath.ParseFile(appPath)
	if err != nil {
		return nil, err
	}
	app := chain.App{
		Name: path.Root,
		Path: appPath,
	}

	c, err := chain.New(app, chain.LogSilent)
	if err != nil {
		return nil, err
	}

	b.ev.Send(events.New(events.StatusOngoing, "Initializing the blockchain"))
	if err := c.Build(ctx); err != nil {
		return nil, err
	}
	b.ev.Send(events.New(events.StatusDone, "Blockchain initialized"))
	if err := c.Init(ctx); err != nil {
		return nil, err
	}

	genesisPath, err := c.GenesisPath()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(genesisPath)
}

// Submit submits Genesis to SPN to announce a new network.
func (b *Builder) Submit(ctx context.Context, genesis []byte) (err error) { return nil }
