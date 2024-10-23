// package ragserver

// import (
// 	"cmp"
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/weaviate/weaviate-go-client/v4/weaviate"
// 	"github.com/weaviate/weaviate/entities/models"
// )

// func initWeaviate(ctx context.Context) (*weaviate.Client, error) {
// 	client,err := weaviate.NewClient(weaviate.Config{
// 		Host:   "localhost:" + cmp.Or(os.Getenv("WVPORT"), "9035"),
// 		Scheme: "http",
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("initializing weaviate: %w", err)
// 	}

// 	cls := &models.Class{
// 		Class:      "Document",
// 		Vectorizer: "none",
// 	}
// 	exists, err := client.Schema().ClassExistenceChecker().WithClassName(cls.Class).Do(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("weaviate error: %w", err)
// 	}
// 	if !exists {
// 		err = client.Schema().ClassCreator().WithClass(cls).Do(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("weaviate error: %w", err)
// 		}
// 	}
// 	return client, nil

// }