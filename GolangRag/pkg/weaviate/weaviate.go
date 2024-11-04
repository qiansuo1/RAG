package weaviate

import (
	"context"
	"errors"
	"fmt"
	"log"


	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)





var (
    ErrInvalidConfig     = errors.New("无效的配置")
    ErrConnectionFailed  = errors.New("连接失败")
    ErrInvalidVector     = errors.New("无效的向量数据")
    ErrSearchFailed      = errors.New("搜索失败")
    ErrSchemaCreateFailed = errors.New("创建Schema失败")
    ErrSchemaExist = errors.New("Schema已存在")
)



// Client Weaviate客户端封装
type Client struct {
    client *weaviate.Client
    ctx    context.Context
}

// Config Weaviate配置
type Config struct {
    WeaviateHost string
    WeaviateKey string
}

// NewClient 创建新的Weaviate客户端
func NewClient(Host,Key string) (*Client, error) {

    var cfg = Config{
        WeaviateHost:Host,
        WeaviateKey:Key,
    }

    config := weaviate.Config{
        Host:   cfg.WeaviateHost,
        Scheme: "https",
        AuthConfig:  auth.ApiKey{Value: cfg.WeaviateKey},
    }

    client, err := weaviate.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("creating Weaviate client failed: %w", err)
    }

   var  myWeaviateClient = Client{
        client: client,
        ctx:    context.Background(),
    }
    // Check the connection
    ready, err := client.Misc().ReadyChecker().Do(myWeaviateClient.ctx)
    if err != nil {
        return nil, fmt.Errorf("checking Weaviate connection failed: %w", err)
    }
    log.Printf("Weaviate ready state: %v\n", ready)
    return &myWeaviateClient, nil
}

// CreateSchema 创建Schema
func (c *Client) CreateSchema(className string) error {

    	   // 检查类是否存在
	   exists, err := c.client.Schema().ClassExistenceChecker().
	   WithClassName(className).
	   Do(c.ctx)
   if err != nil {
	   return err
   }
       // 只在类不存在时创建
	 if !exists {
      
    classObj := &models.Class{
        Class: className,
        Properties: []*models.Property{
            {
                Name:     "pageNumber",
                DataType: []string{"int"},
            },
            {
                Name:     "content",
                DataType: []string{"text"},
            },
        },
        Vectorizer: "none", // 使用自定义向量
    }
    
    err = c.client.Schema().ClassCreator().WithClass(classObj).Do(c.ctx)
    if err != nil {
        return fmt.Errorf("创建Schema失败: %w", err)
    }

    return nil

    }

    return ErrSchemaExist
}



const DefaultClassName = "Document"

// 使用已有的CreateSchema函数
func (c *Client) ensureClassExists() error {
    err := c.CreateSchema(DefaultClassName)
    if err == ErrSchemaExist {
        // class已存在，这是正常情况
        return nil
    }
    if err != nil {
        return err
    }
    return nil
}


func (c *Client) AddEmbedding(pageNum int64, text string, vector []float32) error {
  
    if err := c.ensureClassExists(); err != nil {
        return err
    }

    className := "Document"
    properties := map[string]interface{}{
        "pageNumber": pageNum,
        "content":    text,
    }

    _, err := c.client.Data().Creator().
        WithClassName(className).
        WithProperties(properties).
        WithVector(vector).
        Do(c.ctx)

    if err != nil {
        return fmt.Errorf("添加文档失败: %w", err)
    }

    return nil
}

// SearchSimilar 搜索相似文档
func (c *Client) SearchSimilar(queryText string, limit int) ([]SearchResult, error) {
    className := "Document"
    // 构建查询
    nearText := c.client.GraphQL().NearTextArgBuilder().
        WithConcepts([]string{queryText})

    fields := []graphql.Field{
        {Name: "pageNumber"},
        {Name: "content"},
        {Name: "_additional { certainty }"},
    }

    // 执行查询
    result, err := c.client.GraphQL().Get().
        WithClassName(className).
        WithFields(fields...).
        WithNearText(nearText).
        WithLimit(limit).
        Do(c.ctx)

    if err != nil {
        return nil, fmt.Errorf("搜索文档失败: %w", err)
    }

    return c.parseSearchResults(result)
}

// SearchResult 搜索结果结构
type SearchResult struct {
    PageNumber int64   `json:"pageNumber"`
    Content    string  `json:"content"`
    Certainty  float64 `json:"certainty"`
}

