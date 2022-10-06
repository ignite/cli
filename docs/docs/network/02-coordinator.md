---
sidebar_position: 2
description: Ignite Network commands for coordinator.
---

# Coordinator Guide

Coordinators organize and launch new chains on Ignite Chain

---

## Publish a chain

```shell
ignite n publish https://github.com/ignite/example
```

#### Output

```shell
✔ Source code fetched
✔ Blockchain set up
✔ Chain's binary built
✔ Blockchain initialized
✔ Genesis initialized
✔ Network published
⋆ Launch ID: 3
```

`LaunchID` identifies the published blockchain on Ignite blockchain

## List all published chains

```
ignite n chain list
```

### Output

```
Launch Id 	Chain Id 	Source ...

3 		    example-1 	https://github.com/ignite/example
2 		    spn-10 		https://github.com/tendermint/spn
1 	        example-20 	https://github.com/tendermint/spn
```

---

## Approve validator requests

First, list requests:

```
ignite n request list 3
```

_NOTE: here "3" is specifying the `LaunchID`_

#### Output

```
Id 	Status 		Type 			Content
1 	APPROVED 	Add Genesis Account 	spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 100000000stake
2 	APPROVED 	Add Genesis Validator 	e3d3ca59d8214206839985712282967aaeddfb01@84.118.211.157:26656, spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
3 	PENDING 	Add Genesis Account 	spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
4 	PENDING 	Add Genesis Validator 	b10f3857133907a14dca5541a14df9e8e3389875@84.118.211.157:26656, spn1daefnhnupn85e8vv0yc5epmnkcr5epkqncn2le, 95000000stake
```

Approve the requests. Both syntaxes can be used: `1,2,3,4` and `1-3,4`.

```
ignite n request approve 3 3,4
```

#### Output

```
✔ Source code fetched
✔ Blockchain set up
✔ Requests format verified
✔ Blockchain initialized
✔ Genesis initialized
✔ Genesis built
✔ The network can be started
✔ Request(s) #3, #4 verified
✔ Request(s) #3, #4 approved
```

---

## Initiate the launch of a chain

```
ignite n chain launch 3
```

#### Output

```
✔ Chain 3 will be launched on 2022-10-01 09:00:00.000000 +0200 CEST
```

_This example output shows the launch time of the chain on the network._
