package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/likcoah/irterlan-rp-bot/internal/rng"
)


type Bot struct {
	client			*bot.Bot
	webhookURL		string
	webhookSecret	string
	port			string
	rd				*rng.Randomizer
}


var rollRe = regexp.MustCompile(`^/(\d*)[dдк](\d+)$`)


func New(token, webhookURL, webhookSecret, port string, rd *rng.Randomizer) (*Bot, error) {
	opts := []bot.Option{
		bot.WithWebhookSecretToken(webhookSecret),
		bot.WithDefaultHandler(func(_ context.Context, _ *bot.Bot, _ *models.Update) {}),
	}

	tg, err := bot.New(token, opts...)
	if err != nil { return nil, err }

	b := &Bot{
		client:			tg,
		webhookURL:		webhookURL,
		webhookSecret:	webhookSecret,
		port:			port,
		rd:				rd,
	}

	tg.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommandStartOnly, b.startHandler)
	tg.RegisterHandlerRegexp(bot.HandlerTypeMessageText, rollRe, b.rollHandler)

	return b, nil
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

func (b *Bot) startHandler(ctx context.Context, tg *bot.Bot, update *models.Update) {
	_, err := tg.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		MessageThreadID: update.Message.MessageThreadID,
		Text: "<i>Dice bot created by <a href='https://t.me/likcoah'>likcoah</a> for <a href='https://t.me/UltracruelRP'>Irterlan RP</a></i>\n\n" +
			"You can roll dice using the <code>/[number of dice rolls][d|д|к][number of sides on a die]</code> command\n\n" +
			"Examples:\n" +
			"<code>/d20</code> -> roll a 20-sided die\n" +
			"<code>/4к6</code> -> roll four 6-sided dice",
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
        	IsDisabled: bot.True(),
    	},
	})
	if err != nil { slog.Error("send message failed", "err", err) }
}

func (b *Bot) rollHandler(ctx context.Context, tg *bot.Bot, update *models.Update) {
	matches := rollRe.FindStringSubmatch(update.Message.Text)

	N, _ := strconv.Atoi(matches[2])
	if 2 > N || N > 100 {
		_, err := tg.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			MessageThreadID: update.Message.MessageThreadID,
			Text: "Sorry, but the number of sides must be between 2 and 100\n\n" +
				"Сорян, но количество сторон куба должно быть в диапазоне от 2 до 100",
		})
		if err != nil { slog.Error("send message failed", "err", err) }
		return
	}
	dicesCount := 1
	if matches[1] != "" {
		dicesCount, _ = strconv.Atoi(matches[1])
		if dicesCount < 1 || 100 < dicesCount { dicesCount = 1 }
	}
	result, history := 0, make([]string, dicesCount)

	for i := range dicesCount {
		n := b.rd.Roll(update.Message.From.ID, N)
		result += n
		history[i] = strconv.Itoa(n)
	}

	var text string
	if dicesCount == 1 {
		text = strconv.Itoa(result)
	} else {
		text = fmt.Sprintf("%s = <b>%d</b>", strings.Join(history, " + "), result)
	}

	_, err := tg.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		MessageThreadID: update.Message.MessageThreadID,
		Text: text,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil { slog.Error("send message failed", "err", err) }
}

