package unicode

// Optimized as native bindings in tests.

func IsPrint(r rune) bool
func IsGraphic(r rune) bool
func SimpleFold(r rune) rune
func IsUpper(r rune) bool
