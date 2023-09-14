package main

import (
	"context"
	"encoding/json"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"log"
	"net/url"
	"os"
)

type User struct {
	Name  string   `json:"name"`
	Sex   string   `json:"sex"`
	Age   int      `json:"age"`
	Hobby []string `json:"hobby"`
}

func main() {
	esIP := os.Getenv("es_ip")
	esPort := os.Getenv("es_port")

	ctx := context.Background()
	logger := log.New(os.Stdout, "es_review", log.LstdFlags)
	reqUrl := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", esIP, esPort),
	}

	// 初始化链接
	client, err := es.NewClient(
		es.SetURL(reqUrl.String()),
		es.SetSniff(false),
		es.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}

	q := es.NewMatchQuery("name", "huawei")

	// for print
	src, err := q.Source()
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(src)
	fmt.Println(string(data))

	// 数据搜索
	result, err := client.Search().Index("user").Query(q).Do(ctx)
	if err != nil {
		panic(err)
	}
	total := result.Hits.TotalHits.Value
	fmt.Println(total)

	for _, value := range result.Hits.Hits {

		var user User
		err := json.Unmarshal(value.Source, &user)
		if err != nil {
			log.Fatal("marshal error")
			continue
		}
		fmt.Println(user)

	}
}
