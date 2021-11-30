---
order: 13
description: Chain Simulation
---

# Chain Simulation

The simulator can help you test different scenarios for your chain, simulating the messages, blocks, and accounts. Starport scaffold a base simulation for each module and also dummy simulation methods for each scaffold message.

## Module simulation

Every new module implements the Cosmos SDK simulator. Each new message creates a new file with the simulation methods necessary for the tests. Scaffolding a `CRUD` like a `list` or `map` creates a simulation file with `create`, `update` and `delete` simulation methods into the `x/<module>/simulation` folder and registers this methods into the `x/<module>/module_simulation.go`. Scaffolding a single message creates an empty simulation method to be implemented by the user. Also, the user should maintain the simulation methods for each new modification into the message keeper methods.
Every simulation has your weight. This is because the sender of the operation is assigned randomly. The weight defines how much the simulation will call the message. For better randomizes, the developer can define a random seed, the simulation with the same random seed is deterministic with the same output.

## Scaffolding Simulation

Creates a new chain:
```shell
starport scaffold chain github.com/cosmonaut/mars
```

You can check the empty `x/mars/simulation` folder and the  `x/mars/module_simulation.go` file without any simulation registered. We will scaffold a new message:
```shell
starport scaffold list user address balance:uint state
```

Starport creates the new file `x/mars/simulation/user.go`, and also, the simulation will register with the weight into the `x/mars/module_simulation.go`. The developer should define the proper simulation weight. For this example, we can change the `defaultWeightMsgDeleteUser` to 30 and the `defaultWeightMsgUpdateUser` to 50. 100 is the maximum, and zero is the minimum.

Run the `BenchmarkSimulation` method into `app/simulation_test.go` to run simulation tests for all modules:
```shell
starport chain simulation
```

You can also define some flags provided by the simulation defined by the method `simapp.GetSimulatorFlags()`
```shell
starport chain simulation --NumBlocks 100 --BlockSize 200
```

Wait for the entire simulation to finish and check the result of the messages.

The default go test command works to run the simulation as well:
```shell
go test -v -benchmem -run=^$ -bench ^BenchmarkSimulation -cpuprofile cpu.out ./app -Commit=true
```

### Skip message

The developer can handle a logic to avoid sending a message without returning an error, only return `simtypes.NoOpMsg(...)`  into the simulation message handler.

## Params

Scaffolding a module with params will automatically add it into the  `module_simulaton.go` file:

```shell
starport s module earth --params channel:string,minLaunch:uint,maxLaunch:int
```