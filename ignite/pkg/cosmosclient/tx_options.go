package cosmosclient

// TxOptions contains options for creating a transaction.
// It is used by the CreateTxWithOptions method.
type TxOptions struct {
	// Memo is the memo to be used for the transaction.
	Memo string

	// GasLimit is the gas limit to be used for the transaction.
	// If GasLimit is set to 0, the gas limit will be automatically calculated.
	GasLimit uint64

	// Fees is the fees to be used for the transaction.
	Fees string
}
