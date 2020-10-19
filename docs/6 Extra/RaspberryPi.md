# Starport And The Raspberry Pi

If you make a Starport Chain and push it to GitHub, you'll notice that it generates a raspberry Pi image as a build artifact.

These images allow developers using Starport to provide device images-- or even devices-- to their users with ease.

Additioonally, Starport itself has a GitHub action that builds a Starport development environment that lives on a Raspberry Pi.

## Starport Pi Development Environment

The Starport Pi development environment provides a completely isolated **computer** for you to use when working on starport blockchains.

It's a convenience tool for developers: it contains a perfect and up to date development environment, and if you fork Starport, your fork will create these images too.

## Starport Downstream Pi Images

The Starport Downstream Pi Images produced in the Github actions files are a bit more than images that only run Starport.

They provide a new way of validating Proof of Stake blockchains exclusively from the edge of the network because they natively support the [Zerotier](https://zerotier.com) SD-WAN technology.

As a developer, you can choose to start your own zerotier network, and have devices using your image automatically join the same virtual network at boot time.

By default, Downstream Pi Images do not join a Zerotier network, because it is not clear weather or not joining a global, flat LAN like [`earth`](https://zerotier.atlassian.net/wiki/spaces/SD/pages/7372813/The+Earth+Test+Network).

Historically, one of the largest blockers of the development of peer to peer netwoorks and applications has been NAT (Network Address Translatioon).  Zerotier allows you to sidestep NAT and program hundreds or thousands -- or millions of devices to act in concert with one another.

### How to use Downstream Pi Images

* Buy:
    * [Raspberry Pi 4, 4GB version](https://www.raspberrypi.org/products/raspberry-pi-4-model-b/?variant=raspberry-pi-4-model-b-4gb)
    * [64GB or larger MicroSD Card](https://www.newegg.com/sandisk-64gb-microsdxc/p/N82E16820175006?Description=64gb%20microsd&cm_re=64gb_microsd-_-20-175-006-_-Product) from a reputable brand
    * [Raspberry Pi Power Supply](https://www.raspberrypi.org/products/type-c-power-supply/)
    * [USB MicroSD Card Reader](https://www.newegg.com/iogear-gfr3c11-2-in-1/p/N82E16820283035?Description=usb%20c%20microsd&cm_re=usb_c%20microsd-_-20-283-035-_-Product&quicklink=true)
    * 1 Ethernet Cable for each Pi (wifi is much less convienent, and is sometimes unreliable)
    
NB: If you are developing and want to target Pi, you'll want at least two Raspberry Pi devices.

* Download:
    * [Etcher](https://www.balena.io/etcher/) 
    * [The Pi Image From Github Actions](https://github.com/faddat/clay/actions/runs/262801323)
    
* Install
    * Use Etcher to flash the pi image to the MicroSD Card
    * Put the MicroSD card in the Raspberry Pi and fire it up
    
* Login
    * `ssh ubuntu@yourpiipaddress`
        * password is `ubuntu`

* Check out ZeroTier:
    * `zerotier-cli join e4da7455b26d23be` connects you to Clay's network.
        * This network allows clay to use a wide variety of previouusly unavaiilable P2P technology because it can assume that there is no need to worry about NAT or reaching peers.
        
* Join the network.  You'll find that the yourchaind and yourchaincli commands work perfectly.

If you simply build a chain, you'll get a Raspberry Pi image that has a systemD unit file that's not enabled.  You will want to figure out what the genesis state of your chain should be, and ship that in the device image.  Just have Packer copy it into the device image by mentiooning it in `pibuild.json`.  You will also want to make sure to change the start command in `yourchaind.service` to one that includes a seed node with a public IP address if you're not using ZeroTier.

If you are using Zerotier, you can either have Validators connect directly to one another as persistent peers, or have each validator set up a Zerotier connection to one or more sentry nodes, which ensure that the Pi itself is not accessible.

There are security considerations when using a flat network, depending on what you want to accomplish, flat networks may not be desirable. 

Zerotier can be used for access control, as well, allowing you to create semi-private blockchains that still operate by normal Cosmos style proof of stake validator selection rules.

### How to think about Pi Validators
The device is the key, is the device.  The validator becomes a physical object, and its job is to sign blocks for one or more blockchain networks.  One Pi, One Key.

Keep the Pi Safe, keep your key safe.


### Shout Out
It seems right to leave a shout out here to Clint Nelsen, Guy Lepage and Lukas Etter, who worked with me on similar ideas at [Galaxy](https://github.com/galaxypi/galaxy)
