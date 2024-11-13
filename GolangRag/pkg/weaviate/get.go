package weaviate

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

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
    log.Printf("执行查询: className=%s, fields=%+v", className, fields)
    // 执行查询
    result, err := c.client.GraphQL().Get().
        WithClassName(className).
        WithFields(fields...).
        WithNearText(nearText).
        WithLimit(limit).
        Do(c.ctx)

    if err != nil {
        log.Printf("查询错误: %v", err)
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
        log.Printf("结果为空: result=%v", result)
        return nil, fmt.Errorf("无效的搜索结果")
    }
        // 添加错误检查
        if len(result.Errors) > 0 {
            var errMsgs []string
            for _, err := range result.Errors {
                if err != nil {
                    errMsgs = append(errMsgs, fmt.Sprintf(
                        "消息: %s, 路径: %s", 
                        err.Message,
                        strings.Join(err.Path, "."),
                    ))
                    }
        }
        return nil, fmt.Errorf("GraphQL错误: %s", strings.Join(errMsgs, "; "))
    }
    
    className := "Document"
    data, ok := result.Data["Get"].(map[string]interface{})
    if !ok {
        log.Printf("无法解析Get字段: %+v", result.Data)
        return nil, fmt.Errorf("无效的响应格式: 缺少 Get 字段")
    }

    documents, ok := data[className].([]interface{})
    if !ok {
        log.Printf("无法解析%s字段: %+v", className, data)
        return nil, fmt.Errorf("无效的响应格式: 缺少 %s 字段", className)
    }
    log.Printf("找到 %d 条文档", len(documents))
    var results []SearchResult

    // 解析GraphQL响应
    for i, doc := range documents  {
        
        document, ok := doc.(map[string]interface{})
        if !ok {
            log.Printf("无法解析文档 #%d: %+v", i, doc)
            continue // 跳过无效的文档
        }

        result := SearchResult{}
        
        // 解析页码
        if pageNum, ok := document["pageNumber"].(float64); ok {
            result.PageNumber = int64(pageNum)
        }else {
            log.Printf("文档 #%d 无法解析页码: %v", i, document["pageNumber"])
        }
        
        // 解析内容
        if content, ok := document["sentenceChunk"].(string); ok {
            result.Content = content
        } else {
            log.Printf("文档 #%d 无法解析内容: %v", i, document["sentenceChunk"])
            continue // 如果没有内容，跳过这条记录
        }
    
        
        results = append(results, result)
    }
                
    log.Printf("成功解析 %d 条结果", len(results))

    return results, nil
}



type ListResult struct {
    Id  string `json:"id"`
    PageNumber int64   `json:"pageNumber"`
    SentenceChunk    string  `json:"sentenceChunk"`

}

// ListAll 获取所有文档
func (c *Client) ListAll(limit int) ([]ListResult, error) {
    className := "Document"
    
    // 构建查询字段
    fields := []graphql.Field{
        {Name: "pageNumber"},
        {Name: "sentenceChunk"},
        {
            Name: "_additional",
            Fields: []graphql.Field{
                {Name: "id"},   
            },
        },
    }

    // 执行查询，不使用任何过滤条件
    result, err := c.client.GraphQL().Get().
        WithClassName(className).
        WithFields(fields...).
        WithLimit(limit).
        Do(c.ctx)

    if err != nil {
        log.Printf("查询错误: %v", err)  // 添加日志以便调试
        return nil, fmt.Errorf("获取文档列表失败: %w", err)
    }

    return c.parseListResults(result)
}

