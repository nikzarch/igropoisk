package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

func NewClient() *elasticsearch.Client {
	url := os.Getenv("ELASTIC_URL")

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elastic client: %s", err.Error())
	}
	return es
}
