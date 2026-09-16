package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/likcoah/irterlan-rp-bot/internal/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	c, err := bot.New(os.Getenv("BOT_TOKEN"))
	if err != nil { slog.Error("bot initialization failed, check env tokens", "err", err); return }

	slog.Info("bot started")
	c.Start(ctx)
	slog.Info("bot stopped")
}

