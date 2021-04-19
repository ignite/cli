---
order: 1
---

# Protocol Buffer Files Overview

Protocol buffer files define the data structures used by Cosmos SDK modules. 

Inside the `proto` directory, a directory for each custom module contains `query.proto`, `tx.proto`, `genesis.proto`, and other files.

The `starport serve` command automatically generates Go code from proto files on every file change.

Third-party proto files, including those of Cosmos SDK and Tendermint, are bundled with Starport. To import third-party proto files in your custom proto files:

```proto
import "cosmos/base/query/v1beta1/pagination.proto";
```

You can also manually add third-party proto files. By default, Starport imports proto files from these directories: `third_party/proto` and `proto_vendor`. You can define third-party paths of the import directory in `config.yml`:

```yaml
build:
  proto:
    third_party_paths: ["my_third_party_proto"]
```
