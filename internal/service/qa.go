package service

import (
	"encoding/json"
	"github.com/agnivade/levenshtein"
	"io/ioutil"
	"review-bot/config"
	"strings"
	"unicode"
)

type QaHandler struct {
	Source string
}

type KnowledgeBase struct {
	Questions []struct {
		Variants []string `json:"variants"`
		Answer   string   `json:"answer"`
	} `json:"questions"`
}

func (q *QaHandler) LoadKnowledgeBase() (KnowledgeBase, error) {
	var kb KnowledgeBase
	// load JSON file
	data, err := ioutil.ReadFile(q.Source)
	if err != nil {
		return kb, err
	}
	// Unmarshal json
	err = json.Unmarshal(data, &kb)
	if err != nil {
		return kb, err
	}
	return kb, nil
}

func (q *QaHandler) FindAnswer(question string, kb KnowledgeBase) string {
	// string to lower 将问题转换为小写并去除标点符号
	normalizedQuestion := strings.ToLower(question)
	normalizedQuestion = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, normalizedQuestion)

	var similarQuestions []string
	// 遍历知识库中的问题并查找匹配的变体
	for _, temp := range kb.Questions {
		for _, variant := range temp.Variants {
			normalizedVariant := strings.ToLower(variant)
			normalizedVariant = strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					return r
				}
				return -1
			}, normalizedVariant)
			// if 完全匹配
			if strings.Contains(normalizedQuestion, normalizedVariant) && normalizedVariant != "" {
				return temp.Answer
			}
			// if 不完全匹配
			// 计算动态阈值，这里使用了一个简单的公式，可以根据实际情况调整
			// similarQuestions = append(similarQuestions,  getSimilarQuestions(normalizedQuestion, normalizedVariant))
			if getNlpSimilarQuestions(normalizedQuestion, normalizedVariant, variant) != "" {
				similarQuestions = append(similarQuestions, getNlpSimilarQuestions(normalizedQuestion, normalizedVariant, variant))
			}
		}
	}
	if len(similarQuestions) <= 0 {
		return "Sorry, we couldn't find the relevant answer."
	}
	return "You can still ask me that. eg: \r\n" + strings.Join(similarQuestions, " \r\n")
}

// getSimilarQuestions 方案一： 直接计算相似度
func getSimilarQuestions(normalizedQuestion, normalizedVariant string) string {
	threshold := int(float64(len(normalizedQuestion))*0.3) + 1
	distance := levenshtein.ComputeDistance(normalizedQuestion, normalizedVariant)
	if distance <= threshold {
		return normalizedVariant
	}
	return normalizedVariant
}

// getNlpSimilarQuestions 方案二： 【阿里云】文本相似度（电商）
func getNlpSimilarQuestions(normalizedQuestion, normalizedVariant, variant string) string {
	nlpHandler := AliYunNlpHandler{
		Key:     config.NlpConfig.Key,
		Secret:  config.NlpConfig.Secret,
		Domain:  "alinlp.cn-hangzhou.aliyuncs.com",
		Version: "2020-06-29",
	}
	distance, err := nlpHandler.GetTsChEcom(normalizedQuestion, normalizedVariant)
	if err != nil {
		return ""
	}
	if distance <= 0.5 {
		return ""
	}
	return variant
}
