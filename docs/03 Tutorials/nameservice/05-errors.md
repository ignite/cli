---
order: 5
---

# Errors

Start by navagating to the `errors.go` file within the types folder. Within your `errors.go` file, define errors that are custom to your module along with their codes.

<<< @/nameservice/nameservice/x/nameservice/types/errors.go

You must also add the corresponding method that'll be called at the time of error handling. For instance, let's say we try to delete a name that is not present in the store. In this case, an error should be thrown as the name does not exist.

We'll see later on in the tutorial where this method is called.

Now we move on to writing the Keeper for the module.
