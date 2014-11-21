package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// This handler just returns echo back the path of the incoming request
// in json format with timestamp and host information in the header
func EchoServer(w http.ResponseWriter, req *http.Request) {
	// See if we got authorization headers
	if req.Header["Authorization"] != nil {
		auth := strings.SplitN(req.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		log.Println("Request Payload ", string(payload))
	}

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
func main() {
	http.HandleFunc("/", EchoServer)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