func (c *Client) parseListResults(result *models.GraphQLResponse) ([]ListResult, error) {
    if result == nil || result.Data == nil {
        log.Printf("结果为空: result=%v", result)
        return nil, fmt.Errorf("列表搜索为空,表中无数据")
    }
    if len(result.Errors) > 0 {
        var errMsgs []string
        for _, err := range result.Errors {
            if err != nil {
                errMsgs = append(errMsgs, fmt.Sprintf(
                    "消息: %s, 路径: %s", 
                    err.Message,
                    strings.Join(err.Path, "."),
                ))
                }
    }
    return nil, fmt.Errorf("GraphQL错误: %s", strings.Join(errMsgs, "; "))
}


    className := "Document"
    data, ok := result.Data["Get"].(map[string]interface{})
    if !ok {
        log.Printf("无法解析Get字段: %+v", result.Data)
        return nil, fmt.Errorf("无效的响应格式: 缺少 Get 字段")
    }


    documents, ok := data[className].([]interface{})
    if !ok {
        log.Printf("无法解析%s字段: %+v", className, data)
        return nil, fmt.Errorf("无效的响应格式: 缺少 %s 字段", className)
    }

 var results []ListResult   
 for i,doc := range documents {   
    document, ok := doc.(map[string]interface{})
        if !ok {
            log.Printf("无法解析文档 #%d: %+v", i, doc)
            continue // 跳过无效的文档
        }

        result := ListResult{}
        
        // 解析页码
        if pageNum, ok := document["pageNumber"].(float64); ok {
            result.PageNumber = int64(pageNum)
        }else {
            log.Printf("文档 #%d 无法解析页码: %v", i, document["pageNumber"])
        }
        
        // 解析内容
        if sentenceChunk, ok := document["sentenceChunk"].(string); ok {
            result.SentenceChunk = sentenceChunk
        } else {
            log.Printf("文档 #%d 无法解析内容: %v", i, document["sentenceChunk"])
            continue // 如果没有内容，跳过这条记录
        }

    if id, ok := document["_additional"].(map[string]interface{})["id"].(string); ok {
        result.Id = id
    } else {
        log.Printf("文档 #%d 无法解析ID: %v", i, document["_additional"])
    }   

    results = append(results, result)   
    }
    return results, nil 

}


//相似文本搜索
type NearTextResult struct {
    PageNumber int64   `json:"pageNumber"`
    SentenceChunk    string  `json:"sentenceChunk"`
     Certainty  float64   `json:"certainty,omitempty"`
    Distance   float64   `json:"distance,omitempty"`
}   

