package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
	"github.com/prantoran/go-elastic-textsearch/conf"
	"github.com/prantoran/go-elastic-textsearch/data"
)

// SearchRequest encapsulates Search Request body
type SearchRequest struct {
	Phrase string `json:"phrase"`
}

// LawRecord represents a (lawid, title) pair
type LawRecord struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// SearchResponse represents search response
type SearchResponse struct {
	Laws         []LawRecord `json:"laws"`
	TookInMillis int64       `json:"tookinmillis,omitempty"`
	TotalHits    int64       `json:"totalhits,omitempty"`
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
	req := SearchRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := ParseError{
			Base: err,
		}
		ResponseError(w, parseErr)
		return
	}

	ctx := context.Background()
	err := data.ESConnect(conf.ElasticURL)
	if err != nil {
		res := data.StatusResponse{}
		res.Status = "Error: " + err.Error()
		ServeJSON(w, res)
	}
	q := elastic.NewQueryStringQuery(req.Phrase)

	// all_fields is used by default when the _all field is disabled and
	// no default_field is specified (either in the
	// index settings or in the request body) and no fields are specified.
	// q = q.Field("title").Field("id").Field("ammendments.ammendment") // etc.
	results, err := data.Escon.Client.Search().Index(index).Query(q).Do(ctx)
	// check err and search results here
	if err != nil {
		// Handle error
		panic(err)
	}

	laws := []LawRecord{}

	if results.Hits.TotalHits > 0 {
		for _, hit := range results.Hits.Hits {
			// hit.Index contains the name of the index
			// Deserialize hit.Source into a Laws (could also be just a map[string]interface{}).
			var t data.LawDocument
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}
			u := LawRecord{ID: t.ID, Title: t.Title}
			laws = append(laws, u)
		}
	}

	ServeJSON(w, SearchResponse{Laws: laws, TookInMillis: results.TookInMillis, TotalHits: results.TotalHits()})
}
