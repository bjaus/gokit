package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	}
}

func GracefulShutdown(ctx context.Context, srv *http.Server, timeout time.Duration) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			slog.Error("Shutdown timeout exceeded", "error", err)
		}
	}()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", err)
	}
}
