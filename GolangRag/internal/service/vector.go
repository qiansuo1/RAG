package service

import (
	"github.com/qiansuo1/ragservice/internal/config"
		"github.com/qiansuo1/ragservice/pkg/weaviate"
    
)

type VectorService struct {
    wvClient *weaviate.Client
}

func NewVectorService(cfg *config.Config) (*VectorService, error) {
    client, err := weaviate.NewClient(cfg.WeaviateHost,cfg.WeaviateKey)
    if err != nil {
        return nil, err
    }
    
    return &VectorService{
        wvClient: client,
    }, nil
}

func (s *VectorService) AddVector( pageNum int64, text string, vector []float32) error {
    return s.wvClient.AddEmbedding(pageNum, text, vector)
}

func (s *VectorService) Search(query string, limit int) ([]weaviate.SearchResult, error) {
    return s.wvClient.SearchSimilar(query, limit)
}

func (s *VectorService) ListAll(limit int) ([]weaviate.ListResult, error) {
    if limit <= 0 {
        limit = 10000 // 默认值
    }
    return s.wvClient.ListAll(limit)
}

func (s *VectorService) Delete(ids []string) error {
    return s.wvClient.DeleteDate(ids)
}   

func (s *VectorService) GetNearText(inputVectorOfText []float32, limit int) ([]weaviate.NearTextResult, error) {
    return s.wvClient.GetNearText(inputVectorOfText, limit)
}
