---
sidebar_position: 0
description: Learn module basics by writing and reading blog posts to your chain.
slug: /guide/blog
---

# Build a blog

In this tutorial, we will create a blockchain with a module that allows us to
write and read data from the blockchain. This module will implement the ability
to create and read blog posts, similar to a blogging application. The end user
will be able to submit new blog posts and view a list of existing posts on the
blockchain. This tutorial will guide you through the process of creating and
using this module to interact with the blockchain.

The goal of this tutorial is to provide step-by-step instructions for creating a
feedback loop that allows you to submit data to the blockchain and read that
data back from the blockchain. By the end of this tutorial, you will have
implemented a complete feedback loop and will be able to use it to interact with
the blockchain.

First, create a new `hello` blockchain with Ignite CLI:

```
ignite scaffold chain blog --address-prefix blog
```

The new blockchain is scaffolded with the `--address-prefix` blog flag to use `"blog"` instead of the default `"cosmos"` address prefix.