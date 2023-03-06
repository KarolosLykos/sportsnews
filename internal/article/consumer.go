package article

import (
	"context"
)

type Consumer interface {
	List(ctx context.Context)
}
