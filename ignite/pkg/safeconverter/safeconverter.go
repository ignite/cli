package safeconverter

type SafeToConvertToInt interface {
	uintptr | uint | uint32 | uint64 | int | int32 | int64
}

func ToInt[T SafeToConvertToInt](x T) int {
	return int(x)
}

func ToInt64[T SafeToConvertToInt](x T) int64 {
	return int64(x)
}
