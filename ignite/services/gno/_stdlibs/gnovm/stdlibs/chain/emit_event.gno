package chain

// Emit is a function that constructs a GnoEvent with a specified type and attributes.
// It then forwards this event to the event logger. Each emitted event carries metadata
// such as the event type, the initializing realm, and the provided attributes.
//
// The function takes type and attribute strings. typ (or type) is an arbitrary string that represents
// the type of the event. It plays a role of indexing what kind of event occurred. For example,
// a name like "Foo" or "Bar" can be used to indicating the purpose (nature) of the event.
//
// And the attrs (attributes) accepts an even number of strings and sets them as key-value pairs
// according to the order they are passed. For example, if the attr strings "key1", "value1" are
// passed in, the key is set to "key1" and the value is set to "value1".
//
// The event is dispatched to the EventLogger, which resides in the tm2/pkg/sdk/events.go file.
//
// For more details about the GnoEvent data structure, refer to its definition in the emit_event.go file.
func Emit(typ string, attrs ...string) { emit(typ, attrs) }
func emit(typ string, attrs []string)
