---
sidebar_position: 8
---

# CLI

A command line interface (CLI) lets you interact with your app after it is running on a machine somewhere. Each module has its own namespace within the CLI that gives it the ability to create and sign messages that are destined to be handled by that module. 

The CLI also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI lets you do things like generate keys for a new account or check the status of an interaction you already had with the application.

The CLI for the scavenge module is present in the `tx.go` and `query.go` files in the `x/scavenge/client/cli/` directory.

- The `tx.go` file is for making transactions that contain messages that will ultimately update the state.
- The `query.go` file is for making queries let you read information from the state.

Both files use the [Cobra](https://github.com/spf13/cobra) library.

## The tx.go file

The `tx.go` file contains the `GetTxCmd` standard method that is used in the Cosmos SDK. This method is referenced later in the `module.go` file that describes exactly which attributes a modules has.

This method makes it easier to incorporate different modules for different reasons at the level of the actual application. You are focused on a module now, but later you create an application that uses this module and other modules that are already available within the Cosmos SDK.

## Commit solution

```go
// x/scavenge/client/cli/tx_commit_solution.go

package cli

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"scavenge/x/scavenge/types"
)

func CmdCommitSolution() *cobra.Command {
	cmd := &cobra.Command{
		// pass a solution as the only argument
		Use:   "commit-solution [solution]",
		Short: "Broadcast message commit-solution",
		// set the number of arguments to 1
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			solution := args[0]

			// find a hash of the solution
			solutionHash := sha256.Sum256([]byte(solution))

			// convert the solution hash to string
			solutionHashString := hex.EncodeToString(solutionHash[:])

			// convert a scavenger address to string
			var scavenger = clientCtx.GetFromAddress().String()

			// find the hash of solution and scavenger address
			var solutionScavengerHash = sha256.Sum256([]byte(solution + scavenger))

			// convert the hash to string
			var solutionScavengerHashString = hex.EncodeToString(solutionScavengerHash[:])

			// create a new message
			msg := types.NewMsgCommitSolution(clientCtx.GetFromAddress().String(), string(solutionHashString), string(solutionScavengerHashString))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// broadcast the transaction with the message to the blockchain
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```

Note that this file makes use of the `sha256` library for hashing the plain text solutions into the scrambled hashes. This activity takes place on the client side so the solutions are never leaked to any public entity that might want to sneak a peak and steal the bounty reward associated with the scavenges. You can also notice that the hashes are converted into hexadecimal representation to make them easy to read as strings. Hashes are ultimately stored as hexadecimal representations in the keeper.

## Submit scavenge

```go
// x/scavenge/client/cli/tx_submit_scavenge.go

package cli

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"scavenge/x/scavenge/types"
)

func CmdSubmitScavenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-scavenge [solution] [description] [reward]",
		Short: "Broadcast message submit-scavenge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// find a hash of the solution
			solutionHash := sha256.Sum256([]byte(args[0]))

			// convert the hash to string
			solutionHashString := hex.EncodeToString(solutionHash[:])
			argsDescription := string(args[1])
			argsReward := string(args[2])

			// create a new message
			msg := types.NewMsgSubmitScavenge(clientCtx.GetFromAddress().String(), string(solutionHashString), string(argsDescription), string(argsReward))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			// broadcast the transaction with the message to the blockchain
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```
