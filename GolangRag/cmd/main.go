package main

import (
	"log"
	"github.com/qiansuo1/ragservice/internal/config"
	"github.com/qiansuo1/ragservice/internal/service"
	"github.com/qiansuo1/ragservice/internal/handler"
	"github.com/gin-gonic/gin"
)	
	
// 主程序入口
func main() {
    log.Println("开始加载配置...")
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("配置加载成功")

    log.Println("初始化向量服务...")
    vectorSvc, err := service.NewVectorService(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("向量服务初始化成功")
    log.Println("初始化PDF服务...")
    pdfSvc, err := service.NewPdfService(&cfg, vectorSvc)
    if err != nil {
        log.Fatal(err)
    }
    defer pdfSvc.Close()
    log.Println("PDF服务初始化成功")
    handler := handler.NewHandler(pdfSvc,vectorSvc)
    
    r := gin.Default()
    handler.SetupRoutes(r)

    log.Printf("服务器启动成功! 监听端口: http://localhost:8080")
    log.Printf("可用的API接口:")
    log.Printf("- POST /api/pdf/upload     : 上传PDF文件")
    log.Printf("- POST /api/pdf/process    : 处理本地PDF文件")
    log.Printf("- POST /api/vectors/search : 向量搜索")
    log.Printf("- GET /api/vectors/list : 数据列表")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("服务器启动失败:", err)
    }

}


