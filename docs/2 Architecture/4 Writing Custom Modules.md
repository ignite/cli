# Writing custom modules

Starport allows you to jump directly into creating your own module. With the before described `type` function you can add new transaction types to your application. Under the hood, starport creates a handler, types and messages for you. 

Without using starport, you would need to manipulate these functions yourself. Here is what starport does when you add a `type`. Understanding what starport does might help in order to either add more complex structures or debug in case something does not work as it should.

## Proto

When using the type command a Protobuffer definition is created for you in the `proto` directory.

It contains messages for full CRUD (Create, Read, Update, Delete) operations for your created transaction type.

## Module

Once you have created your starport blockchain application, you will have your own module resident of `yourapp/x/yourmodule`, it comes predefined with a couple of files and folders which define types, functions and messages of your module.

## Types

The `types` folder defines structures of your golang blockchain application. Here you can define your basic or more advanced types which will later be data and functions usable on your blockchain.

The message types are defined in the file `types/messages_type`, or other functions that you are planning to use.

## Client

There is a `rest` folder that takes care of the rest API that your module exposes.

The `cli` folder with the contents take care of the Command line interface commands that will be available to the user.

## Frontend

Currently starport provides a basic Vue User-Interface that you can get inspired by or build ontop on. The source code is available in the `vue` folder. Written in JavaScript you can hop directly into writing the frontend for your application.

[Learn more about Vue.](https://vuejs.org/)

## Summary

- Starport bootstraps a module for you.
- You can change a module by modifying the files in `yourapp/x/` or the `proto` diretory.
- Starport has a Vue frontend where you can start to work immediately.
