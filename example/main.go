package main

import (
	"fmt"
	"log"
	"time"

	elasticconvert "github.com/suiguo/elastic-convert"
	"github.com/suiguo/elastic-convert/example/client"
)

type TestStruct1 struct {
	Id  int    `json:"id"`
	Tag string `json:"tag"`
}

type TestStruct2 struct {
	Id  int       `json:"id"`
	Tag time.Time `json:"tag"`
}

func main() {
	cfg := &client.ElasticCfg{
		Host: []string{"http://127.0.0.1:9200"},
	}
	cli, err := client.GetInstanceElastic(cfg, true)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		err := cli.InsertNewRcord("index_1", i, &TestStruct1{
			Id:  i,
			Tag: fmt.Sprintf("idx:%d", i),
		})
		if err != nil {
			log.Println(err)
		}
	}

	for i := 0; i < 100; i++ {
		err := cli.InsertNewRcord("index_2", i, &TestStruct2{
			Id:  i,
			Tag: time.Now(),
		})
		if err != nil {
			log.Println(err)
		}
	}
	datas, err := cli.Search([]string{"index_1", "index_2"}, client.WithSort("id", client.Desc))
	if err != nil {
		panic(err)
	}
	for _, d := range datas {
		out, err := elasticconvert.Result(d, []string{"index_1", "index_2"}, []any{&TestStruct1{}, &TestStruct2{}})
		if err != nil {
			log.Println(err)
			continue
		}
		if out == nil {
			continue
		}
		for _, o := range out.Hits.Hits {
			log.Println(o)
		}
	}
}
