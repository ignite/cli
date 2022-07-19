package foo

type Foo[T any] struct {
	FooKeeper foo.keeper
	i         T
}

func (f Foo[T]) RegisterAPIRoutes()         {}
func (f Foo[T]) RegisterTxService()         {}
func (f Foo[T]) RegisterTendermintService() {}
func (f Foo[T]) Name() string               { return app.BaseApp.Name() }
func (f Foo[T]) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (f Foo[T]) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
