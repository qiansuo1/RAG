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

    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    vectorSvc, err := service.NewVectorService(&cfg)
    if err != nil {
        log.Fatal(err)
    }

    pdfSvc, err := service.NewPdfService(&cfg, vectorSvc)
    if err != nil {
        log.Fatal(err)
    }

    handler := handler.NewHandler(pdfSvc,vectorSvc)
    
    r := gin.Default()
    handler.SetupRoutes(r)

    r.Run(":8080")
}


