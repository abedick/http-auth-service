package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func initServer() {

	r := mux.NewRouter()

	r.HandleFunc("/", handleStatus).Methods("GET")
	r.HandleFunc("/status", handleStatus).Methods("GET")
	r.HandleFunc("/status", isAuth(handleStatus)).Methods("POST")
	r.HandleFunc("/update", isAuth(handleUpdate)).Methods("POST")
	r.HandleFunc("/validate", isAuth(isConfig(handleValidate))).Methods("POST")
	r.HandleFunc("/issue", isAuth(isConfig(handleIssue))).Methods("POST")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		http.ListenAndServe(c.Port, universalLogger(r))
		wg.Done()
	}()
	wg.Wait()
}

func errorResponse(w http.ResponseWriter, err error) {
	response := publicResponse{
		Status:  "error",
		Message: err.Error(),
	}
	c.Logger.Message("error", err.Error())
	output, _ := json.Marshal(response)
	w.Write([]byte(output))
}

func successResponse(w http.ResponseWriter, response publicResponse, msg string) {
	response.Status = "success"
	c.Logger.Message("success", msg)
	output, _ := json.Marshal(response)
	w.Write([]byte(output))
}
