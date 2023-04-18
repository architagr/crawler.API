package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/connection"
	"jobcrawler.api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	_ "jobcrawler.api/middleware"
)

var conn connection.IConnection
var env *config.Config

type Message struct {
	To      string `json:"+919950458542"`
	From    string `json:"test"`
	Message string `json:"message"`
}

func setupDB() {
	var err error
	conn, err = connection.InitConnection(env.GetDatabaseConnectionString(), 10)
	if err != nil {
		log.Fatalf("error in conncting to mongo %+v", err)
	}
	err = conn.ValidateConnection()
	if err != nil {
		log.Fatalf("error in conncting to mongo %+v", err)
	}
}

func main() {
	env = config.GetConfig()
	setupDB()
	//defer conn.Disconnect()
	app := fiber.New()

	app.Get("/test/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString(fmt.Sprintf("Hello, %s!", name))
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

	app.Get("/getJobs", func(c *fiber.Ctx) error {
		var pageSize, pageNumber int64 = 10, 0 // todo: get this from querystring

		jobservice, err := service.GetJobServiceObj()
		if err != nil {
			return err
		}
		response, err := jobservice.GetJobs(pageSize, pageNumber)
		if err != nil {
			return err
		}
		return c.JSON(response)
	})

	app.Post("/sms", func(c *fiber.Ctx) error {
		fmt.Println("SMS API initiated")

		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//get request param

		phoneNo := ""
		if err := c.BodyParser(phoneNo); err != nil {
			return err
		}
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
			return err1
		} else {
			fmt.Println("SMS sent successfully")
		}

		return c.JSON("Message sent successfully")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		fmt.Println("Process for login initiated")
		detail := new(models.LoginDetails)
		//get request param
		err := c.BodyParser(detail)
		if err != nil {
			return err
		}
		jsonStr, err := json.Marshal(detail)
		fmt.Printf("OTP received: " + string(jsonStr))

		authService, err := service.AuthServiceObj()
		token, err := authService.GetJWT(detail)

		userService, err := service.UserServiceObj()
		result, err := userService.Login(detail)
		fmt.Printf("user created: " + result.(string))
		//validate OTP
		return c.JSON(token)
	})

	log.Fatal(app.Listen(":8080"))
}
