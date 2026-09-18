package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)


type Bot struct {
	client			*bot.Bot
	webhookURL		string
	webhookSecret	string
	port			string
}


func New(token, webhookURL, webhookSecret, port string) (*Bot, error) {
	opts := []bot.Option{
		bot.WithWebhookSecretToken(webhookSecret),
	}

	b, err := bot.New(token, opts...)
	if err != nil { return nil, err }

	b.RegisterHandler(bot.HandlerTypeMessageText, "/ping", bot.MatchTypeExact, ping)

	return &Bot{
		client:			b,
		webhookURL:		webhookURL,
		webhookSecret:	webhookSecret,
		port:			port,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	_, err := b.client.SetWebhook(ctx, &bot.SetWebhookParams{
		URL:			b.webhookURL,
		SecretToken:	b.webhookSecret,
	})
	if err != nil { return fmt.Errorf("set webhook failed: %w", err) }

	go b.client.StartWebhook(ctx)

	server := &http.Server{
		Addr:				":" + b.port,
		Handler:			b.client.WebhookHandler(),
		ReadHeaderTimeout:	5*time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("webhook HTTP server shutdown failed", "err", err)
		}
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("webhook HTTP server failed: %w", err)
	}

	return nil
}

func ping(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "pong",
	})
}

