package v1

// CommandPath returns the absolute command path including the binary name as prefix.
func (h *Hook) CommandPath() string {
	return ensureFullCommandPath(h.PlaceHookOn)
}
