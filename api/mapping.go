package api

import (
	"fmt"
	"net/http"
)

// DefaultResponse encapsulate simple responses
type DefaultResponse struct {
	Status string `json:"status"`
}

// SetDefaultMapping sets the default mapping for elasticsearch
func SetDefaultMapping(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SetDefaultMapping reached\n")
	res := DefaultResponse{Status: "ok"}

	ServeJSON(w, res)
}
