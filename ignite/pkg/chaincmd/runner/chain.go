package chaincmdrunner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// Start starts the blockchain.
func (r Runner) Start(ctx context.Context, args ...string) error {
	return r.run(
		ctx,
		runOptions{wrappedStdErrMaxLen: 50000},
		r.chainCmd.StartCommand(args...),
	)
}

// Init inits the blockchain.
func (r Runner) Init(ctx context.Context, moniker string) error {
	return r.run(ctx, runOptions{}, r.chainCmd.InitCommand(moniker))
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

var gentxRe = regexp.MustCompile(`(?m)"(.+?)"`)

func (r Runner) InPlace(ctx context.Context, newChainID, newOperatorAddress string, options ...chaincmd.InPlaceOption) error {
	runOptions := runOptions{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	return r.run(
		ctx,
		runOptions,
		r.chainCmd.TestnetInPlaceCommand(newChainID, newOperatorAddress, options...),
	)
}

// Initialize config directories & files for a multi-validator testnet locally.
func (r Runner) MultiNode(ctx context.Context, options ...chaincmd.MultiNodeOption) error {
	runOptions := runOptions{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	return r.run(
		ctx,
		runOptions,
		r.chainCmd.TestnetMultiNodeCommand(options...),
	)
}

// Gentx generates a genesis tx carrying a self delegation.
func (r Runner) Gentx(
	ctx context.Context,
	validatorName,
	selfDelegation string,
	options ...chaincmd.GentxOption,
) (gentxPath string, err error) {
	b := newBuffer()

	if err := r.run(ctx, runOptions{
		stdout: b,
		stderr: b,
		stdin:  os.Stdin,
	}, r.chainCmd.GentxCommand(validatorName, selfDelegation, options...)); err != nil {
		return "", err
	}

	return gentxRe.FindStringSubmatch(b.String())[1], nil
}

// CollectGentxs collects gentxs.
func (r Runner) CollectGentxs(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.chainCmd.CollectGentxsCommand())
}

// ValidateGenesis validates genesis.
func (r Runner) ValidateGenesis(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.chainCmd.ValidateGenesisCommand())
}

// UnsafeReset resets the blockchain database.
func (r Runner) UnsafeReset(ctx context.Context) error {
	return r.run(ctx, runOptions{}, r.chainCmd.UnsafeResetCommand())
}

// ShowNodeID shows node id.
func (r Runner) ShowNodeID(ctx context.Context) (nodeID string, err error) {
	b := newBuffer()
	err = r.run(ctx, runOptions{stdout: b}, r.chainCmd.ShowNodeIDCommand())
	nodeID = strings.TrimSpace(b.String())
	return
}

// NodeStatus keeps info about node's status.
type NodeStatus struct {
	ChainID string
}

// Status returns the node's status.
func (r Runner) Status(ctx context.Context) (NodeStatus, error) {
	b := newBuffer()

	if err := r.run(ctx, runOptions{stdout: b, stderr: b}, r.chainCmd.StatusCommand()); err != nil {
		return NodeStatus{}, err
	}

	var chainID string

	data, err := b.JSONEnsuredBytes()
	if err != nil {
		return NodeStatus{}, err
	}

	version := r.chainCmd.SDKVersion()
	switch {
	case version.GTE(cosmosver.StargateFortyVersion):
		out := struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"NodeInfo"`
		}{}

		if err := json.Unmarshal(data, &out); err != nil {
			return NodeStatus{}, err
		}

		chainID = out.NodeInfo.Network
	default:
		out := struct {
			NodeInfo struct {
				Network string `json:"network"`
			} `json:"node_info"`
		}{}

		if err := json.Unmarshal(data, &out); err != nil {
			return NodeStatus{}, err
		}

		chainID = out.NodeInfo.Network
	}

	return NodeStatus{
		ChainID: chainID,
	}, nil
}

// BankSend sends amount from fromAccount to toAccount.
func (r Runner) BankSend(ctx context.Context, fromAccount, toAccount, amount string, options ...chaincmd.BankSendOption) (string, error) {
	b := newBuffer()
	opt := []step.Option{
		r.chainCmd.BankSendCommand(fromAccount, toAccount, amount, options...),
	}

	if r.chainCmd.KeyringPassword() != "" {
		input := newBuffer()
		_, err := fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(input, r.chainCmd.KeyringPassword())
		if err != nil {
			return "", err
		}
		opt = append(opt, step.Write(input.Bytes()))
	}

	if err := r.run(ctx, runOptions{stdout: b}, opt...); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			return "", errors.New("account doesn't have any balances")
		}

		return "", err
	}

	txResult, err := decodeTxResult(b)
	if err != nil {
		return "", err
	}

	if txResult.Code > 0 {
		return "", errors.Errorf("cannot send tokens (SDK code %d): %s", txResult.Code, txResult.RawLog)
	}

	return txResult.TxHash, nil
}

// WaitTx waits until a tx is successfully added to a block and can be queried.
func (r Runner) WaitTx(ctx context.Context, txHash string, retryDelay time.Duration, maxRetry int) error {
	retry := 0

	// retry querying the request
	checkTx := func() error {
		b := newBuffer()
		if err := r.run(ctx, runOptions{stdout: b}, r.chainCmd.QueryTxCommand(txHash)); err != nil {
			// filter not found error and check for max retry
			if !strings.Contains(err.Error(), "not found") {
				return backoff.Permanent(err)
			}
			retry++
			if retry == maxRetry {
				return backoff.Permanent(errors.Errorf("can't retrieve tx %s", txHash))
			}
			return err
		}

		// parse tx and check code
		txResult, err := decodeTxResult(b)
		if err != nil {
			return backoff.Permanent(err)
		}
		if txResult.Code != 0 {
			return backoff.Permanent(errors.Errorf("tx %s failed: %s", txHash, txResult.RawLog))
		}

		return nil
	}
	return backoff.Retry(checkTx, backoff.WithContext(backoff.NewConstantBackOff(retryDelay), ctx))
}

// Export exports the state of the chain into the specified file.
func (r Runner) Export(ctx context.Context, exportedFile string) error {
	// Make sure the path exists
	dir := filepath.Dir(exportedFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	stdout, stderr := newBuffer(), newBuffer()
	if err := r.run(ctx, runOptions{stdout: stdout, stderr: stderr}, r.chainCmd.ExportCommand()); err != nil {
		return err
	}

	// Exported genesis is written on stderr from Cosmos-SDK v0.44.0
	var exportedState []byte
	if stdout.Len() > 0 {
		exportedState = stdout.Bytes()
	} else {
		exportedState = stderr.Bytes()
	}

	// Save the new state
	return os.WriteFile(exportedFile, exportedState, 0o600)
}

// EventSelector is used to query events.
type EventSelector struct {
	typ   string
	attr  string
	value string
}

// NewEventSelector creates a new event selector.
func NewEventSelector(typ, addr, value string) EventSelector {
	return EventSelector{typ, addr, value}
}

// Event represents a TX event.
type Event struct {
	Type       string
	Attributes []EventAttribute
	Time       time.Time
}

// EventAttribute holds event's attributes.
type EventAttribute struct {
	Key   string
	Value string
}

// QueryTxByEvents queries tx events by event selectors.
func (r Runner) QueryTxByEvents(
	ctx context.Context,
	selectors ...EventSelector,
) ([]Event, error) {
	if len(selectors) == 0 {
		return nil, errors.New("event selector list should be greater than zero")
	}
	list := make([]string, len(selectors))
	for i, event := range selectors {
		list[i] = fmt.Sprintf("%s.%s=%s", event.typ, event.attr, event.value)
	}
	query := strings.Join(list, "&")
	return r.QueryTx(ctx, r.chainCmd.QueryTxEventsCommand(query))
}

// QueryTxByQuery queries tx events by event selectors.
func (r Runner) QueryTxByQuery(
	ctx context.Context,
	selectors ...EventSelector,
) ([]Event, error) {
	if len(selectors) == 0 {
		return nil, errors.New("event selector list should be greater than zero")
	}
	list := make([]string, len(selectors))
	for i, query := range selectors {
		list[i] = fmt.Sprintf("%s.%s='%s'", query.typ, query.attr, query.value)
	}
	query := strings.Join(list, " AND ")
	return r.QueryTx(ctx, r.chainCmd.QueryTxQueryCommand(query))
}

// QueryTx queries tx events/query selectors.
func (r Runner) QueryTx(
	ctx context.Context,
	option ...step.Option,
) ([]Event, error) {
	// execute the command and parse the output.
	b := newBuffer()
	if err := r.run(ctx, runOptions{stdout: b}, option...); err != nil {
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
			TimeStamp string `json:"timestamp"`
		} `json:"txs"`
	}{}

	data, err := b.JSONEnsuredBytes()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
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

				txTime, err := time.Parse(time.RFC3339, tx.TimeStamp)
				if err != nil {
					return nil, err
				}

				events = append(events, Event{
					Type:       e.Type,
					Attributes: attrs,
					Time:       txTime,
				})
			}
		}
	}

	return events, nil
}
