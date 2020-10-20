# Starport And The Raspberry Pi

If you make a Starport Chain and push it to GitHub, you'll notice that it generates a raspberry Pi image as a build artifact.

These images allow developers using Starport to provide device images-- or even devices-- to their users with ease, which increases node count, [blockchain health](https://www.notion.so/Blockchain-App-TCO-3a86ae028a7f4c589efdb3a538d19bf2), and network

Additionally, Starport itself has a GitHub action that builds a Starport development environment that lives on a Raspberry Pi.


## Starport Downstream Pi Images

The Starport Downstream Pi Images produced in the Github actions files are a bit more than images that only run Starport.

They provide a new way of validating Proof of Stake blockchains exclusively from the edge of the network because they natively support the [Zerotier](https://zerotier.com) SD-WAN technology.

As a developer, you can choose to start your own Zerotier network, and have devices using your image automatically join the same virtual network at boot time.

By default, join a Zerotier network flat LAN, [`earth`](https://zerotier.atlassian.net/wiki/spaces/SD/pages/7372813/The+Earth+Test+Network).  This will provide you with the initial connectivity necessary to perform an absolutely silent Genesis ceremony.

One of the largest blockers to the development of peer to peer netwoorks and applications has been NAT (Network Address Translation).  Zerotier allows you to sidestep NAT and program hundreds, thousands, millions of devices to act in concert with one another.


### How to use Downstream Pi Images

* **Buy**:
    * [Raspberry Pi 4, 4GB version](https://www.raspberrypi.org/products/raspberry-pi-4-model-b/?variant=raspberry-pi-4-model-b-4gb)
    * [64GB or larger MicroSD Card](https://www.newegg.com/sandisk-64gb-microsdxc/p/N82E16820175006?Description=64gb%20microsd&cm_re=64gb_microsd-_-20-175-006-_-Product) from a reputable brand
    * [Raspberry Pi Power Supply](https://www.raspberrypi.org/products/type-c-power-supply/)
    * [USB MicroSD Card Reader](https://www.newegg.com/iogear-gfr3c11-2-in-1/p/N82E16820283035?Description=usb%20c%20microsd&cm_re=usb_c%20microsd-_-20-283-035-_-Product&quicklink=true)
    * 1 Ethernet Cable for each Pi ([wifi](https://medium.com/@huobur/how-to-setup-wifi-on-raspberry-pi-4-with-ubuntu-20-04-lts-64-bit-arm-server-ceb02303e49b) is much less convienent, and is sometimes unreliable.)

**NB**: If you are developing and want to target Pi, you'll want at least two Raspberry Pi devices because it will be much easier and faster to swap out the images as you test.

* **Download**:
    * [Etcher](https://www.balena.io/etcher/)
    * The Pi Image From Github Actions, [**for example**](https://github.com/faddat/clay/actions/runs/262801323)
    * [Fing](https://www.fing.com/) For android / ios

* **Install**
    * Use Etcher to flash the pi image to the MicroSD Card
    * Put the MicroSD card in the Raspberry Pi and fire it up

* **Find your device**
    * Open Fing and use it to get your raspberry pi's IP address
    * Open the web UI at http://raspiipaddress:8080

* **Login**
    * `ssh ubuntu@raspiipaddress`
        * password is `ubuntu`
        * you'll be prompted to change your password at first login

* **Check out ZeroTier**:
    * By default, you'll be connected to [`earth`](https://zerotier.atlassian.net/wiki/spaces/SD/pages/7372813/The+Earth+Test+Network)
        * This network allows clay to use a wide variety of previouusly unavailable P2P technology because it can assume that there is no need to worry about NAT or reaching peers.
        * Earth isn't exactly safe.  You're on a LAN with all other peers.  As your blockchain evolves, you may wish to drop Zerotier, or even set up your own Zerotier [controller](https://github.com/key-networks/ztncui).  The [AUR](https://aur.archlinux.org/) also has an excellent Zerotier controller you may wish to use.


## Usage Modes
The Pi images have two usage modes: Pre-Genesis and Post-Genesis

You can find examples of both usage modes in [Clay](https://github.com/faddat/clay).


### Pre-Genesis
If you simply build a chain, you'll get a Raspberry Pi image that has a systemD unit file that's not enabled, and a vue.service that is enabled.  Got a later time, we may create a user interface that allows for simple, *optionally disclosable* genesis ceremonies.

You will want to figure out what the genesis state of your chain should be, and ship that in later device images as described in the post-genesis section.  Just have Packer copy it into the device image by mentioning it in `pibuild.json`.

Zerotier can be used for access control, as well, allowing you to create semi-private blockchains that still operate by normal Cosmos style proof of stake validator selection rules.


### Post-Genesis
After Genesis, these images can be used to distribute your blockchain application, UI and all.

You will also want to make sure to change the start command in `yourchaind.service` to one that includes a seed node with a public IP address if you're not using ZeroTier.

Post-genesis is still an evolving design pattern, but the intent is to deliver plug and play blockchain nodes that can serve as:

* Full nodes
* Validators
* Application Interfaces
* Network Gateways

Depending on your desired design, any one node can perform one or all of the roles above.


### How to think about Pi Validators
The device is the key, is the device.  The validator becomes a physical object, and its job is to sign blocks for one or more blockchain networks.  One Pi, One Key.

Keep the Pi Safe, keep your key safe.

* Security Considerations
    * Flat Networks
There are security considerations when using a flat network, depending on what you want to accomplish, flat networks may not be desirable.  Because of how effectively Zerotier penetrates network address translation, you may wish to have a more typical network layout with Sentries. Your sentries could be on a public network with real public IP addresses or on another virtual private network.  You could have your validators operate each on their own vlan with no public IP address, just as described in the conventional [Sentry Node Architecture](https://forum.cosmos.network/t/sentry-node-architecture-overview/454).
    * If all validators are on the same virtual network, there is risk of ddos attacks from one validator to another, not to mention other nastiness.
    * If all nodes except for validators are on the same virtual network, you can architect most of the interesting peer to peer application patterns that these images support, safely.


## Hardware build targets
The Raspberry Pi is our first build target in terms of hardware company is not the last. With only minor modifications, you should be able to build for a wide array of devices today:

* bananapi-r1 (Archlinux ARM)
* beaglebone-black (Archlinux ARM, Debian)
* jetson-nano (Ubuntu)
* odroid-u3 (Archlinux ARM)
* odroid-xu4 (Archlinux ARM, Ubuntu)
* parallella (Ubuntu)
* raspberry-pi-3 (Archlinux ARM (armv8))
* raspberry-pi-4 (Archlinux ARM (armv7), Ubuntu 20.04 LTS))
* wandboard (Archlinux ARM)

It is important to note that you should avoid 32 bit device images.  There is a [known security issue](https://github.com/tendermint/tendermint/blob/master/docs/tendermint-core/running-in-production.md#validator-signing-on-32-bit-architectures-or-arm) with Tendermint on 32 bit CPUs, so please avoid armv7 devices.


## Shout Outs
* @clintnelsen, @guylepage3 and @lukasetter, who worked on similar ideas with @faddat at [Galaxy](https://github.com/galaxypi/galaxy).
* @mkaczanowski built and documented [packer-builder-arm](https://github.com/mkaczanowski/packer-builder-arm), Which drastically sped up making device images.
