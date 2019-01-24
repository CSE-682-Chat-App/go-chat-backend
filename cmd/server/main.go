package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GenericResponse is the structure for the error response messages
type GenericResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func tmpRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenericResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("%s %s", "Hello", "World"),
	})
}

//HandleRequest acts as a generic middleware for handling the request
func HandleRequest(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps)
	}
}

func main() {

	router := httprouter.New()

	//Authentication endpoints
	router.GET("/", HandleRequest(tmpRoute))

	log.Info("Starting Server on port 9090")
	log.Fatal(http.ListenAndServe(":9090", router))
}
