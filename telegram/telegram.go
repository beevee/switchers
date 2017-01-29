package telegram

import (
	"errors"
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
	outboxText       chan outgoingText
	outboxForward    chan outgoingForward
	telebot          *telebot.Bot
	tomb             tomb.Tomb
}

type outgoingText struct {
	chat telebot.Chat
	text string
}

type outgoingForward struct {
	chat    telebot.Chat
	message telebot.Message
}

// Start initializes Telegram API connections
func (b *Bot) Start() error {
	var err error
	b.telebot, err = telebot.NewBot(b.TelegramToken)
	if err != nil {
		return err
	}

	inbox := make(chan telebot.Message)
	b.telebot.Listen(inbox, 1*time.Second)
	b.tomb.Go(func() error {
		for {
			select {
			case message := <-inbox:
				b.Logger.Log("msg", "Telegram message received", "firstname", message.Sender.FirstName,
					"lastname", message.Sender.LastName, "username", message.Sender.Username,
					"chatid", message.Chat.ID, "message", message.Text)

				b.GameProcessor.ExecuteCommand(strconv.Itoa(message.ID), message.Text, "telegram_"+strconv.FormatInt(message.Chat.ID, 10))
			case <-b.tomb.Dying():
				b.Logger.Log("msg", "aborted Telegram message receiving goroutine")
				return nil
			}
		}
	})

	b.outboxText = make(chan outgoingText, 1000)
	b.outboxForward = make(chan outgoingForward, 1000)
	b.tomb.Go(func() error {
		for {
			select {
			case message := <-b.outboxText:
				if err = b.telebot.SendMessage(message.chat, message.text, nil); err != nil {
					b.Logger.Log("msg", "failed to send Telegram message", "chatid", message.chat.ID, "message", message.text, "error", err)
				} else {
					b.Logger.Log("msg", "Telegram message sent", "chatid", message.chat.ID, "message", message.text)
				}
				time.Sleep(200 * time.Millisecond)
			case message := <-b.outboxForward:
				if err = b.telebot.ForwardMessage(message.chat, message.message); err != nil {
					b.Logger.Log("msg", "failed to forward Telegram message", "chatid", message.chat.ID, "messageid", message.message.ID, "ownerid", message.message.Sender.ID, "error", err)
				} else {
					b.Logger.Log("msg", "Telegram message forwarded", "chatid", message.chat.ID, "messageid", message.message.ID, "ownerid", message.message.Sender.ID)
				}
				time.Sleep(200 * time.Millisecond)
			case <-b.tomb.Dying():
				b.Logger.Log("msg", "aborted Telegram message sending goroutine")
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
	chatID, err := b.parseUserID(playerID)
	if err != nil {
		b.Logger.Log("msg", "cannot parse player id", "playerid", playerID, "message", message, "error", err)
		return
	}

	b.outboxText <- outgoingText{
		chat: telebot.Chat{ID: chatID},
		text: message,
	}
}

// ForwardMessage forwards message to Telegram user
func (b *Bot) ForwardMessage(playerID string, messageText string, messageID string, messageOwnerID string) {
	chatID, err := b.parseUserID(playerID)
	if err != nil {
		b.Logger.Log("msg", "cannot parse player id", "playerid", playerID, "messageowner", messageOwnerID, "error", err)
		return
	}

	ownerID, err := b.parseUserID(messageOwnerID)
	if err != nil {
		b.Logger.Log("msg", "cannot parse message owner id", "playerid", playerID, "messageowner", messageOwnerID, "error", err)
		return
	}

	msgID, err := strconv.Atoi(messageID)
	if err != nil {
		b.Logger.Log("msg", "cannot parse original message id", "playerid", playerID, "messageowner", messageOwnerID, "error", err)
		return
	}

	b.outboxForward <- outgoingForward{
		chat:    telebot.Chat{ID: chatID},
		message: telebot.Message{ID: msgID, Sender: telebot.User{ID: int(ownerID)}},
	}
}

func (b *Bot) parseUserID(id string) (int64, error) {
	if !strings.HasPrefix(id, "telegram_") {
		return 0, errors.New("non-Telegram user")
	}

	return strconv.ParseInt(id[9:], 10, 64)
}
