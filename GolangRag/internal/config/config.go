package config
import (
	"fmt"

	"os"



	"gopkg.in/yaml.v3"
)


type Config struct {
	GrpcAddress  string `yaml:"grpc_address"`  //
	
	WeaviateHost string `yaml:"weaviate_host"`
	WeaviateKey string `yaml:"weaviate_key"`
}

func LoadConfig() (Config, error) {
	var config Config
	file, err := os.ReadFile("../configs/config.yaml")
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("解析YAML失败: %v", err)
	}
	
	return config, nil
}

// type textInfo struct{
// 	pageNumber int//第几页
// 	pageCharCount int//页面的字符数
// 	pageWordCount int//页面的单词数
// 	PageSentenceCountRaw int//页面的句子数
// 	pageTokenCount int//页面的token数
// 	text string//页面的文本	
// }
