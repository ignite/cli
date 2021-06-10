# Starport CLI Action
This action makes the `starport` cli available as a Github Action.

## Quick start

Add a new workflow to your repo:

```yml
on: push

jobs:
  help:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2

      - name: Print Help 
        uses: tendermint/starport/actions/cli@develop
        with:
          args: -h 
```
