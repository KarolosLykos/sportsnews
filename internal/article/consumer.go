package article

import (
	"context"
)

type Consumer interface {
	Consume(ctx context.Context)
}
