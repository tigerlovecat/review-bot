package utils

import (
	"encoding/xml"
	"fmt"
)

type XMLData struct {
	ToUserName   CDATAText `xml:"ToUserName"`
	FromUserName CDATAText `xml:"FromUserName"`
	CreateTime   int64     `xml:"CreateTime"`
	MsgType      CDATAText `xml:"MsgType"`
	Event        CDATAText `xml:"Event"`
	EventKey     CDATAText `xml:"EventKey"`
	Ticket       CDATAText `xml:"Ticket"`
}

type CDATAText struct {
	Text string `xml:",cdata"`
}

type ReplyTextMessage struct {
	XMLName      xml.Name  `xml:"xml"`
	ToUserName   CDATAText `xml:"ToUserName"`
	FromUserName CDATAText `xml:"FromUserName"`
	CreateTime   int64     `xml:"CreateTime"`
	MsgType      CDATAText `xml:"MsgType"`
	Content      CDATAText `xml:"Content"`
}

type XMLPostData struct {
	ToUserName   CDATAText `xml:"ToUserName"`
	FromUserName CDATAText `xml:"FromUserName"`
	CreateTime   int64     `xml:"CreateTime"`
	MsgType      CDATAText `xml:"MsgType"`
	MediaId      CDATAText `xml:"MediaId"`
	ThumbMediaId CDATAText `xml:"ThumbMediaId"`
	MsgId        int64     `xml:"MsgId"`
	MsgDataId    string    `xml:"MsgDataId"`
	Idx          string    `xml:"Idx"`
	Event        CDATAText `xml:"Event"`
	EventKey     CDATAText `xml:"EventKey"`
	Ticket       CDATAText `xml:"Ticket"`
	Content      CDATAText `xml:"Content"`
}

func ParseXMLData(xmlContent []byte) (XMLPostData, error) {
	var data XMLPostData
	err := xml.Unmarshal(xmlContent, &data)
	if err != nil {
		// 处理错误
		fmt.Println("无法解析 XML 数据")
		return data, err
	}
	return data, nil
}
