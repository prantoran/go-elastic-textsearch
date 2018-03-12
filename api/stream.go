package api

import (
	"context"
	"encoding/json"
	"errors"
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
		err := data.ESError{
			Base: errors.New("Could not connect to ES"),
		}
		ResponseError(w, err)
		return
	}

	createIndex, err := data.Escon.Client.CreateIndex(req.Index).BodyJson(req.Mapping).Do(ctx)
	if err != nil {
		err := data.ESError{
			Base: errors.New("Could not create index in ES"),
		}
		ResponseError(w, err)
		return
	}
	if !createIndex.Acknowledged {
		err := data.ESError{
			Base: errors.New("Creating index in ES not acknoledged"),
		}
		ResponseError(w, err)
		return
	}

	ServeJSON(w, data.StatusResponse{Status: "Index created"})
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
		err := data.ESError{
			Base: errors.New("Could not connect to ES"),
		}
		ResponseError(w, err)
		return
	}
	ctx := context.Background()
	exists, err := data.Escon.Client.IndexExists(index).Do(ctx)
	if err != nil {
		err := data.ESError{
			Base: errors.New("Could not check whether index exists"),
		}
		ResponseError(w, err)
		return
	}

	if !exists {
		err := NotFoundError{
			Base: errors.New("Document not found"),
		}
		ResponseError(w, err)
	}
	ServeJSON(w, data.StatusResponse{Status: "Index Exists"})
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
		err := data.ESError{
			Base: errors.New("Could not connect to ES"),
		}
		ResponseError(w, err)
		return
	}

	ctx := context.Background()
	del, err := data.Escon.Client.DeleteIndex(index).Do(ctx)
	if err != nil {
		err := data.ESError{
			Base: errors.New("Could not delete index from ES"),
		}
		ResponseError(w, err)
		return
	}
	if !del.Acknowledged {
		err := data.ESError{
			Base: errors.New("Index not deleted from ES"),
		}
		ResponseError(w, err)
		return
	}

	ServeJSON(w, data.StatusResponse{Status: "Index deleted"})
}
