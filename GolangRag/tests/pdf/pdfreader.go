package pdf

import (
	"fmt"

	"os"



	"gopkg.in/yaml.v3"
)


type Config struct {
	PdfPath string `yaml:"pdf_path"`
}

func LoadConfig() (Config, error) {
	var config Config
	file, err := os.ReadFile("./tests/pdf/config.yaml")
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("解析YAML失败: %v", err)
	}
	
	return config, nil
}
type textInfo struct{
	pageNumber int//第几页
	pageCharCount int//页面的字符数
	pageWordCount int//页面的单词数
	PageSentenceCountRaw int//页面的句子数
	pageTokenCount int//页面的token数
	text string//页面的文本	
}

