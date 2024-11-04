// 从 ragserver/weaviate.go 移出 HTTP 处理相关代码
package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiansuo1/ragservice/internal/service"
)
    
type Handler struct {
    pdfSvc    *service.PdfService
    vectorSvc *service.VectorService
}

func NewHandler(pdfSvc *service.PdfService, vectorSvc *service.VectorService) *Handler {
    return &Handler{
        pdfSvc:    pdfSvc,
        vectorSvc: vectorSvc,
    }
}


// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(r *gin.Engine) {
    // 创建上传文件的临时目录
    if err := os.MkdirAll("temp", 0755); err != nil {
        log.Fatal("创建临时目录失败:", err)
    }

    api := r.Group("/api")
    {
        api.POST("/pdf/upload", h.HandlePDFUpload)
        api.POST("/vectors/search", h.HandleVectorSearch)
    }
}

const (
    maxRetries    = 3        // 最大重试次数
    retryInterval = time.Second * 2  // 重试间隔
)


// HandlePDFUpload 处理PDF文件上传
func (h *Handler) HandlePDFUpload(c *gin.Context) {
    // 获取上传的文件
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "文件上传失败",
            "detail": err.Error(),
        })
        return
    }

    // 验证文件类型
    if filepath.Ext(file.Filename) != ".pdf" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "只支持PDF文件",
        })
        return
    }

    // 保存文件到临时目录
    tempPath := filepath.Join("temp", file.Filename)
    if err := c.SaveUploadedFile(file, tempPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "保存文件失败",
            "detail": err.Error(),
        })
        return
    }

    // 处理PDF文件
    if err := h.processWithRetry(tempPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "处理PDF失败",
            "detail": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "PDF处理成功",
        "filename": file.Filename,
    })
}

// processWithRetry 带重试机制的PDF处理
func (h *Handler) processWithRetry(filepath string) error {
    var lastErr error
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        // 如果不是第一次尝试，等待一段时间
        if attempt > 0 {
            time.Sleep(retryInterval * time.Duration(attempt))
            log.Printf("第 %d 次重试处理文件: %s", attempt+1, filepath)
        }

        err := h.pdfSvc.ProcessPdf(filepath)
        if err == nil {
            if attempt > 0 {
                log.Printf("重试成功: %s", filepath)
            }
            return nil
        }

        lastErr = err
        log.Printf("处理失败 (尝试 %d/%d): %v", attempt+1, maxRetries, err)
    }

    return fmt.Errorf("处理失败，已重试%d次: %w", maxRetries, lastErr)
}

// SearchRequest 搜索请求结构
type SearchRequest struct {
    Query string `json:"query" binding:"required"`
    Limit int    `json:"limit,omitempty"`
}

// HandleVectorSearch 处理向量搜索
func (h *Handler) HandleVectorSearch(c *gin.Context) {
    var req SearchRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "无效的请求参数",
            "detail": err.Error(),
        })
        return
    }

    // 设置默认限制
    if req.Limit <= 0 {
        req.Limit = 10
    }

    // 执行搜索
    results, err := h.vectorSvc.Search(req.Query, req.Limit)   
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "搜索失败",
            "detail": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "results": results,
    })
}