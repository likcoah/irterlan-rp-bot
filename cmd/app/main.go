package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/likcoah/irterlan-rp-bot/internal/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	c, err := bot.New(os.Getenv("BOT_TOKEN"), os.Getenv("BOT_WEBHOOK_URL"),
		os.Getenv("BOT_WEBHOOK_SECRET"), os.Getenv("PORT"))
	if err != nil {
		slog.Error("bot initialization failed, check env tokens", "err", err)
		return
	}

	slog.Info("bot started")
	err = c.Start(ctx)
	if err != nil { slog.Error("bot runtime error", "err", err); return }
	slog.Info("bot stopped")
}

