package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jobcrawler.api/config"
	"jobcrawler.api/models"
	"jobcrawler.api/repository/connection"
	"jobcrawler.api/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

//#region init db connection
//#endregin

//#region init elastic connection
//#endregin

//#region init cache
//#endregin

//#region init repo objects passing the required connection
//#endregin

//#region init service objects using repo objects and cache
//#endregin

//#region init controller objects using service ojects
//#endregin

//#region init router using controller
//#endregin

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

		saveImagetoAWS(src, filename, file.Size)

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

		result.ImagePath, err = getUserImageURL(result.Id + ".png")
		if err != nil {
			return err
		}

		return c.JSON(result)
	})

	log.Fatal(app.Listen(":8080"))
}

func saveImagetoAWS(_file multipart.File, fileName string, size int64) {
	// Specify your AWS region and S3 bucket name
	//region := "Asia Pacific (Mumbai) ap-south-1"
	bucketName := "jobcrawler.portalimages"

	// Create an AWS session
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(region)},
	// )
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an S3 client
	svc := s3.New(sess)

	// Create an S3 object with the specified bucket and key (filename)
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          _file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("image/jpeg"), // Specify the correct content type
	})
	if err != nil {
		fmt.Println("Failed to upload image:", err)
		return
	}

	fmt.Println("Image uploaded successfully!")
}

func getUserImageURL(filename string) (string, error) {
	bucketName := "jobcrawler.portalimages"

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an S3 client
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(time.Duration(2 * time.Hour)) // Set link expiration time
	if err != nil {
		return "", err
	}
	return url, nil
}
