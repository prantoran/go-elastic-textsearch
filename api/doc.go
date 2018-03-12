package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/gorilla/mux"
	"github.com/prantoran/go-elastic-textsearch/conf"
	"github.com/prantoran/go-elastic-textsearch/data"
)

// InsertBulk inserts records from request body in bulk
// using ES's Bulk Processor
func InsertBulk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	mappingtype, ok := vars["type"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	fmt.Printf("EnterBulk()\n")

	res := data.StatusResponse{}

	err := data.ESConnect(conf.ElasticURL)

	fmt.Printf("escon err: %v\n", err)

	if err != nil {
		res.Status = "could not connect ot escon\n"
		ServeJSON(w, res)
	}

	ctx := context.Background()

	p, err := data.Escon.Client.BulkProcessor().
		Name("MyBackgroundWorker-1").
		Workers(2).
		BulkActions(1000).               // commit if # requests >= 1000
		BulkSize(2 << 20).               // commit if size of requests >= 2 MB
		FlushInterval(30 * time.Second). // commit every 30s
		Do(ctx)

	if err != nil {
	}

	// ... Do some work here

	req := []data.LawDocument{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := ParseError{
			Base: err,
		}
		ResponseError(w, parseErr)
		return
	}

	for _, u := range req {
		fmt.Printf("ID: %v Title: %v\n", u.ID, u.Title)
		r := elastic.NewBulkIndexRequest().Index(index).Type(mappingtype).Id(u.ID).Doc(u)
		p.Add(r)
	}
	fmt.Printf("passed\n")
	// Stop the bulk processor and do some cleanup
	err = p.Close()
	if err != nil {
		fmt.Printf("pclose: %v\n", err)
	}

	res.Status = fmt.Sprintf("Inserted law with id: %v", req[0].ID)
	ServeJSON(w, res)

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

	mappingtype, ok := vars["type"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	fmt.Printf("InsertSingle()\n")
	req := data.LawDocument{}

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
		Type(mappingtype).
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

// GetSingle retrieves a single record
func GetSingle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	mappingtype, ok := vars["type"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	id, ok := vars["id"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	ctx := context.Background()

	err := data.ESConnect(conf.ElasticURL)

	fmt.Printf("escon err: %v\n", err)
	defaultres := data.StatusResponse{}
	if err != nil {
		defaultres.Status = "could not connect ot escon\n"
		ServeJSON(w, defaultres)
	}

	get1, err := data.Escon.Client.Get().
		Index(index).
		Type(mappingtype).
		Id(id).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		bytes, err := json.Marshal(get1.Source)
		if err != nil {
			defaultres.Status = err.Error()
			ServeJSON(w, defaultres)

		}
		fmt.Printf("Got document %s in version %d from index %s, type %s\nUID: %v\nrouting: %v\nParent: %v\n\nsource: %v\n\nfields: %v\\n",
			get1.Id, get1.Version, get1.Index, get1.Type, get1.Uid, get1.Routing, get1.Parent, string(bytes), get1.Fields)
		res := data.LawDocument{}
		err = json.Unmarshal(bytes, &res)
		if err != nil {
			defaultres.Status = err.Error()
			ServeJSON(w, defaultres)

		}
		ServeJSON(w, res)

	} else {
		fmt.Printf("Document not found\n")

		defaultres.Status = "Document not found"
		ServeJSON(w, defaultres)
	}
}

// DeleteSingle deletes a single record
func DeleteSingle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	mappingtype, ok := vars["type"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}

	id, ok := vars["id"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}
	ctx := context.Background()

	err := data.ESConnect(conf.ElasticURL)
	res := data.StatusResponse{}
	if err != nil {
		res.Status = "Error: " + err.Error()
		ServeJSON(w, res)
	}
	// Delete tweet with specified ID
	delres, err := data.Escon.Client.Delete().
		Index(index).
		Type(mappingtype).
		Id(id).
		Do(ctx)
	if err != nil {
		// Handle error
		res.Status = "Error: " + err.Error()
		ServeJSON(w, res)

	}
	if delres.Found {
		res.Status = fmt.Sprintf("Document with id %v deleted", id)

	} else {
		res.Status = fmt.Sprintf("Deletion not complete")

	}
	ServeJSON(w, res)
}
