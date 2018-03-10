package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/prantoran/go-elastic-textsearch/data"

	"github.com/gorilla/mux"
	"github.com/prantoran/go-elastic-textsearch/conf"
)

// Stream represents index of ES
type Stream struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TimestampKey string `json:"timestamp_key"`
}

// SetMappingRequest represents the request body
type SetMappingRequest struct {
	Index   string                 `json:"index"`
	Type    string                 `json:"type"`
	Mapping map[string]interface{} `json:"mapping"`
}

// SetMapping sets the default mapping for elasticsearch
func SetMapping(w http.ResponseWriter, r *http.Request) {
	req := SetMappingRequest{}

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
		fmt.Printf("could not connect ot escon\n")
	}

	fmt.Printf("Escon: %v\n", data.Escon)
	createIndex, err := data.Escon.Client.CreateIndex(req.Index).BodyJson(req.Mapping).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		fmt.Printf("Not Acknowledged\n")
	}

	res := data.StatusResponse{Status: "ok"}

	ServeJSON(w, res)
}

// IndexExists checks whether index exists
func IndexExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}
	err := data.ESConnect(conf.ElasticURL)

	if err != nil {
		fmt.Printf("could not connect ot escon\n")
	}
	ctx := context.Background()
	exists, err := data.Escon.Client.IndexExists(index).Do(ctx)
	res := data.StatusResponse{}
	if exists {
		res.Status = "Index exists"
	} else {
		res.Status = "Index does not exist"
	}
	ServeJSON(w, res)
}

// DeleteIndex deletes an index
func DeleteIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, ok := vars["index"]
	if ok == false {
		err := data.InvalidIDError{
			Base: errors.New("ID parameter does not exist"),
		}
		ResponseError(w, err)
		return
	}
	err := data.ESConnect(conf.ElasticURL)

	if err != nil {
		fmt.Printf("could not connect ot escon\n")
	}
	res := data.StatusResponse{}

	ctx := context.Background()
	deleteIndex, err := data.Escon.Client.DeleteIndex(index).Do(ctx)
	if err != nil {
		res.Status = "Error:" + err.Error()
		ServeJSON(w, res)
	}
	if deleteIndex.Acknowledged {
		res.Status = "Index deleted"
	} else {
		res.Status = "Error: Index was not deleted"
	}
	ServeJSON(w, res)
}
