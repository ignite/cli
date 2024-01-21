# Creating the structure

Create a new blockchain with the following command:

```
ignite scaffold chain blog
```

This will create a new directory called `blog/` containing the necessary files
and directories for your [blockchain
application](https://docs.cosmos.network/main/basics/app-anatomy). Next,
navigate to the newly created directory by running:

```
cd blog
```

Since your app will be storing and operating with blog posts, you will need to
create a `Post` type to represent these posts. You can do this using the
following Ignite CLI command:

```
ignite scaffold type post title body creator id:uint
```

This will create a `Post` type with four fields: `title`, `body`, `creator`, all
of type `string`, and `id` of type `uint`.

It is a good practice to commit your changes to a version control system like
Git after using Ignite's code scaffolding commands. This will allow you to
differentiate between changes made automatically by Ignite and changes made
manually by developers, and also allow you to roll back changes if necessary.
You can commit your changes to Git with the following commands:

```
git add .
git commit -am "ignite scaffold type post title body"
```

### Creating messages

Next, you will be implementing CRUD (create, read, update, and delete)
operations for your blog posts. Since create, update, and delete operations
change the state of the application, they are considered write operations. In
Cosmos SDK blockchains, state is changed by broadcasting
[transactions](https://docs.cosmos.network/main/basics/tx-lifecycle) that
contain messages that trigger state transitions. To create the logic for
broadcasting and handling transactions with a "create post" message, you can use
the following Ignite CLI command:

```
ignite scaffold message create-post title body --response id:uint
```

This will create a "create post" message with two fields: `title` and `body`,
both of which are of type `string`. Posts will be stored in the key-value store
in a list-like data structure, where they are indexed by an incrementing integer
ID. When a new post is created, it will be assigned an ID integer. The
`--response` flag is used to return `id` of type `uint` as a response to the
"create post" message.

To update a specific blog post in your application, you will need to create a
message called "update post" that accepts three arguments: `title`, `body`, and
`id`. The `id` argument of type `uint` is necessary to specify which blog post
you want to update. You can create this message using the Ignite CLI command:

```
ignite scaffold message update-post title body id:uint
```

To delete a specific blog post in your application, you will need to create a
message called "delete post" that accepts only the `id` of the post to be
deleted. You can create this message using the Ignite CLI command:

```
ignite scaffold message delete-post id:uint
```

### Creating queries

[Queries](https://docs.cosmos.network/main/basics/query-lifecycle) allow users
to retrieve information from the blockchain state. In your application, you will
have two queries: "show post" and "list post". The "show post" query will allow
users to retrieve a specific post by its ID, while the "list post" query will
return a paginated list of all stored posts.

To create the "show post" query, you can use the following Ignite CLI command:

```
ignite scaffold query show-post id:uint --response post:Post
```

This query will accept `id` of type `uint` as an argument, and will return a
`post` of type `Post` as a response.

To create the "list post" query, you can use the following Ignite CLI command:

```
ignite scaffold query list-post --response post:Post --paginated
```

This query will return a post of type Post in a paginated output. The
`--paginated` flag indicates that the query should return its results in a
paginated format, allowing users to retrieve a specific page of results at a
time.

## Summary

Congratulations on completing the initial setup of your blockchain application!
You have successfully created a "post" data type and generated the necessary
code for handling three types of messages (create, update, and delete) and two
types of queries (list and show posts).

However, at this point, the messages you have created will not trigger any state
transitions, and the queries you have created will not return any results. This
is because Ignite only generates the boilerplate code for these features, and it
is up to you to implement the necessary logic to make them functional.

In the next chapters of the tutorial, you will learn how to implement the
message handling and query logic to complete your blockchain application. This
will involve writing code to process the messages and queries you have created
and use them to modify or retrieve data from the blockchain's state. By the end
of this process, you will have a fully functional blog application on a Cosmos
SDK blockchain.