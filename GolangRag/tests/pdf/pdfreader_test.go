package pdf

import (
	"testing"
	"os"
)



func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if config.PdfPath == "" {
		t.Error("PDF路径不应为空")
	}

	// 检查文件是否存在
	_, err = os.Stat(config.PdfPath)
	if os.IsNotExist(err) {
		t.Errorf("配置中的PDF文件不存在: %s", config.PdfPath)
	}
}