package service

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"strconv"
)

type AliYunNlpHandler struct {
	Key     string
	Secret  string
	Domain  string
	Version string
}

type GetTsChEcomResponse struct {
	RequestId string `json:"RequestId"`
	Data      string `json:"Data"`
}

type GetTsChEcomResponseData struct {
	Result []struct {
		Score string `json:"score"`
		Flag  bool   `json:"flag"`
	} `json:"result"`
	Success bool `json:"success"`
}

func (n *AliYunNlpHandler) GetTsChEcom(OriginT, OriginQ string) (score float64, err error) {
	AccessKeyId := n.Key
	AccessKeySecret := n.Secret
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", AccessKeyId, AccessKeySecret)
	if err != nil {
		return 0, err
	}
	request := requests.NewCommonRequest()
	request.Domain = n.Domain
	request.Version = n.Version
	// 因为是RPC接口，因此需指定ApiName(Action)
	request.ApiName = "GetTsChEcom"
	request.QueryParams["ServiceCode"] = "alinlp"
	request.QueryParams["OriginT"] = OriginT
	request.QueryParams["OriginQ"] = OriginQ
	request.QueryParams["Type"] = "similarity"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return 0, err
	}
	var (
		temp     GetTsChEcomResponse
		tempData GetTsChEcomResponseData
	)
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &temp)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal([]byte(temp.Data), &tempData)
	if err != nil {
		return 0, err
	}
	if tempData.Result[0].Flag {
		score, err = strconv.ParseFloat(tempData.Result[0].Score, 64)
	}
	return score, err
}
