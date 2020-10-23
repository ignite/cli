# Writing custom modules

Starport allows you to jump directly into creating your own module. With the before described `type` function you can add new types to your application. Under the hood, starport creates a handler, types and messages for you. 

Without using starport, you would need to manipulate these functions yourself. Here is what starport does when you add a `type`. Understanding what starport does might help in order to either add more complex structures or debug in case something does not work as it should.

Once you have created your starport blockchain application, you will have your own module resident of `yourapp/x/yourmodule`, it comes predefined with a couple of files and folders which deifne types, functions and messages of your module.

## Types

The `types` folder defines structures of your golang blockchain application. Here you can define your basic or more advanced types which will later be data and functions usable on your blockchain.

Your basic `types` reside in `types/TypeModule.go`, define the fields and functions of your application.

The message types are defined in the file `types/MsgCreateX`, or other functions that you are planning to use.

There is a `rest` folder that takes care of the rest API that your module exposes.

The `cli` folder with the contents take care of the Command line interface commands that will be available to the user.

## Frontend

Currently, Starport provides a basic Vue user interface that you can take inspiration from or use it as a starting point for your app. The source code is available inside the `vue` directory. Written in JavaScript you can hop directly into writing the frontend for your application.

[Learn more about Vue.](https://vuejs.org/)

## Summary

- Starport bootstraps a module for you.
- You can change a module by modifying the files in `yourapp/x/`.
- Starport has a Vue frontend where you can start to work immediately.
