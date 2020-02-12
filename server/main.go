package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type PortList struct {
	Msg   string `json:"msg"`
	Ports []int  `json:"ports"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ports/{challengeId}/{key}", servePortList).Methods("GET")
	router.HandleFunc("/", connectionTest)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func connectionTest(resp http.ResponseWriter, req *http.Request) {
	_, err := resp.Write([]byte("flagProxy"))
	if err != nil {
		log.Print(err)
	}
}

func servePortList(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	challengeId := params["challengeId"]
	key := params["key"]

	// TODO: connect to source & get auth data
	cId := "abcdefghijk"
	k := "testkeyforchallenge0"

	resp.Header().Set("Content-Type", "application/json")

	// authenticate
	if strings.Compare(challengeId, cId) != 0 || strings.Compare(key, k) != 0 { // fails
		var result PortList
		result.Msg = "auth error"
		result.Ports = []int{}
		if err := json.NewEncoder(resp).Encode(&result); err != nil {
			fmt.Println(err)
		}
	} else { // succeeds
		// TODO: get ports data
		var result PortList
		result.Msg = "success"
		result.Ports = []int{1001, 1002, 1003, 1004}
		if err := json.NewEncoder(resp).Encode(&result); err != nil {
			fmt.Println(err)
		}
	}
}
