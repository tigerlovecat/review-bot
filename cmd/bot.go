package cmd

import "review-bot/internal/service"

func TgBotSetup() {
	tgHandler := service.TgBotHandler{}
	tgHandler.BotUpdates()
}
