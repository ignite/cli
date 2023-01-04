# Vue frontend

Welcome to this tutorial on using Ignite to develop a web application for your
blockchain with Vue 3. Ignite is a tool that simplifies the process of building
a blockchain application by providing a set of templates and generators that can
be used to get up and running quickly.

One of the key features of Ignite is its support for Vue 3, a popular JavaScript
framework for building user interfaces. In this tutorial, you will learn how to
use Ignite to create a new blockchain and scaffold a Vue frontend template. This
will give you a basic foundation for your web application and make it easier to
get started building out the rest of your application.

Once you have your blockchain and Vue template set up, the next step is to
generate an API client. This will allow you to easily interact with your
blockchain from your web application, allowing you to retrieve data and make
transactions. By the end of this tutorial, you will have a fully functional web
application that is connected to your own blockchain.

Create a new blockchain project:

```
ignite scaffold chain example
```

To create a Vue frontend template, go to the `example` directory and run the
following command:

```
ignite scaffold vue
```

This will create a new Vue project in the `vue` directory. This project can be
used with any blockchain, but it depends on an API client to interact with the
blockchain. To generate an API client, run the following command in the
`example` directory:

```
ignite generate composables
```

This command generates two directories:

* `ts-client`: a framework-agnostic TypeScript client that can be used to
  interact with your blockchain. You can learn more about how to use this client
  in the [TypeScript client tutorial](/clients/typescript).
* `vue/src/composables`: a collection of Vue 3
  [composables](https://vuejs.org/guide/reusability/composables.html) that wrap
  the TypeScript client and make it easier to interact with your blockchain from
  your Vue application.

To start your Vue application, go to the `vue` directory and run the following
command:

```
npm install && npm run dev
```

It is recommended to run `npm install` before starting your app with `npm run
dev` to ensure that all dependencies are installed (including the ones that the
API client has, see `vue/postinstall.js`).

In a separate terminal window in the `example` directory run the following
command to start your blockchain:

```
ignite chain serve
```
