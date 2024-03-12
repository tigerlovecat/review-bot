package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"review-bot/config"
	"review-bot/internal/model"
)

type TgBotHandler struct {
	Client *tgbotapi.BotAPI
}

func (t *TgBotHandler) InitClient() {
	bot, err := tgbotapi.NewBotAPI(config.TelegramConfig.Token)
	if err != nil {
		fmt.Println(" - -- InitClient - -- ")
		fmt.Printf(" - -- err:{%v} - -- ", err.Error())
		fmt.Println(" - -- InitClient - -- ")
		t.Client = nil
		return
	}
	t.Client = bot
	return
}

func (t *TgBotHandler) BotUpdates() {
	t.InitClient()
	if t.Client == nil {
		fmt.Println("InitClient err .")
		return
	}
	// set debug
	t.Client.Debug = true
	botConfig := tgbotapi.NewUpdate(0)
	botConfig.Timeout = 0
	// get updates
	updateList := t.Client.GetUpdatesChan(botConfig)
	// reply
	for update := range updateList {
		go handleUpdate(t.Client, update)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	sourceFile := "./config/knowledge_base.json"
	qaHandler := QaHandler{
		Source: sourceFile,
	}
	// load knowledge base 加载知识库
	kb, err := qaHandler.LoadKnowledgeBase()
	if err != nil {
		fmt.Printf("load knowledge base err: {%v} ", err)
		return
	}
	// find answer 查找答案
	answer := qaHandler.FindAnswer(text, kb)
	replyMsg := tgbotapi.NewMessage(chatID, answer)
	_, _ = bot.Send(replyMsg)
	// add records
	chatResponse := update.Message.Chat
	fmt.Printf("chat data: {%+v} \r\n", chatResponse)
	fmt.Printf("chat UserName: {%v} \r\n", &chatResponse.LastName)
	recordModel := model.CustomerRecordsModel{
		Message:     text,
		Username:    update.Message.Chat.LastName,
		Reply:       answer,
		MessageTime: int64(update.Message.Date),
	}
	err = recordModel.AddRecord()
	if err != nil {
		fmt.Printf("add record err: {%v} \r\n", err.Error())
	} else {
		fmt.Println("add record success.")
	}
	return
}
