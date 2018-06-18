package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"./syntax"
)

type Message struct {
	Msg string `json:"msg,omitempty"`
}

func runParser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// get code
		body, _ := ioutil.ReadAll(r.Body)
		code := string(body[:])
		// make message
		var message Message
		message.Msg = syntax.Validate(code)

		json.NewEncoder(w).Encode(&message)
		return
	}
	http.Redirect(w, r, "/", 200)
}

func main() {
	port := ":3000"
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/parser", runParser)
	log.Println("Server on port 3000")
	log.Fatal(http.ListenAndServe(port, nil))
}
