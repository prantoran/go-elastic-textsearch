package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	es "gopkg.in/olivere/elastic.v5"
)

const (
	// ElasticURL to connect to, elasticsearch provides a relative urlpath within docker
	ElasticURL   = "http://elasticsearch:9200"
	indexName    = "applications"
	docType      = "log"
	appName      = "myApp"
	indexMapping = `{
                        "mappings" : {
                            "log" : {
                                "properties" : {
                                    "app" : { "type" : "string", "index" : "not_analyzed" },
                                    "message" : { "type" : "string", "index" : "not_analyzed" },
                                    "time" : { "type" : "date" }
                                }
                            }
                        }
                    }`
)

type Log struct {
	App     string    `json:"app"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func main() {
	err := ESConnect(ElasticURL)
	if err != nil {
		log.Println("Connecting ElasticSearch@: ", ElasticURL)
		log.Fatal("Elasticsearch Error: ", err)
	} else {
		log.Print("Connected to ELASTIC")
	}
	err = createIndexWithLogsIfDoesNotExist(Escon.Client)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})

	fmt.Println("Listening on :6969")
	log.Fatal(http.ListenAndServe(":6969", nil))

}

func createIndexWithLogsIfDoesNotExist(client *es.Client) error {
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("exists: %v\n", exists)

	if exists {
		return nil
	}

	res, err := client.CreateIndex(indexName).
		Body(indexMapping).
		Do(context.Background())

	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct")
	}

	return addLogsToIndex(client)
}

func addLogsToIndex(client *es.Client) error {
	for i := 0; i < 10; i++ {
		l := Log{
			App:     "myApp",
			Message: fmt.Sprintf("message %d", i),
			Time:    time.Now(),
		}

		_, err := client.Index().
			Index(indexName).
			Type(docType).
			BodyJson(l).
			Do(context.Background())

		if err != nil {
			return err
		}
	}

	return nil
}

func findAndPrintAppLogs(client *es.Client) error {
	termQuery := es.NewTermQuery("app", appName)

	res, err := client.Search(indexName).
		Index(indexName).
		Query(termQuery).
		Sort("time", true).
		Do(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("Logs found:")
	var l Log
	for _, item := range res.Each(reflect.TypeOf(l)) {
		l := item.(Log)
		fmt.Printf("time: %s message: %s\n", l.Time, l.Message)
	}

	return nil
}

// data: elasticsearch.go

// Escon is the connetion variable
var Escon ESConnection

// ESConnection represents a connection to elasticsearch
type ESConnection struct {
	Client *es.Client
}

// ESConnect makes a connection with the given URL
func ESConnect(url string) error {
	f, _ := os.Create("elastic_search.log")
	client, err := es.NewClient(es.SetURL(url), es.SetTraceLog(log.New(f, "", 0)), es.SetSniff(false), es.SetHealthcheck(false))
	if err != nil {
		return err
	}
	Escon = ESConnection{Client: client}
	return nil
}
