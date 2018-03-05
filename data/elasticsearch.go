package data

import (
	"log"
	"os"

	es "gopkg.in/olivere/elastic.v5"
)

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
