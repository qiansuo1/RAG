package weaviate



import (
	"fmt"

	"github.com/weaviate/weaviate/entities/models"
)


// CreateCollection 
// PropertyDefinition 
type PropertyDefinition struct {
    Name     string   `json:"name"`
    DataType []string `json:"dataType"`
}

// CollectionConfig 定义集合配置
type CollectionConfig struct {
    Name       string              `json:"name"`
    Properties []PropertyDefinition `json:"properties"`
    Vectorizer string              `json:"vectorizer"` // 可以是 "none", "text2vec-openai" 等
}

// CreateSchema 创建Schema
func (c *Client) CreateCollection(config CollectionConfig) error {
    if config.Name == "" {
        return fmt.Errorf("collection name is empty")
    }
    if len(config.Properties) == 0 {
        return fmt.Errorf("at least one property is required")
    }
    	   // 检查类是否存在
	   exists, err := c.client.Schema().ClassExistenceChecker().
	   WithClassName(config.Name).
	   Do(c.ctx)
   if err != nil {
	   return err
   }
   if exists {
    return fmt.Errorf("%s: %w", config.Name, ErrCollectionExist)
}
properties := make([]*models.Property, len(config.Properties))
for i, prop := range config.Properties {
    properties[i] = &models.Property{
        Name:     prop.Name,
        DataType: prop.DataType,
    }
}

    // 设置默认值
    if config.Vectorizer == "" {
        config.Vectorizer = "none"
    }

       // 只在类不存在时创建

    classObj := &models.Class{
        Class: config.Name,
        Properties: properties,
        Vectorizer: config.Vectorizer, // 使用自定义向量
    }
    
    err = c.client.Schema().ClassCreator().WithClass(classObj).Do(c.ctx)
    if err != nil {
        return fmt.Errorf("创建Schema失败: %w", err)
    }

    return nil

    

}



func (c *Client) AddEmbedding(pageNum int64, text string, vector []float32) error {
  
    if err := c.ensureCollectionExists(); err != nil {
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
