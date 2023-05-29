package controller

import (
	"UserAPI/models"
	"UserAPI/service"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type IUserController interface {
	SaveUserProfile(c *fiber.Ctx) error
	SaveUserImage(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type userController struct {
	service service.IUserProfileService
}

var userControllerObj IUserController

func InitUserController(serviceObj service.IUserProfileService) IUserController {
	if userControllerObj == nil {
		userControllerObj = &userController{
			service: serviceObj,
		}
	}
	return userControllerObj
}

func (ctlr *userController) SaveUserProfile(c *fiber.Ctx) error {
	detail := new(models.UserDetail)
	//get request param
	err := c.BodyParser(detail)
	if err != nil {
		return err
	}

	result, err := ctlr.service.SaveUserProfile(detail)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (ctrl *userController) SaveUserImage(c *fiber.Ctx) error {
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

	ctrl.service.SaveImagetoAWS(src, filename, file.Size)

	return c.JSON(fiber.Map{"message": "Image uploaded successfully"})
}

func (ctrl *userController) GetUserProfile(c *fiber.Ctx) error {
	email := c.Params("userId")

	result, err := ctrl.service.GetUserProfile(email)
	if err != nil {
		return err
	}

	result.ImagePath, err = ctrl.service.GetUserImageURL(result.Id + ".png")
	if err != nil {
		return err
	}

	return c.JSON(result)
}
