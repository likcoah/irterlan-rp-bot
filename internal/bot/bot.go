package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)


type Bot struct {
	client *bot.Bot
}


func New(token string) (*Bot, error) {
	opts := []bot.Option{
	}

	b, err := bot.New(token, opts...)
	if err != nil { return nil, err }

	b.RegisterHandler(bot.HandlerTypeMessageText, "/ping", bot.MatchTypeExact, ping)

	return &Bot{client: b}, nil
}

func (b *Bot) Start(ctx context.Context){
	b.client.Start(ctx)
}

func ping(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "pong",
	})
}

