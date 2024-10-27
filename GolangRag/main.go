package main

import (
	"context"
	//"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	pb	 "github.com/qiansuo1/ragservice/grpcclient" 
	tests "github.com/qiansuo1/ragservice/tests/pdf"
	    "google.golang.org/grpc/credentials/insecure"
	ragserver "github.com/qiansuo1/ragservice/ragserver"	
)

func Pdfclient() {
	config, err := tests.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	conn, err := grpc.NewClient("localhost:50051",grpc.WithTransportCredentials(insecure.NewCredentials())) // 替换为您的服务器地址
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()	
	client := pb.NewPdfServiceClient(conn)

	request := &pb.PdfRequest{
		FilePath: config.PdfPath, // 替换为您实际的PDF文件路径
	}

	stream, err := client.ExtractText(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to extract text: %v", err)
	}
	ctx := context.Background()
	wvClient, err := ragserver.InitWeaviate(ctx)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {	
			log.Fatalf("Failed to receive response: %v", err)
		}
		// fmt.Printf("Page Number: %d\n", resp.PageNumber)
        // fmt.Printf("Sentence Chunk: %s\n", resp.SentenceChunk)
        // fmt.Printf("Chunk Char Count: %d\n", resp.ChunkCharCount)
        // fmt.Printf("Chunk Word Count: %d\n", resp.ChunkWordCount)
        // fmt.Printf("Chunk Token Count: %d\n", resp.ChunkTokenCount)
        // fmt.Printf("Embedding: [%d elements]\n", len(resp.Embedding))
		ragserver.AddToWeaviate(wvClient, resp.PageNumber, resp.SentenceChunk, resp.Embedding)	
        // 可以选择打印嵌入向量的前几个元素
        // if len(resp.Embedding) > 5 {
        //   fmt.Printf("First 5 embedding values: %v\n", resp.Embedding[:5])
        // }
	}	
}

func main(){

	Pdfclient()

}

