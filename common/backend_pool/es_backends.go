package backend_pool

import (
	"github.com/olivere/elastic"
	"github.com/open-falcon/falcon-plus/common/model"
	"context"
	"log"
	"errors"
)

//EsClient
type EsClient struct {
	cli  *elastic.Client
	name string
	callTimeout int
}

func CreateEsClient(connTimeout int, callTimeout int, urls []string) *EsClient {
	client, err := elastic.NewClient(elastic.SetURL(urls...))
	if err != nil {
		panic(err)
		return nil
	}
	return &EsClient{cli: client, name: "EsConnection", callTimeout: callTimeout }
}

func (t *EsClient) Send(index string, esType string, items []*model.EsItem) (err error) {
	client := t.cli

	ctx := context.Background()

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		log.Printf("Insert to elasticsearch failed, err %v.\n", err)
		return err
	}
	if !exists {
		err = errors.New("ES index "+index+" not exist")
		log.Print(err)
		return err
	}
	count := len(items)
	for i := 0; i < count; i++ {
		result, err := client.Index().
			Index(index).
			Type(esType).
			BodyJson(items[i]).
			Do(ctx)
		//TODO 修改为bulk
		if err != nil {
			log.Printf("Insert to elasticsearch failed, err %v.\n", err)
			continue
		}
		log.Printf("Insert to index %s, id %s\n", result.Index, result.Id)
	}
	log.Printf("Once Insert to es comple\n")
	return err;
}




