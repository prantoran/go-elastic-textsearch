package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prantoran/go-elastic-textsearch/conf"
	"github.com/prantoran/go-elastic-textsearch/data"
)

// Section represents a single law section
type Section struct {
	Detail string `json:"detail"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

// Ammendment represents a single law ammendment
type Ammendment struct {
	Ammendment string   `json:"ammendment"`
	Atags      []string `json:"atags,omitempty"`
	ID         int      `json:"id"`
}

// SingleInsertRequest encapsulates the generic structure of single law record
type SingleInsertRequest struct {
	CreatedAt   string       `json:"created_at,omitempty"`
	Sections    []Section    `json:"sections,omitempty"`
	Ammendments []Ammendment `json:"ammendments,omitempty"`
	Act         string       `json:"act,omitempty"`
	ID          string       `json:"id"`
	Preamble    []string     `json:"preamble,omitempty"`
	Title       string       `json:"title,omitempty"`
}

// InsertSingle inserts a single document
func InsertSingle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	Mappingtype, ok := vars["type"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	fmt.Printf("InsertSingle()\n")
	req := SingleInsertRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := ParseError{
			Base: err,
		}
		ResponseError(w, parseErr)
		return
	}

	fmt.Printf("CreatedAt: %v\nSections: %v\nAmmendments: %v\nAct: %v\nID: %v\nPreamble: %v\nTitle: %v\n",
		req.CreatedAt, req.Sections, req.Ammendments, req.Act, req.ID, req.Preamble, req.Title)

	res := data.StatusResponse{}

	err := data.ESConnect(conf.ElasticURL)

	fmt.Printf("escon err: %v\n", err)

	if err != nil {
		res.Status = "could not connect ot escon\n"
		ServeJSON(w, res)
	}

	ctx := context.Background()

	put1, err := data.Escon.Client.Index().
		Index(index).
		Type(Mappingtype).
		Id(req.ID).
		BodyJson(req).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	res.Status = fmt.Sprintf("Inserted law with id: %v", req.ID)
	ServeJSON(w, res)

}
