# How to use Starport

Starport creates a blockchain for you in Go. Requirements for this is to have Go installed. You can get all the information here https://golang.org/doc/install.

Starport installation instructions can be found here: https://github.com/tendermint/starport#install

To the machine you are executing on there are not many requirements. It runs on Linux or Mac Operating Systems and can be run from a Raspberry Pi.

## Your blockchain application

To create a blockchain application we use the command `app`

```bash
starport app github.com/username/myapp
```

This will create the folder `myapp` and is a usable blockchain blueprint. If you want to dive directly into looking at the details of your blockchain you can run it with entering your `myapp` folder and use the command `serve` to initialise your blockchain and start it.

`starport serve`

The output of the `serve` command will already indicate that you find helpful information, guides and analytics about your blockchain, development tips as well as an interaction user interface on http://localhost:12345/.

The first step of your own blockchain is already done. Using the default settings, a blockchain that has networking, consensus protocol with an own token is hereby established. From here on, you can implement logic that makes your own blockchain unique. 

## The Key-Value Store (KV)
### How to use types

In the SDK, data are stored in the multistore. As Key-Value pairs those are saved in the KVStores. Multiple stores can be created and managed at the same time. We will use the store to save our data to the blockchain.
Starport assists us in setting up the Key-Value Store with the command `type`. 
In order to use `type` we should give our type a fitting `typeName` with the intended fields that we want to use. If we wanted to store user with username and age, we would use the command

`starport type user username age:int` 

Now a Key-Value Store for the user with fields username and age is created. We can create a new user with the command

`myappcli tx myapp create-user "my-first-username" 35`

Which creates the user with username `my-first-username` and age of `35`. 

These are the basic commands navigating starport. From creating a first blockchain to adding your own data types and accessing the User Interface. In the next two chapters, we will be looking closer at the initial setup for starport and how to configure it. Afterwards, we will be looking into more complex usecases, where each of the commands and more will be explained in more detail.

## Summary

- With the command `starport add` a new blockchain can be initialised.
- A combination `starport add` and `starport serve` already let's you manage your blockchain out of the box.
- The default blockchain includes networking and a consensus protocol with your own token.
- Data is managed with the Key-Value Store and data types can be added with `starport type`.

[◀️ Previous - Development Mode](../../01%20introduction/03_development_mode/03_development_mode.md) | [▶️ Next - Genesis File](../../02%20using%20starport/02_genesis_file/02_genesis_file.md)  