func (c *Client) GetNearText(inputVectorOfText []float32, limit int) ([]NearTextResult, error) {
    
    exists, err := c.client.Schema().ClassExistenceChecker().
    WithClassName("Document").
        Do(c.ctx)
    if err != nil {
        return nil, fmt.Errorf("检查集合是否存在失败: %w", err)
    }   
    if !exists {
        return nil, fmt.Errorf("集合 Document 不存在")
    }
// //如果没有配置向量化，则使用grpc调用python端的向量化模型
// hasVectorizer, err := c.checkVectorizer("Document")
// if err != nil {
//     return nil, fmt.Errorf("检查向量化器失败: %w", err)
// }
// if !hasVectorizer {
//     log.Printf("集合 Document 未配置向量化器")

//    //使用grpc调用python端的向量化模型

   


// }
   
nearText := c.client.GraphQL().NearVectorArgBuilder().
WithVector(inputVectorOfText)

    fields := []graphql.Field{
        {Name: "pageNumber"},
        {Name: "sentenceChunk"},
        {
            Name: "_additional",
            Fields: []graphql.Field{
                {Name: "certainty"},
                {Name: "distance"},
            },
        },
    }   
  
    response, err := c.client.GraphQL().Get().
        WithClassName("Document").
        WithFields(fields...).
        WithNearVector(nearText).
        WithLimit(limit).
        Do(c.ctx)   
    if err != nil { 
        log.Printf("查询错误: %v", err)
        return nil, fmt.Errorf("获取文档列表失败: %w", err)
    }
    if len(response.Errors) > 0 {
        var errMsgs []string
        for _, graphqlErr := range response.Errors {
            errMsgs = append(errMsgs, fmt.Sprintf(
                "GraphQL错误: 消息=%s, 路径=%v",
                graphqlErr.Message,
                graphqlErr.Path,
               
            ))
        }
        log.Printf("GraphQL查询失败: %s", strings.Join(errMsgs, "; "))
        return nil, fmt.Errorf("GraphQL查询失败: %s", strings.Join(errMsgs, "; "))
    }
    log.Printf("Weaviate响应: %+v", response)
    return c.parseNearTextResults(response)
}
// parseSearchResults 解析结果时处理两种相似度
func (c *Client) parseNearTextResults(result *models.GraphQLResponse) ([]NearTextResult, error) {
        if result == nil || result.Data == nil {
        log.Printf("结果为空: result=%v", result)
        return nil, fmt.Errorf("列表搜索为空,表中无数据")
    }
    // 错误检查
    if len(result.Errors) > 0 {
        var errMsgs []string
        for _, err := range result.Errors {
            if err != nil {
                errMsgs = append(errMsgs, fmt.Sprintf(
                    "消息: %s, 路径: %s", 
                    err.Message,
                    strings.Join(err.Path, "."),
                ))
            }
        }
        return nil, fmt.Errorf("GraphQL错误: %s", strings.Join(errMsgs, "; "))
    }


    className := "Document"
    data, ok := result.Data["Get"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("无效的响应格式: 缺少 Get 字段")
    }

    documents, ok := data[className].([]interface{})
    if !ok {
        return nil, fmt.Errorf("无效的响应格式: 缺少 %s 字段", className)
    }

    var results []NearTextResult
    for i, doc := range documents {
        document, ok := doc.(map[string]interface{})
        if !ok {
            log.Printf("无法解析文档 #%d: %+v", i, doc)
            continue
        }

        result := NearTextResult{}
        
        // 解析页码
        if pageNum, ok := document["pageNumber"].(float64); ok {
            result.PageNumber = int64(pageNum)
        } else {
            log.Printf("文档 #%d 无法解析页码: %v", i, document["pageNumber"])
        }
        
        // 解析内容
        if sentenceChunk, ok := document["sentenceChunk"].(string); ok {
            result.SentenceChunk = sentenceChunk
        } else {
            log.Printf("文档 #%d 无法解析内容: %v", i, document["sentenceChunk"])
            continue
        }

        // 解析额外信息
        if additional, ok := document["_additional"].(map[string]interface{}); ok {
                // 解析 certainty
            if certainty, ok := additional["certainty"].(float64); ok {
                result.Certainty = certainty
            }
            
            // 解析 distance
            if distance, ok := additional["distance"].(float64); ok {
                result.Distance = distance
            }
        }

        results = append(results, result)
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("未找到任何相似内容")
    }

    // 按相似度排序（默认使用 certainty）
    sort.Slice(results, func(i, j int) bool {
        return results[i].Certainty > results[j].Certainty
    })

    return results, nil
}


// checkVectorizer 检查集合是否配置了向量化器
func (c *Client) checkVectorizer(className string) (bool, error) {
    schema, err := c.client.Schema().ClassGetter().
    WithClassName(className).
    Do(c.ctx)
if err != nil {
    return false, fmt.Errorf("获取集合配置失败: %w", err)
}

// 打印完整配置以便调试
log.Printf("集合配置: %+v", schema)

// 检查是否配置了向量化器
if schema.Vectorizer == "" {
    return false, nil
}

// 类型断言检查 ModuleConfig
if schema.ModuleConfig == nil {
    return false, nil
}

moduleConfig, ok := schema.ModuleConfig.(map[string]interface{})
if !ok {
    return false, fmt.Errorf("无效的 ModuleConfig 类型: %T", schema.ModuleConfig)
}

// 检查是否存在 text2vec-openai 配置
_, hasOpenAI := moduleConfig["text2vec-openai"]

log.Printf("集合 %s 的向量化器配置: vectorizer=%s, hasOpenAI=%v", 
    className, schema.Vectorizer, hasOpenAI)

return hasOpenAI, nil
}