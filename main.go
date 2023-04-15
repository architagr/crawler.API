package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"users/gauravagarwal/CRAWLER.API/repository/config"
	"users/gauravagarwal/CRAWLER.API/repository/connection"

	"github.com/gorilla/mux"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var conn connection.IConnection
var env *config.Config

type Message struct {
	To      string `json:"+919950458542"`
	From    string `json:"test"`
	Message string `json:"message"`
}

func setupDB() {
	env = config.GetConfig()
	conn = connection.InitConnection(env.GetDatabaseConnectionString(), 10)
	err := conn.ValidateConnection()
	if err != nil {
		log.Fatalf("error in conncting to mongo %+v", err)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, p *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// r.HandleFunc("/set", func(w http.ResponseWriter, p *http.Request) {
	// 	query := p.URL.Query()
	// 	data, present := query["id"]
	// 	if !present || len(data) == 0 {
	// 		fmt.Println("data not present")
	// 	}
	// 	cacheservice.Set(data[0])
	// })

	// r.HandleFunc("/get", func(w http.ResponseWriter, p *http.Request) {
	// 	var res = cacheservice.Get("1")
	// 	fmt.Fprint(w, res)
	// })

	r.HandleFunc("/sms", func(w http.ResponseWriter, p *http.Request) {
		fmt.Println("SMS API initiated")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//get request param
		body, err := ioutil.ReadAll(p.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer p.Body.Close()
		phoneNo := string(body)
		fmt.Println("mobileno: " + phoneNo)

		accountSaid := "AC67b72d1e3041953fc3fb10da36a0c5a0"
		authToken := "92a71371a777646a65d102c469f70bb4"
		client := twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSaid,
			Password: authToken,
		})
		randno := rand.Intn(9999)
		params := &openapi.CreateMessageParams{}
		params.SetTo(phoneNo)
		params.SetFrom("+14344045914")
		params.SetBody("OTP received to login JobPortal is " + strconv.Itoa(randno))

		_, err1 := client.Api.CreateMessage(params)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("SMS sent successfully")
		}

		json.NewEncoder(w).Encode("Message sent successfully")
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, p *http.Request) {
		fmt.Println("Process for login initiated")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		//get request param
		body, err := ioutil.ReadAll(p.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer p.Body.Close()
		otp := string(body)
		fmt.Printf("OTP received: " + otp)

		//validate OTP
		fmt.Fprint(w, true)
	}).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
