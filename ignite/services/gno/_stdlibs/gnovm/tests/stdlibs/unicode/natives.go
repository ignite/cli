package unicode

import "unicode"

func IsPrint(r rune) bool    { return unicode.IsPrint(r) }
func IsGraphic(r rune) bool  { return unicode.IsGraphic(r) }
func SimpleFold(r rune) rune { return unicode.SimpleFold(r) }
func IsUpper(r rune) bool    { return unicode.IsUpper(r) }
