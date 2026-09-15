package runtime

import gno "github.com/gnolang/gno/gnovm/pkg/gnolang"

func GC(m *gno.Machine) {
	_, ok := m.GarbageCollect()
	if !ok {
		panic("should not happen, allocation limit exceeded while gc.")
	}
}

func MemStats(m *gno.Machine) string {
	return m.Alloc.MemStats()
}
