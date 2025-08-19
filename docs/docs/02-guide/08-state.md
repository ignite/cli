---
description: Learn how Cosmos SDK modules manage state with collections
title: State Management
---

# State Management in Modules

In blockchain applications, state refers to the current data stored on the blockchain at a specific point in time. Handling state is usually the core of any blockchain application. The Cosmos SDK provides powerful tools for state management, with the `collections` package being the recommended approach for modern applications.

## Collections Package

Ignite scaffolds using the [`collections`](http://pkg.go.dev/cosmossdk.io/collections) package for module code. This package provides a type-safe and efficient way to set and query values from the module store.

### Key Features of Collections

- **Type Safety**: Collections are type-safe, reducing the risk of runtime errors.
- **Simplified API**: Easy-to-use methods for common operations like Get, Set, and Has.
- **Performance**: Optimized for performance with minimal overhead.
- **Integration**: Seamlessly integrates with the Cosmos SDK ecosystem.

## Understand keeper field

Ignite creates all the necessary boilerplate for collections in the `x/<module>/keeper/keeper.go` file. The `Keeper` struct contains fields for each collection you define in your module. Each field is an instance of a collection type, such as `collections.Map`, `collections.Item`, or `collections.List`.

```go
type Keeper struct {
	// ...

	Params   collections.Item[Params]
	Counters collections.Map[string, uint64]
	Profiles collections.Map[sdk.AccAddress, Profile]
}
```

## Common State Operations

### Reading State

To read values from state, use the `Get` method:

```go
// getting a single item
params, err := k.Params.Get(ctx)
if err != nil {
	// handle error
	// collections.ErrNotFound is returned when an item doesn't exist
}

// getting a map entry
counter, err := k.Counters.Get(ctx, "my-counter")
if err != nil {
	// handle error
}
```

### Writing State

To write values to state, use the `Set` method:

```go
// setting a single item
err := k.Params.Set(ctx, params)
if err != nil {
	// handle error
}

// setting a map entry
err = k.Counters.Set(ctx, "my-counter", 42)
if err != nil {
	// handle error
}
```

### Checking Existence

Use the `Has` method to check if a value exists without retrieving it:

```go
exists, err := k.Counters.Has(ctx, "my-counter")
if err != nil {
	// handle error
}
if exists {
	// value exists
}
```

### Removing State

To remove values from state, use the `Remove` method:

```go
err := k.Counters.Remove(ctx, "my-counter")
if err != nil {
	// handle error
}
```

## Implementing Business Logic in Messages

Messages in Cosmos SDK modules modify state based on user transactions. Here's how to implement business logic in a message handler using collections:

```go
func (k msgServer) CreateProfile(ctx context.Context, msg *types.MsgCreateProfile) (*types.MsgCreateProfileResponse, error) {
	// validate message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// parse sender address
	senderBz, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, err
	}
	sender := sdk.AccAddress(senderBz)

	// check if profile already exists
	exists, err := k.Profiles.Has(ctx, sender)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, sdkerrors.Wrap(types.ErrProfileExists, "profile already exists")
	}

	// create new profile
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	profile := types.Profile{
		Name:      msg.Name,
		Bio:       msg.Bio,
		CreatedAt: sdkCtx.BlockTime().Unix(),
	}

	// store the profile
	err = k.Profiles.Set(ctx, sender, profile)
	if err != nil {
		return nil, err
	}

	// increment profile counter
	counter, err := k.Counters.Get(ctx, "profiles")
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	// set the counter (adding 1)
	err = k.Counters.Set(ctx, "profiles", counter+1)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateProfileResponse{}, nil
}
```

## Implementing Queries

Queries allow users to read state without modifying it. Here's how to implement a query handler using collections:

```go
func (q queryServer) GetProfile(ctx context.Context, req *types.QueryGetProfileRequest) (*types.QueryGetProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// parse address
	addressBz, err := k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}
	address := sdk.AccAddress(addressBz)

	// get profile
	profile, err := q.k.Profiles.Get(ctx, address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "profile not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetProfileResponse{Profile: profile}, nil
}
```

## Error Handling with Collections

When working with collections, proper error handling is essential:

```go
// example from a query function
params, err := q.k.Params.Get(ctx)
if err != nil && !errors.Is(err, collections.ErrNotFound) {
	return nil, status.Error(codes.Internal, "internal error")
}
```

In the snippet above, it uses the `Get` method to get a collection item. A `collections.ErrNotFound` can be a valid error when the collection is empty, whereas any other error is considered an internal error that should be handled appropriately.

## Iterating Over Collections

Collections also support iteration:

```go
// iterate over all profiles
err := k.Profiles.Walk(ctx, nil, func(key sdk.AccAddress, value types.Profile) (bool, error) {
	// process each profile
	// return true to stop iteration, false to continue
	return false, nil
})
if err != nil {
	// handle error
}

// iterate over a range of counters
startKey := "a"
endKey := "z"
err = k.Counters.Walk(ctx, collections.NewPrefixedPairRange[string, uint64](startKey, endKey), func(key string, value uint64) (bool, error) {
	// process each counter in the range
	return false, nil
})
if err != nil {
	// handle error
}
```

## Conclusion

The `collections` package provides a powerful and type-safe way to manage state in Cosmos SDK modules. By understanding how to use collections effectively, you can build robust and efficient blockchain applications that handle state transitions reliably.

When developing with Ignite CLI, you are already taking advantage of collections which significantly simplify the state management code and reduce the potential for errors.
