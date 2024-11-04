package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/qiansuo1/ragservice/internal/config"
	grpcclient "github.com/qiansuo1/ragservice/pkg/grpc/pdfservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PdfService struct {
	conn       *grpc.ClientConn
    grpcClient grpcclient.PdfServiceClient
    vectorSvc  *VectorService
}

func NewPdfService(cfg *config.Config, vectorSvc *VectorService) (*PdfService, error) {
	//creat gprc connection
	conn, err := grpc.NewClient(
        cfg.GrpcAddress,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("创建gRPC连接失败: %w", err)
    }
	defer conn.Close()
	
	//Creating a grpc client
    client := grpcclient.NewPdfServiceClient(conn)

    // Creating a pdf service instance
    service := &PdfService{
		conn: conn,
        grpcClient: client,
        vectorSvc:  vectorSvc,
    }

    return service, nil
}

const (
    maxRetries    = 3        // 最大重试次数
    retryInterval = time.Second * 2  // 重试间隔
)

func (s *PdfService) ProcessPdf(filePath string) error {
	
	request := &grpcclient.PdfRequest{
        FilePath: filePath,
    }

    stream, err := s.grpcClient.ExtractText(context.Background(), request)
    if err != nil {
        return fmt.Errorf("提取文本失败: %w", err)
    }

    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("接收响应失败: %w", err)
        }

        // 将数据保存到向量数据库
        if err := s.addVectorWithRetry(resp); err != nil {
            return fmt.Errorf("保存向量数据失败: %w", err)
        }
    }

    return nil
}	
// addVectorWithRetry 带重试机制的向量添加
func (s *PdfService) addVectorWithRetry(resp *grpcclient.PdfResponse) error {

    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        if attempt > 0 {
            time.Sleep(retryInterval * time.Duration(attempt))
            log.Printf("第 %d 次重试添加向量，页码: %d", attempt+1, resp.PageNumber)
        }

        err := s.vectorSvc.AddVector(
			int64(resp.PageNumber),
            resp.SentenceChunk,
            resp.Embedding,
        )
        
        if err == nil {
            return nil
        }

        lastErr = err
        log.Printf("添加向量失败 (尝试 %d/%d): %v", attempt+1, maxRetries, err)
    }

    return fmt.Errorf("添加向量失败，已重试%d次: %w", maxRetries, lastErr)

}
// Close 关闭连接
func (s *PdfService) Close() error {
    if s.conn != nil {
        return s.conn.Close()
    }
    return nil
}