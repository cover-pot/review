package main

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"log"
	"os"
)

type Goods struct {
	Name   string   `json:"name"`
	Price  float64  `json:"price"`
	Images []string `json:"images"`
}

func main() {
	esIP := os.Getenv("es_ip")
	esPort := os.Getenv("es_port")

	ctx := context.Background()
	logger := log.New(os.Stdout, "es_review", log.LstdFlags)

	host := fmt.Sprintf("http://%s:%s", esIP, esPort)
	fmt.Println(host)
	// 初始化链接
	client, err := es.NewClient(
		es.SetURL(host),
		es.SetSniff(false),
		es.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}

	q := es.NewMatchQuery("name", "huawei")

	// to print
	src, err := q.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(src)
	fmt.Println(string(data))

	// 数据搜索
	result, err := client.Search().Index("goods").Query(q).Do(ctx)
	if err != nil {
		panic(err)
	}
	total := result.Hits.TotalHits.Value
	fmt.Println(total)

	for _, value := range result.Hits.Hits {

		var goods Goods
		err := json.Unmarshal(value.Source, &goods)
		if err != nil {
			log.Fatal("marshal error")
			continue
		}
		fmt.Println(goods)

	}
}
