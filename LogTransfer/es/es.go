package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
)

var client *elastic.Client

type LogData struct {
	Data string `json:"log_data"`
}

func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return err
	}
	fmt.Println("connect to es succeed!")
	return
}

func SendToES(topic string, logData []byte) {
	d := LogData{
		Data: string(logData),
	}
	//	把数据发送给ES
	response, err := client.Index().Index(topic).BodyJson(d).Do(context.Background())
	if err != nil {
		fmt.Println("SendToES failed,err:", err)
		return
	}
	fmt.Println("SendToES ----", "Index:", response.Index, " Type:", response.Type, "ID:", response.Id, "Result:", response.Result)
}
