package tsrelayer

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

// Call calls a method in the ts relayer wrapper lib with args and fills reply from the returned value.
func Call(ctx context.Context, method string, args, reply interface{}) error {
	command, cleanup, err := nodetime.Command(nodetime.CommandXRelayer)
	if err != nil {
		return err
	}
	defer cleanup()

	req, err := json2.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	res := &bytes.Buffer{}

	err = cmdrunner.New().Run(
		ctx,
		step.New(
			step.Exec(command[0], command[1:]...),
			step.Write(req),
			step.Stdout(res),
		),
	)
	if err != nil {
		return err
	}

	// regular logs can be printed to the stdout by the other process before a jsonrpc response is emitted.
	// differentiate both kinds and simulate printing regular logs if there are any.
	sc := bufio.NewScanner(res)
	for sc.Scan() {
		err = json2.DecodeClientResponse(bytes.NewReader(sc.Bytes()), reply)

		var e *json2.Error
		if errors.As(err, &e) || errors.Is(err, json2.ErrNullResult) { // jsonrpc returned with a server-side error.
			return err
		}

		if err != nil { // a line printed to the stdout by the other process.
			fmt.Println(sc.Text())
		}
	}

	return sc.Err()
}
