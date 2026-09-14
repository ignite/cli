package time

import (
	"time"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/internal/execctx"
)

func X_now(m *gno.Machine) (sec int64, nsec int32, mono int64) {
	if m == nil || m.Context == nil {
		return 0, 0, 0
	}

	ctx := execctx.GetContext(m)
	return ctx.Timestamp, int32(ctx.TimestampNano), ctx.Timestamp*int64(time.Second) + ctx.TimestampNano
}
