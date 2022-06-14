package adapter

import "github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/query"

const (
	EntityTX query.Entity = iota
	EntityEvent
)

const (
	FieldTXHash query.Field = iota
	FieldTXIndex
	FieldTXBlockHeight
	FieldTXBlockTime
	FieldTXCreateTime
	FieldEventID
	FieldEventTXHash
	FieldEventType
	FieldEventIndex
	FieldEventAttrName
	FieldEventAttrValue
	FieldEventCreateTime
)
