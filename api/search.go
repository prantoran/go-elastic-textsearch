package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
	"github.com/prantoran/go-elastic-textsearch/conf"
	"github.com/prantoran/go-elastic-textsearch/data"
)

// SearchRequest encapsulates Search Request body
type SearchRequest struct {
	Body json.RawMessage `json:"body"`
}

// QuerySearchQuery handles search query, using ElasticSearch's QueryStringQuery
func QuerySearchQuery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}
	fmt.Printf("index: %v\n", index)
	req := SearchRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := ParseError{
			Base: err,
		}
		ResponseError(w, parseErr)
		return
	}

	// Search with a term query
	ctx := context.Background()

	err := data.ESConnect(conf.ElasticURL)
	res := data.StatusResponse{}
	if err != nil {
		res.Status = "Error: " + err.Error()
		ServeJSON(w, res)
	}
	fmt.Printf("ok\n")
	q := elastic.NewQueryStringQuery("The es 134 of")
	fmt.Printf("ok2\n")
	q = q.Field("title").Field("id").Field("ammendments.ammendment") // etc.
	results, err := data.Escon.Client.Search().Index(index).Query(q).Do(ctx)
	// check err and search results here
	fmt.Printf("resuerr: %v\n", err)

	if err != nil {
		// Handle error
		panic(err)
	}

	// results is of type results and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", results.TookInMillis)
	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over the process, see below.
	var ttyp data.LawDocument
	for _, item := range results.Each(reflect.TypeOf(ttyp)) {
		t := item.(data.LawDocument)
		fmt.Printf("Laws by %s: %s\n", t.ID, t.Title)
	}

	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d Laws\n", results.TotalHits())

	// Here's how you iterate through the search results with full control over each step.
	if results.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d Laws\n", results.Hits.TotalHits)

		// Iterate through results
		for _, hit := range results.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Laws (could also be just a map[string]interface{}).
			var t data.LawDocument
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with Laws
			fmt.Printf("Law by %s: %s\n", t.ID, t.Title)
		}
	} else {
		// No hits
		fmt.Print("Found no Laws\n")
	}
	ServeJSON(w, data.StatusResponse{Status: "ok"})
}
