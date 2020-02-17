package main

import (
	"database/sql"
	"encoding/json"
	"flagProxy/client/swaper"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ChallengeId string
	Key         string
)

const (
	host     = "localhost"
	dbPort   = 5432
	user     = "Tp0tOj"
	password = "Tp0tOj-dev"
	dbName   = "Tp0tOj"
)

type PortList struct {
	Msg   string `json:"msg"`
	Ports []int  `json:"ports"`
}

func init() {
	// TODO: connect to data source & get auth data
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("server connected")

	ChallengeId = "abcdefghijk"
	Key = "testkeyforchallenge0"
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ports/{challengeId}/{key}", servePortList).Methods("GET")
	router.HandleFunc("/flagByPort/{challengeId}/{key}/{port}", serveFlagByPort).Methods("GET")
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

func servePortList(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	challengeId := params["challengeId"]
	key := params["key"]

	writer.Header().Set("Content-Type", "application/json")

	// authenticate
	if strings.Compare(challengeId, ChallengeId) != 0 || strings.Compare(key, Key) != 0 { // fails
		var result PortList
		result.Msg = "auth error"
		result.Ports = []int{}
		if err := json.NewEncoder(writer).Encode(&result); err != nil {
			fmt.Println(err)
		}
	} else { // succeeds
		// TODO: get ports data
		var result PortList
		result.Msg = "success"
		result.Ports = []int{10001, 10002, 10003, 10004}
		if err := json.NewEncoder(writer).Encode(&result); err != nil {
			fmt.Println(err)
		}
	}
}

func serveFlagByPort(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	challengeId := params["challengeId"]
	key := params["key"]
	port := params["port"]

	writer.Header().Set("Content-Type", "application/json")

	// authenticate
	if strings.Compare(challengeId, ChallengeId) != 0 || strings.Compare(key, Key) != 0 { // fails
		var flag swaper.FlagByPort
		flag.Msg = "auth error"
		flag.Flag = ""
		if err := json.NewEncoder(writer).Encode(&flag); err != nil {
			fmt.Println(err)
		}
	} else { // succeeds

		// TODO: get real flag from data source
		fmt.Println("query by port :", port)

		var flag swaper.FlagByPort
		flag.Msg = "success"
		flag.Flag = "flag{test_flag_for_flag_proxy}"
		if err := json.NewEncoder(writer).Encode(&flag); err != nil {
			fmt.Println(err)
		}
	}
}
