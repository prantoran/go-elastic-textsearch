package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/prantoran/go-elastic-textsearch/conf"
	"github.com/prantoran/go-elastic-textsearch/data"
)

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "showenv":
			log.Printf("Environment")
			env := os.Environ()
			sort.Strings(env)
			for _, e := range env {
				log.Printf("- %s", e)
			}
			os.Exit(0)
		}
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := data.ESConnect(conf.ElasticURL)
	if err != nil {
		log.Println("Connecting ElasticSearch@: ", conf.ElasticURL)
		log.Fatal("Elasticsearch Error: ", err)
	} else {
		log.Println("Connecting ElasticSearch@: ", conf.ElasticURL)

		log.Print("Connected to ELASTIC")
	}

	// utilities.LaunchESConnectionTest()

	router := NewRouter()

	fmt.Println("Listening on :6969")
	log.Fatal(http.ListenAndServe(":6968", router))

}
