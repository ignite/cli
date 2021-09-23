package types

import "encoding/binary"

type (
	// OracleScriptID is the type-safe unique identifier type for oracle scripts.
	OracleScriptID uint64

	// OracleRequestID is the type-safe unique identifier type for data requests.
	OracleRequestID int64
)

// int64ToBytes convert int64 to a byte slice
func int64ToBytes(num int64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, uint64(num))
	return result
}
