# Protocol Buffer Files Overview

Protocol buffer files define the data structures used by Cosmos SDK modules. 

Inside the `proto` directory, a directory for each custom module contains `query.proto`, `tx.proto`, `genesis.proto`, and other files.

`starport serve` automatically generates Go code from proto files on every file change.

Third-party proto files, including those of Cosmos SDK and Tendermint are bundled with Starport. You can import them in your custom proto files like so:

```proto
import "cosmos/base/query/v1beta1/pagination.proto";
```

Additional third-party proto files can be added manually. By default Starport imports proto files from two directories: `third_party/proto` and `proto_vendor`. This can be customised in `config.yml`:

```yaml
build:
  proto:
    third_party_paths: ["my_third_party_proto"]
```
