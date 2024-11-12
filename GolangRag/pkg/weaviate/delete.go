package weaviate

import (
        "fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
)

// Delete a Collection
func (c *Client) DeleteCollection(className string) error {
    
    err := c.client.Schema().ClassDeleter().WithClassName(className).Do(c.ctx)
    if err != nil {
        return fmt.Errorf("删除Schema失败: %w", err)
    }

    return nil
}

//根据ID删除

func (c *Client) DeleteDate(ids []string) error{
  
    operands := make([]*filters.WhereBuilder, len(ids))
    for i, id := range ids {
        operands[i] = filters.Where().
            WithPath([]string{"id"}).
            WithOperator(filters.Equal).
            WithValueText(id)  // 每个 ID 单独作为一个条件
    }

    filter := filters.Where().
    WithOperator(filters.Or).
    WithOperands(operands)

    
    _,err := c.client.Batch().ObjectsBatchDeleter().
    WithClassName("Document").
    WithOutput("minimal").
        WithWhere(filter).
        Do(c.ctx)
        if err != nil {
        return fmt.Errorf("删除数据失败: %w", err)
        }
   
   
        
    
    return nil  
}


