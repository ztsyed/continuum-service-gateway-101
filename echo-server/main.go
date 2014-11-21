package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/echo", EchoServer)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// This handler just returns echo back the path of the incoming request
// in json format with timestamp and host information in the header
func EchoServer(w http.ResponseWriter, req *http.Request) {
	//Echo Back query string
	responseStruct := struct {
		Response  string    `json:"response"`
		Status    string    `json:"status"`
		TimeStamp time.Time `json:"timestamp"`
	}{
		"Echo Server says: " + req.URL.String(),
		"OK",
		time.Now(),
	}
	responseJson, _ := json.Marshal(responseStruct)
	io.WriteString(w, string(responseJson))
	w.Header().Set("Host", req.RemoteAddr+":")
}
