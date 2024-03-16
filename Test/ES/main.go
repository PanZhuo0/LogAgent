package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// Es insert data demo
type student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func main() {
	//	初始化连接
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to es succeed!")
	s1 := student{Name: "张三", Age: 22, Married: false}
	//链式操作
	put1, err := client.Index().Index("student").BodyJson(s1).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed student %s to index:%s,type:%s\n", put1.Id, put1.Index, put1.Type)
}
