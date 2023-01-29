package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Event struct {
	Id    string
	AppId string
	Name  string
	Data  map[string]any
}

type EventReqBody struct {
	AppId string
	Name  string
	Data  map[string]any
}

type EventStore []Event

var store = EventStore{}

func main() {
	http.HandleFunc("/", HandleEvents)
	http.ListenAndServe(":8080", nil)
}

/*
*	Handle events should take in post requests and publish the events
* to a Kafka broker. As for now they are simply stored in `store`.
 */
func HandleEvents(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	// GET will become unnecessary as this interface is for handling POSTs
	case "GET":
		resp := store
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error occured in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	case "POST":
		var body EventReqBody
		if r.Body == nil {
			http.Error(w, "A request body must be sent", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		event := Event{Id: uuid.New().String(), AppId: body.AppId, Name: body.Name, Data: body.Data}
		store = append(store, event)
		jsonResp, err := json.Marshal(event)
		if err != nil {
			log.Fatalf("Error occured in JSON marshal. Err: %s", err)
		}
		fmt.Printf("jsonResp: %s", string(jsonResp))
		w.Write(jsonResp)
	}
}
