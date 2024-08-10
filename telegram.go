package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"strings"
)

const chunkSize = 4096

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

// NewTelegramBot return telegram bot instance
func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = viper.GetBool("TELEGRAM_DEBUG")
	return &TelegramBot{
		Bot: bot,
	}, nil
}

// SendMessage send message to telegram chat
func (t *TelegramBot) SendMessage(chatId int64, message string) error {
	for _, chunk := range chunkString(message, chunkSize) {
		msg := tgbotapi.NewMessage(chatId, chunk)
		_, err := t.Bot.Send(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return err
		}
	}
	return nil
}

// chunkString chunk message because of telegram limit
func chunkString(message string, chunkSize int) []string {
	var chunks []string
	for len(message) > chunkSize {
		// find index of last "\n" in string
		cutIndex := strings.LastIndex(message[:chunkSize], "\n")
		// if no "\n" then split exactly by chunk size
		if cutIndex == -1 {
			cutIndex = chunkSize
		}
		chunks = append(chunks, message[:cutIndex])
		message = message[cutIndex:]
		// remove "\n" from begin of the string
		message = strings.TrimLeft(message, "\n")
	}
	// add rest of the message
	if len(message) > 0 {
		chunks = append(chunks, message)
	}

	return chunks
}
