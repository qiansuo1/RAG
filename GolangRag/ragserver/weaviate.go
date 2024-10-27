package ragserver

import (

	"context"
	"fmt"
	"os"

	    "github.com/weaviate/weaviate-go-client/v4/weaviate"
    "github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	 "github.com/weaviate/weaviate/entities/models"
)

func InitWeaviate(ctx context.Context) (*weaviate.Client, error) {
	// client,err := weaviate.NewClient(weaviate.Config{
	// 	Host:   "localhost:" + cmp.Or(os.Getenv("WVPORT"), "9035"),
	// 	Scheme: "http",
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("initializing weaviate: %w", err)
	// }

	// cls := &models.Class{
	// 	Class:      "Document",
	// 	Vectorizer: "none",
	// }
	// exists, err := client.Schema().ClassExistenceChecker().WithClassName(cls.Class).Do(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("weaviate error: %w", err)
	// }
	// if !exists {
	// 	err = client.Schema().ClassCreator().WithClass(cls).Do(ctx)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("weaviate error: %w", err)
	// 	}
	// }
	// return client, nil
	cfg := weaviate.Config{
        Host:        os.Getenv("WCD_HOSTNAME"),
        Scheme:     "https",
        AuthConfig:  auth.ApiKey{Value: os.Getenv("WCD_API_KEY")},
    }

    client, err := weaviate.NewClient(cfg)
    if err != nil {
        fmt.Println(err)
    }


	classObj := &models.Class{
        Class:      "Document",
        Vectorizer: "none",
        
    }
    err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
    if err != nil {
        panic(err)
    }

    // Check the connection
    ready, err := client.Misc().ReadyChecker().Do(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("%v", ready)

	return client, nil

}

func AddToWeaviate(client *weaviate.Client, pageNumber int32, sentenceChunk string, embedding []float32) error {
	// 创建对象
	object := &models.Object{
		Class: "Document",
		Properties: map[string]interface{}{
			"pageNumber":    pageNumber,
			"sentenceChunk": sentenceChunk,
		},
		Vector: embedding,
	}

	// 上传对象
	_, err := client.Data().Creator().
		WithClassName("Document").
		WithProperties(object.Properties).
		WithVector(object.Vector).
		Do(context.Background())

	return err
}