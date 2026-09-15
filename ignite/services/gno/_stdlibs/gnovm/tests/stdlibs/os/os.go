package os

import (
	"time"

	"github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/tests/stdlibs/chain/runtime"
)

func X_write(m *gnolang.Machine, p []byte, isStderr bool) int {
	if isStderr {
		if w, ok := m.Output.(interface{ StderrWrite(p []byte) (int, error) }); ok {
			n, _ := w.StderrWrite(p)
			return n
		}
	}
	n, _ := m.Output.Write(p)
	return n
}

func X_sleep(m *gnolang.Machine, duration int64) {
	arg0 := m.LastBlock().GetParams1(m.Store).TV
	d := arg0.GetInt64()
	sec := d / int64(time.Second)
	nano := d % int64(time.Second)
	ctx := m.Context.(*runtime.TestExecContext)
	ctx.Timestamp += sec
	ctx.TimestampNano += nano
	if ctx.TimestampNano >= int64(time.Second) {
		ctx.Timestamp += 1
		ctx.TimestampNano -= int64(time.Second)
	}

	m.Context = ctx
}
