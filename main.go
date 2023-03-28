package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Message struct {
	To      string `json:"+919950458542"`
	From    string `json:"test"`
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.HandleFunc("/sms", func(w http.ResponseWriter, p *http.Request) {
		fmt.Println("SMS API initiated")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		accountSaid := "AC67b72d1e3041953fc3fb10da36a0c5a0"
		authToken := "f5f3e764930b0224d35f1b8439ba9a7a"
		client := twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSaid,
			Password: authToken,
		})

		params := &openapi.CreateMessageParams{}
		params.SetTo("+919950458542")
		params.SetFrom("+14344045914")
		params.SetBody("Message")

		_, err := client.Api.CreateMessage(params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("SMS sent successfully")
		}

		json.NewEncoder(w).Encode("Message sent successfully")
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