// parseSearchResults 解析搜索结果
func (c *Client) parseSearchResults(result *models.GraphQLResponse) ([]SearchResult, error) {
    if result == nil || result.Data == nil {
        return nil, fmt.Errorf("无效的搜索结果")
    }
    className := "Document"
    var results []SearchResult

    // 解析GraphQL响应
    if data, ok := result.Data["Get"].(map[string]interface{}); ok {
        if documents, ok := data[className].([]interface{}); ok {
            for _, doc := range documents {
                if document, ok := doc.(map[string]interface{}); ok {
                    result := SearchResult{}
                    
                    // 解析页码
                    if pageNum, ok := document["pageNumber"].(float64); ok {
                        result.PageNumber = int64(pageNum)
                    }
                    
                    // 解析内容
                    if content, ok := document["content"].(string); ok {
                        result.Content = content
                    }
                    
                    // 解析相似度
                    if additional, ok := document["_additional"].(map[string]interface{}); ok {
                        if certainty, ok := additional["certainty"].(float64); ok {
                            result.Certainty = certainty
                        }
                    }
                    
                    results = append(results, result)
                }
            }
        }
    }

    return results, nil
}



// DeleteSchema 删除Schema
func (c *Client) DeleteSchema(className string) error {
    
    err := c.client.Schema().ClassDeleter().WithClassName(className).Do(c.ctx)
    if err != nil {
        return fmt.Errorf("删除Schema失败: %w", err)
    }

    return nil
}

// func initWeaviate(ctx context.Context) (*weaviate.Client, error) {


//     }

//     client, err := weaviate.NewClient(cfg)
//     if err != nil {
//         return nil, fmt.Errorf("init weaviate fail: %w", err)
//     }
// 	   // 检查类是否存在
// 	   exists, err := client.Schema().ClassExistenceChecker().
// 	   WithClassName("Document").
// 	   Do(ctx)
//    if err != nil {
// 	   return nil,  err
//    }
//        // 只在类不存在时创建
// 	   if !exists {
//         classObj := &models.Class{
//             Class:      "Document",
//             Vectorizer: "none",
//         }
//         if err := client.Schema().ClassCreator().WithClass(classObj).Do(ctx); err != nil {
//             return nil, fmt.Errorf("create class fail: %w", err)
//         }
//     }


//     // Check the connection
//     ready, err := client.Misc().ReadyChecker().Do(context.Background())
//     if err != nil {
//         return nil, err
//     }
// 	fmt.Printf("Weaviate ready state: %v\n", ready)
// 	return client, nil

// }

// type ragServer struct {
// 	ctx      context.Context
// 	wvClient *weaviate.Client
// }

// func InitRagServer() (*ragServer, error) {
// 	ctx := context.Background()	
// 	wvClient, err := initWeaviate(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	server := &ragServer{ctx: ctx, wvClient: wvClient}

// 	router := gin.Default()
// 	router.POST("/add", server.AddToWeaviate)
// 	router.GET("/query", server.QueryWeaviate)
// 	router.Run(":8080")	
// 	return server, nil	
// }



// type AddReq struct{
// 	PageNumber int64 `json:"page_number"`
// 	SentenceChunk string `json:"sentence_chunk"`
// 	Embedding []float32 `json:"embedding"`
// }

// func (s *ragServer) AddToWeaviate(c *gin.Context) {
// 	// 创建对象
// 	var req AddReq
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
	

	
// 	object := &models.Object{
	
// 		Properties: map[string]interface{}{
// 			"pageNumber":    req.PageNumber,
// 			"sentenceChunk": req.SentenceChunk,
// 		},
// 		Vector: req.Embedding,
// 	}

// 	// 上传对象
// 	_, err = s.wvClient.Data().Creator().
// 		WithClassName("Document").
// 		WithProperties(object.Properties).
// 		WithVector(object.Vector).
// 		Do(context.Background())

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Object added successfully"})
// }



// type QueryReq struct {
// 	ClassName string `json:"className"`
// 	Fields []string `json:"fields"`
// 	Limit    int      `json:"limit"`
//     Include  struct {
//         Vector bool `json:"vector"`
//     } 				`json:"include"`
// }	
// //查询
// func (s *ragServer) QueryWeaviate(c *gin.Context){
// 	var req QueryReq
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fields := make([]graphql.Field, len(req.Fields))
//     for i, fieldName := range req.Fields {
//         fields[i] = graphql.Field{
//             Name: fieldName,
//         }
//     }
// 	    // 根据需要添加向量字段
// 		if req.Include.Vector {
// 			fields = append(fields, graphql.Field{
// 				Name: "_additional",
// 				Fields: []graphql.Field{
// 					{
// 						Name: "vector",
// 					},
// 				},
// 			})
// 		}
//     // 设置默认限制
//     limit := 10
//     if req.Limit > 0 {
//         limit = req.Limit
//     }

// 		response, err := s.wvClient.GraphQL().Get().
//         WithClassName(req.ClassName).
//         WithFields(fields...).
//         WithLimit(limit).
//         Do(context.Background())

//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // 处理响应
//     fmt.Printf("查询结果: %+v\n", response)
 


// }	