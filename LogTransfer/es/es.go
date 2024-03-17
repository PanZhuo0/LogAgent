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

type Log struct {
	topic string
	msg   string
}

var ch chan Log

func Init(address string, chanSize int, nums int) (err error) {
	ch = make(chan Log, chanSize) //怎么大应该够用了
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		return err
	}
	fmt.Println("connect to es succeed!")
	//启动这个goroutine，让其一直获取数据
	for i := 0; i < nums; i++ {
		go sendToES()
	}
	return
}

func SendToESChan(topic string, msg []byte) {
	//存入chan
	ch <- Log{
		topic: topic,
		msg:   string(msg),
	}
}

func sendToES() {
	//从chan中读出
	for v := range ch {
		d := LogData{
			Data: v.msg,
		}
		//	把数据发送给ES
		response, err := client.Index().Index(v.topic).BodyJson(d).Do(context.Background())
		if err != nil {
			fmt.Println("SendToES failed,err:", err)
			return
		}
		fmt.Println("SendToES ----", "Index:", response.Index, " Type:", response.Type, "ID:", response.Id, "Result:", response.Result)
	}

}
