package telegram

import (
	"strconv"
	"strings"
	"time"

	"github.com/tucnak/telebot"
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

// Bot handles interactions with Telegram users
type Bot struct {
	TelegramToken    string
	TrumpCode        string
	PlayerRepository switchers.PlayerRepository
	GameProcessor    switchers.GameProcessor
	Logger           switchers.Logger
	telebot          *telebot.Bot
	tomb             tomb.Tomb
}

// Start initializes Telegram API connections
func (b *Bot) Start() error {
	var err error
	b.telebot, err = telebot.NewBot(b.TelegramToken)
	if err != nil {
		return err
	}

	messages := make(chan telebot.Message)
	b.telebot.Listen(messages, 1*time.Second)

	b.tomb.Go(func() error {
		for {
			select {
			case message := <-messages:
				b.Logger.Log("msg", "Telegram message received", "firstname", message.Sender.FirstName,
					"lastname", message.Sender.LastName, "username", message.Sender.Username,
					"chatid", message.Chat.ID, "message", message.Text)

				b.GameProcessor.ExecuteCommand(message.Text, "telegram_"+strconv.FormatInt(message.Chat.ID, 10))
			case <-b.tomb.Dying():
				return nil
			}
		}
	})

	return nil
}

// Stop gracefully stops Telegram API connections
func (b *Bot) Stop() error {
	b.tomb.Kill(nil)
	return b.tomb.Wait()
}

// SendMessage sends message to Telegram user
func (b *Bot) SendMessage(ID string, message string) {
	if !strings.HasPrefix(ID, "telegram_") {
		b.Logger.Log("msg", "cannot send messages to a non-Telegram user", "id", ID, "message", message)
		return
	}

	chatID, err := strconv.ParseInt(ID[9:], 10, 64)
	if err != nil {
		b.Logger.Log("msg", "unexpected non-numeric Telegram chat id", "id", ID, "message", message, "error", err)
		return
	}

	if err = b.telebot.SendMessage(telebot.Chat{ID: chatID}, message, nil); err != nil {
		b.Logger.Log("msg", "failed to send message to user", "id", ID, "message", message, "error", err)
		return
	}

	b.Logger.Log("msg", "Telegram message sent", "chatid", ID, "message", message)
}
