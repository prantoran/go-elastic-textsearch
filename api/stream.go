package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prantoran/go-elastic-textsearch/conf"

	"git.meghdut.io/meghdut/dataengine/data"
)

const (
	indexMapping = `"mappings" : {
                            "law_details" : {
                                "properties" : {
									"created_at": {
										"type": "string",
										"index" : "not_analyzed"
									},
									"sections": {
										"type": "nested",
										"properties": {
											"details": {
												"type": "string",
												"index": "analyzed"
											},
											"id": {
												"type": "integer"
											},
											"title": {
												"type": "string",
												"index": "analyzed"
											}
										}
									},
									"ammendments": {
										"type": "nested",
										"properties": {
											"ammendment": {
												"type": "string",
												"index": "analyzed"
											},
											"atags": {
												"type": "nested"
											}
										}
									},
									"act": {
										"type": "string",
										"index": "not_analyzed"
									},
									"id": {
										"type": "string",
										"index": "not_analyzed"
									},
									"preamble": {
										"type": "nested"
									},
									"title": {
										"type": "string",
										"index": "analyzed"
									} 
                                }
                            }
                        }
                    }`
)

type Stream struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	TimestampKey string        `json:"timestamp_key"`
	Mappings     StreamMapping `json:"mappings"`
}

// SetMappingRequest represents the request body
type SetMappingRequest struct {
	Index   string                 `json:"index"`
	Type    string                 `json:"type"`
	Mapping map[string]interface{} `json:"mapping"`
}

// SetMappingResponse encapsulate simple responses
type SetMappingResponse struct {
	Status string `json:"status"`
}

// SetMapping sets the default mapping for elasticsearch
func SetMapping(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SetMapping reached\n")
	req := SetMappingRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parseErr := ParseError{
			Base: err,
		}
		ResponseError(w, parseErr)
		return
	}
	fmt.Printf("Index: %v Type: %v Mapping:\n%v\n", req.Index, req.Type, req.Mapping)
	fmt.Printf("mappings type: %T\n", req.Mapping)

	ctx := context.Background()
	fmt.Printf("defmap type: %T \n", indexMapping)

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

	res := SetMappingResponse{Status: "ok"}

	ServeJSON(w, res)
}
