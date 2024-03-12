package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"review-bot/internal/model"
	"review-bot/internal/service"
	"review-bot/pkg/app"
	"review-bot/pkg/app/response"
	"review-bot/pkg/app/utils"
)

// QaTest question reply
func QaTest(c *gin.Context) {
	sourceFile := "./config/knowledge_base.json"
	qaHandler := service.QaHandler{
		Source: sourceFile,
	}
	// load knowledge base 加载知识库
	kb, err := qaHandler.LoadKnowledgeBase()
	if err != nil {
		fmt.Printf("load knowledge base err: {%v} ", err)
		return
	}
	// get question 客户数据
	question := c.Query("question")
	// find answer 查找答案
	answer := qaHandler.FindAnswer(question, kb)
	fmt.Printf("answer: {%v} \r\n", answer)

	app.XOK(c, answer, "success")
	return
}

// QaListTest question reply list
func QaListTest(c *gin.Context) {
	customerModel := model.CustomerRecordsModel{}
	list, err := customerModel.List()
	if err != nil {
		app.XError(c, "ParamsHandlerErr")
		return
	}
	var result []response.QaListTestResponse
	for _, temp := range list {
		result = append(result, response.QaListTestResponse{
			Time:     utils.TimeInt64ToString(temp.MessageTime),
			UserName: temp.Username,
			Message:  temp.Message,
			Reply:    temp.Reply,
		})
	}

	app.XOK(c, result, "success")
	return
}
