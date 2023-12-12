# Creating a structure of the application

To create a structure for a blockchain application that enables users to lend
and borrow digital assets from each other, use the Ignite CLI to generate the
necessary code.

First, create a new blockchain called `loan` by running the following command:

```
ignite scaffold chain loan --no-module
```

The `--no-module` flag tells Ignite not to create a default module. Instead, you
will create the module yourself in the next step.

Next, change the directory to `loan/`:

```
cd loan
```

Create a module with a dependency on the standard Cosmos SDK `bank` module by
running the following command:

```
ignite scaffold module loan --dep bank
```

Create a `loan` model with a list of properties.

```
ignite scaffold list loan amount fee collateral deadline state borrower lender --no-message
```

The `--no-message` flag tells Ignite not to generate Cosmos SDK messages for
creating, updating, and deleting loans. Instead, you will generate the code for
custom messages.


To generate the code for handling the messages for requesting, approving,
repaying, liquidating, and cancelling loans, run the following commands:

```
ignite scaffold message request-loan amount fee collateral deadline
```

```
ignite scaffold message approve-loan id:uint
```

```
ignite scaffold message repay-loan id:uint
```

```
ignite scaffold message liquidate-loan id:uint
```

```
ignite scaffold message cancel-loan id:uint
```

Great job! By using a few simple commands with Ignite CLI, you have successfully
set up the foundation for your blockchain application. You have created a loan
model and included keeper methods to allow interaction with the store. In
addition, you have also implemented message handlers for five custom messages.

Now that the basic structure is in place, it's time to move on to the next phase
of development. In the coming sections, you will be focusing on implementing the
business logic within the message handlers you have created. This will involve
writing code to define the specific actions and processes that should be carried
out when each message is received.