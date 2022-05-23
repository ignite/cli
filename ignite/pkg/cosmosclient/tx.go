package cosmosclient

// TX contains transaction information.
type TX struct {
	// Hash contains the transaction hash.
	Hash string

	// Height of the block that contains the transaction.
	Height int64

	// BlockTime contains the timestamp of the block that contains the transaction.
	BlockTime string

	// EventLog contains the events emitted during the execution.
	// The value is a JSON string containing the list of events.
	EventLog string
}
