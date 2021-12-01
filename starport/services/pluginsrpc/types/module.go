package plugintypes

import "context"

type Module interface {
	Init(context.Context) error
}
