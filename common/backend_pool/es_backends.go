package backend_pool

import (
	"github.com/olivere/elastic"
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"context"
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

	done := make(chan error, 1)
	ctx := context.Background()
	go func() {
		exists, err := client.IndexExists(index).Do(ctx)
		if err != nil {
			// Handle error
			done <- err
			return
		}
		if !exists {
			err := fmt.Errorf("ES index %s not exist!", index)
			done <- err
			return
		}
		count := len(items)
		for i := 0; i < count; i++ {
			result, err := client.Index().
				Index(index).
				Type(esType).
				BodyJson(items[i]).
				Do(ctx)
			if err != nil {
				// Handle error
				done <- err
				continue
			}
			fmt.Printf("Insert to index %s, id %s\n", result.Index, result.Id)
		}
	}()

	select {
	//case <-time.After(time.Duration(t.callTimeout) * time.Millisecond):
	//	return fmt.Errorf("elasticsearch call timeout.")
	case err = <-done:
		if err != nil {
			err = fmt.Errorf("insert to elasticsearch failed, err %v.", err)
		}
		return err
	}
}




