package chaincmdrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// Start starts the blockchain.
func (r Runner) Start(ctx context.Context, args ...string) error {
	return r.run(ctx, runOptions{longRunning: true}, r.cc.StartCommand(args...))
}

// LaunchpadStartRestServer start launchpad rest server.
func (r Runner) LaunchpadStartRestServer(ctx context.Context, apiAddress, rpcAddress string) error {
	return r.run(ctx, runOptions{longRunning: true}, r.cc.LaunchpadRestServerCommand(apiAddress, rpcAddress))
}

// Init inits the blockchain.
func (r Runner) Init(ctx context.Context, moniker string) error {
	return r.run(ctx, runOptions{}, r.cc.InitCommand(moniker))
}

// KV holds a key, value pair.
type KV struct {
	key   string
	value string
}

// NewKV returns a new key, value pair.
func NewKV(key, value string) KV {
	return KV{key, value}
}

// LaunchpadSetConfigs updates configurations for a launchpad app.
func (r Runner) LaunchpadSetConfigs(ctx context.Context, kvs ...KV) error {
	for _, kv := range kvs {
		if err := r.run(ctx, runOptions{}, r.cc.LaunchpadSetConfigCommand(kv.key, kv.value)); err != nil {
			return err
		}
	}
	return nil
}

var gentxRe = regexp.MustCompile(`(?m)"(.+?)"`)

// Gentx generates a genesis tx carrying a self delegation.
func (r Runner) Gentx(ctx context.Context, validatorName, selfDelegation string, options ...chaincmd.GentxOption) (gentxPath string, err error) {
	b := &bytes.Buffer{}

	// note: launchpad outputs from stderr.
	if err := r.run(ctx, runOptions{stdout: b, stderr: b}, r.cc.GentxCommand(validatorName, selfDelegation, options...)); err != nil {
		return "", err
	}

	return gentxRe.FindStringSubmatch(b.String())[1], nil
}

// CollectGentxs collects gentxs.
func (r Runner) CollectGentxs(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.cc.CollectGentxsCommand())
}

// ValidateGenesis validates genesis.
func (r Runner) ValidateGenesis(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.cc.ValidateGenesisCommand())
}

// ShowNodeID shows node id.
func (r Runner) ShowNodeID(ctx context.Context) (nodeID string, err error) {
	b := &bytes.Buffer{}
	err = r.run(ctx, runOptions{stdout: b}, r.cc.ShowNodeIDCommand())
	nodeID = strings.TrimSpace(b.String())
	return
}

type NodeStatus struct {
	ChainID string
}

func (r Runner) Status(ctx context.Context) (NodeStatus, error) {
	b := &bytes.Buffer{}

	if err := r.run(ctx, runOptions{stdout: b, stderr: b}, r.cc.StatusCommand()); err != nil {
		return NodeStatus{}, err
	}

	var chainID string

	if r.cc.SDKVersion().Major().Is(cosmosver.Stargate) {
		out := struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"NodeInfo"`
		}{}

		if err := json.NewDecoder(b).Decode(&out); err != nil {
			return NodeStatus{}, err
		}

		chainID = out.NodeInfo.Network
	} else {
		out := struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"node_info"`
		}{}

		if err := json.NewDecoder(b).Decode(&out); err != nil {
			return NodeStatus{}, err
		}

		chainID = out.NodeInfo.Network
	}

	return NodeStatus{
		ChainID: chainID,
	}, nil
}

//BankSend sends amount from fromAccount to toAccount.
func (r Runner) BankSend(ctx context.Context, fromAccount, toAccount, amount string) error {
	b := &bytes.Buffer{}

	if err := r.run(ctx, runOptions{stdout: b}, r.cc.BankSendCommand(fromAccount, toAccount, amount)); err != nil {
		if strings.Contains(err.Error(), "key not found") || // stargate
			strings.Contains(err.Error(), "unknown address") || // launchpad
			strings.Contains(b.String(), "item could not be found") { // launchpad
			return errors.New("account doesn't have any balances")
		}

		return err
	}

	out := struct {
		Code  int    `json:"code"` // TODO LP with code 19.
		Error string `json:"raw_log"`
	}{}

	if err := json.NewDecoder(b).Decode(&out); err != nil {
		return err
	}

	if out.Code > 0 {
		return errors.New(out.Error)
	}

	return nil
}

type EventSelector struct {
	typ   string
	attr  string
	value string
}

func NewEventSelector(typ, addr, value string) EventSelector {
	return EventSelector{typ, addr, value}
}

type Event struct {
	Type       string
	Attributes []EventAttribute
}

type EventAttribute struct {
	Key   string
	Value string
}

func (r Runner) QueryTxEvents(ctx context.Context, event EventSelector, moreEvents ...EventSelector) ([]Event, error) {
	// prepare the slector.
	var list []string

	eventsSelectors := append([]EventSelector{event}, moreEvents...)

	for _, event := range eventsSelectors {
		list = append(list, fmt.Sprintf("%s.%s=%s", event.typ, event.attr, event.value))
	}

	query := strings.Join(list, "&")

	// execute the commnd and parse the output.
	b := &bytes.Buffer{}

	if err := r.run(ctx, runOptions{stdout: b}, r.cc.QueryTxEventsCommand(query)); err != nil {
		return nil, err
	}

	out := struct {
		Txs []struct {
			Logs []struct {
				Events []struct {
					Type  string `json:"type"`
					Attrs []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
					} `json:"attributes"`
				} `json:"events"`
			} `json:"logs"`
		} `json:"txs"`
	}{}

	if err := json.NewDecoder(b).Decode(&out); err != nil {
		return nil, err
	}

	var events []Event

	for _, tx := range out.Txs {
		for _, log := range tx.Logs {
			for _, e := range log.Events {
				var attrs []EventAttribute
				for _, attr := range e.Attrs {
					attrs = append(attrs, EventAttribute{
						Key:   attr.Key,
						Value: attr.Value,
					})
				}

				events = append(events, Event{
					Type:       e.Type,
					Attributes: attrs,
				})
			}
		}
	}

	return events, nil
}
