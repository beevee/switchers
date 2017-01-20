package telegram

import (
	"time"

	"github.com/tucnak/telebot"
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

// Bot handles interactions with Telegram users
type Bot struct {
	TelegramToken string
	Logger        switchers.Logger
	telebot       *telebot.Bot
	tomb          tomb.Tomb
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
				b.handleMessage(message)
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

func (b *Bot) handleMessage(message telebot.Message) {
	b.Logger.Log("msg", "message received", "firstname", message.Sender.FirstName,
		"lastname", message.Sender.LastName, "username", message.Sender.Username,
		"chatid", message.Chat.ID, "command", message.Text)
}
