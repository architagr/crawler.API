package main

import (
	"jobcrawler.api/config"
	"jobcrawler.api/repository"
	"jobcrawler.api/service"
)

var envVariables config.IConfig
var jobDetailsRepoObj repository.IJobDetailsRepository
var jobDetailsService service.IJobService

func main() {
	initConfig()
	initRepository()
	intitServices()
}
func initConfig() {
	envVariables = config.GetConfig()
}
func initRepository() {
	mongodbConnection, err := repository.InitConnection(envVariables.GetDatabaseConnectionString(), 10)
	if err != nil {
		panic(err)
	}

	jobDetailsRepoObj, err = repository.InitJobDetailsRepo(mongodbConnection, envVariables.GetDatabaseName(), envVariables.GetCollectionName())
	if err != nil {
		panic(err)
	}
}

func intitServices() {
	jobDetailsService = service.InitJobService(jobDetailsRepoObj)
}

/*


import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/connection"
	"jobcrawler.api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

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

	app.Post("/getJobs", func(c *fiber.Ctx) error {
		//var pageSize, pageNumber int64 = 10, 0 // todo: get this from querystring

		//get params from body
		filter := new(models.JobFilter)
		if err := c.BodyParser(filter); err != nil {
			return err
		}
		_filter, err := json.Marshal(filter)
		fmt.Println("filter: " + string(_filter))
		jobservice, err := service.GetJobServiceObj()
		if err != nil {
			return err
		}
		response, err := jobservice.GetJobs(filter, filter.PageSize, filter.PageNumber)
		if err != nil {
			return err
		}
		return c.JSON(response)
	})

	app.Get("/getJobDetail/:id", func(c *fiber.Ctx) error {
		//get params from body
		jobId := c.Params("id")
		jobservice, err := service.GetJobServiceObj()
		if err != nil {
			return err
		}
		response, err := jobservice.GetJobDetail(jobId)
		if err != nil {
			return err
		}
		return c.JSON(response)
	})

	app.Post("/sms", func(c *fiber.Ctx) error {
		fmt.Println("SMS API initiated")

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

		userService, err := service.UserServiceObj("userDetails")
		result, err := userService.Login(detail)
		fmt.Printf("user created: " + strconv.FormatBool(result))
		var res = new(models.Response)
		if result {
			res.Status = "OK"
			res.Data = token
		} else {
			res.Status = "Failed"
			res.Data = "Login Failed"
		}
		return c.JSON(res)
	})

	app.Post("/saveUserProfile", func(c *fiber.Ctx) error {
		detail := new(models.UserDetail)
		//get request param
		err := c.BodyParser(detail)
		if err != nil {
			return err
		}

		userService, err := service.UserProfileServiceObj("userProfile")
		result, err := userService.SaveUserProfile(detail)
		if err != nil {
			return err
		}
		return c.JSON(result)
	})

	app.Post("/saveuserImage/:id", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to retrieve the image"})
		}

		files := form.File["image"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No image file found"})
		}

		file := files[0]
		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open the image"})
		}
		defer src.Close()

		userId := c.Params("id")

		// Create a unique filename
		ext := filepath.Ext(file.Filename)
		filename := userId + ext

		// Specify the folder path to save the uploaded file
		currentDir, err := os.Getwd()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get current directory"})
		}

		// Specify the folder path to save the uploaded file
		uploadDir := filepath.Join(currentDir, "uploads")
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create upload directory"})
		}

		filePath := filepath.Join(uploadDir, filename)

		// Create the destination file
		dst, err := os.Create(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save the image"})
		}
		defer dst.Close()

		// Copy the uploaded file to the destination
		if _, err := io.Copy(dst, src); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save the image"})
		}

		return c.JSON(fiber.Map{"message": "Image uploaded successfully"})
	})

	app.Get("/image/:id", func(c *fiber.Ctx) error {
		imagePath, err := os.Getwd() // Update with the actual path to your saved image
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get current directory"})
		}

		uploadDir := filepath.Join(imagePath, "uploads")

		filename := c.Params("id") + ".png"
		filePath := filepath.Join(uploadDir, filename)

		// Check if the image file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("Image not found")
		}

		// Get the file's content type
		contentType := "image/jpeg" // Update with the appropriate content type of your image

		// Set the response headers
		c.Set(fiber.HeaderContentType, contentType)
		c.Set(fiber.HeaderCacheControl, "max-age=31536000") // Optional: Cache the image for a year

		// Send the image file as the response
		return c.SendFile(filePath)
	})

	app.Get("/getUserProfile/:userId", func(c *fiber.Ctx) error {
		email := c.Params("userId")

		userService, err := service.UserProfileServiceObj("userProfile")
		if err != nil {
			return err
		}
		result, err := userService.GetUserProfile(email)
		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	log.Fatal(app.Listen(":8080"))
}
*/
