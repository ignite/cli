---
order: 8
---

# CLI

A Command Line Interface (CLI) will help us interact with our app once it is running on a machine somewhere. Each Module has it's own namespace within the CLI that gives it the ability to create and sign Messages destined to be handled by that module. It also comes with the ability to query the state of that module. When combined with the rest of the app, the CLI will let you do things like generate keys for a new account or check the status of an interaction you already had with the application.

The CLI for our module is broken into two files called `tx.go` and `query.go` which are located in `./x/scavenge/client/cli/`. One file is for making transactions that contain messages which will ultimately update our state. The other is for making queries which will give us the ability to read information from our state. Both files utilize the [Cobra](https://github.com/spf13/cobra) library.

## `tx.go`

The `tx.go` file contains `GetTxCmd` which is a standard method within the Cosmos SDK. It is referenced later in the `module.go` file which describes exactly which attributes a modules has. This makes it easier to incorporate different modules for different reasons at the level of the actual application. After all, we are focusing on a module at this point, but later we will create an application that utilizes this module as well as other modules which are already available within the Cosmos SDK.

## Commit Solution

```go
 // x/scavenge/client/cli/tx_commit_solution.go
func CmdCommitSolution() *cobra.Command {
	cmd := &cobra.Command{
    // pass a solution as the only argument
		Use:   "commit-solution [solution]",
		Short: "Broadcast message commit-solution",
    // set the number of arguments to 1
		Args:  cobra.ExactArgs(1),
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

Note that this file makes use of the `sha256` library for hashing our plain text solutions into the scrambled hashes. This activity takes place on the client side so the solutions are never leaked to any public entity which might want to sneak a peak and steal the bounty reward associated with the scavenges. You can also notice that the hashes are converted into hexadecimal representation to make them easy to read as strings (which is how they are ultimately stored in the keeper).


## Create Scavenge

```go
// x/scavenge/client/cli/tx_create_scavenge.go
func CmdCreateScavenge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-scavenge [solution] [description] [reward]",
		Short: "Broadcast message create-scavenge",
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
			msg := types.NewMsgCreateScavenge(clientCtx.GetFromAddress().String(), string(solutionHashString), string(argsDescription), string(argsReward))
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