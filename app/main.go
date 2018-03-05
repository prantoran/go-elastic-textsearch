package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"time"

	"github.com/prantoran/go-elastic-textsearch/data"
	"github.com/prantoran/go-elastic-textsearch/utilities"
)

const (
	// ElasticURL to connect to, elasticsearch provides a relative urlpath within docker
	ElasticURL = "http://linkesdb:9200"
	// ElasticURL   = "elasticsearch://linkesdb:9200" // no es node found
	indexName    = "yo"
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

// Log you
type Log struct {
	App     string    `json:"app"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

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

	cmd := exec.Command("pwd")
	log.Printf("Running command and waiting for it to finish...")
	if err := cmd.Run(); err != nil {
		log.Printf("Command finished with error: %v", err)

	} else {
		b, _ := cmd.Output()
		fmt.Printf("shell output: pwd: %v\n", string(b))
	}
	cmd = exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	err = data.ESConnect(ElasticURL)
	if err != nil {
		log.Println("Connecting ElasticSearch@: ", ElasticURL)
		log.Fatal("Elasticsearch Error: ", err)
	} else {
		log.Println("Connecting ElasticSearch@: ", ElasticURL)

		log.Print("Connected to ELASTIC")
	}

	utilities.LaunchESConnectionTest()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})
	fmt.Println("Listening on :6969")
	log.Fatal(http.ListenAndServe(":6969", nil))

}
