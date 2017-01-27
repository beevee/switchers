package telegram

import (
	"strconv"
	"strings"
	"time"
	"errors"
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

				b.GameProcessor.ExecuteCommand(message.Text, message.ID, "telegram_"+strconv.FormatInt(message.Chat.ID, 10))
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
func (b *Bot) SendMessage(playerID string, message string) {

	chatID, err := b.ParseUserID(playerID)
	if err != nil {
		b.Logger.Log("Cannot send message", "playerid", playerID, "message", message, "error", err)
		return
	}

	if err = b.telebot.SendMessage(telebot.Chat{ID: chatID}, message, nil); err != nil {
		b.Logger.Log("msg", "failed to send message to user", "playerid", playerID, "message", message, "error", err)
		return
	}

	b.Logger.Log("msg", "Telegram message sent", "chatid", chatID, "message", message)
}

func (b *Bot) ParseUserID(id string) (int64, error) {
	if !strings.HasPrefix(id, "telegram_") {
		return 0, errors.New(" non-Telegram user")
	}

	return strconv.ParseInt(id[9:], 10, 64);
}

func (b *Bot) ForwardMessage(playerID string, messageID int, messageOwnerID string) {

	chatID, err := b.ParseUserID(playerID)
	if err != nil {
		b.Logger.Log("Cannot forward message", "playerid", playerID, "messageOwner", messageOwnerID, "error", err)
		return
	}

	ownerID, err := b.ParseUserID(playerID)
	if err != nil {
		b.Logger.Log("Cannot forward message", "playerid", playerID, "messageOwner", messageOwnerID, "error", err)
		return
	}

	message := telebot.Message{ID: messageID, Sender: telebot.User{ID: int(ownerID)}}
	if err = b.telebot.ForwardMessage(telebot.Chat{ID: chatID}, message); err != nil {
		b.Logger.Log("msg", "failed to send message to user", "playerid", playerID, "message", message, "error", err)
		return
	}

	b.Logger.Log("msg", "Telegram message sent", "chatid", chatID, "message", message)
}
