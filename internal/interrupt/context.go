package interrupt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Context returns a context that is canceled on SIGINT and SIGTERM.
func Context() (context.Context, func()) {
	return WrappedContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

// WrappedContext returns a new context wrapping the provided context and
// canceling it on the provided signals.
func WrappedContext(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}
