// package ragserver

// import (
// 	"cmp"
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/sashabaranov/go-openai"
// 	"github.com/weaviate/weaviate-go-client/v4/weaviate"
// 	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
// 	"github.com/weaviate/weaviate/entities/models"
// )	
// func raginit(){
// 	ctx := context.Background()
// 	wvClient, err := initWeaviate(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	openaiClient := openai.NewClient(apiKey)
// 	server := &ragServer{
// 		ctx: ctx,
// 		wvClient: wvClient,
// 		openaiClient: openaiClient,
// 	}	
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("POST/add/", server.addDocumentsHandler)
// 	mux.HandleFunc("POST/query/", server.queryHandler)

// 	port := cmp.Or(os.Getenv("SERVERPORT"), "9020")
// 	address := "localhost:" + port
// 	log.Println("listening on", address)
// 	log.Fatal(http.ListenAndServe(address, mux))
// }

// type ragServer struct {
// 	ctx      context.Context
// 	wvClient *weaviate.Client
// 	openaiClient *openai.Client
// }

// func (s *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
// 	type document struct{
// 		Text string
// 	}
// 	type addRequest struct {
// 		Documents []document
// 	}
// 	ar := &addRequest{}
// 	err := readRequestJSON(req, ar)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}	
	
// 	// 使用 OpenAI 的 API 批量生成嵌入
// 	var texts []string
// 	for _, doc := range ar.Documents {
// 		texts = append(texts, doc.Text)
// 	}
	
// 	log.Printf("为 %v 个文档生成嵌入", len(texts))
// 	embeddings, err := s.openaiClient.CreateEmbeddings(s.ctx, openai.EmbeddingRequest{
// 		Input: texts,
// 		Model: openai.AdaEmbeddingV2,
// 	})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if len(embeddings.Data) != len(ar.Documents) {
// 		http.Error(w, "嵌入批次大小不匹配", http.StatusInternalServerError)
// 		return
// 	}

// 	// 准备 Weaviate 对象
// 	objects := make([]*models.Object, len(ar.Documents))
// 	for i, doc := range ar.Documents {
// 		objects[i] = &models.Object{
// 			Class: "Document",
// 			Properties: map[string]interface{}{
// 				"text": doc.Text,
// 			},
// 			Vector: embeddings.Data[i].Embedding,
// 		}
// 	}

// 	// 将文档和嵌入存储到 Weaviate
// 	log.Printf("正在向 Weaviate 存储 %v 个对象", len(objects))
// 	_, err = s.wvClient.Batch().ObjectsBatcher().WithObjects(objects...).Do(s.ctx)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("所有文档已处理并存储"))

// }

// func (s *ragServer) queryHandler(w http.ResponseWriter, req *http.Request) {
// 	type queryRequest struct {
// 		Content string `json:"content"`
// 	}
// 	qr := &queryRequest{}
// 	err := readRequestJSON(req, qr)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	// Embed the query contents.
// 	//调用OpenAI的API来将请求内容向量化，为后续使用向量数据库查询做准备
// 	rsp, err := s.openaiClient.CreateEmbeddings(s.ctx, openai.EmbeddingRequest{
// 		Input: []string{qr.Content},
// 		Model: openai.AdaEmbeddingV2,
// 	})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	

// // 使用Weaviate进行向量搜索
// 	gql := s.wvClient.GraphQL()
// 	result, err := gql.Get().
// 		WithNearVector(
// 			gql.NearVectorArgBuilder().WithVector(rsp.Data[0].Embedding)).
// 		WithClassName("Document").
// 		WithFields(graphql.Field{Name: "text"}).
// 		WithLimit(3).
// 		Do(s.ctx)

// 	if err != nil {
// 		http.Error(w, "搜索相关文档失败", http.StatusInternalServerError)
// 		return
// 	}
// 	contents, err := decodeGetResults(result)
// 	if err != nil {
// 		http.Error(w, fmt.Errorf("reading weaviate response: %w", err).Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Creata a RAG query for the LLM with the most relevant documents as
// 	// context.
// 	ragQuery := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(contents, "\n"))
// 	resp, err := s.openaiClient.CreateChatCompletion(s.ctx, openai.ChatCompletionRequest{
// 		Model: openai.GPT3Dot5Turbo,
// 		Messages: []openai.ChatCompletionMessage{
// 			{
// 				Role: openai.ChatMessageRoleUser,
// 				Content: ragQuery,
// 			},
// 		},
// 	})

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if len(resp.Choices) > 0 {
// 		answer := resp.Choices[0].Message.Content
// 		log.Printf("回答: %s\n", answer)
// 	} else {	
// 		log.Println("API 响应中没有选择")
// 	}	

// 	var respTexts []string
// 	for _, part := range resp.Choices[0].Message.Content {
// 		respTexts = append(respTexts, string(part))
// 	}		
// 	renderJSON(w, strings.Join(respTexts, "\n"))

// }

// func decodeGetResults(result *models.GraphQLResponse) ([]string, error) {
// 	data, ok := result.Data["Get"]
// 	if !ok {
// 		return nil, fmt.Errorf("get key not found in result")
// 	}
// 	doc, ok := data.(map[string]any)
// 	if !ok {
// 		return nil, fmt.Errorf("get key unexpected type")
// 	}
// 	slc, ok := doc["Document"].([]any)
// 	if !ok {
// 		return nil, fmt.Errorf("document is not a list of results")
// 	}

// 	var out []string
// 	for _, s := range slc {
// 		smap, ok := s.(map[string]any)
// 		if !ok {
// 			return nil, fmt.Errorf("invalid element in list of documents")
// 		}
// 		s, ok := smap["text"].(string)
// 		if !ok {
// 			return nil, fmt.Errorf("expected string in list of documents")
// 		}
// 		out = append(out, s)
// 	}
// 	return out, nil
// }


// const ragTemplateStr = `
// I will ask you a question and will provide some additional context information.
// Assume this context information is factual and correct, as part of internal
// documentation.
// If the question relates to the context, answer it using the context.
// If the question does not relate to the context, answer it as normal.

// For example, let's say the context has nothing in it about tropical flowers;
// then if I ask you about tropical flowers, just answer what you know about them
// without referring to the context.

// For example, if the context does mention minerology and I ask you about that,
// provide information from the context along with general knowledge.

// Question:
// %s

// Context:
// %s
// `