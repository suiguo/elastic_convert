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
	Id  string    `json:"id"`
	Tag time.Time `json:"tag"`
}

func main() {
	cfg := &client.ElasticCfg{
		Host: []string{"127.0.0.1:9200"},
	}
	cli, err := client.GetInstanceElastic(cfg, true)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		cli.InsertNewRcord("test1_index1", i, &TestStruct1{
			Id:  i,
			Tag: fmt.Sprintf("idx:%d", i),
		})
	}

	for i := 0; i < 100; i++ {
		cli.InsertNewRcord("test1_index2", i, &TestStruct2{
			Id:  fmt.Sprintf("idx:[%d]", i),
			Tag: time.Now(),
		})
	}
	datas, err := cli.Search([]string{"test1_index1", "test1_index2"}, client.WithSort("id", client.Desc))
	if err != nil {
		panic(err)
	}
	for _, d := range datas {
		out, err := elasticconvert.Result(d, []string{"test1_index1", "test1_index2"}, []any{&TestStruct1{}, &TestStruct2{}})
		if err != nil {
			log.Println(err)
			continue
		}
		if out == nil {
			continue
		}
		for _, o := range out.Hits.Hits {
			log.Panicln(o)
		}
	}
}
