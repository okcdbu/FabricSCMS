package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gin-gonic/gin"
)

type elasticClient struct {
	es        *elasticsearch.Client
	IndexName string
}

type ESResponse struct {
	Took int64
	Hits struct {
		Total struct {
			Value int64
		}
		Hits []*ESHit
	}
}

type ESHit struct {
	Score   float64 `json:"_score"`
	Index   string  `json:"_index"`
	Type    string  `json:"_type"`
	Version int64   `json:"_version,omitempty"`

	Source Article `json:"_source"`
}

type Article struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Author          string              `json:"author"`
	UploadDate      string              `json:"uploaded"`
	Description     string              `json:"description"`
	Platform        string              `json:"platform"`
	SignaturePolicy string              `json:"signature_policy"`
	CCLanguages     []*CCLanguage       `json:"cc_languages"`
	AppLanguages    []map[string]string `json:"app_languages"`
	Versions        []map[string]string `json:"versions"`
}

type CCLanguage struct {
	Language     string            `json:"language"`
	Link         string            `json:"link"`
	AssetStruct  map[string]string `json:"asset_struct"`
	Dependencies map[string]string `json:"dependencies"`
}

var (
	esClient = &elasticClient{}
)

func init() {
	cfg := esClientConfig()
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating the client: %s", err)
	} else {
		log.Println(elasticsearch.Version)
		log.Println(es.Info())
	}
	esClient.es = es
	esClient.IndexName = "smart_contract"
}

func esClientConfig() elasticsearch.Config {
	cfg := elasticsearch.Config{
		Addresses:              []string{"https://localhost:9200"},
		APIKey:                 "ajVuZzE0SUJlczZ5ZEtIV2FEaDI6Wjh0ZXZySmNSaG1HbTBkakpGcC1pdw==",
		CertificateFingerprint: "39342712d2129a0a4bb9c835452777108e30588860ae260f6b11a7e53aae7659",
	}
	//password : XEGwyYour=Xi*wdhYIRl
	return cfg
}

func ESSearchAll(c *gin.Context) {

	searchString := c.Query("filter")
	log.Println(searchString)

	var searchRequest map[string]interface{}
	err := json.Unmarshal([]byte(searchString), &searchRequest)
	if err != nil {
		log.Println(err)
	}

	if searchRequest["q"] == nil {
		searchRequest["q"] = "*"
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": searchRequest["q"],
			},
		},
	}

	res, err := esClient.es.Search(
		esClient.es.Search.WithIndex(esClient.IndexName),
		esClient.es.Search.WithBody(esutil.NewJSONReader(query)),
		esClient.es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer res.Body.Close()

	var sr ESResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		log.Printf("Error: %s\n", err)
	}

	var scs []Article

	for _, h := range sr.Hits.Hits {
		scs = append(scs, h.Source)
	}
	c.Header("content-range", fmt.Sprintf("%d", len(scs)))

	if scs == nil {
		scs = make([]Article, 0)
	}
	c.IndentedJSON(http.StatusOK, scs)
}

func ESSearchWithLanguage(c *gin.Context) {
	lang := c.Query("language")

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"cc_languages.language": lang,
			},
		},
	}

	res, err := esClient.es.Search(
		esClient.es.Search.WithIndex("smart_contract"),
		esClient.es.Search.WithBody(esutil.NewJSONReader(query)),
		esClient.es.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Error: %s", err)
	}
	defer res.Body.Close()

	var sr ESResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		log.Printf("Error: %s\n", err)
	}

	var scs []Article

	for _, h := range sr.Hits.Hits {
		scs = append(scs, h.Source)
	}
	log.Println(scs)
	c.IndentedJSON(http.StatusOK, scs)
}

func EsDocumentByID(c *gin.Context) {
	ccId := c.Param("id")
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": ccId,
			},
		},
	}

	res, err := esClient.es.Search(
		esClient.es.Search.WithIndex("smart_contract"),
		esClient.es.Search.WithBody(esutil.NewJSONReader(query)),
		esClient.es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer res.Body.Close()

	var sr ESResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		log.Printf("Error: %s\n", err)
	}

	c.IndentedJSON(http.StatusOK, sr.Hits.Hits[0].Source)
}

func AddDocumentToES(item *Article) (string, error) {
	payload, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return "", err
	}
	ctx := context.Background()
	req := esapi.IndexRequest{
		Index:      esClient.IndexName,
		DocumentID: string(item.ID),
		Body:       bytes.NewReader(payload),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, esClient.es)
	if err != nil {
		log.Fatalf("Error getting rsponse: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Print("payload : ")
			log.Println(e)
			log.Println(err)
			return "", err
		}
		log.Print(e)
		return "", fmt.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	return "Contract successfully added to search index", nil
}
