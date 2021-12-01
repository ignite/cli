---
order: 13
description: Test different scenarios for your chain. 

---

# Chain Simulation

The Starport chain simulator can help you test different scenarios for your chain, simulating the messages, blocks, and accounts. You can scaffold a base simulation for each module along with dummy simulation methods for each scaffolded message.

## Module Simulation

Every new module that is scaffolded with Starport implements the Cosmos SDK [Module Simulation](https://docs.cosmos.network/master/building-modules/simulator.html). 

- Each new message creates a file with the simulation methods required for the tests. 
- Scaffolding a `CRUD` like a `list` or `map` creates a simulation file with `create`, `update`, and `delete` simulation methods in the `x/<module>/simulation` folder and registers these methods in `x/<module>/module_simulation.go`. 
- Scaffolding a single message creates an empty simulation method to be implemented by the user. 

We recommend that you maintain the simulation methods for each new modification into the message keeper methods.

Every simulation is weighted because the sender of the operation is assigned randomly. The weight defines how much the simulation calls the message. 

For better randomizations, you can define a random seed. The simulation with the same random seed is deterministic with the same output.

## Scaffold a Simulation

To create a new chain:

```shell
starport scaffold chain github.com/cosmonaut/mars
```

Review the empty `x/mars/simulation` folder and the `x/mars/module_simulation.go` file to see that a simulation is not registered. 

Now, scaffold a new message:

```shell
starport scaffold list user address balance:uint state
```

A new file `x/mars/simulation/user.go` is created and is registered with the weight in the `x/mars/module_simulation.go` file. 

Be sure to define the proper simulation weight with a minimum weight of 0 and a maximum weight of 100. 

For this example, change the `defaultWeightMsgDeleteUser` to 30 and the `defaultWeightMsgUpdateUser` to 50. 

Run the `BenchmarkSimulation` method into `app/simulation_test.go` to run simulation tests for all modules:

```shell
starport chain simulation
```

You can also define flags that are provided by the simulation. Flags are defined by the method `simapp.GetSimulatorFlags()`:

```shell
starport chain simulation --NumBlocks 100 --BlockSize 200
```

Wait for the entire simulation to finish and check the result of the messages.

The default `go test` command works to run the simulation:

```shell
go test -v -benchmem -run=^$ -bench ^BenchmarkSimulation -cpuprofile cpu.out ./app -Commit=true
```

### Skip Message

Use logic to avoid sending a message without returning an error. Return only `simtypes.NoOpMsg(...)` into the simulation message handler.

## Params

Scaffolding a module with params automatically adds the module in the  `module_simulaton.go` file:

```shell
starport s module earth --params channel:string,minLaunch:uint,maxLaunch:int
```