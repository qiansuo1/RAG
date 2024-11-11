package weaviate

import (
"fmt"
//"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
)
// Delete a Collection
func (c *Client) DeleteCollection(className string) error {
    
    err := c.client.Schema().ClassDeleter().WithClassName(className).Do(c.ctx)
    if err != nil {
        return fmt.Errorf("删除Schema失败: %w", err)
    }

    return nil
}


