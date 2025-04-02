---
sidebar_position: 10
description: Test different scenarios for your chain.
---

# Chain simulation

The Ignite CLI chain simulator can help you to run your chain based in
randomized inputs for you can make fuzz testing and also benchmark test for your
chain, simulating the messages, blocks, and accounts. You can scaffold a
template to perform simulation testing in each module along with a boilerplate
simulation methods for each scaffolded message.

## Module simulation

Every new module that is scaffolded with Ignite CLI implements the Cosmos SDK
[Module
Simulation](https://docs.cosmos.network/main/building-modules/simulator).

- Each new message creates a file with the simulation methods required for the
  tests.
- Scaffolding a `CRUD` type like a `list` or `map` creates a simulation file
  with `create`, `update`, and `delete` simulation methods in the
  `x/<module>/simulation` folder and registers these methods in
  `x/<module>/module_simulation.go`.
- Scaffolding a single message creates an empty simulation method to be
  implemented by the user.

We recommend that you maintain the simulation methods for each new modification
into the message keeper methods.

Every simulation is weighted because the sender of the operation is assigned
randomly. The weight defines how much the simulation calls the message.

For better randomizations, you can define a random seed. The simulation with the
same random seed is deterministic with the same output.

## Scaffold a simulation

To create a new chain:

```
ignite scaffold chain mars
```

Review the empty `x/mars/simulation` folder and the
`x/mars/module_simulation.go` file to see that a simulation is not registered.

Now, scaffold a new message:

```
ignite scaffold list user address balance:uint state
```

A new file `x/mars/simulation/user.go` is created and is registered with the
weight in the `x/mars/module_simulation.go` file.

Be sure to define the proper simulation weight with a minimum weight of 0 and a
maximum weight of 100.

For this example, change the `defaultWeightMsgDeleteUser` to 30 and the
`defaultWeightMsgUpdateUser` to 50.

Run the `BenchmarkSimulation` method into `app/simulation_test.go` to run
simulation tests for all modules:

```
ignite chain simulate
```

You can also define flags that are provided by the simulation. Flags are defined
by the method `simapp.GetSimulatorFlags()`:

```
ignite chain simulate -v --numBlocks 200 --blockSize 50 --seed 33
```

Wait for the entire simulation to finish and check the result of the messages.

The default `go test` command works to run the simulation:

```
go test -v -benchmem -run=^$ -bench ^BenchmarkSimulation -cpuprofile cpu.out ./app -Commit=true
```

### Skip message

Use logic to avoid sending a message without returning an error. Return only
`simtypes.NoOpMsg(...)` into the simulation message handler.

## Params

Scaffolding a module with params automatically adds the module in the
`module_simulaton.go` file:

```
ignite s module earth --params channel:string,minLaunch:uint,maxLaunch:int
```

After the parameters are scaffolded, change the
`x/<module>/module_simulation.go` file to set the random parameters into the
`RandomizedParams` method. The simulation will change the params randomly
according to call the function.
