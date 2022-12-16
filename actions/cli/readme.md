# Ignite CLI Action

This action makes the `ignite` CLI available as a Github Action.

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
        uses: ignite/cli/actions/cli@main
        with:
          args: -h 
```
